package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/canrozanes/lenslocked/models"
	"github.com/canrozanes/lenslocked/write"
)

type Users struct {
	UserService *models.UserService
}

type userResponse struct {
	User *models.User `json:"user"`
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {

	var data struct {
		Email    string
		Password string
	}

	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err.Error())
		json.NewEncoder(w).Encode(&write.ErrorResponse{Error: "Something went wrong."})

		return
	}

	user, err := u.UserService.Create(data.Email, data.Password)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err.Error())
		json.NewEncoder(w).Encode(&write.ErrorResponse{Error: "Something went wrong."})

		return
	}

	json.NewEncoder(w).Encode(userResponse{User: user})
}

func (u Users) ProcessSignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Password string
	}

	err := json.NewDecoder(r.Body).Decode(&data)

	fmt.Println(data)
	fmt.Println("we are past the decode")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err.Error())
		json.NewEncoder(w).Encode(&write.ErrorResponse{Error: "Something went wrong."})
		return
	}

	user, err := u.UserService.Authenticate(data.Email, data.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Println(err.Error())
		json.NewEncoder(w).Encode(&write.ErrorResponse{Error: "Not authorized."})
		return
	}

	fmt.Println(user)

	json.NewEncoder(w).Encode(userResponse{User: user})
}
