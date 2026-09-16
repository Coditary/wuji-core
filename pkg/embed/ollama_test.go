package embed_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/coditary/wuji-core/pkg/data"
	"github.com/coditary/wuji-core/pkg/driver"
	"github.com/coditary/wuji-core/pkg/embed"
)

func TestOllamaProduceDataEmbed(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/embed" {
			http.NotFound(w, r)
			return
		}
		var req struct {
			Input []string `json:"input"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatal(err)
		}
		vecs := make([][]float64, len(req.Input))
		for i := range req.Input {
			vecs[i] = []float64{0.1, 0.2, 0.3}
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"embeddings": vecs})
	}))
	t.Cleanup(srv.Close)

	prod := embed.NewOllama(srv.URL, "nomic-embed-text", "test")
	shape, err := prod.ProduceData(context.Background(), driver.DataRequest{
		Task: driver.DataTaskEmbed, Texts: []string{"hello", "world"},
	})
	if err != nil {
		t.Fatal(err)
	}
	batch, ok := shape.(*data.VectorBatch)
	if !ok {
		t.Fatalf("got %T", shape)
	}
	if batch.Count() != 2 || batch.Dims != 3 {
		t.Fatalf("count=%d dims=%d", batch.Count(), batch.Dims)
	}
}

func TestOllamaProduceDataFallback(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unavailable", http.StatusServiceUnavailable)
	}))
	t.Cleanup(srv.Close)

	prod := embed.NewOllama(srv.URL, "nomic-embed-text", "test")
	shape, err := prod.ProduceData(context.Background(), driver.DataRequest{
		Task: driver.DataTaskEmbed, Text: "fallback",
	})
	if err != nil {
		t.Fatal(err)
	}
	batch, ok := shape.(*data.VectorBatch)
	if !ok {
		t.Fatalf("got %T", shape)
	}
	if batch.MetaData.Extra["embed_fallback"] != "pseudo" {
		t.Fatalf("expected pseudo fallback, meta=%v", batch.MetaData.Extra)
	}
}
