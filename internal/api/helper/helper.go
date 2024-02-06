package helper

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Errors []string `json:"errors"`
}

func ResponseErrors(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	messages := []string{err.Error()}

	json.NewEncoder(w).Encode(ErrorResponse{Errors: messages}) //nolint:errcheck,errchkjson
}

func ResponseInternalError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	messages := []string{err.Error()}

	json.NewEncoder(w).Encode(ErrorResponse{Errors: messages}) //nolint:errcheck,errchkjson
}

func ResponseOk(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			ResponseErrors(w, err)
		}
	}
}

func ResponseNoContent(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func ResponseNotFound(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		ResponseErrors(w, err)
	}
}
