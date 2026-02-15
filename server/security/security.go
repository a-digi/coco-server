package security

import (
	"net/http"

	"github.com/a-digi/coco-server/server/di"
)

// SecurityLayer defines an interface for security checks before route handling.
type SecurityLayer interface {
	// Authorize is called before the route handler. Return an error to block the request.
	Authorize(w http.ResponseWriter, r *http.Request, ctx di.Context) error
}
