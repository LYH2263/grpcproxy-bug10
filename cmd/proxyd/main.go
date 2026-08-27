package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LYH2263/go-grpcproxy"
	"github.com/LYH2263/go-grpcproxy/console"
)

func main() {
	grpcAddr := flag.String("grpc", ":8232", "gRPC admin listen address")
	httpAddr := flag.String("http", ":9232", "console HTTP address")
	flag.Parse()

	kit := grpcproxy.New()
	defer kit.Close()

	cs := console.New(*httpAddr, func() string { return grpcproxy.InspectSummary(kit) })
	go func() {
		if err := cs.ListenAndServe(); err != nil {
			fmt.Fprintln(os.Stderr, "console:", err)
		}
	}()

	fmt.Printf("proxyd grpc stub on %s (kit ready)\n", *grpcAddr)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = shutdown
}
