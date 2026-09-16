package grpc_test

import (
	"context"
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/coditary/wuji-core/pkg/data"
	"github.com/coditary/wuji-core/pkg/driver"
	"github.com/coditary/wuji-core/pkg/driver/dummy"
	grpcdriver "github.com/coditary/wuji-core/pkg/driver/grpc"
)

func TestRemoteRunRAGIndexQuery(t *testing.T) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := lis.Addr().String()
	_ = lis.Close()

	go func() { _ = grpcdriver.Serve(dummy.New(), grpcdriver.HostOptions{Addr: addr}) }()

	remote, err := grpcdriver.Connect(context.Background(), addr)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = remote.Close() })

	root := t.TempDir()
	doc := filepath.Join(root, "remote.txt")
	if err := os.WriteFile(doc, []byte("remote rag indexing works"), 0o644); err != nil {
		t.Fatal(err)
	}

	indexShape, err := remote.ProduceRAG(context.Background(), driver.RAGRequest{
		Task: driver.RAGTaskIndex, StoreRoot: root, Collection: "remote",
		SourcePaths: []string{doc},
	})
	if err != nil {
		t.Fatal(err)
	}
	if indexShape.Kind() != data.ShapeRecordSet {
		t.Fatalf("index kind=%s", indexShape.Kind())
	}

	queryShape, err := remote.ProduceRAG(context.Background(), driver.RAGRequest{
		Task: driver.RAGTaskQuery, StoreRoot: root, Collection: "remote",
		Query: "indexing", TopK: 2,
	})
	if err != nil {
		t.Fatal(err)
	}
	hits, err := driver.RAGQueryFromShape(queryShape)
	if err != nil {
		t.Fatal(err)
	}
	if len(hits.Hits) == 0 {
		t.Fatal("expected hits")
	}
}

func TestRemoteRunRAGDryRunAndFilter(t *testing.T) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := lis.Addr().String()
	_ = lis.Close()

	go func() { _ = grpcdriver.Serve(dummy.New(), grpcdriver.HostOptions{Addr: addr}) }()

	remote, err := grpcdriver.Connect(context.Background(), addr)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = remote.Close() })

	root := t.TempDir()
	doc := filepath.Join(root, "keep.txt")
	if err := os.WriteFile(doc, []byte("filter me source keep"), 0o644); err != nil {
		t.Fatal(err)
	}

	dryShape, err := remote.ProduceRAG(context.Background(), driver.RAGRequest{
		Task: driver.RAGTaskIndex, StoreRoot: root, Collection: "proto",
		SourcePaths: []string{doc}, DryRun: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	rs, ok := dryShape.(*data.RecordSet)
	if !ok {
		t.Fatalf("got %T", dryShape)
	}
	foundDry := false
	for _, rec := range rs.Records {
		if v, ok := rec.Fields["dry_run"]; ok && v.B {
			foundDry = true
		}
	}
	if !foundDry {
		t.Fatal("expected dry_run in remote index response")
	}

	_, err = remote.ProduceRAG(context.Background(), driver.RAGRequest{
		Task: driver.RAGTaskIndex, StoreRoot: root, Collection: "proto",
		SourcePaths: []string{doc},
	})
	if err != nil {
		t.Fatal(err)
	}

	queryShape, err := remote.ProduceRAG(context.Background(), driver.RAGRequest{
		Task: driver.RAGTaskQuery, StoreRoot: root, Collection: "proto",
		Query: "filter", TopK: 5,
		Filter: map[string]string{"source": doc}, IncludeScores: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	hits, err := driver.RAGQueryFromShape(queryShape)
	if err != nil {
		t.Fatal(err)
	}
	if len(hits.Hits) == 0 {
		t.Fatal("expected filtered hits")
	}
}
