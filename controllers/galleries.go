package controllers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/canrozanes/lenslocked/context"
	"github.com/canrozanes/lenslocked/errors"
	"github.com/canrozanes/lenslocked/models"
	"github.com/canrozanes/lenslocked/write"
	"github.com/go-chi/chi/v5"
)

type Galleries struct {
	GalleryService *models.GalleryService
}

type galleryResponse struct {
	Gallery *models.Gallery `json:"gallery"`
}

func (g Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var data struct {
		UserID int
		Title  string
	}

	err := json.NewDecoder(r.Body).Decode(&data)
	data.UserID = context.User(r.Context()).ID

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		err = errors.Public(err, "Request is missing title.")
		errors.WriteErrorResponse(w, err)
		return
	}

	gallery, err := g.GalleryService.Create(data.Title, data.UserID)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		err = errors.Public(err, "Could not create gallery.")
		errors.WriteErrorResponse(w, err)
		return
	}

	json.NewEncoder(w).Encode(galleryResponse{Gallery: gallery})
}

func (g Galleries) Update(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}

	var data struct {
		UserID int
		Title  string
	}

	err = json.NewDecoder(r.Body).Decode(&data)
	data.UserID = context.User(r.Context()).ID

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err.Error())
		err = errors.Public(err, "Bad request")
		errors.WriteErrorResponse(w, err)
		return
	}

	gallery.Title = data.Title
	err = g.GalleryService.Update(gallery)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = errors.Public(err, "Something went wrong")
		errors.WriteErrorResponse(w, err)
		return
	}

	json.NewEncoder(w).Encode(galleryResponse{Gallery: gallery})
}

// Used to be called Index
func (g Galleries) GetAllGalleries(w http.ResponseWriter, r *http.Request) {
	type Gallery struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
	}
	var data struct {
		Galleries []Gallery
	}

	user := context.User(r.Context())

	galleries, err := g.GalleryService.ByUserID(user.ID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errors.WriteErrorResponse(w, fmt.Errorf("Something went wrong"))

		return
	}

	data.Galleries = []Gallery{}
	for _, gallery := range galleries {
		data.Galleries = append(data.Galleries, Gallery{
			ID:    gallery.ID,
			Title: gallery.Title,
		})
	}

	type galleriesResponse struct {
		Galleries []Gallery `json:"galleries"`
	}

	json.NewEncoder(w).Encode(galleriesResponse{Galleries: data.Galleries})
}

// Used to be called Show
func (g Galleries) GetGallery(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}

	// TODO Delete this when we implemented images
	type Gallery struct {
		ID     int      `json:"id"`
		Title  string   `json:"title"`
		Images []string `json:"images"`
		UserId int      `json:"user_id"`
	}

	type galleryResponse struct {
		Gallery `json:"gallery"`
	}

	data := Gallery{
		ID:     gallery.ID,
		Title:  gallery.Title,
		UserId: gallery.UserID,
	}

	for i := 0; i < 20; i++ {
		w, h := rand.Intn(500)+200, rand.Intn(500)+200
		catImageURL := fmt.Sprintf("https://placehold.co/%dx%d", w, h)
		data.Images = append(data.Images, catImageURL)
	}

	json.NewEncoder(w).Encode(galleryResponse{Gallery: data})
}

func (g Galleries) Delete(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}
	err = g.GalleryService.Delete(gallery.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errors.WriteErrorResponse(w, fmt.Errorf("something went wrong"))

		return
	}

	write.Success(w)

}

// 1. Define the galleryOpt type, which is really a function
// that takes in response writer, request, and gallery and
// will return an error if something isn't okay.
type galleryOpt func(http.ResponseWriter, *http.Request, *models.Gallery) error

func (g Galleries) galleryByID(w http.ResponseWriter, r *http.Request, opts ...galleryOpt) (*models.Gallery, error) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errors.WriteErrorResponse(w, fmt.Errorf("invalid ID"))

		return nil, err
	}
	gallery, err := g.GalleryService.ByID(id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			errors.WriteErrorResponse(w, fmt.Errorf("gallery not found"))

			return nil, err
		}
		w.WriteHeader(http.StatusInternalServerError)
		errors.WriteErrorResponse(w, fmt.Errorf("something went wrong"))

		return nil, err
	}
	// 3. Add a for loop to iterate over all of our functional
	// options, calling each and returning if there is an error.
	for _, opt := range opts {
		err = opt(w, r, gallery)
		if err != nil {
			return nil, err
		}
	}
	return gallery, nil
}

func userMustOwnGallery(w http.ResponseWriter, r *http.Request, gallery *models.Gallery) error {
	user := context.User(r.Context())
	if user.ID != gallery.UserID {
		w.WriteHeader(http.StatusForbidden)
		errors.WriteErrorResponse(w, fmt.Errorf("user does not have access to this gallery"))
		return fmt.Errorf("user does not have access to this gallery")
	}
	return nil
}
