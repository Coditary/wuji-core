package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/coditary/wuji-core/pkg/config"
	"github.com/coditary/wuji-core/pkg/core"
	grpcserver "github.com/coditary/wuji-core/internal/server/grpc"
)

func main() {
	addr := flag.String("addr", "", "core listen endpoint (default: unix socket under .wuji/run/drivers/core.sock)")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if strings.TrimSpace(*addr) == "" {
		*addr = config.DefaultCoreEndpoint(cfg.Root)
	}

	c, err := core.New(core.Config{AppConfig: cfg, Lazy: true})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer c.Close()

	srv, err := grpcserver.NewServer(c, *addr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "Wuji core daemon listening on %s\n", srv.Addr())

	errCh := make(chan error, 1)
	go func() { errCh <- srv.Start() }()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case <-sigCh:
		fmt.Fprintln(os.Stderr, "Shutting down...")
		srv.Stop()
	}
}
