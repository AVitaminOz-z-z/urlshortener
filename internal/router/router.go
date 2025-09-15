package router

import (
	"github.com/AVitaminOz-z-z/urlshortener.git/internal/handler"
	mwLogger "github.com/AVitaminOz-z-z/urlshortener.git/internal/logger"
	"github.com/AVitaminOz-z-z/urlshortener.git/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
)

func NewURLRouter(storage *storage.URLStorage, logger *slog.Logger) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(mwLogger.NewMiddlewareLogger(logger))
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Get("/{id}", handler.RedirectToFullURL(storage))
	r.Post("/", handler.CreateShortURL(storage))
	r.Post("/api/shorten", handler.CreateAPIShortURL(storage))
	return r
}
