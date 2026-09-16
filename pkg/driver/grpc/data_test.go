package grpc_test

import (
	"context"
	"net"
	"testing"

	"github.com/coditary/wuji-core/pkg/data"
	"github.com/coditary/wuji-core/pkg/driver"
	"github.com/coditary/wuji-core/pkg/driver/dummy"
	grpcdriver "github.com/coditary/wuji-core/pkg/driver/grpc"
)

func TestRemoteProduceDataEmbed(t *testing.T) {
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

	shape, err := remote.ProduceData(context.Background(), driver.DataRequest{
		Task: driver.DataTaskEmbed,
		Text: "remote",
	})
	if err != nil {
		t.Fatal(err)
	}
	if shape.Kind() != data.ShapeVector {
		t.Fatalf("kind=%s", shape.Kind())
	}
}
