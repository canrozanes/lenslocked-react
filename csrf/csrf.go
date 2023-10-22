package customcsrf

import (
	"fmt"
	"net/http"

	"github.com/canrozanes/lenslocked/errors"
	"github.com/gorilla/csrf"
)

type ErrorHandler struct{}

func (c *ErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
	err := csrf.FailureReason(r)
	fmt.Println(err.Error())
	errors.WriteErrorResponse(w, err)
}

func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

func RefreshCSRFToken(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie := &http.Cookie{
			Name:     "csrf",
			Value:    csrf.Token(r),
			Path:     "/",
			HttpOnly: false,
		}

		http.SetCookie(w, cookie)
		h.ServeHTTP(w, r)
	})
}
