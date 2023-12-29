package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
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

	type Image struct {
		GalleryID       int    `json:"gallery_id"`
		Filename        string `json:"filename"`
		FilenameEscaped string `json:"filename_escaped"`
	}

	// TODO Delete this when we implemented images
	type Gallery struct {
		ID     int     `json:"id"`
		Title  string  `json:"title"`
		Images []Image `json:"images"`
		UserId int     `json:"user_id"`
	}

	type galleryResponse struct {
		Gallery `json:"gallery"`
	}

	data := Gallery{
		ID:     gallery.ID,
		Title:  gallery.Title,
		UserId: gallery.UserID,
	}

	images, err := g.GalleryService.Images(gallery.ID)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		errors.WriteErrorResponse(w, fmt.Errorf("Something went wrong"))
		return
	}

	for _, image := range images {
		data.Images = append(data.Images, Image{
			GalleryID:       image.GalleryID,
			Filename:        image.Filename,
			FilenameEscaped: url.PathEscape(image.Filename),
		})
	}

	json.NewEncoder(w).Encode(galleryResponse{Gallery: data})
}

func (g Galleries) Image(w http.ResponseWriter, r *http.Request) {
	filename := g.filename(w, r)
	galleryID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusNotFound)
		return
	}

	image, err := g.GalleryService.Image(galleryID, filename)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Image not found", http.StatusNotFound)
			return
		}
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, image.Path)
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

func (g Galleries) DeleteImage(w http.ResponseWriter, r *http.Request) {
	filename := g.filename(w, r)
	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}
	err = g.GalleryService.DeleteImage(gallery.ID, filename)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errors.WriteErrorResponse(w, fmt.Errorf("something went wrong"))
		return
	}
	write.Success(w)
}

func (g Galleries) UploadImage(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}

	err = r.ParseMultipartForm(5 << 20) // 5mb
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errors.WriteErrorResponse(w, fmt.Errorf("something went wrong"))
		return
	}
	fileHeaders := r.MultipartForm.File["files"]
	for _, fileHeader := range fileHeaders {
		file, err := fileHeader.Open()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			errors.WriteErrorResponse(w, fmt.Errorf("something went wrong"))
			return
		}
		defer file.Close()

		err = g.GalleryService.CreateImage(gallery.ID, fileHeader.Filename, file)
		if err != nil {
			// Add some extra error handling.
			var fileErr models.FileError
			if errors.As(err, &fileErr) {
				msg := fmt.Sprintf("%v has an invalid content type or extension. Only png, gif, and jpg files can be uploaded.", fileHeader.Filename)
				w.WriteHeader(http.StatusInternalServerError)
				errors.WriteErrorResponse(w, fmt.Errorf(msg))
				return
			}
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			errors.WriteErrorResponse(w, fmt.Errorf("something went wrong"))
			return
		}
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

func (g Galleries) filename(w http.ResponseWriter, r *http.Request) string {
	filename := chi.URLParam(r, "filename")
	filename = filepath.Base(filename)
	return filename
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
