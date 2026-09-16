package rag

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/coditary/wuji-core/pkg/driver"
	"github.com/coditary/wuji-core/pkg/ragstore"
)

// Index ingests documents into a collection with real or fallback embeddings.
func Index(ctx context.Context, req driver.RAGRequest) (*driver.RAGIndexResponse, error) {
	ctx = ragStoreCtx(ctx, req)
	store := ragstore.ResolveBackend(ctx)
	docs, err := loadDocs(req)
	if err != nil {
		return nil, err
	}
	if len(docs) == 0 {
		return nil, fmt.Errorf("no documents to index")
	}
	chunkSize := req.ChunkSize
	if chunkSize <= 0 {
		chunkSize = 800
	}
	newChunks := ragstore.SplitDocs(docs, chunkSize, req.ChunkOverlap)
	if len(newChunks) == 0 {
		return nil, fmt.Errorf("no chunks produced from input")
	}
	applyIndexMetadata(newChunks, req.IndexMetadata)

	var vecs [][]float32
	dims := 0
	if !req.DryRun {
		emb := ResolveEmbedder(ctx)
		texts := make([]string, len(newChunks))
		for i, c := range newChunks {
			texts[i] = c.Text
		}
		var err error
		vecs, dims, err = emb(ctx, req.EmbedModel, texts)
		if err != nil {
			return nil, fmt.Errorf("embed chunks: %w", err)
		}
		if err := ragstore.AttachEmbeddings(newChunks, vecs, dims); err != nil {
			return nil, err
		}
	}

	incomingSources := uniqueDocSources(docs)
	var kept []ragstore.Chunk
	if store.Exists(ctx, req.StoreRoot, req.Collection) && !req.DryRun {
		col, err := store.Load(ctx, req.StoreRoot, req.Collection)
		if err != nil {
			return nil, err
		}
		if err := checkEmbedCompatibility(col.Manifest, req.EmbedModel, dims, req.Force); err != nil {
			return nil, err
		}
		switch req.IndexModeOrDefault() {
		case driver.RAGIndexAppend:
			kept = col.Chunks
		default:
			kept = ragstore.ChunksWithoutSources(col.Chunks, incomingSources)
		}
	}

	chunks := append(kept, newChunks...)
	totalChunks := len(chunks)
	var bytes int64
	for _, doc := range docs {
		bytes += int64(len(doc.Content))
	}

	resp := &driver.RAGIndexResponse{
		Collection: req.Collection, ChunksAdded: len(newChunks), TotalChunks: totalChunks,
		BytesIndexed: bytes, Sources: ragstore.MergeSources(kept, incomingSources),
		DryRun: req.DryRun,
	}
	if req.DryRun {
		resp.Message = fmt.Sprintf("dry-run: would index %d new chunks (%d total)", len(newChunks), totalChunks)
		return resp, nil
	}

	vectors := ragstore.FlattenEmbeddings(chunks)
	manifest := ragstore.Manifest{
		Collection: req.Collection,
		EmbedModel: req.EmbedModel,
		Dims:       dims,
		ChunkSize:  chunkSize,
		Overlap:    req.ChunkOverlap,
		Sources:    resp.Sources,
	}
	if err := store.Save(ctx, req.StoreRoot, manifest, chunks, vectors); err != nil {
		return nil, err
	}
	resp.Message = fmt.Sprintf("indexed %d new chunks (%d total)", len(newChunks), totalChunks)
	return resp, nil
}

// Query retrieves top matching chunks for a query.
func Query(ctx context.Context, req driver.RAGRequest) (*driver.RAGQueryResponse, error) {
	ctx = ragStoreCtx(ctx, req)
	store := ragstore.ResolveBackend(ctx)
	col, err := store.Load(ctx, req.StoreRoot, req.Collection)
	if err != nil {
		return nil, err
	}
	dims := col.Manifest.Dims
	if dims <= 0 {
		return nil, fmt.Errorf("collection %q has no embedding dims", req.Collection)
	}
	emb := ResolveEmbedder(ctx)
	vecs, gotDims, err := emb(ctx, req.EmbedModel, []string{req.Query})
	if err != nil {
		return nil, fmt.Errorf("embed query: %w", err)
	}
	if len(vecs) == 0 {
		return nil, fmt.Errorf("empty query embedding")
	}
	if gotDims != dims {
		return nil, fmt.Errorf("query embedding dims %d != collection dims %d (embed model changed?)", gotDims, dims)
	}
	pairs := ragstore.SelectTopK(vecs[0], col.Vectors, dims, req.SearchK(), req.FinalTopK(), req.Diverse)
	hits := pairsToHits(col.Chunks, pairs, req)
	return &driver.RAGQueryResponse{Query: req.Query, Collection: req.Collection, Hits: hits}, nil
}

func pairsToHits(chunks []ragstore.Chunk, pairs []ragstore.ScorePair, req driver.RAGRequest) []driver.RAGHit {
	hits := make([]driver.RAGHit, 0, len(pairs))
	for _, p := range pairs {
		if p.Index >= len(chunks) {
			continue
		}
		if req.MinScore > 0 && p.Score < req.MinScore {
			continue
		}
		if !matchFilter(chunks[p.Index], req.Filter) {
			continue
		}
		c := chunks[p.Index]
		hits = append(hits, driver.RAGHit{
			ID: c.ID, Score: p.Score, Text: c.Text, Source: c.Source, ChunkIdx: c.ChunkIdx, Metadata: c.Metadata,
		})
	}
	return hits
}

