package rest

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	"log"
	"net/http"
	v1 "proto-parking-services/pkg/api/v1"
)

// RunServer http/rest gateway
func RunServer(ctx context.Context, gRPCPort, httpPort string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	if err := v1.RegisterParkingServicesHandlerFromEndpoint(ctx, mux, "localhost:"+gRPCPort, opts); err != nil {
		log.Fatalf("Failed to start HTTP gateway: %v", err)
	}

	srv := &http.Server{
		Addr:    ":" + httpPort,
		Handler: mux,
	}

	// gracefull shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
		}

		_, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		_ = srv.Shutdown(ctx)
	}()

	log.Println("starting HTTP/REST gateway...")

	return srv.ListenAndServe()
}
