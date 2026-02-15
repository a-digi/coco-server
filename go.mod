module github.com/a-digi/coco-server

go 1.25.0

require (
	github.com/a-digi/coco-logger v0.1.0
	github.com/a-digi/coco-orm v0.1.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/mattn/go-sqlite3 v1.14.34 // indirect
)

replace github.com/a-digi/coco-server/di => ./server/di

replace github.com/a-digi/coco-server/routing => ./server/routing
