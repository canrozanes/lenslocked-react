package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/canrozanes/lenslocked/context"
	"github.com/canrozanes/lenslocked/models"
	"github.com/canrozanes/lenslocked/write"
)

type Users struct {
	UserService          *models.UserService
	SessionService       *models.SessionService
	PasswordResetService *models.PasswordResetService
	EmailService         *models.EmailService
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

	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		// TODO: Long term, we should show a warning about not being able to sign the user in.
		// This is not likely to happen
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&write.ErrorResponse{Error: "User was created but couldn't be logged in."})
		return
	}

	setCookie(w, CookieSession, session.Token)

	json.NewEncoder(w).Encode(userResponse{User: user})

}

func (u Users) ProcessSignIn(w http.ResponseWriter, r *http.Request) {
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

	user, err := u.UserService.Authenticate(data.Email, data.Password)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Println(err.Error())
		json.NewEncoder(w).Encode(&write.ErrorResponse{Error: "No signed in user found."})
		return
	}

	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		json.NewEncoder(w).Encode(&write.ErrorResponse{Error: "Something went wrong."})
		return
	}

	setCookie(w, CookieSession, session.Token)

	json.NewEncoder(w).Encode(userResponse{User: user})
}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())

	json.NewEncoder(w).Encode(userResponse{User: user})
}

func (u Users) ProcessSignOut(w http.ResponseWriter, r *http.Request) {
	token, err := readCookie(r, CookieSession)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&write.ErrorResponse{Error: "No signed in user found."})
		return
	}
	err = u.SessionService.Delete(token)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&write.ErrorResponse{Error: "Couldn't delete token"})
		return
	}
	deleteCookie(w, CookieSession)

	write.Success(w)
}

func (u Users) ProcessForgotPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}

	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err.Error())
		json.NewEncoder(w).Encode(&write.ErrorResponse{Error: "Request is missing email."})
		return
	}

	fmt.Println(data.Email)
	pwReset, err := u.PasswordResetService.Create(data.Email)
	if err != nil {
		// TODO: Handle other cases in the future. For instance, if a user does not exist with that email address.
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&write.ErrorResponse{Error: "Couldn't find the user"})
		return
	}
	vals := url.Values{
		"token": {pwReset.Token},
	}
	resetURL := "https://www.lenslocked.com/reset-pw?" + vals.Encode()
	err = u.EmailService.ForgotPassword(data.Email, resetURL)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&write.ErrorResponse{Error: "Couldn't send the email"})
		return
	}
	// Don't render the reset token here! We need the user to confirm they have
	// access to the email account to verify their identity.

	write.Success(w)
}

type UserMiddleware struct {
	SessionService *models.SessionService
}

func (umw UserMiddleware) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := readCookie(r, CookieSession)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		user, err := umw.SessionService.User(token)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithUser(r.Context(), user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)

	})
}

func (umw UserMiddleware) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())

		if user == nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(&write.ErrorResponse{Error: "No signed in user found."})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (u Users) ProcessResetPassword(w http.ResponseWriter, r *http.Request) {

	var data struct {
		Token    string
		Password string
	}
	data.Token = r.FormValue("token")
	data.Password = r.FormValue("password")

	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err.Error())
		json.NewEncoder(w).Encode(&write.ErrorResponse{Error: "Request is missing token or password."})
		return
	}

	user, err := u.PasswordResetService.Consume(data.Token)
	if err != nil {
		fmt.Println(err)
		// TODO: Distinguish between server errors and invalid token errors.
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&write.ErrorResponse{Error: "Token is wrong."})
		return

	}

	err = u.UserService.UpdatePassword(user.ID, data.Password)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&write.ErrorResponse{Error: "Token is wrong."})
		return
	}

	// Sign the user in now that they have reset their password.
	// Any errors from this point onward should redirect to the sign in page.
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&write.ErrorResponse{Error: "Password Reset but login failed."})
		return
	}
	setCookie(w, CookieSession, session.Token)
	write.Success(w)
}
