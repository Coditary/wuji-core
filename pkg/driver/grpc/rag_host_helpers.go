package grpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/data"
	"github.com/coditary/wuji-core/pkg/driver"
	"github.com/coditary/wuji-core/pkg/rag"
)

func runRAGOnHost(ctx context.Context, d driver.Driver, req driver.RAGRequest) (data.DataShape, error) {
	ctx = rag.WithEmbedder(ctx, rag.DriverEmbedder(d))
	if req.Task == driver.RAGTaskAnswer {
		return runRAGAnswerOnHost(ctx, d, req)
	}
	return driver.RunRAGTask(ctx, d, req)
}

func runRAGAnswerOnHost(ctx context.Context, d driver.Driver, req driver.RAGRequest) (data.DataShape, error) {
	qreq := req
	qreq.Task = driver.RAGTaskQuery
	shape, err := driver.RunRAGTask(ctx, d, qreq)
	if err != nil {
		return nil, err
	}
	hits, err := driver.RAGQueryFromShape(shape)
	if err != nil {
		return nil, err
	}
	gen, err := driver.As[driver.TextGenerator](d, capability.TextGeneration)
	if err != nil {
		if ans, ok := d.(driver.RAGAnswerer); ok {
			resp, ansErr := ans.Answer(ctx, req)
			if ansErr == nil {
				return driver.RAGAnswerToRecordSet(resp), nil
			}
		}
		return nil, fmt.Errorf("answer requires text capability on driver host: %w", err)
	}
	contextBlock := formatHostRAGContext(hits.Hits, req.SystemPrompt, req.Cite)
	textResp, err := gen.GenerateText(ctx, driver.TextRequest{
		Prompt:       req.Query,
		SystemPrompt: contextBlock,
		Model:        req.TextModel,
		MaxTokens:    req.MaxTokensOrDefault(),
	})
	if err != nil {
		if ans, ok := d.(driver.RAGAnswerer); ok {
			resp, ansErr := ans.Answer(ctx, req)
			if ansErr == nil {
				return driver.RAGAnswerToRecordSet(resp), nil
			}
		}
		return nil, err
	}
	return driver.RAGAnswerToRecordSet(&driver.RAGAnswerResponse{
		Answer: textResp.Text, Collection: req.Collection, Query: req.Query,
		Hits: hits.Hits, TokensUsed: textResp.TokensUsed,
	}), nil
}

func formatHostRAGContext(hits []driver.RAGHit, extra string, cite bool) string {
	var b strings.Builder
	b.WriteString("Answer using only the retrieved context below. If the answer is not in the context, say you do not know.\n\n")
	for i, hit := range hits {
		fmt.Fprintf(&b, "[%d] (score=%.3f source=%s)\n%s\n\n", i+1, hit.Score, hit.Source, hit.Text)
	}
	if cite {
		b.WriteString("Cite supporting passages using [1], [2], ... matching the chunk numbers in the context.\n")
	}
	if strings.TrimSpace(extra) != "" {
		return extra + "\n\n" + b.String()
	}
	return b.String()
}
