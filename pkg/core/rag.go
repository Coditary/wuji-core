package core

import (
	"context"
	"fmt"
	"strings"

	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/data"
	"github.com/coditary/wuji-core/pkg/driver"
)

// RunRAG executes a RAG operation and returns WDD output.
func (c *Core) RunRAG(ctx context.Context, driverID string, req driver.RAGRequest) (data.DataShape, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if req.Task == driver.RAGTaskAnswer {
		return scheduleGenerate(c, ctx, driverID, req.TextModel, capability.RAG, func(ctx context.Context) (data.DataShape, error) {
			return c.runRAGAnswer(ctx, driverID, req)
		})
	}
	return scheduleGenerate(c, ctx, driverID, req.EmbedModel, capability.RAG, func(ctx context.Context) (data.DataShape, error) {
		d, err := c.resolveForCapability(driverID, capability.RAG)
		if err != nil {
			return nil, err
		}
		ctx = c.ragContext(ctx, driverID)
		return driver.RunRAGTask(ctx, d, req)
	})
}

func (c *Core) runRAGAnswer(ctx context.Context, driverID string, req driver.RAGRequest) (data.DataShape, error) {
	ragD, err := c.resolveForCapability(driverID, capability.RAG)
	if err != nil {
		return nil, err
	}
	qreq := req
	qreq.Task = driver.RAGTaskQuery
	ctx = c.ragContext(ctx, driverID)
	shape, err := driver.RunRAGTask(ctx, ragD, qreq)
	if err != nil {
		return nil, err
	}
	hits, err := driver.RAGQueryFromShape(shape)
	if err != nil {
		return nil, err
	}
	contextBlock := formatRAGContext(hits.Hits, req.ContextMaxChars)
	textResp, err := c.GenerateText(ctx, c.ResolveDriverID(driverID, capability.TextGeneration), driver.TextRequest{
		Prompt:       req.Query,
		SystemPrompt: buildRAGSystemPrompt(contextBlock, req.SystemPrompt, req.Cite),
		Model:        req.TextModel,
		MaxTokens:    req.MaxTokensOrDefault(),
		Temperature:  req.Temperature,
		TopP:         req.TopP,
	})
	if err != nil {
		if ans, ok := ragD.(driver.RAGAnswerer); ok {
			resp, ansErr := ans.Answer(ctx, req)
			if ansErr == nil {
				return driver.RAGAnswerToRecordSet(resp), nil
			}
		}
		return nil, err
	}
	return driver.RAGAnswerToRecordSet(&driver.RAGAnswerResponse{
		Answer: textResp.Text, Collection: req.Collection, Query: req.Query,
		Hits: hits.Hits, TokensUsed: textResp.TokensUsed, NoContext: req.NoContext,
	}), nil
}

func buildRAGSystemPrompt(contextBlock, extra string, cite bool) string {
	base := "Answer using only the retrieved context below. If the answer is not in the context, say you do not know.\n\n" + contextBlock
	if cite {
		base += "\n\nCite supporting passages using [1], [2], ... matching the chunk numbers in the context."
	}
	if strings.TrimSpace(extra) == "" {
		return base
	}
	return extra + "\n\n" + base
}

func formatRAGContext(hits []driver.RAGHit, maxChars int) string {
	var b strings.Builder
	for i, hit := range hits {
		fmt.Fprintf(&b, "[%d] (score=%.3f source=%s)\n%s\n\n", i+1, hit.Score, hit.Source, hit.Text)
		if maxChars > 0 && b.Len() >= maxChars {
			return b.String()[:maxChars]
		}
	}
	out := b.String()
	if maxChars > 0 && len(out) > maxChars {
		return out[:maxChars]
	}
	return out
}
