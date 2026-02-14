package response_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/a-digi/coco-server/server/response"
)

func TestSuccessResponse(t *testing.T) {
	rw := httptest.NewRecorder()
	response.SuccessResponse(rw, http.StatusOK, "ok")
	if rw.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rw.Code)
	}
	var body map[string]interface{}
	if err := json.NewDecoder(bytes.NewReader(rw.Body.Bytes())).Decode(&body); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if !body["success"].(bool) || body["message"] != "ok" {
		t.Errorf("unexpected body: %+v", body)
	}
}

func TestErrorResponse(t *testing.T) {
	rw := httptest.NewRecorder()
	response.ErrorResponse(rw, 400, "fail")
	if rw.Code != 400 {
		t.Errorf("expected status 400, got %d", rw.Code)
	}
	var body map[string]interface{}
	if err := json.NewDecoder(bytes.NewReader(rw.Body.Bytes())).Decode(&body); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if !body["error"].(bool) || body["message"] != "fail" {
		t.Errorf("unexpected body: %+v", body)
	}
}

func TestNotFoundResponse(t *testing.T) {
	rw := httptest.NewRecorder()
	response.NotFoundResponse(rw, "not found")
	if rw.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rw.Code)
	}
	var body map[string]interface{}
	if err := json.NewDecoder(bytes.NewReader(rw.Body.Bytes())).Decode(&body); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if !body["error"].(bool) || body["message"] != "not found" {
		t.Errorf("unexpected body: %+v", body)
	}
}
