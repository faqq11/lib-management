package helper

import (
	"encoding/json"
	"net/http"
)

func SuccessResponse(writer http.ResponseWriter, statusCode int, payload interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	json.NewEncoder(writer).Encode(payload)
}