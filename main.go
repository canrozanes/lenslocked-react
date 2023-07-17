package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/canrozanes/lenslocked/spa"
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

	// r.Get("/", controllers.StaticHandler(views.Must(
	// 	views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))))

	// r.Get("/contact", controllers.StaticHandler(views.Must(
	// 	views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))))

	// r.Get("/faq", controllers.FAQ(
	// 	views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))))

	r.Mount("/api", getApiRouter())

	// we want all routes besides /api to go to the SPA, hence we use the NotFound handler
	r.NotFound(spa.SpaHandler)

	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", r)
}
