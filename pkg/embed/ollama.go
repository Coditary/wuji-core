package embed

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/coditary/wuji-core/pkg/data"
	"github.com/coditary/wuji-core/pkg/driver"
	"github.com/coditary/wuji-core/pkg/rag"
)

// Ollama produces embeddings via the Ollama HTTP API.
type Ollama struct {
	BaseURL      string
	DefaultModel string
	DriverID     string
	HTTPClient   *http.Client
}

// NewOllama returns an Ollama embedding producer with sensible defaults.
func NewOllama(baseURL, defaultModel, driverID string) *Ollama {
	if baseURL == "" {
		baseURL = "http://127.0.0.1:11434"
	}
	if defaultModel == "" {
		defaultModel = "nomic-embed-text"
	}
	if driverID == "" {
		driverID = "local-rag"
	}
	return &Ollama{
		BaseURL:      strings.TrimRight(baseURL, "/"),
		DefaultModel: defaultModel,
		DriverID:     driverID,
		HTTPClient:   &http.Client{Timeout: 5 * time.Minute},
	}
}

func (o *Ollama) ProduceData(ctx context.Context, req driver.DataRequest) (data.DataShape, error) {
	task := req.TaskOrDefault()
	switch task {
	case driver.DataTaskEmbed, driver.DataTaskVector:
	default:
		return nil, fmt.Errorf("%s driver: unsupported data task %q", o.DriverID, task)
	}

	inputs := append([]string(nil), req.Texts...)
	if t := strings.TrimSpace(req.Text); t != "" {
		inputs = append(inputs, t)
	}
	if req.ImagePath != "" {
		inputs = append(inputs, "image:"+req.ImagePath)
	}
	if len(inputs) == 0 {
		return nil, fmt.Errorf("embed requires --text, batch texts, or --image")
	}

	model := strings.TrimSpace(req.Model)
	if model == "" {
		model = o.DefaultModel
	}

	rows, dims, err := o.embedTexts(ctx, model, inputs)
	if err != nil {
		rows, dims, fbErr := rag.PseudoEmbedder(ctx, model, inputs)
		if fbErr != nil {
			return nil, err
		}
		_ = err
		batch, bErr := data.NewVectorBatch(inputs, dims, rows)
		if bErr != nil {
			return nil, bErr
		}
		batch.MetaData = embedMeta(o.DriverID, model, task, true)
		return batch, nil
	}

	batch, err := data.NewVectorBatch(inputs, dims, rows)
	if err != nil {
		return nil, err
	}
	batch.MetaData = embedMeta(o.DriverID, model, task, false)
	return batch, nil
}

func embedMeta(driverID, model string, task driver.DataTask, fallback bool) data.Meta {
	meta := data.Meta{
		Capability: "data",
		Task:       string(task),
		Driver:     driverID,
		Model:      model,
	}
	if fallback {
		meta.Extra = map[string]string{"embed_fallback": "pseudo"}
	}
	return meta
}

func (o *Ollama) embedTexts(ctx context.Context, model string, texts []string) ([][]float32, int, error) {
	rows, dims, err := o.embedBatch(ctx, model, texts)
	if err == nil {
		return rows, dims, nil
	}
	rows = make([][]float32, 0, len(texts))
	dims = 0
	for _, text := range texts {
		vec, d, err := o.embedOne(ctx, model, text)
		if err != nil {
			return nil, 0, err
		}
		if dims == 0 {
			dims = d
		} else if d != dims {
			return nil, 0, fmt.Errorf("embedding dims changed: got %d want %d", d, dims)
		}
		rows = append(rows, vec)
	}
	return rows, dims, nil
}

func (o *Ollama) embedBatch(ctx context.Context, model string, texts []string) ([][]float32, int, error) {
	body, err := json.Marshal(map[string]any{
		"model": model,
		"input": texts,
	})
	if err != nil {
		return nil, 0, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, o.BaseURL+"/api/embed", bytes.NewReader(body))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.HTTPClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, 0, fmt.Errorf("ollama embed: %s: %s", resp.Status, strings.TrimSpace(string(msg)))
	}

	var out struct {
		Embeddings [][]float64 `json:"embeddings"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, 0, err
	}
	if len(out.Embeddings) != len(texts) {
		return nil, 0, fmt.Errorf("ollama embed: got %d vectors for %d inputs", len(out.Embeddings), len(texts))
	}
	rows := make([][]float32, len(out.Embeddings))
	dims := 0
	for i, row := range out.Embeddings {
		if len(row) == 0 {
			return nil, 0, fmt.Errorf("ollama embed: empty vector at index %d", i)
		}
		if dims == 0 {
			dims = len(row)
		} else if len(row) != dims {
			return nil, 0, fmt.Errorf("ollama embed: row %d has %d dims, want %d", i, len(row), dims)
		}
		rows[i] = float64To32(row)
	}
	return rows, dims, nil
}

func (o *Ollama) embedOne(ctx context.Context, model, text string) ([]float32, int, error) {
	body, err := json.Marshal(map[string]string{
		"model":  model,
		"prompt": text,
	})
	if err != nil {
		return nil, 0, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, o.BaseURL+"/api/embeddings", bytes.NewReader(body))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.HTTPClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, 0, fmt.Errorf("ollama embeddings: %s: %s", resp.Status, strings.TrimSpace(string(msg)))
	}

	var out struct {
		Embedding []float64 `json:"embedding"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, 0, err
	}
	if len(out.Embedding) == 0 {
		return nil, 0, fmt.Errorf("ollama embeddings: empty vector")
	}
	return float64To32(out.Embedding), len(out.Embedding), nil
}

func float64To32(in []float64) []float32 {
	out := make([]float32, len(in))
	for i, v := range in {
		out[i] = float32(v)
	}
	return out
}
