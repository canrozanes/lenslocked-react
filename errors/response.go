package errors

import (
	"encoding/json"
	"errors"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func WriteErrorResponse(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&ErrorResponse{Error: errMessage(err)})
}

func errMessage(err error) string {
	var pubErr publicError
	if errors.As(err, &pubErr) {
		return pubErr.Public()
	} else {
		return "Something went wrong."
	}
}
