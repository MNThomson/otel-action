package main

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/credentials"
)

const (
	instrumentationName    = "github.com/MNThomson/otel-action"
	instrumentationVersion = "dev"
)

func Resource(conf configType) *resource.Resource {
	attributes := []attribute.KeyValue{
		attribute.String("telemetry.sdk.language", "go"),
		attribute.String("service.version", instrumentationVersion),
		attribute.String("service.name", conf.serviceName),
	}

	return resource.NewWithAttributes(
		semconv.SchemaURL,
		attributes...,
	)
}

func setupOTEL(ctx context.Context, conf configType) (func(context.Context) error, error) {
	stdoutexp, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, fmt.Errorf("creating stdout exporter: %w", err)
	}

	client := otlptracegrpc.NewClient(
		otlptracegrpc.WithEndpoint(conf.otelEndpoint),
		otlptracegrpc.WithHeaders(conf.otelHeaders),
		otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, "")),
	)
	otelexp, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("creating otlp exporter: %w", err)
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(stdoutexp),
		sdktrace.WithBatcher(otelexp),
		sdktrace.WithResource(Resource(conf)),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	otel.SetTracerProvider(tracerProvider)

	tracer = otel.GetTracerProvider().Tracer(
		instrumentationName,
		trace.WithInstrumentationVersion(instrumentationVersion),
		trace.WithSchemaURL(semconv.SchemaURL),
	)

	return tracerProvider.Shutdown, nil
}
