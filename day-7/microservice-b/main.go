package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var (
	requestCounter        metric.Int64Counter
	requestDuration       metric.Float64Histogram
	activeRequestsCounter metric.Int64UpDownCounter
)

func initProvider() (func(context.Context) error, error) {
	ctx := context.Background()

	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Using environment variables 👌")
	}

	// Read the OTEL collector endpoint from environment variable
	otelEndpoint := os.Getenv("OTEL_COLLECTOR_ENDPOINT")
	if otelEndpoint == "" {
		otelEndpoint = "localhost:4318" // Default endpoint
	}

	// Create a resource with the service name
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("microservice-b"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create OTLP trace exporter over HTTP with custom endpoint
	traceExporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(otelEndpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// Create OTLP metric exporter over HTTP with custom endpoint
	metricExporter, err := otlpmetrichttp.New(ctx,
		otlpmetrichttp.WithEndpoint(otelEndpoint),
		otlpmetrichttp.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create metric exporter: %w", err)
	}

	// Create trace provider with the exporter and resource
	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter),
		trace.WithResource(res),
	)

	// Create metric reader and meter provider with the resource
	metricReader := sdkmetric.NewPeriodicReader(metricExporter)
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(metricReader),
		sdkmetric.WithResource(res),
	)

	// Set global providers
	otel.SetTracerProvider(tracerProvider)
	otel.SetMeterProvider(meterProvider)

	return func(ctx context.Context) error {
		err := tracerProvider.Shutdown(ctx)
		if err != nil {
			return err
		}
		err = meterProvider.Shutdown(ctx)
		if err != nil {
			return err
		}
		return nil
	}, nil
}

// Basic Hello Handler
func hello(c *gin.Context) {
	startTime := time.Now()
	ctx := c.Request.Context()

	// Increment active requests
	activeRequestsCounter.Add(ctx, 1)
	defer activeRequestsCounter.Add(ctx, -1)

	c.JSON(http.StatusOK, gin.H{
		"message": "👋 Hello from microservice-b",
	})

	duration := time.Since(startTime).Milliseconds()

	requestCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("endpoint", "/hello-b")))
	requestDuration.Record(ctx, float64(duration), metric.WithAttributes(attribute.String("endpoint", "/hello-b")))
}

// Call Service A Handler
func callA(c *gin.Context) {
	startTime := time.Now()
	ctx := c.Request.Context()

	activeRequestsCounter.Add(ctx, 1)
	defer activeRequestsCounter.Add(ctx, -1)

	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Using environment variables 👌")
	}

	SVC_A_URI := os.Getenv("SVC_A_URI")
	if SVC_A_URI == "" {
		SVC_A_URI = "http://localhost:8080" // Default URI for service-A
	}

	// Create a new HTTP client with OpenTelemetry instrumentation
	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	// Create a new request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/hello-a", SVC_A_URI), nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request to service-A"})
		return
	}

	// Use the context from Gin
	req = req.WithContext(ctx)

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reach service-A"})
		return
	}
	defer resp.Body.Close()

	resBody, _ := ioutil.ReadAll(resp.Body)

	c.JSON(http.StatusOK, gin.H{
		"message":  "🥳 Response from service-A",
		"response": string(resBody),
	})

	duration := time.Since(startTime).Milliseconds()

	requestCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("endpoint", "/call-a")))
	requestDuration.Record(ctx, float64(duration), metric.WithAttributes(attribute.String("endpoint", "/call-a")))
}

// Get Coffee Handler
func getMeCoffee(c *gin.Context) {
	startTime := time.Now()
	ctx := c.Request.Context()

	activeRequestsCounter.Add(ctx, 1)
	defer activeRequestsCounter.Add(ctx, -1)

	// Create a new HTTP client with OpenTelemetry instrumentation
	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	// Create a new request
	req, err := http.NewRequest("GET", "https://api.sampleapis.com/coffee/iced", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request to coffee API"})
		return
	}

	// Use the context from Gin
	req = req.WithContext(ctx)

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch coffee"})
		return
	}
	defer resp.Body.Close()

	resBody, _ := ioutil.ReadAll(resp.Body)

	c.JSON(http.StatusOK, gin.H{
		"message":  "🍵 Here is your coffee",
		"response": string(resBody),
	})

	duration := time.Since(startTime).Milliseconds()

	requestCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("endpoint", "/getme-coffee")))
	requestDuration.Record(ctx, float64(duration), metric.WithAttributes(attribute.String("endpoint", "/getme-coffee")))
}

func main() {
	ctx := context.Background()
	shutdown, err := initProvider()
	if err != nil {
		log.Fatalf("Failed to initialize OpenTelemetry: %v", err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatalf("Error shutting down provider: %v", err)
		}
	}()

	router := gin.Default()

	// Use OpenTelemetry middleware for Gin
	router.Use(otelgin.Middleware("microservice-b"))

	// Initialize the Meter
	meter := otel.GetMeterProvider().Meter("microservice-b")

	// Initialize instruments using the Meter interface methods
	requestCounter, err = meter.Int64Counter(
		"request_count",
		metric.WithDescription("Counts the number of requests received"),
	)
	if err != nil {
		log.Fatalf("Failed to create counter: %v", err)
	}

	requestDuration, err = meter.Float64Histogram(
		"request_duration_ms",
		metric.WithDescription("Records the duration of requests in milliseconds"),
	)
	if err != nil {
		log.Fatalf("Failed to create histogram: %v", err)
	}

	activeRequestsCounter, err = meter.Int64UpDownCounter(
		"active_requests",
		metric.WithDescription("Counts the number of active requests"),
	)
	if err != nil {
		log.Fatalf("Failed to create up-down counter: %v", err)
	}

	router.GET("/hello-b", hello)
	router.GET("/call-a", callA)
	router.GET("/getme-coffee", getMeCoffee)

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "80" // Default URI for service-B
	}

	// Start the server
	router.Run(fmt.Sprintf(":%s", PORT))
}
