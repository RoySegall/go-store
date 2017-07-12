package api

import (
	"net/http"
	"encoding/json"
)

// Display an error easily.
func WriteError(writer http.ResponseWriter, errorMessage string) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusBadRequest)

	response, _ := json.Marshal(map[string] string {
		"data": errorMessage,
	})

	writer.Write(response)
}
