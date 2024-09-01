'use strict';

const { NodeTracerProvider } = require('@opentelemetry/sdk-trace-node'); // Updated import
const { JaegerExporter } = require('@opentelemetry/exporter-jaeger');
const { registerInstrumentations } = require('@opentelemetry/instrumentation');
const { Resource } = require('@opentelemetry/resources');
const { SemanticResourceAttributes } = require('@opentelemetry/semantic-conventions');
const { SimpleSpanProcessor } = require('@opentelemetry/sdk-trace-base');
const { HttpInstrumentation } = require('@opentelemetry/instrumentation-http');
const { ExpressInstrumentation } = require('@opentelemetry/instrumentation-express');

// Initialize the provider
const provider = new NodeTracerProvider({
  resource: new Resource({
    [SemanticResourceAttributes.SERVICE_NAME]: 'service-b',
  }),
});

const JAEGER_ENDPOINT =  process.env.OTEL_EXPORTER_JAEGER_ENDPOINT

// Setup the exporter
const exporter = new JaegerExporter({
  endpoint: JAEGER_ENDPOINT, // Replace with the appropriate Jaeger collector endpoint
});

// Add the exporter to the provider
provider.addSpanProcessor(new SimpleSpanProcessor(exporter));

// Initialize the provider and instrumentations
provider.register();

registerInstrumentations({
  instrumentations: [
    new HttpInstrumentation({
      applyCustomAttributesOnSpan: (span, request, response) => {
        span.setAttribute('custom-attribute', 'custom-value');
      },
    }),
    new ExpressInstrumentation(), // Add this for Express.js instrumentation
  ],
});

console.log('Tracing initialized');
