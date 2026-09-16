package driver_test

import (
	"encoding/json"
	"testing"

	"github.com/coditary/wuji-core/pkg/driver"
)

func TestRAGRequestJSONFromTaijiPayload(t *testing.T) {
	raw := `{"task":"index","store_root":"/tmp","collection":"default","source_paths":["."],"recursive":true}`
	var r driver.RAGRequest
	if err := json.Unmarshal([]byte(raw), &r); err != nil {
		t.Fatal(err)
	}
	if r.Task != driver.RAGTaskIndex {
		t.Fatalf("task=%q", r.Task)
	}
	if r.StoreRoot != "/tmp" {
		t.Fatalf("store_root=%q", r.StoreRoot)
	}
	if len(r.SourcePaths) != 1 || r.SourcePaths[0] != "." {
		t.Fatalf("source_paths=%v", r.SourcePaths)
	}
	if err := r.Validate(); err != nil {
		t.Fatalf("validate: %v", err)
	}
}

func TestRAGRequestProtoRoundTrip(t *testing.T) {
	orig := driver.RAGRequest{
		Task: driver.RAGTaskQuery, StoreRoot: "/tmp", Collection: "docs", Query: "hello",
		SourcePaths: []string{"a.txt"}, UseStdin: false, Recursive: true,
		ChunkSize: 512, ChunkOverlap: 64, EmbedModel: "e1", TextModel: "t1",
		SystemPrompt: "be brief", TopK: 8, MinScore: 0.2, MaxTokens: 256,
		Filter:      map[string]string{"source": "a.txt"},
		PurgeSource: "old.txt", IndexMode: driver.RAGIndexAppend, Force: true, DryRun: true,
		IndexMetadata: map[string]string{"team": "hr"}, Glob: "*.md", Exclude: []string{"drafts"},
		URL: "http://example.com", QueryFile: "q.txt", QueryStdin: false,
		RerankTop: 20, Diverse: true, ContextMaxChars: 4000, Cite: true,
		Temperature: 0.7, TopP: 0.9, NoContext: false, IncludeScores: true,
		RenameTo: "new", ExportPath: "/out", ImportPath: "/in",
	}
	pb := driver.RAGRequestToProto(orig)
	back, err := driver.RAGRequestFromProto(pb)
	if err != nil {
		t.Fatal(err)
	}
	if back.Task != orig.Task || back.Collection != orig.Collection || back.Query != orig.Query {
		t.Fatalf("basic fields mismatch: %+v", back)
	}
	if back.Filter["source"] != "a.txt" || back.IndexMetadata["team"] != "hr" {
		t.Fatalf("maps mismatch filter=%v meta=%v", back.Filter, back.IndexMetadata)
	}
	if back.IndexMode != driver.RAGIndexAppend || !back.Force || !back.DryRun || !back.Diverse || !back.Cite {
		t.Fatalf("flags mismatch: %+v", back)
	}
	if back.RerankTop != 20 || back.ExportPath != "/out" || back.RenameTo != "new" {
		t.Fatalf("extended fields mismatch: %+v", back)
	}
}

func TestRAGTasksProtoIncludesNewTasks(t *testing.T) {
	tasks := driver.AllRAGTasks()
	pb := driver.RAGTasksToProto(tasks)
	found := map[string]bool{}
	for _, task := range pb {
		found[task] = true
	}
	for _, want := range []driver.RAGTask{driver.RAGTaskStats, driver.RAGTaskExport, driver.RAGTaskImport, driver.RAGTaskRename} {
		if !found[string(want)] {
			t.Fatalf("missing task %q in %v", want, pb)
		}
	}
}

func TestProduceDataProtoTextsRoundTrip(t *testing.T) {
	req := driver.DataRequest{
		Task: driver.DataTaskEmbed, Model: "m", Texts: []string{"a", "b", "c"},
	}
	pb, err := driver.ProduceDataRequestToProto(req)
	if err != nil {
		t.Fatal(err)
	}
	back, err := driver.ProduceDataRequestFromProto(pb)
	if err != nil {
		t.Fatal(err)
	}
	if len(back.Texts) != 3 || back.Texts[1] != "b" {
		t.Fatalf("texts=%v", back.Texts)
	}
}
