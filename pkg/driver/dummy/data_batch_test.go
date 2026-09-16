package dummy_test

import (
	"context"
	"testing"

	"github.com/coditary/wuji-core/pkg/data"
	"github.com/coditary/wuji-core/pkg/driver"
	"github.com/coditary/wuji-core/pkg/driver/dummy"
)

func TestProduceDataEmbedBatch(t *testing.T) {
	d := dummy.New()
	shape, err := d.ProduceData(context.Background(), driver.DataRequest{
		Task:  driver.DataTaskEmbed,
		Texts: []string{"one", "two", "three"},
	})
	if err != nil {
		t.Fatal(err)
	}
	batch, ok := shape.(*data.VectorBatch)
	if !ok {
		t.Fatalf("got %T", shape)
	}
	if batch.Count() != 3 || batch.Dims != 8 {
		t.Fatalf("count=%d dims=%d", batch.Count(), batch.Dims)
	}
}
