# In-Depth Metrics & Instrumentation Analysis: Service A

This document provides a detailed, technical analysis of `service-a`. It explains how metrics, tracing, and logging are implemented, detailing the logic, libraries, and middleware patterns used.

---

## 🗺️ Project Architecture Overview

`service-a` is a Node.js web application built with the **Express** framework. It serves as an entry point for user requests, handles basic routing, generates logs, interacts with downstream services (like `service-b`), and is fully instrumented for observability:
- **Metrics collection & exposure** via `prom-client` (Prometheus format).
- **Distributed tracing** via **OpenTelemetry SDK** (exported to Jaeger).
- **Structured logging** via `pino`.

---

## 📦 Key Dependencies

The primary observability dependencies are:

| Package | Purpose |
| :--- | :--- |
| `prom-client` | Standard Prometheus client library to define, register, and serve metrics. |
| `@opentelemetry/sdk-node` | OpenTelemetry Node.js SDK to bootstrap tracing. |
| `@opentelemetry/exporter-jaeger` | Exporter to ship spans directly to Jaeger collector. |
| `@opentelemetry/instrumentation-express` / `http` | Automatic hook interception to trace Express routes and HTTP requests. |
| `pino` | High-performance JSON-based logger. |

---

## 📈 Metrics Instrumentation Implementation

The metrics instrumentation is implemented using `prom-client` inside [index.js]. Three primary Prometheus metric types are utilized: **Counter**, **Histogram**, **Summary**, and **Gauge**.

### 1. Defined Metric Types & Declarations

#### A. Request Counter (`http_requests_total`)
Tracks the cumulative count of requests incoming to the server.
```javascript
const httpRequestCounter = new promClient.Counter({
    name: 'http_requests_total',
    help: 'Total number of HTTP requests',
    labelNames: ['method', 'path', 'status_code'],
});
```
*   **Dimensions (Labels):** `method` (GET/POST), `path` (endpoint URL), and `status_code` (e.g., 200, 404, 500). This allows slicing and dicing the error rates and endpoint traffic.

#### B. Request Duration Histogram (`http_request_duration_seconds`)
Measures the distribution of request durations using pre-configured latency buckets.
```javascript
const requestDurationHistogram = new promClient.Histogram({
    name: 'http_request_duration_seconds',
    help: 'Duration of HTTP requests in seconds',
    labelNames: ['method', 'path', 'status_code'],
    buckets: [0.1, 0.5, 1, 5, 10], // Buckets in seconds
});
```
*   **Buckets:** Count boundaries at $100\text{ms}$, $500\text{ms}$, $1\text{s}$, $5\text{s}$, and $10\text{s}$. Useful for computing quantiles (like p95/p99) on the Prometheus server using `histogram_quantile`.

#### C. Request Duration Summary (`http_request_duration_summary_seconds`)
Calculates streaming quantiles directly on the client side.
```javascript
const requestDurationSummary = new promClient.Summary({
    name: 'http_request_duration_summary_seconds',
    help: 'Summary of the duration of HTTP requests in seconds',
    labelNames: ['method', 'path', 'status_code'],
    percentiles: [0.5, 0.9, 0.99],
});
```
*   **Percentiles:** Configured to calculate the 50th (median), 90th, and 99th percentiles in-memory.

#### D. Active Task Gauge (`node_gauge_example`)
Tracks a value that can go up and down (in this case, measuring execution times of async tasks).
```javascript
const gauge = new promClient.Gauge({
    name: 'node_gauge_example',
    help: 'Example of a gauge tracking async task duration',
    labelNames: ['method', 'status']
});
```

---

## ⚙️ Metrics Collection & Scraping Logic

### 1. Global Request Tracking Middleware
To record request metrics automatically without manual instrumentation in every route, a custom Express middleware is applied globally (lines 66-77):

```javascript
app.use((req, res, next) => {
    const start = Date.now();
    res.on('finish', () => {
        const duration = (Date.now() - start) / 1000; // Convert ms to seconds
        const { method, url } = req;
        const statusCode = res.statusCode;
        
        // Record Counter, Histogram, and Summary
        httpRequestCounter.labels({ method, path: url, status_code: statusCode }).inc();
        requestDurationHistogram.labels({ method, path: url, status_code: statusCode }).observe(duration);
        requestDurationSummary.labels({ method, path: url, status_code: statusCode }).observe(duration);
    });
    next();
});
```