func matchFilter(chunk ragstore.Chunk, filter map[string]string) bool {
	if len(filter) == 0 {
		return true
	}
	for k, want := range filter {
		got := chunk.Metadata[k]
		if got == "" && k == "source" {
			got = chunk.Source
		}
		if got != want {
			return false
		}
	}
	return true
}

func ListCollections(ctx context.Context, req driver.RAGRequest) ([]driver.RAGCollectionInfo, error) {
	ctx = ragStoreCtx(ctx, req)
	store := ragstore.ResolveBackend(ctx)
	items, err := store.ListCollections(ctx, req.StoreRoot)
	if err != nil {
		return nil, err
	}
	out := make([]driver.RAGCollectionInfo, len(items))
	for i, m := range items {
		out[i] = manifestToInfo(m)
	}
	return out, nil
}

func CollectionInfo(ctx context.Context, req driver.RAGRequest) (driver.RAGCollectionInfo, error) {
	ctx = ragStoreCtx(ctx, req)
	store := ragstore.ResolveBackend(ctx)
	col, err := store.Load(ctx, req.StoreRoot, req.Collection)
	if err != nil {
		return driver.RAGCollectionInfo{}, err
	}
	return manifestToInfo(col.Manifest), nil
}

func DeleteCollection(ctx context.Context, req driver.RAGRequest) error {
	ctx = ragStoreCtx(ctx, req)
	return ragstore.ResolveBackend(ctx).Delete(ctx, req.StoreRoot, req.Collection)
}

func PurgeSource(ctx context.Context, req driver.RAGRequest) error {
	ctx = ragStoreCtx(ctx, req)
	store := ragstore.ResolveBackend(ctx)
	col, err := store.Load(ctx, req.StoreRoot, req.Collection)
	if err != nil {
		return err
	}
	filtered := ragstore.ChunksWithoutSources(col.Chunks, []string{req.PurgeSource})
	if len(filtered) == len(col.Chunks) {
		return fmt.Errorf("source %q not found in collection", req.PurgeSource)
	}
	col.Manifest.Sources = ragstore.UniqueSources(filtered)
	vectors := ragstore.FlattenEmbeddings(filtered)
	return store.Save(ctx, req.StoreRoot, col.Manifest, filtered, vectors)
}

func ExportCollection(ctx context.Context, req driver.RAGRequest) error {
	ctx = ragStoreCtx(ctx, req)
	return ragstore.ResolveBackend(ctx).ExportCollection(ctx, req.StoreRoot, req.Collection, req.ExportPath)
}

func ImportCollection(ctx context.Context, req driver.RAGRequest) error {
	ctx = ragStoreCtx(ctx, req)
	return ragstore.ResolveBackend(ctx).ImportCollection(ctx, req.StoreRoot, req.Collection, req.ImportPath)
}

func RenameCollection(ctx context.Context, req driver.RAGRequest) error {
	ctx = ragStoreCtx(ctx, req)
	return ragstore.ResolveBackend(ctx).RenameCollection(ctx, req.StoreRoot, req.Collection, req.RenameTo)
}

func ragStoreCtx(ctx context.Context, req driver.RAGRequest) context.Context {
	if dir := strings.TrimSpace(req.StoreDir); dir != "" {
		return ragstore.WithStoreDir(ctx, dir)
	}
	return ctx
}

func loadDocs(req driver.RAGRequest) ([]ragstore.SourceDoc, error) {
	if req.UseStdin {
		doc, err := ragstore.ReadStdin(os.Stdin)
		if err != nil {
			return nil, err
		}
		return []ragstore.SourceDoc{doc}, nil
	}
	return ragstore.ReadSources(req.SourcePaths, ragstore.ReadOptions{
		Recursive: req.Recursive,
		Glob:      req.Glob,
		Exclude:   req.Exclude,
	})
}

func applyIndexMetadata(chunks []ragstore.Chunk, meta map[string]string) {
	if len(meta) == 0 {
		return
	}
	for i := range chunks {
		if chunks[i].Metadata == nil {
			chunks[i].Metadata = map[string]string{}
		}
		for k, v := range meta {
			chunks[i].Metadata[k] = v
		}
	}
}

func checkEmbedCompatibility(m ragstore.Manifest, model string, dims int, force bool) error {
	if force {
		return nil
	}
	if m.Dims > 0 && dims > 0 && m.Dims != dims {
		return fmt.Errorf("embedding dims mismatch: collection has %d, new embeddings have %d (use --force or --delete)", m.Dims, dims)
	}
	if m.EmbedModel != "" && model != "" && !strings.EqualFold(m.EmbedModel, model) {
		return fmt.Errorf("embed model mismatch: collection uses %q, request uses %q (use --force or --delete)", m.EmbedModel, model)
	}
	return nil
}

func uniqueDocSources(docs []ragstore.SourceDoc) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, d := range docs {
		if _, ok := seen[d.Path]; ok {
			continue
		}
		seen[d.Path] = struct{}{}
		out = append(out, d.Path)
	}
	return out
}

func manifestToInfo(m ragstore.Manifest) driver.RAGCollectionInfo {
	updated := ""
	if !m.UpdatedAt.IsZero() {
		updated = m.UpdatedAt.UTC().Format(time.RFC3339)
	}
	return driver.RAGCollectionInfo{
		Collection: m.Collection, EmbedModel: m.EmbedModel, Dims: m.Dims,
		ChunkCount: m.ChunkCount, ChunkSize: m.ChunkSize, Overlap: m.Overlap,
		Sources: m.Sources, UpdatedAt: updated,
	}
}
