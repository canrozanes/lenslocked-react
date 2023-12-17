package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
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
	"github.com/joho/godotenv"
)

type config struct {
	PSQL models.PostgresConfig
	SMTP models.SMTPConfig
	CSRF struct {
		Key    string
		Secure bool
	}
	Server struct {
		Address string
	}

	CORS struct {
		Host string
	}
}

func loadEnvConfig() (config, error) {
	var cfg config
	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}

	// TODO: PSQL
	cfg.PSQL = models.DefaultPostgresConfig()

	// TODO: SMTP
	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	cfg.SMTP.Port, err = strconv.Atoi(portStr)
	if err != nil {
		panic(err)
	}
	cfg.SMTP.Username = os.Getenv("SMTP_USERNAME")
	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")

	// TODO: CSRF
	cfg.CSRF.Key = "gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX"
	cfg.CSRF.Secure = false

	// TODO: Server
	cfg.Server.Address = ":3000"

	return cfg, nil
}

func getApiRouter(db *sql.DB, cfg config) chi.Router {
	r := chi.NewRouter()

	csrfErrorHandler := customcsrf.NewErrorHandler()
	csrfMw := csrf.Protect(
		[]byte(cfg.CSRF.Key),
		csrf.Secure(cfg.CSRF.Secure),
		csrf.ErrorHandler(csrfErrorHandler),
		csrf.SameSite(csrf.SameSiteStrictMode),
	)

	userService := &models.UserService{
		DB: db,
	}

	sessionService := &models.SessionService{
		DB: db,
	}

	pwResetService := &models.PasswordResetService{
		DB: db,
	}

	emailService := models.NewEmailService(cfg.SMTP)

	umw := controllers.UserMiddleware{
		SessionService: sessionService,
	}

	r.Use(csrfMw)
	r.Use(customcsrf.RefreshCSRFToken)
	r.Use(umw.SetUser)

	// Setup our controllers
	usersC := controllers.Users{
		UserService:          userService,
		SessionService:       sessionService,
		PasswordResetService: pwResetService,
		EmailService:         emailService,
	}

	galleryService := &models.GalleryService{
		DB: db,
	}

	galleriesC := controllers.Galleries{
		GalleryService: galleryService,
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
	r.Post("/forgot-pw", usersC.ProcessForgotPassword)
	r.Post("/reset-pw", usersC.ProcessResetPassword)

	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
	})

	r.Route("/galleries", func(r chi.Router) {
		r.Get("/{id}", galleriesC.GetGallery) // used to be called Show
		r.Group(func(r chi.Router) {
			r.Use(umw.RequireUser)
			r.Get("/", galleriesC.GetAllGalleries)
			r.Post("/", galleriesC.Create)
			r.Post("/{id}", galleriesC.Update)
			r.Post("/{id}/delete", galleriesC.Delete)
		})
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Route not found", http.StatusNotFound)
	})

	return r
}

func main() {
	r := chi.NewRouter()

	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}
	db, err := models.Open(cfg.PSQL)

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
				return strings.HasPrefix(origin, cfg.CORS.Host)
			},
		),
		handlers.AllowedHeaders([]string{"X-Csrf-Token"}),
		handlers.ExposedHeaders([]string{"X-Csrf-Token"}),
	)

	r.Use(middleware.Logger)
	r.Use(corsMw)

	r.Mount("/api", getApiRouter(db, cfg))

	// we want all routes besides /api to go to the SPA, hence we use the NotFound handler
	r.NotFound(spa.SpaHandler)

	fmt.Printf("Starting the server on %s...\n", cfg.Server.Address)
	err = http.ListenAndServe(":3000", r)

	if err != nil {
		panic(err)
	}
}
