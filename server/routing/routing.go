package routing

import (
	"net/http"

	serverdi "github.com/a-digi/coco-server/server/di"
	"github.com/a-digi/coco-server/server/security"
	"github.com/a-digi/coco-server/server/response"
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

func (rb *RoutingBuilder) authorizeRequest(w http.ResponseWriter, r *http.Request) bool {
	if rb.SecurityLayer == nil {
		return true
	}

	if err := rb.SecurityLayer.Authorize(w, r, rb.Context); err != nil {
		return false
	}

	return true
}

func (rb *RoutingBuilder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		// Handle CORS preflight
		response.BuildHeaders(w)
		w.WriteHeader(http.StatusOK)
		return
	}
	if !rb.authorizeRequest(w, r) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

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