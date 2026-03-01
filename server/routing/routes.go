package routing

import (
	"bytes"
	"net/http"
	"strings"

	"github.com/a-digi/coco-logger/logger"
	serverdi "github.com/a-digi/coco-server/server/di"
	"github.com/a-digi/coco-server/server/request"
	"github.com/a-digi/coco-server/server/response"
	"github.com/a-digi/coco-server/server/security"
	"gopkg.in/yaml.v3"
)

type HandlerInterface interface {
	ServeHTTP(ctx request.RequestContext)
}

type RouteConfig struct {
	Path        string        `yaml:"path"`
	Method      string        `yaml:"method"`
	Executor    string        `yaml:"executor"`
	ContentType string        `yaml:"content_type"`
	Security    string        `yaml:"security"`
	Scopes      []string      `yaml:"scopes"`
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

	var currentScopes []string
	if len(parentConfig.Scopes) > 0 {
		currentScopes = append(currentScopes, parentConfig.Scopes...)
	}
	if len(currentItem.Scopes) > 0 {
		currentScopes = append(currentScopes, currentItem.Scopes...)
	}

	currentSecurity := currentItem.Security
	if currentSecurity == "" {
		currentSecurity = parentConfig.Security
	}

	if currentItem.Method != "" && currentItem.Executor != "" {
		route := currentItem
		route.Path = currentPath
		route.ContentType = currentContentType
		route.Security = currentSecurity
		route.Scopes = currentScopes
		routes = append(routes, route)
	}

	if len(currentItem.Children) > 0 {
		for _, child := range currentItem.Children {
			// Pass currentItem as parentConfig to propagate inherited scopes
			childCopy := currentItem
			childCopy.Scopes = currentScopes
			childRoutes := flattenRoutes(child, currentPath, childCopy)
			routes = append(routes, childRoutes...)
		}
	}

	return routes
}

func createSecRoute(reqPath string, method string, flatRoutes []RouteConfig) *security.Route {
	for _, fRoute := range flatRoutes {
		if fRoute.Method == method && MatchPath(fRoute.Path, reqPath) {
			exact := true
			pSegs := strings.Split(fRoute.Path, "/")
			rSegs := strings.Split(reqPath, "/")
			for i := 0; i < len(pSegs) && i < len(rSegs); i++ {
				if pSegs[i] != rSegs[i] {
					if strings.HasPrefix(pSegs[i], "{res:") && strings.HasPrefix(rSegs[i], "{res:") {
						exact = false
						break
					}
				}
			}
			if exact {
				return &security.Route{
					Path:     fRoute.Path,
					Method:   fRoute.Method,
					Security: fRoute.Security,
					Scopes:   fRoute.Scopes,
				}
			}
		}
	}
	return nil
}

type Route struct {
	Path        string
	Method      string
	Executor    string
	ContentType string
}

type RouteBuilder struct {
	routes        []Routes
	HandlerMap    map[string]HandlerInterface
	Context       serverdi.Context
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
		response.BuildHeaders(w)
		w.WriteHeader(http.StatusOK)
		return
	}
	// Security check moved inside route match loop
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
			if !MatchPath(rc.Path, r.URL.Path) || r.Method != rc.Method {
				continue
			}

			handler, ok := route.HandlerMap[rc.Executor]
			if !ok {
				continue
			}

			reqCtx := request.NewContext(w, r, rb.Context)
			reqCtx.GetURI().ExtractPathVariables(rc.Path)
			rb.Context.GetLogger().Info("Request success: %s %s -> %s", r.Method, r.URL.Path, rc.Executor)
			if rb.SecurityLayer == nil {
				handler.ServeHTTP(reqCtx)
				return
			}

			secRoute := createSecRoute(reqCtx.GetURI().Path, r.Method, flatRoutes)
			if secRoute == nil {
				rb.Context.GetLogger().Warning("Request denied: %s %s", r.Method, r.URL.Path)
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}

			if err := rb.SecurityLayer.Authorize(w, r, rb.Context, secRoute); err != nil {
				rb.Context.GetLogger().Warning("Request denied: %s %s", r.Method, r.URL.Path)
				return
			}

			handler.ServeHTTP(reqCtx)
			return
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
			rb.AddRoute(route.Method, route.Path, route.Security, route.Scopes, func(w http.ResponseWriter, r *http.Request) {
				if route.ContentType != "" && r.Header.Get("Content-Type") != route.ContentType {
					http.Error(w, "Content-Type must be "+route.ContentType, http.StatusUnsupportedMediaType)
					return
				}
				reqCtx := request.NewContext(w, r, ctx)
				reqCtx.GetURI().ExtractPathVariables(route.Path)
				handler.ServeHTTP(reqCtx)
			})
		}
	}

	return rb
}
