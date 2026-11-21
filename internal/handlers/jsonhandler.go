package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func UnmarshallJson(rbody io.ReadCloser, dest any) error {
	baits, err := io.ReadAll(rbody)
	if err != nil {
		return fmt.Errorf("Invalid request Body")
	}
	err = json.Unmarshal(baits, dest)
	if err != nil {
		return fmt.Errorf("Unable to unmarshall request Body")
	}
	return nil
}

func WriteJsonResponse(w http.ResponseWriter, statusCode int, jsonBody any) error {
	w.WriteHeader(statusCode)
	data, err := json.Marshal(jsonBody)
	if err != nil {
		return fmt.Errorf("Error writting the response")
	}
	w.Write(data)
	return nil
}

func WriteError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	w.Header().Add("Content-Type", "text/html")
	w.Header().Add("Content-Type", "charset=utf-8")
	w.Write([]byte(message))
}