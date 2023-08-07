package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/canrozanes/lenslocked/controllers"
	customcsrf "github.com/canrozanes/lenslocked/csrf"
	"github.com/canrozanes/lenslocked/models"
	"github.com/canrozanes/lenslocked/spa"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
)

type Resource struct {
	Route string `json:"route"`
}

func getApiRouter(db *sql.DB) chi.Router {
	r := chi.NewRouter()

	// Setup our model services
	userService := models.UserService{
		DB: db,
	}

	sessionService := models.SessionService{
		DB: db,
	}

	// Setup our controllers
	usersC := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{ "message": "bar" }`))
	})

	r.Post("/users", usersC.Create)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Post("/signout", usersC.ProcessSignOut)
	r.Get("/users/me", usersC.CurrentUser)

	r.Get("/{resource}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		resource := chi.URLParam(r, "resource")

		response := Resource{resource}
		json.NewEncoder(w).Encode(response)

	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Route not found", http.StatusNotFound)
	})

	return r
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	r.Mount("/api", getApiRouter(db))

	// we want all routes besides /api to go to the SPA, hence we use the NotFound handler
	r.NotFound(spa.SpaHandler)

	csrfKey := "gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX"
	csrfErrorHandler := customcsrf.NewErrorHandler()
	csrfMw := csrf.Protect(
		[]byte(csrfKey),
		// TODO: Fix this before deploying
		csrf.Secure(false),
		csrf.ErrorHandler(csrfErrorHandler),
	)

	skipCsrf := customcsrf.NewSkipper()

	fmt.Println("Starting the server on :3000...")

	http.ListenAndServe(":3000", skipCsrf(csrfMw(r)))
}
