package routing_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/a-digi/coco-server/server/routing"
)

func TestRoutingBuilder_AddRouteAndServe(t *testing.T) {
	rb := routing.NewRoutingBuilder()
	rb.AddRoute("GET", "/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		w.Write([]byte("ok"))
	})
	r := httptest.NewRequest("GET", "/test", nil)
	rw := httptest.NewRecorder()
	rb.ServeHTTP(rw, r)
	if rw.Code != 201 || rw.Body.String() != "ok" {
		t.Errorf("unexpected response: %d, %s", rw.Code, rw.Body.String())
	}
}

func TestRoutingBuilder_NotFound(t *testing.T) {
	rb := routing.NewRoutingBuilder()
	r := httptest.NewRequest("GET", "/notfound", nil)
	rw := httptest.NewRecorder()
	rb.ServeHTTP(rw, r)
	if rw.Code != 404 {
		t.Errorf("expected 404, got %d", rw.Code)
	}
}
