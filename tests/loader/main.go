// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Example using OTLP exporters + collector + third-party backends. For
// information about using the exporter, see:
// https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp?tab=doc#example-package-Insecure
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"google.golang.org/grpc"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

var token = `Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwczovL0RFVi1pbnRlZy1rM2QtcmNhLmF1dGguc2VhbGlnaHRzLmlvLyIsImp3dGlkIjoiREVWLWludGVnLWszZC1yY2EsbmVlZFRvUmVtb3ZlLEFQSUdXLWI3ZmQxNzRiLWY1ZTgtNDFiZi04YzI3LTY2ZWYyYmQ0MmQwOSwxNjcxNDUxNDIxNDEzIiwic3ViamVjdCI6IlNlYUxpZ2h0c0BhZ2VudCIsImF1ZGllbmNlIjpbImFnZW50cyJdLCJ4LXNsLXJvbGUiOiJhZ2VudCIsIngtc2wtc2VydmVyIjoiaHR0cHM6Ly9kZXYtaW50ZWctazNkLXJjYS1ndy5kZXYuc2VhbGlnaHRzLmNvL2FwaSIsInNsX2ltcGVyX3N1YmplY3QiOiIiLCJpYXQiOjE2NzE0NTE0MjF9.Kqa6g6NVZ_h7cMNcqQHNSUkF2dfwCFW35xEzv5y493vLjncBls1OBWwM54_nNIWcAbBGOGX8Wx8qJVURWt5fW6o33wo5FFpQMcQRRGkz1s1T2ciGdoS3M-lnpDT0BUBlWmLqhwwhb2LOWRY4aRbUAHOOKPInqgWoDw5GS3qxMsyboPf5h9udNv3TdkxB3kAHJnsvfxeZ5WTcKJCsqocPa3hrysPVHm0KNZOpZawNg3O4PKhFC6GKSpU7cWJhr3BK_4aHD0o5pq_o2hqg3pgnzH6FR1yZwAL68wyMGquvR5p9f3-hw9jCaWo0KyAKYKBUIXDbiHM6U1tGB-4wNa1kh8pw0eBaPQVj1zS5TKgESDgPM5CYzGvjWxys0wJCpzfqQiAwqnMZbiMWXX0f0ks0A1ms6brk1vmg-uqC8YyPFo3lT0wASdT66MoNfhOUP3IKV85vkygcWwDVOIKLP9-3o3hueSMkYDVjLiB1qF8cm7y1fEO9MF6Cv5CC8od7-4JKNPQ_or6zuXIylRk7O5vXveMW6vzgm702SvJcLLxqwrpNkXuHPZruu9o8l23-LN5O_qe_Mdl0urrLMFe2r2QwHg7ZNg9sBHPb61RUVYKLiTFm2KQzA-NlUD-jCTfrF-yGybH6M1X016dfW6rr4R_0kls3WBRhmxIq8DdbhUEPH-E`

// Initializes an OTLP exporter, and configures the corresponding trace and
// metric providers.
func initProvider() (func(context.Context) error, error) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceNameKey.String("test-service"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// If the OpenTelemetry Collector is running on a local cluster (minikube or
	// microk8s), it should be accessible through the NodePort service at the
	// `localhost:30080` endpoint. Otherwise, replace `localhost` with the
	// endpoint of your cluster. If you run the app inside k8s, then you can
	// probably connect directly to the service through dns.
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	//conn, err := grpc.DialContext(ctx, "localhost:4317",
	//	// Note the use of insecure transport here. TLS is recommended in production.
	//	grpc.WithTransportCredentials(insecure.NewCredentials()),
	//	grpc.WithBlock(),
	//)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	//}

	headers := map[string]string{"Authorization": token} // Replace xtraceToken with the authentication token obtained in the Prerequisites section.
	traceClient := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint("localhost:4317"), // Replace otelAgentAddr with the endpoint obtained in the Prerequisites section.
		otlptracegrpc.WithHeaders(headers),
		otlptracegrpc.WithDialOption(grpc.WithBlock()))
	log.Println("start to connect to server")
	// Set up a trace exporter
	traceExporter, err := otlptrace.New(ctx, traceClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := sdktrace.NewSimpleSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Shutdown will flush any remaining spans and shut down the exporter.
	return tracerProvider.Shutdown, nil
}

func main() {
	log.Printf("Waiting for connection...")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	shutdown, err := initProvider()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	tracer := otel.Tracer("test-tracer")

	// Attributes represent additional key-value descriptors that can be bound
	// to a metric observer or recorder.
	commonAttrs := []attribute.KeyValue{
		attribute.String("attrA", "chocolate"),
		attribute.String("attrB", "raspberry"),
		attribute.String("attrC", "vanilla"),
	}

	// work begins
	ctx, span := tracer.Start(
		ctx,
		"CollectorExporter-Example",
		trace.WithAttributes(commonAttrs...))
	defer span.End()
	for i := 0; i < 10; i++ {
		_, iSpan := tracer.Start(ctx, fmt.Sprintf("Sample-%d", i))
		log.Printf("Doing really hard work (%d / 10)\n", i+1)

		<-time.After(time.Second)
		iSpan.End()
	}

	log.Printf("Done!")
}
