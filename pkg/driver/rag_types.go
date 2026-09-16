package driver

import (
	"fmt"
	"strings"
)

// RAGIndexMode controls how new documents merge into an existing collection.
type RAGIndexMode string

const (
	RAGIndexAuto    RAGIndexMode = "auto"
	RAGIndexReplace RAGIndexMode = "replace"
	RAGIndexAppend  RAGIndexMode = "append"
)

// RAGOutputShape selects answer output projection.
type RAGOutputShape string

const (
	RAGOutputFull       RAGOutputShape = "full"
	RAGOutputAnswerOnly RAGOutputShape = "answer"
)

// RAGRequest is the unified CLI/core input for RAG operations.
type RAGRequest struct {
	Task         RAGTask           `json:"task"`
	StoreRoot    string            `json:"store_root"`
	StoreDir     string            `json:"store_dir,omitempty"`
	Collection   string            `json:"collection"`
	Query        string            `json:"query"`
	SourcePaths  []string          `json:"source_paths"`
	UseStdin     bool              `json:"use_stdin"`
	Recursive    bool              `json:"recursive"`
	ChunkSize    int               `json:"chunk_size"`
	ChunkOverlap int               `json:"chunk_overlap"`
	EmbedModel   string            `json:"embed_model"`
	TextModel    string            `json:"text_model"`
	SystemPrompt string            `json:"system_prompt"`
	TopK         int               `json:"top_k"`
	MinScore     float32           `json:"min_score"`
	MaxTokens    int               `json:"max_tokens"`
	Filter       map[string]string `json:"filter"`
	PurgeSource  string            `json:"purge_source"`

	// Index options
	IndexMode     RAGIndexMode      `json:"index_mode"`
	Force         bool              `json:"force"`
	DryRun        bool              `json:"dry_run"`
	IndexMetadata map[string]string `json:"index_metadata"`
	Glob          string            `json:"glob"`
	Exclude       []string          `json:"exclude"`
	URL           string            `json:"url"`

	// Query / answer options
	QueryFile       string  `json:"query_file"`
	QueryStdin      bool    `json:"query_stdin"`
	RerankTop       int     `json:"rerank_top"`
	Diverse         bool    `json:"diverse"`
	ContextMaxChars int     `json:"context_max_chars"`
	Cite            bool    `json:"cite"`
	Temperature     float32 `json:"temperature"`
	TopP            float32 `json:"top_p"`
	NoContext       bool    `json:"no_context"`
	IncludeScores   bool    `json:"include_scores"`

	// Collection management
	RenameTo   string `json:"rename_to"`
	ExportPath string `json:"export_path"`
	ImportPath string `json:"import_path"`
}

// RAGHit is one retrieval result.
type RAGHit struct {
	ID       string
	Score    float32
	Text     string
	Source   string
	ChunkIdx int
	Metadata map[string]string
}

// RAGIndexResponse summarizes indexing.
type RAGIndexResponse struct {
	Collection   string
	ChunksAdded  int
	TotalChunks  int
	BytesIndexed int64
	Sources      []string
	Message      string
	DryRun       bool
}

// RAGQueryResponse holds retrieval hits.
type RAGQueryResponse struct {
	Query      string
	Collection string
	Hits       []RAGHit
}

// RAGAnswerResponse combines hits and generated text.
type RAGAnswerResponse struct {
	Answer     string
	Collection string
	Query      string
	Hits       []RAGHit
	TokensUsed int
	NoContext  bool
}

// RAGCollectionInfo describes a stored collection.
type RAGCollectionInfo struct {
	Collection string
	EmbedModel string
	Dims       int
	ChunkCount int
	ChunkSize  int
	Overlap    int
	Sources    []string
	UpdatedAt  string
}

// RAGManageResponse is returned for list/info/delete/purge and port ops.
type RAGManageResponse struct {
	Task        RAGTask
	Message     string
	Collections []RAGCollectionInfo
	Collection  RAGCollectionInfo
}

func (r RAGRequest) Validate() error {
	switch r.Task {
	case RAGTaskIndex:
		if strings.TrimSpace(r.Collection) == "" {
			return fmt.Errorf("--collection is required for --index")
		}
		if r.URL != "" {
			return fmt.Errorf("URL ingestion is not implemented yet; use --file or --dir")
		}
		if len(r.SourcePaths) == 0 && !r.UseStdin {
			return fmt.Errorf("--index requires --dir/--file or --stdin")
		}
	case RAGTaskQuery, RAGTaskAnswer:
		if strings.TrimSpace(r.Collection) == "" {
			return fmt.Errorf("--collection is required")
		}
		if strings.TrimSpace(r.Query) == "" {
			return fmt.Errorf("query text is required (positional, --text, --query-file, or --query-stdin)")
		}
	case RAGTaskInfo, RAGTaskDelete, RAGTaskStats:
		if strings.TrimSpace(r.Collection) == "" {
			return fmt.Errorf("--collection is required")
		}
	case RAGTaskPurge:
		if strings.TrimSpace(r.Collection) == "" || strings.TrimSpace(r.PurgeSource) == "" {
			return fmt.Errorf("--collection and --source are required for --purge")
		}
	case RAGTaskRename:
		if strings.TrimSpace(r.Collection) == "" || strings.TrimSpace(r.RenameTo) == "" {
			return fmt.Errorf("--collection and --rename-to are required for --rename")
		}
	case RAGTaskExport:
		if strings.TrimSpace(r.Collection) == "" || strings.TrimSpace(r.ExportPath) == "" {
			return fmt.Errorf("--collection and --to are required for --export")
		}
	case RAGTaskImport:
		if strings.TrimSpace(r.Collection) == "" || strings.TrimSpace(r.ImportPath) == "" {
			return fmt.Errorf("--collection and --from are required for --import")
		}
	case RAGTaskList:
	default:
		return fmt.Errorf("unknown rag task %q", r.Task)
	}
	if r.StoreRoot == "" {
		return fmt.Errorf("store root is required")
	}
	return nil
}

func (r RAGRequest) MaxTokensOrDefault() int {
	if r.MaxTokens > 0 {
		return r.MaxTokens
	}
	return 1024
}

func (r RAGRequest) IndexModeOrDefault() RAGIndexMode {
	switch r.IndexMode {
	case "", RAGIndexAuto, RAGIndexReplace:
		return RAGIndexReplace
	case RAGIndexAppend:
		return RAGIndexAppend
	default:
		return RAGIndexReplace
	}
}

func (r RAGRequest) SearchK() int {
	k := r.TopK
	if k <= 0 {
		k = 5
	}
	if r.RerankTop > k {
		return r.RerankTop
	}
	return k
}

func (r RAGRequest) FinalTopK() int {
	k := r.TopK
	if k <= 0 {
		return 5
	}
	return k
}
