package spa

import (
	"embed"
	"errors"
	"io"
	"mime"
	"net/http"
	"path"
	"path/filepath"
)

var ErrDir = errors.New("path is dir")

//go:embed dist/*
var assets embed.FS

func tryRead(fs embed.FS, prefix, requestedPath string, w http.ResponseWriter) error {
	f, err := fs.Open(path.Join(prefix, requestedPath))
	if err != nil {
		return err
	}
	defer f.Close()

	stat, _ := f.Stat()
	if stat.IsDir() {
		return ErrDir
	}

	contentType := mime.TypeByExtension(filepath.Ext(requestedPath))
	w.Header().Set("Content-Type", contentType)
	_, err = io.Copy(w, f)
	return err
}

func SpaHandler(w http.ResponseWriter, r *http.Request) {
	err := tryRead(assets, "dist", r.URL.Path, w)
	if err == nil {
		return
	}

	err = tryRead(assets, "dist", "index.html", w)
	if err != nil {
		panic(err)
	}
}
