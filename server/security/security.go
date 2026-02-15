package security

import (
	"net/http"

	"github.com/a-digi/coco-server/server/di"
)

type SecurityLayer interface {
	Authorize(w http.ResponseWriter, r *http.Request, ctx di.Context) error
}
