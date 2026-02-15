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

## Example

```go
package main

import (
    "github.com/a-digi/coco-server/server"
)

func main() {
    srv := server.New()
    srv.GET("/", func(ctx *server.Context) {
        ctx.String(200, "Hello, world!")
    })
    srv.Run(8080)
}
```

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

## Documentation

Full documentation can be found at [pkg.go.dev](https://pkg.go.dev/github.com/a-digi/coco-server).

## License

MIT License – see [LICENSE](LICENSE)
