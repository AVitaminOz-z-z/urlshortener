package router

import (
	"github.com/AVitaminOz-z-z/urlshortener.git/internal/handler"
	"github.com/AVitaminOz-z-z/urlshortener.git/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewURLRouter(urlStorage *storage.URLStorage) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Get("/{id}", handler.ManageGET(urlStorage))
	r.Post("/", handler.ManagePOST(urlStorage))
	return r
}
