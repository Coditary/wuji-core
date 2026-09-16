package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/coditary/wuji-core/pkg/config"
	grpcdriver "github.com/coditary/wuji-core/pkg/driver/grpc"
	"github.com/coditary/wuji-core/pkg/driver/local"
)

func main() {
	var (
		addr     = flag.String("addr", "", "gRPC listen address (default: unix socket under .wuji/run/drivers/)")
		coreAddr = flag.String("core", "", "Wuji core address to register with (optional)")
	)
	flag.Parse()

	listenAddr, err := config.ResolveListenAddr(local.DriverID, *addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "config: %v\n", err)
		os.Exit(1)
	}

	appCfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config: %v\n", err)
		os.Exit(1)
	}
	ragCfg := appCfg.ResolvedLocalRAG()
	fmt.Fprintf(os.Stderr, "local-rag embed model: %s\n", ragCfg.DefaultEmbedModel)
	fmt.Fprintf(os.Stderr, "Ollama API: %s\n", ragCfg.OllamaAPI)

	drv := local.NewCompositeFromConfig(appCfg)
	if err := grpcdriver.Serve(drv, grpcdriver.HostOptions{
		Addr:     listenAddr,
		CoreAddr: *coreAddr,
	}); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
