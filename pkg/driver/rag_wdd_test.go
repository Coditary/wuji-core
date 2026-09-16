package driver_test

import (
	"testing"

	"github.com/coditary/wuji-core/pkg/driver"
)

func TestRAGQueryShapeRoundtrip(t *testing.T) {
	orig := &driver.RAGQueryResponse{
		Query: "hello", Collection: "docs",
		Hits: []driver.RAGHit{
			{ID: "c1", Score: 0.9, Text: "chunk one", Source: "a.txt", ChunkIdx: 1},
		},
	}
	shape := driver.RAGQueryToRecordSet(orig)
	got, err := driver.RAGQueryFromShape(shape)
	if err != nil {
		t.Fatal(err)
	}
	if got.Query != orig.Query || got.Collection != orig.Collection {
		t.Fatalf("meta query=%q collection=%q", got.Query, got.Collection)
	}
	if len(got.Hits) != 1 || got.Hits[0].Text != "chunk one" {
		t.Fatalf("hits=%+v", got.Hits)
	}
}
