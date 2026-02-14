package response

import (
	"encoding/json"
	"net/http"
)

func SuccessResponse(w http.ResponseWriter, status int, resp interface{}) {
	buildHeaders(w)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
        "success":   true,
        "message": resp,
    })
}

func ErrorResponse(w http.ResponseWriter, status int, message string) {
	buildHeaders(w)
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":   true,
		"message": message,
	})
}

func NotFoundResponse(w http.ResponseWriter, message string) {
	buildHeaders(w)
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":   true,
		"message": message,
	})
}

func buildHeaders(w http.ResponseWriter) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Content-Type", "application/json")
}