#### Middleware Lifecycle Logic:
1.  **Start Timestamp:** Captured at the entry of the request via `Date.now()`.
2.  **`res.on('finish')` Hook:** Express processes the routes asynchronously. Listening to the `finish` event ensures that the metrics are logged **only after** the response has been completely written and flushed to the network. This guarantees:
    *   The correct HTTP status code (`res.statusCode`) is captured.
    *   The total duration includes the time taken by the controller and route handlers.
3.  **Observation Recording:**
    *   `.inc()` increments the counter by 1.
    *   `.observe(duration)` records the execution time into both the Histogram and Summary buckets.

### 2. Manual Timer Tracking (Gauge)
For specific long-running or asynchronous operations (like the `/example` route), a timer-based gauge is used:
```javascript
app.get('/example', async (req, res) => {
    const endGauge = gauge.startTimer({ method: req.method, status: res.statusCode });
    await simulateAsyncTask(); // Simulates a 0-5 second delay
    endGauge(); // Stops the timer and records the duration
    res.send('Async task completed');
});
```
*   `gauge.startTimer()` initiates a clock.
*   Calling the returned `endGauge()` function stops the timer and automatically updates the gauge with the elapsed duration.

### 3. Exposing the Scrape Endpoint
Prometheus collects metrics by periodically pulling (scraping) an HTTP endpoint. This is exposed at `/metrics`:
```javascript
app.get('/metrics', async (req, res) => {
    res.set('Content-Type', promClient.register.contentType);
    res.end(await promClient.register.metrics());
});
```
*   Sets the Content-Type header to `text/plain; version=0.0.4` (standard Prometheus format).
*   Retrieves all registered metrics serialized into the Prometheus text format via `promClient.register.metrics()`.

---

## 🔍 Tracing Instrumentation (OpenTelemetry)

Distributed tracing tracks the request flow across service boundaries. Tracing in `service-a` is orchestrated via [tracing.js].

### 1. Initialization Order (CRITICAL)
In `index.js`, the very first import block is:
```javascript
require('dotenv').config();
require('./tracing'); // Initialize tracing immediately!
```
> [!IMPORTANT]
> OpenTelemetry depends on **monkey-patching** core Node.js modules (`http`, `https`) and third-party modules (like `express` and `axios`). Therefore, `tracing.js` **must** run before any package imports (like `express` or `axios`) occur. If imported later, auto-instrumentation will fail to intercept outgoing/incoming requests.

### 3. Tracing Setup Details
*   **Tracer Provider:** `NodeTracerProvider` bootstraps tracing, configured with a service name resource attribute of `service-a`.
*   **Exporter:** Spans are collected using a `JaegerExporter` pointing to the collector endpoint defined by `process.env.OTEL_EXPORTER_JAEGER_ENDPOINT`.
*   **Span Processor:** A `SimpleSpanProcessor` passes spans to the exporter immediately as they end (ideal for development/testing).
*   **Auto-instrumentation:**
    ```javascript
    registerInstrumentations({
      instrumentations: [
        new HttpInstrumentation({
          applyCustomAttributesOnSpan: (span, request, response) => {
            span.setAttribute('custom-attribute', 'custom-value');
          },
        }),
        new ExpressInstrumentation(),
      ],
    });
    ```
    *   `HttpInstrumentation` traces outgoing HTTP calls (like `axios.get` calls to `service-b`) and incoming HTTP traffic.
    *   `ExpressInstrumentation` tracks Express routes and middleware execution times as sub-spans.

---

## 🚦 Endpoints & Testing Logic

1.  **`/` and `/healthy` (200 OK):** Used to test basic success flows.
2.  **`/serverError` (500 Internal Server Error) & `/notFound` (404 Not Found):** Verifies that the metrics middleware successfully logs non-200 HTTP statuses.
3.  **`/logs`:** Invokes custom logger statements to test integration with log collection backends.
4.  **`/crash`:** Exits the process (`process.exit(1)`), which tests container restart policies and system resilience.
5.  **`/call-service-b`:** Traces a multi-tier downstream call. When `axios` calls `service-b`, OpenTelemetry injects trace context headers (`traceparent`) into the HTTP header, which `service-b`'s tracing library reads, correlating the logs and traces between both microservices.
