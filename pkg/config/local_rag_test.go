package config_test

import (
	"testing"

	"github.com/coditary/wuji-core/pkg/config"
)

func TestResolvedLocalRAGDefaults(t *testing.T) {
	cfg := &config.Config{}
	rag := cfg.ResolvedLocalRAG()
	if rag.DefaultEmbedModel != "nomic-embed-text" {
		t.Fatalf("embed model=%q", rag.DefaultEmbedModel)
	}
	if rag.OllamaAPI != "http://127.0.0.1:11434" {
		t.Fatalf("ollama api=%q", rag.OllamaAPI)
	}
}

func TestResolvedLocalRAGMerge(t *testing.T) {
	cfg := &config.Config{
		Drivers: config.DriversConfig{
			LocalRAG: &config.LocalRAGConfig{
				DefaultEmbedModel: "mxbai-embed-large",
				OllamaAPI:         "http://127.0.0.1:11500",
			},
		},
	}
	rag := cfg.ResolvedLocalRAG()
	if rag.DefaultEmbedModel != "mxbai-embed-large" {
		t.Fatalf("embed model=%q", rag.DefaultEmbedModel)
	}
	if rag.OllamaAPI != "http://127.0.0.1:11500" {
		t.Fatalf("ollama api=%q", rag.OllamaAPI)
	}
}

func TestResolvedLocalRAGLegacyYAMLKey(t *testing.T) {
	cfg := &config.Config{
		Drivers: config.DriversConfig{
			RAG: &config.LocalRAGConfig{
				DefaultEmbedModel: "legacy-model",
			},
		},
	}
	rag := cfg.ResolvedLocalRAG()
	if rag.DefaultEmbedModel != "legacy-model" {
		t.Fatalf("embed model=%q", rag.DefaultEmbedModel)
	}
}
