package routing

import (
	"net/http"
	"bytes"

	"gopkg.in/yaml.v3"
	"github.com/a-digi/coco-logger/logger"
	serverdi "github.com/a-digi/coco-server/server/di"
	"github.com/a-digi/coco-server/server/security"
	"github.com/a-digi/coco-server/server/response"
)

type HandlerInterface interface {
	ServeHTTP(http.ResponseWriter, *http.Request, serverdi.Context)
}

type RouteConfig struct {
	Path        string        `yaml:"path"`
	Method      string        `yaml:"method"`
	Executor    string        `yaml:"executor"`
	ContentType string        `yaml:"content_type"`
	Children    []RouteConfig `yaml:"children"`
}

type RoutesYAML struct {
	APIPrefix string        `yaml:"api_prefix"`
	Routes    []RouteConfig `yaml:"routes"`
}

func flattenRoutes(currentItem RouteConfig, prefix string, parentConfig RouteConfig) []RouteConfig {
	var routes []RouteConfig

	currentPath := prefix
	if len(currentPath) > 0 && currentPath[len(currentPath)-1] == '/' {
		currentPath = currentPath[:len(currentPath)-1]
	}
	if len(currentItem.Path) > 0 && currentItem.Path[0] != '/' {
		currentPath += "/" + currentItem.Path
	} else {
		currentPath += currentItem.Path
	}

	currentContentType := currentItem.ContentType

	if currentContentType == "" {
		currentContentType = parentConfig.ContentType
	}

    if currentItem.Method != "" && currentItem.Executor != "" {
        route := currentItem
        route.Path = currentPath
        route.ContentType = currentContentType
        routes = append(routes, route)
    }

	if len(currentItem.Children) > 0 {
		for _, child := range currentItem.Children {
			childRoutes := flattenRoutes(child, currentPath, parentConfig)
			routes = append(routes, childRoutes...)
		}
	}

	return routes
}

type Route struct {
	Path        string
	Method      string
	Executor    string
	ContentType string
}

type RouteBuilder struct {
	routes     []Routes
	HandlerMap map[string]HandlerInterface
	Context    serverdi.Context
	SecurityLayer security.SecurityLayer
}

func (rb *RouteBuilder) AddContext(ctx serverdi.Context) {
	rb.Context = ctx
}

func (rb *RouteBuilder) AddRoute(route Routes) {
	rb.routes = append(rb.routes, route)
}

func (rb *RouteBuilder) Build(log logger.Logger) http.Handler {
	return rb
}

func (rb *RouteBuilder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		// Handle CORS preflight
		response.BuildHeaders(w)
		w.WriteHeader(http.StatusOK)
		return
	}
	if rb.SecurityLayer != nil {
		if err := rb.SecurityLayer.Authorize(w, r, rb.Context); err != nil {
			// Authorization failed, response already handled or block
			return
		}
	}
	for _, route := range rb.routes {

		if len(route.YamlContent) == 0 || route.HandlerMap == nil {
            rb.Context.GetLogger().Warning("Route skipped: YamlContent empty or HandlerMap nil")
			continue
		}

		var yamlConfig RoutesYAML
		decoder := yaml.NewDecoder(bytes.NewReader(route.YamlContent))

		if err := decoder.Decode(&yamlConfig); err != nil {
            rb.Context.GetLogger().Warning("YAML could not be parsed: %v", err)
			continue
		}

		flatRoutes := []RouteConfig{}

		for _, parent := range yamlConfig.Routes {
			flatRoutes = append(flatRoutes, flattenRoutes(parent, "", parent)...)
		}

		for _, rc := range flatRoutes {
			if r.URL.Path == rc.Path && r.Method == rc.Method {
				handler, ok := route.HandlerMap[rc.Executor]
				if ok {
					rb.Context.GetLogger().Info("Request success: %s %s -> %s", r.Method, r.URL.Path, rc.Executor)
					handler.ServeHTTP(w, r, rb.Context)
					return
				}
			}
		}
	}

    rb.Context.GetLogger().Warning("Request not found: %s %s", r.Method, r.URL.Path)
	http.NotFound(w, r)
}

func (rb *RouteBuilder) SetSecurityLayer(layer security.SecurityLayer) {
	rb.SecurityLayer = layer
}

type Routes struct {
	YamlContent []byte
	HandlerMap  map[string]HandlerInterface
}

var GlobalRouteBuilder = &RouteBuilder{}

func RegisterRoutes(routeConfigs Routes, log logger.Logger, ctx serverdi.Context) *RoutingBuilder {
	if len(routeConfigs.YamlContent) == 0 {
		log.Error("No YAML content provided for routes")
		return NewRoutingBuilder()
	}
	var yamlConfig RoutesYAML
	decoder := yaml.NewDecoder(bytes.NewReader(routeConfigs.YamlContent))

	if err := decoder.Decode(&yamlConfig); err != nil {
		log.Error("Failed to parse routes.yaml: %v", err)
		return NewRoutingBuilder()
	}

	rb := NewRoutingBuilder()
	rb.AddContext(ctx)

	for _, parent := range yamlConfig.Routes {
		flatRoutes := flattenRoutes(parent, "", parent)
		for _, route := range flatRoutes {
			handler, ok := routeConfigs.HandlerMap[route.Executor]
			if !ok {
				log.Error("Unknown executor: %s", route.Executor)
				continue
			}
			rb.AddRoute(route.Method, route.Path, func(w http.ResponseWriter, r *http.Request) {
				if route.ContentType != "" && r.Header.Get("Content-Type") != route.ContentType {
					http.Error(w, "Content-Type must be "+route.ContentType, http.StatusUnsupportedMediaType)
					return
				}
				handler.ServeHTTP(w, r, ctx)
			})
		}
	}

	return rb
}
