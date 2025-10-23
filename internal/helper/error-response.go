package helper

import (
	"encoding/json"
	"net/http"
)

func ErrorResponse(writer http.ResponseWriter, statusCode int, message string) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)

	res := map[string]string{
		"message": message,
	}

	json.NewEncoder(writer).Encode(res)
}