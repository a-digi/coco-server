# coco-server

[![Go Reference](https://pkg.go.dev/badge/github.com/a-digi/coco-server.svg)](https://pkg.go.dev/github.com/a-digi/coco-server)

A lightweight, modular Go server framework.

## Installation

```bash
go get github.com/a-digi/coco-server
```

## Features
- Simple server configuration
- Routing and middleware support
- Request and response handling
- Dependency injection
- Logging
- Security layer (request authorization)

## Example

```go
package main

import (
    "github.com/a-digi/coco-server/server"
    "github.com/a-digi/coco-server/server/security"
)

type MySecurityLayer struct{}

func (s *MySecurityLayer) Authorize(w http.ResponseWriter, r *http.Request, ctx *server.Context) error {
    // Implement your authorization logic here
    // Return an error to deny access
    return nil // Allow all requests
}

func main() {
    srv := server.New()
    srv.SecurityLayer = &MySecurityLayer{}
    srv.GET("/", func(ctx *server.Context) {
        ctx.String(200, "Hello, world!")
    })
    srv.Run(8080)
}
```

## Security Layer

The security layer allows you to add custom request authorization logic. Implement the `Authorize` method and assign your security layer to the server. If `Authorize` returns an error, the request is denied with HTTP 403 Forbidden.

## Directory Structure

- `server/` – Main framework logic
  - `config.go` – Configuration
  - `log.go` – Logging
  - `pid.go` – Process management
  - `port.go` – Port management
  - `server.go` – Server logic
  - `di/` – Dependency injection
  - `request/` – Request handling
  - `response/` – Response handling
  - `routing/` – Routing and route management
  - `security/` – Security layer logic

## Documentation

Full documentation can be found at [pkg.go.dev](https://pkg.go.dev/github.com/a-digi/coco-server).

## License

MIT License – see [LICENSE](LICENSE)
