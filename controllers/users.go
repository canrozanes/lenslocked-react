package controllers

import (
	"encoding/json"
	"net/http"
)

type Users struct{}

func (u Users) New(w http.ResponseWriter, r *http.Request) {

	var data struct {
		Email    string
		Password string
	}

	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(data)
}
