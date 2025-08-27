package utils

import (
	"encoding/json"
	"net/http"
	"os"
)

func SendJSONResponse(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if os.Getenv("ENV") != "production" {
		// prettify JSON for easier debugging
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		encoder.Encode(data)
	} else {
		// in production, compact JSON for smaller payload
		json.NewEncoder(w).Encode(data)
	}
}

func SendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	response := map[string]any{
		"error": message,
	}
	SendJSONResponse(w, statusCode, response)
}

func SendNoContentResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
