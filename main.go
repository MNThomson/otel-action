package main

import (
	"context"
	"log"

	"go.opentelemetry.io/otel/trace"
)

var (
	tracer trace.Tracer
)

func main() {
	ctx := context.Background()

	// Get env vars
	conf, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Registers a tracer Provider globally.
	shutdown, err := setupOTEL(ctx, conf)
	if err != nil {
		log.Fatal(err)
	}
	defer shutdown(ctx)

	// Create traces from workflow data
	err = createTraces(ctx, conf)
	if err != nil {
		log.Fatal(err)
	}
}
