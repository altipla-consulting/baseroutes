package baseroutes

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/altipla-consulting/env"
	"libs.altipla.consulting/routing"
)

type RegisterOption func(r *routing.Router)

func WithFrontend(root string) RegisterOption {
	return func(r *routing.Router) {
		if env.IsLocal() {
			r.ServeFiles("/images", http.Dir(filepath.Join(root, "images")))
		}
	}
}

func Register(r *routing.Router, baseTemplate string, options ...RegisterOption) {
	go func() {
		if baseTemplate != "" {
			// Touch template to reload the page every time we change the Go implementation.
			_ = os.Chtimes(baseTemplate, time.Now(), time.Now())
		}
	}()

	r.Get("/robots.txt", fileHandler("robots.txt"))
	r.Get("/favicon.ico", fileHandler("favicon.ico"))
	r.Get("/apple-touch-icon.png", fileHandler("apple-touch-icon.png"))
	r.Get("/apple-touch-icon-precomposed.png", fileHandler("apple-touch-icon.png"))
	r.Get("/apple-touch-icon-120x120.png", fileHandler("apple-touch-icon.png"))
	r.Get("/apple-touch-icon-120x120-precomposed.png", fileHandler("apple-touch-icon.png"))
	r.Get("/apple-touch-icon-152x152.png", fileHandler("apple-touch-icon.png"))
	r.Get("/apple-touch-icon-152x152-precomposed.png", fileHandler("apple-touch-icon.png"))

	for _, opt := range options {
		opt(r)
	}
}

func fileHandler(path string) routing.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if env.IsLocal() {
			http.ServeFile(w, r, filepath.Join("..", "public", path))
		} else {
			http.ServeFile(w, r, filepath.Join("public", path))
		}
		return nil
	}
}
