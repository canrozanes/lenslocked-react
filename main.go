package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/canrozanes/lenslocked/controllers"
	"github.com/canrozanes/lenslocked/models"
	"github.com/canrozanes/lenslocked/templates"
	"github.com/canrozanes/lenslocked/views"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Resource struct {
	Route string `json:"route"`
}

func getApiRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{ "message": "bar" }`))
	})

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

	r.Get("/", controllers.StaticHandler(views.Must(
		views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))))

	r.Get("/contact", controllers.StaticHandler(views.Must(
		views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))))

	r.Get("/faq", controllers.FAQ(
		views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))))

	// Setup a database connection
	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	userService := models.UserService{
		DB: db,
	}

	usersC := controllers.Users{
		UserService: &userService,
	}
	usersC.Templates.New = views.Must(views.ParseFS(templates.FS, "signup.gohtml", "tailwind.gohtml"))
	usersC.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "signin.gohtml", "tailwind.gohtml"))

	r.Get("/signup", usersC.New)
	r.Post("/users", usersC.Create)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)

	r.Mount("/api", getApiRouter())

	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", r)
}
