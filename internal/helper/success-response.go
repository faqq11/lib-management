package helper

import (
	"encoding/json"
	"net/http"
)

func SuccessResponse(writer http.ResponseWriter, statusCode int, payload map[string]interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	json.NewEncoder(writer).Encode(payload)
}