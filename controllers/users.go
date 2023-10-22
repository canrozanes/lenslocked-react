package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/canrozanes/lenslocked/context"
	"github.com/canrozanes/lenslocked/errors"
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
		errors.WriteErrorResponse(w, err)

		return
	}

	user, err := u.UserService.Create(data.Email, data.Password)

	if err != nil {
		if errors.Is(err, models.ErrEmailTaken) {
			err = errors.Public(err, "That email address is already associated with an account.")
			w.WriteHeader(http.StatusConflict) // 409
		} else {
			w.WriteHeader(http.StatusBadRequest) // 500
		}
		fmt.Println(err.Error())
		errors.WriteErrorResponse(w, err)
		return
	}

	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		// TODO: Long term, we should show a warning about not being able to sign the user in.
		// This is not likely to happen
		w.WriteHeader(http.StatusInternalServerError)
		err = errors.Public(err, "User was created but couldn't be logged in.")
		errors.WriteErrorResponse(w, err)

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
		errors.WriteErrorResponse(w, err)
		return
	}

	user, err := u.UserService.Authenticate(data.Email, data.Password)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Println(err.Error())
		err = errors.Public(err, "Invalid email/password")
		errors.WriteErrorResponse(w, err)
		return
	}

	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		errors.WriteErrorResponse(w, err)
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
		err = errors.Public(err, "No signed in user found.")
		errors.WriteErrorResponse(w, err)
		return
	}
	err = u.SessionService.Delete(token)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		err = errors.Public(err, "Couldn't delete token")
		errors.WriteErrorResponse(w, err)
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
		err = errors.Public(err, "Request is missing email.")
		errors.WriteErrorResponse(w, err)
		return
	}

	pwReset, err := u.PasswordResetService.Create(data.Email)
	if err != nil {
		// TODO: Handle other cases in the future. For instance, if a user does not exist with that email address.
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		err = errors.Public(err, "Couldn't find the user")
		errors.WriteErrorResponse(w, err)
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
		err = errors.Public(err, "Couldn't send the email")
		errors.WriteErrorResponse(w, err)
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
			// TODO: Public requires an error. But we don't have an error here so we create a dummy one
			err := errors.Public(errors.New(""), "No signed in user found.")
			errors.WriteErrorResponse(w, err)
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
		err = errors.Public(err, "Request is missing token or password.")
		errors.WriteErrorResponse(w, err)
		return
	}

	user, err := u.PasswordResetService.Consume(data.Token)
	if err != nil {
		fmt.Println(err)
		// TODO: Distinguish between server errors and invalid token errors.
		w.WriteHeader(http.StatusInternalServerError)
		err = errors.Public(err, "Token is wrong.")
		errors.WriteErrorResponse(w, err)
		return

	}

	err = u.UserService.UpdatePassword(user.ID, data.Password)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		err = errors.Public(err, "Token is wrong.")
		errors.WriteErrorResponse(w, err)
		return
	}

	// Sign the user in now that they have reset their password.
	// Any errors from this point onward should redirect to the sign in page.
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		err = errors.Public(err, "Password Reset but login failed.")
		errors.WriteErrorResponse(w, err)
		return
	}
	setCookie(w, CookieSession, session.Token)
	write.Success(w)
}
