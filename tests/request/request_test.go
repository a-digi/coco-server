package request_test

import (
	"net/http"
	"net/url"
	"testing"
	"github.com/a-digi/coco-server/server/request"
)

type testStruct struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestMapFormToStruct(t *testing.T) {
	form := url.Values{}
	form.Set("name", "Max")
	form.Set("age", "42")
	r, _ := http.NewRequest("POST", "/", nil)
	r.Form = form
	var ts testStruct
	err := request.MapFormToStruct(r, &ts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ts.Name != "Max" || ts.Age != 42 {
		t.Errorf("unexpected struct: %+v", ts)
	}
}

func TestMapFormToStruct_InvalidDest(t *testing.T) {
	r, _ := http.NewRequest("POST", "/", nil)
	err := request.MapFormToStruct(r, nil)
	if err == nil {
		t.Error("expected error for nil dest")
	}
}
