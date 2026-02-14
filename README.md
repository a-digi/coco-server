# coco-server

[![Go Reference](https://pkg.go.dev/badge/github.com/a-digi/coco-server.svg)](https://pkg.go.dev/github.com/a-digi/coco-server)

Ein leichtgewichtiges, modular aufgebautes Go-Server-Framework.

## Installation

```bash
go get github.com/a-digi/coco-server
```

## Features
- Einfache Server-Konfiguration
- Routing und Middleware-Unterstützung
- Request- und Response-Handling
- Dependency Injection
- Logging

## Beispiel

```go
package main

import (
    "github.com/a-digi/coco-server/server"
)

func main() {
    srv := server.New()
    srv.GET("/", func(ctx *server.Context) {
        ctx.String(200, "Hallo, Welt!")
    })
    srv.Run(8080)
}
```

## Verzeichnisstruktur

- `server/` – Hauptlogik des Frameworks
  - `config.go` – Konfiguration
  - `log.go` – Logging
  - `pid.go` – Prozessmanagement
  - `port.go` – Portverwaltung
  - `server.go` – Server-Logik
  - `di/` – Dependency Injection
  - `request/` – Request-Handling
  - `response/` – Response-Handling
  - `routing/` – Routing und Routenverwaltung

## Dokumentation

Die vollständige Dokumentation findest du auf [pkg.go.dev](https://pkg.go.dev/github.com/a-digi/coco-server).

## Lizenz

MIT License – siehe [LICENSE](LICENSE)
