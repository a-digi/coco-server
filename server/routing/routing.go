package routing

import (
	"net/http"

	serverdi "github.com/a-digi/coco-server/di"
)

type RouteHandler struct {
	Method  string
	Pattern string
	Handler http.HandlerFunc
}

type RoutingBuilder struct {
	routes   []RouteHandler
	Context  serverdi.Context
}

func NewRoutingBuilder() *RoutingBuilder {
	return &RoutingBuilder{
		routes: make([]RouteHandler, 0),
	}
}

func (rb *RoutingBuilder) AddContext(ctx serverdi.Context) {
	rb.Context = ctx
}

func (rb *RoutingBuilder) AddRoute(method, pattern string, handler http.HandlerFunc) {
	rb.routes = append(rb.routes, RouteHandler{Method: method, Pattern: pattern, Handler: handler})
}

func (rb *RoutingBuilder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range rb.routes {
		if r.URL.Path == route.Pattern && r.Method == route.Method {
			route.Handler(w, r)
			return
		}
	}

	http.NotFound(w, r)
}

func (rb *RoutingBuilder) Build() http.Handler {
	return rb
}