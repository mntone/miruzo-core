package response

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(
	responseWriter http.ResponseWriter,
	statusCode int,
	response any,
) error {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(statusCode)

	return json.NewEncoder(responseWriter).Encode(response)
}

func WriteJSONText(
	responseWriter http.ResponseWriter,
	statusCode int,
	response string,
) error {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(statusCode)

	_, err := responseWriter.Write([]byte(response))
	return err
}
