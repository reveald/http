# reveald-http

A Go HTTP library that exposes the [reveald](https://github.com/reveald/reveald) search backend via HTTP endpoints. Build powerful search APIs with minimal configuration.

## Features

- **Simple HTTP Wrapper**: Expose reveald search functionality through HTTP endpoints
- **Flexible Routing**: Configure multiple routes with different search capabilities
- **Feature-Rich Search**: Built-in support for pagination, sorting, filtering, and histograms
- **Middleware Support**: Standard HTTP middleware pattern for authentication, logging, and more
- **Elasticsearch Backend**: Powered by Elasticsearch through the reveald abstraction layer
- **OpenTelemetry Ready**: Auto-instrumentation support for observability
- **Customizable**: Configure parameter readers, loggers, and backend options

## Installation

```bash
go get github.com/reveald/http/v2
```

## Quick Start

```go
package main

import (
    "fmt"
    "net/http"

    revealdhttp "github.com/reveald/http/v2"
    "github.com/reveald/reveald/featureset"
    "github.com/reveald/reveald/v2"
)

func main() {
    // Create an Elasticsearch backend
    backend, err := reveald.NewElasticBackend([]string{"http://localhost:9200"})
    if err != nil {
        panic(err)
    }

    // Create the HTTP API
    api := revealdhttp.New(backend)

    // Configure a search route
    api.Route("/search",
        revealdhttp.WithIndex("my-index"),
        revealdhttp.WithFeatures(
            featureset.NewPaginationFeature(),
            featureset.NewSortingFeature("sort"),
        ))

    // Start the server
    http.ListenAndServe(":8080", api)
}
```

Visit `http://localhost:8080/search` to see search results.

## API Reference

### Creating an API

```go
api := revealdhttp.New(backend, options...)
```

**Options:**
- `WithLogger(logger)` - Set a custom logger
- `WithParamReader(readers...)` - Set custom parameter readers for extracting search parameters from HTTP requests

### Configuring Routes

```go
api.Route(pattern, options...)
```

**Route Options:**
- `WithIndex(index)` - Specify which Elasticsearch index to search (can be called multiple times)
- `WithFeatures(features...)` - Add reveald features (pagination, sorting, filtering, etc.)
- `WithMiddleware(handler)` - Add HTTP middleware for the route

### Available Features

Features are provided by the `github.com/reveald/reveald/featureset` package:

- **Pagination**: `featureset.NewPaginationFeature()`
- **Sorting**: `featureset.NewSortingFeature(paramName, options...)`
- **Static Filtering**: `featureset.NewStaticFilterFeature(options...)`
  - `WithRequiredProperty(field)` - Filter by field presence
  - `WithRequiredValue(field, value)` - Filter by exact value
- **Dynamic Filtering**: `featureset.NewDynamicFilterFeature(field)` - User-controlled filtering via query parameters
- **Histograms**: `featureset.NewHistogramFeature(field, options...)`

## Examples

### Advanced Search Endpoint

```go
api.Route("/products/search",
    // Add custom middleware
    revealdhttp.WithMiddleware(authMiddleware),

    // Search in multiple indices
    revealdhttp.WithIndex("products"),
    revealdhttp.WithIndex("product-variants"),

    // Add search features
    revealdhttp.WithFeatures(
        featureset.NewPaginationFeature(),
        featureset.NewSortingFeature("sort",
            featureset.WithSortOption("by-price", "price", false),
            featureset.WithSortOption("by-date", "created_at", true)),
        featureset.NewStaticFilterFeature(
            featureset.WithRequiredValue("status", "active")),
        featureset.NewDynamicFilterFeature("category"),
        featureset.NewDynamicFilterFeature("brand"),
        featureset.NewHistogramFeature("price",
            featureset.WithInterval(100)),
    ))
```

### Custom Middleware

```go
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

### Custom Logger

```go
type CustomLogger struct{}

func (l *CustomLogger) Errorf(format string, args ...any) {
    log.Printf("ERROR: "+format, args...)
}

api := revealdhttp.New(backend,
    revealdhttp.WithLogger(&CustomLogger{}))
```

## Development

### Prerequisites

- Go 1.25+
- Docker (for running Elasticsearch locally)

### Local Development

```bash
# Start local Elasticsearch
./examples/local-elasticsearch.sh

# Seed with test data
./examples/seed-elasticsearch.sh

# Run the example server
go run examples/main.go

# Run tests
go test ./...

# Build
go build ./...
```

### Running the Example

The example server demonstrates a complete setup with pagination, sorting, filtering, and histograms:

```bash
go run examples/main.go
```

Then navigate to:
- `http://localhost:8080/search` - View search results
- `http://localhost:8080/search?text_field=Third` - Filtered results
- Add `Authorization` header to see middleware in action

## Architecture

The library follows a clean separation of concerns:

- **HTTPAPI** (`server.go`): Main HTTP server managing routes and middleware
- **Route** (`route.go`): Route configuration with indices, features, and middleware
- **Handler** (`handler.go`): Creates reveald endpoints and executes searches
- **ParamReader** (`reader.go`): Extracts search parameters from HTTP requests
- **Result** (`result.go`): Formats search results as JSON responses

## Dependencies

- [github.com/reveald/reveald](https://github.com/reveald/reveald) - Core search functionality
- [github.com/elastic/go-elasticsearch/v8](https://github.com/elastic/go-elasticsearch) - Elasticsearch client

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## Related Projects

- [reveald](https://github.com/reveald/reveald) - The core search backend library
