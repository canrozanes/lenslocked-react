package frontend

import (
	"embed"
	"errors"
	"html/template"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"path"
	"path/filepath"
)

var ErrDir = errors.New("path is dir")

//go:embed dist/*
var dist embed.FS

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

func tryReadHtml(efs embed.FS, prefix, requestedPath string, w http.ResponseWriter, r *http.Request) error {
	if requestedPath != "index.html" {
		return errors.New("path is not index.html")
	}

	indexHtmlBits, err := fs.ReadFile(efs, path.Join(prefix, requestedPath))

	if err != nil {
		return err
	}

	tpl, err := template.New("index.html").Parse(string(indexHtmlBits))
	if err != nil {
		return err
	}

	err = tpl.Execute(w, nil)

	return err
}

func SpaHandler(w http.ResponseWriter, r *http.Request) {
	// dist/assets
	err := tryRead(dist, "dist", r.URL.Path, w)
	if err == nil {
		return
	}

	// dist/index.html
	err = tryReadHtml(dist, "dist", "index.html", w, r)
	if err != nil {
		panic(err)
	}
}
