package routing

import (
	"net/http"

	serverdi "github.com/a-digi/coco-server/server/di"
	"github.com/a-digi/coco-server/server/security"
)

type RouteHandler struct {
	Method  string
	Pattern string
	Handler http.HandlerFunc
}

type RoutingBuilder struct {
	routes        []RouteHandler
	Context       serverdi.Context
	SecurityLayer security.SecurityLayer
}

func NewRoutingBuilder() *RoutingBuilder {
	return &RoutingBuilder{
		routes: make([]RouteHandler, 0),
	}
}

func (rb *RoutingBuilder) AddContext(ctx serverdi.Context) {
	rb.Context = ctx
}

func (rb *RoutingBuilder) SetSecurityLayer(layer security.SecurityLayer) {
	rb.SecurityLayer = layer
}

func (rb *RoutingBuilder) AddRoute(method, pattern string, handler http.HandlerFunc) {
	rb.routes = append(rb.routes, RouteHandler{Method: method, Pattern: pattern, Handler: handler})
}

func (rb *RoutingBuilder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range rb.routes {
		if r.URL.Path == route.Pattern && r.Method == route.Method {
			if rb.SecurityLayer != nil {
				if err := rb.SecurityLayer.Authorize(w, r, rb.Context); err != nil {
					// Authorization failed, response already handled or block
					return
				}
			}
			route.Handler(w, r)
			return
		}
	}

	http.NotFound(w, r)
}

func (rb *RoutingBuilder) Build() http.Handler {
	return rb
}