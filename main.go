package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/canrozanes/lenslocked/controllers"
	customcsrf "github.com/canrozanes/lenslocked/csrf"
	"github.com/canrozanes/lenslocked/migrations"
	"github.com/canrozanes/lenslocked/models"
	"github.com/canrozanes/lenslocked/spa"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"github.com/gorilla/handlers"
)

func getApiRouter(db *sql.DB) chi.Router {
	r := chi.NewRouter()

	csrfKey := "gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX"
	csrfErrorHandler := customcsrf.NewErrorHandler()
	csrfMw := csrf.Protect(
		[]byte(csrfKey),
		// TODO: Fix this before deploying
		csrf.Secure(false),
		csrf.ErrorHandler(csrfErrorHandler),
		csrf.SameSite(csrf.SameSiteStrictMode),
	)

	userService := models.UserService{
		DB: db,
	}

	sessionService := models.SessionService{
		DB: db,
	}

	umw := controllers.UserMiddleware{
		SessionService: &sessionService,
	}

	r.Use(csrfMw)
	r.Use(customcsrf.RefreshCSRFToken)
	r.Use(umw.SetUser)

	// Setup our controllers
	usersC := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{ "message": "API running" }`))
	})

	r.Get("/csrf", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		cookie := &http.Cookie{
			Name:     "csrf",
			Value:    csrf.Token(r),
			Path:     "/",
			HttpOnly: false,
		}

		http.SetCookie(w, cookie)

		// TODO: Use a struct
		w.Write([]byte(`{ "success": "true" }`))
	})

	r.Post("/users", usersC.Create)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Post("/signout", usersC.ProcessSignOut)

	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Route not found", http.StatusNotFound)
	})

	return r
}

func main() {
	r := chi.NewRouter()

	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	if err != nil {
		panic(err)
	}

	defer db.Close()

	// migrations are the folder
	err = models.MigrateFS(db, migrations.FS, ".")

	if err != nil {
		panic(err)
	}

	corsMw := handlers.CORS(
		handlers.AllowCredentials(),
		handlers.AllowedOriginValidator(
			func(origin string) bool {
				return strings.HasPrefix(origin, "http://localhost")
			},
		),
		handlers.AllowedHeaders([]string{"X-Csrf-Token"}),
		handlers.ExposedHeaders([]string{"X-Csrf-Token"}),
	)

	r.Use(middleware.Logger)
	r.Use(corsMw)

	r.Mount("/api", getApiRouter(db))

	// we want all routes besides /api to go to the SPA, hence we use the NotFound handler
	r.NotFound(spa.SpaHandler)

	fmt.Println("Starting the server on :3000...")

	http.ListenAndServe(":3000", r)
}
