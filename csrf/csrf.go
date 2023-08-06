package customcsrf

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/canrozanes/lenslocked/write"
	"github.com/gorilla/csrf"
)

type ErrorHandler struct{}

func (c *ErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(&write.ErrorResponse{Error: "CSRF token invalid."})
}

func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

type SkipCSRF struct {
	h http.Handler
}

func NewSkipper() func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		sk := &SkipCSRF{
			h: h,
		}
		return sk
	}
}

// TODO: Modify skipCsrf mw after employing environment variables.
// Code should skip csrf during development and keep it for production
func (sr *SkipCSRF) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.Referer(), "http://localhost:5173/") {
		fmt.Println("csrf by-passed due to local dev")
		r = csrf.UnsafeSkipCheck(r)
	}
	sr.h.ServeHTTP(w, r)
}
