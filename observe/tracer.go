package observe

import (
	"context"
	"fmt"
	"strings"

	"github.com/bushubdegefu/m-playground/configs"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var AppTracer = otel.Tracer(fmt.Sprintf("cli-server-%v", configs.AppConfig.GetOrDefault("APP_NAME", "blue-app")))

func InitTracer() *sdktrace.TracerProvider {
	traceExporter := configs.AppConfig.Get("TRACE_EXPORTER")
	tracerHost := configs.AppConfig.Get("TRACER_HOST")
	tracerPort := configs.AppConfig.GetOrDefault("TRACER_PORT", "4317") // Default to gRPC port

	// Create resource with service name
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(configs.AppConfig.Get("APP_NAME")),
	)

	// Create TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
	)

	// Set global tracer provider and propagator
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	if strings.ToLower(traceExporter) == "jaeger" && tracerHost != "" {
		exporter, err := otlptracegrpc.New(context.Background(),
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(fmt.Sprintf("%s:%s", tracerHost, tracerPort)),
		)
		if err != nil {
			panic(fmt.Sprintf("Failed to create Jaeger exporter: %v", err))
		}

		// Register batch span processor
		batcher := sdktrace.NewBatchSpanProcessor(exporter)
		tp.RegisterSpanProcessor(batcher)
	}

	return tp
}

func EchoAppSpanner(ctx echo.Context, span_name string) (context.Context, oteltrace.Span) {
	gen, _ := uuid.NewV7()
	id := gen.String()

	//  getting request body
	trace, span := AppTracer.Start(ctx.Request().Context(), span_name,
		oteltrace.WithAttributes(attribute.String("id", id)),
	)
	return trace, span
}

type RouteTracer struct {
	Tracer context.Context
	Span   oteltrace.Span
}
