# go-http-middleware

A collection of composable HTTP middleware for Go `net/http` servers.

## Middleware

- **Logger** — Structured request/response logging with duration tracking
- **Recovery** — Panic recovery with configurable error responses
- **CORS** — Cross-origin resource sharing with fine-grained control
- **RateLimit** — Token bucket rate limiting per client IP

## Installation

```bash
go get github.com/marwantech91/go-http-middleware
```

## Usage

```go
package main

import (
    "net/http"
    mw "github.com/marwantech91/go-http-middleware"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("ok"))
    })

    handler := mw.Chain(
        mux,
        mw.Logger(),
        mw.Recovery(),
        mw.CORS(mw.CORSOptions{AllowOrigins: []string{"*"}}),
        mw.RateLimit(100, 60),
    )

    http.ListenAndServe(":8080", handler)
}
```

## License

MIT
