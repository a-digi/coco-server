package security

import (
	"net/http"

	"github.com/a-digi/coco-server/server/di"
)

type Route struct {
	Path          string
	Method        string
	Security      string
	Scopes        []string
	PathVariables map[string]string
}

type SecurityLayer interface {
	Authorize(w http.ResponseWriter, r *http.Request, ctx di.Context, route *Route) error
}
