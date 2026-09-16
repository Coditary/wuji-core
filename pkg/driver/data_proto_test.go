package driver_test

import (
	"testing"

	wujiv1 "github.com/coditary/wuji-core/api/proto/v1"
	"github.com/coditary/wuji-core/pkg/data"
	"github.com/coditary/wuji-core/pkg/driver"
)

func TestProduceDataProtoRoundTrip(t *testing.T) {
	batch, err := data.NewVectorBatch([]string{"hello"}, 4, [][]float32{{0.1, 0.2, 0.3, 0.4}})
	if err != nil {
		t.Fatal(err)
	}
	req := driver.DataRequest{
		Task:  driver.DataTaskEmbed,
		Text:  "hello",
		Input: batch,
	}
	pb, err := driver.ProduceDataRequestToProto(req)
	if err != nil {
		t.Fatal(err)
	}
	back, err := driver.ProduceDataRequestFromProto(pb)
	if err != nil {
		t.Fatal(err)
	}
	if back.Task != driver.DataTaskEmbed || back.Input == nil {
		t.Fatalf("req=%+v", back)
	}

	resp, err := driver.ProduceDataResponseToProto(batch)
	if err != nil {
		t.Fatal(err)
	}
	shape, err := driver.ProduceDataResponseFromProto(resp)
	if err != nil {
		t.Fatal(err)
	}
	if shape.Kind() != data.ShapeVector {
		t.Fatalf("kind=%s", shape.Kind())
	}
}

func TestDataTasksProto(t *testing.T) {
	tasks := driver.AllDataTasks()
	pb := driver.DataTasksToProto(tasks)
	if len(pb) != len(tasks) {
		t.Fatalf("len=%d", len(pb))
	}
	back := driver.DataTasksFromProto(pb)
	if len(back) != len(tasks) {
		t.Fatalf("back=%d", len(back))
	}
	_ = wujiv1.ProduceDataRequest{}
}
