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

func WriteJSONBytes(
	responseWriter http.ResponseWriter,
	statusCode int,
	response []byte,
) error {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(statusCode)

	_, err := responseWriter.Write(response)
	return err
}
