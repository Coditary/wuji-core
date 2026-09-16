package driver

import (
	wujiv1 "github.com/coditary/wuji-core/api/proto/v1"
	"github.com/coditary/wuji-core/pkg/data"
)

func RAGTasksToProto(tasks []RAGTask) []string {
	if len(tasks) == 0 {
		return nil
	}
	out := make([]string, 0, len(tasks))
	for _, task := range tasks {
		out = append(out, string(task))
	}
	return out
}

func RAGTasksFromProto(items []string) []RAGTask {
	if len(items) == 0 {
		return nil
	}
	out := make([]RAGTask, 0, len(items))
	for _, item := range items {
		task, err := ParseRAGTask(item)
		if err != nil {
			continue
		}
		out = append(out, task)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func RAGRequestToProto(req RAGRequest) *wujiv1.RAGRequestMsg {
	return &wujiv1.RAGRequestMsg{
		Task: req.Task.String(), StoreRoot: req.StoreRoot, Collection: req.Collection, Query: req.Query,
		SourcePaths: req.SourcePaths, UseStdin: req.UseStdin, Recursive: req.Recursive,
		ChunkSize: int32(req.ChunkSize), ChunkOverlap: int32(req.ChunkOverlap),
		EmbedModel: req.EmbedModel, TextModel: req.TextModel, SystemPrompt: req.SystemPrompt,
		TopK: int32(req.TopK), MinScore: req.MinScore, MaxTokens: int32(req.MaxTokens),
		PurgeSource: req.PurgeSource,
		Filter:      req.Filter, IndexMetadata: req.IndexMetadata,
		IndexMode: string(req.IndexMode), Force: req.Force, DryRun: req.DryRun,
		Glob: req.Glob, Exclude: req.Exclude, Url: req.URL,
		QueryFile: req.QueryFile, QueryStdin: req.QueryStdin,
		RerankTop: int32(req.RerankTop), Diverse: req.Diverse,
		ContextMaxChars: int32(req.ContextMaxChars), Cite: req.Cite,
		Temperature: req.Temperature, TopP: req.TopP,
		NoContext: req.NoContext, IncludeScores: req.IncludeScores,
		RenameTo: req.RenameTo, ExportPath: req.ExportPath, ImportPath: req.ImportPath,
	}
}

func RAGRequestFromProto(p *wujiv1.RAGRequestMsg) (RAGRequest, error) {
	if p == nil {
		return RAGRequest{}, nil
	}
	task, err := ParseRAGTask(p.GetTask())
	if err != nil {
		return RAGRequest{}, err
	}
	return RAGRequest{
		Task: task, StoreRoot: p.GetStoreRoot(), Collection: p.GetCollection(), Query: p.GetQuery(),
		SourcePaths: p.GetSourcePaths(), UseStdin: p.GetUseStdin(), Recursive: p.GetRecursive(),
		ChunkSize: int(p.GetChunkSize()), ChunkOverlap: int(p.GetChunkOverlap()),
		EmbedModel: p.GetEmbedModel(), TextModel: p.GetTextModel(), SystemPrompt: p.GetSystemPrompt(),
		TopK: int(p.GetTopK()), MinScore: p.GetMinScore(), MaxTokens: int(p.GetMaxTokens()),
		PurgeSource: p.GetPurgeSource(),
		Filter:      cloneStringMap(p.GetFilter()), IndexMetadata: cloneStringMap(p.GetIndexMetadata()),
		IndexMode: RAGIndexMode(p.GetIndexMode()), Force: p.GetForce(), DryRun: p.GetDryRun(),
		Glob: p.GetGlob(), Exclude: append([]string(nil), p.GetExclude()...), URL: p.GetUrl(),
		QueryFile: p.GetQueryFile(), QueryStdin: p.GetQueryStdin(),
		RerankTop: int(p.GetRerankTop()), Diverse: p.GetDiverse(),
		ContextMaxChars: int(p.GetContextMaxChars()), Cite: p.GetCite(),
		Temperature: p.GetTemperature(), TopP: p.GetTopP(),
		NoContext: p.GetNoContext(), IncludeScores: p.GetIncludeScores(),
		RenameTo: p.GetRenameTo(), ExportPath: p.GetExportPath(), ImportPath: p.GetImportPath(),
	}, nil
}

func cloneStringMap(in map[string]string) map[string]string {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func (t RAGTask) String() string { return string(t) }

func RAGShapeToProto(shape data.DataShape) (*wujiv1.RAGResponseMsg, error) {
	if shape == nil {
		return &wujiv1.RAGResponseMsg{}, nil
	}
	payload, err := data.MarshalMsgpack(shape)
	if err != nil {
		return nil, err
	}
	return &wujiv1.RAGResponseMsg{
		Payload: payload, PayloadFormat: DataPayloadMsgpack, Kind: string(shape.Kind()),
	}, nil
}

func RAGShapeFromProto(p *wujiv1.RAGResponseMsg) (data.DataShape, error) {
	if p == nil {
		return nil, nil
	}
	format := p.GetPayloadFormat()
	if format == "" {
		format = DataPayloadMsgpack
	}
	return dataShapeFromPayload(p.GetPayload(), format)
}
