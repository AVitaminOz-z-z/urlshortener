package router

import (
	mwCompress "github.com/AVitaminOz-z-z/urlshortener.git/internal/compress"
	mwCookie "github.com/AVitaminOz-z-z/urlshortener.git/internal/cookieman"
	"github.com/AVitaminOz-z-z/urlshortener.git/internal/handler"
	mwLogger "github.com/AVitaminOz-z-z/urlshortener.git/internal/logger"
	"github.com/AVitaminOz-z-z/urlshortener.git/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
)

func NewURLRouter(storage *storage.URLStorage, logger *slog.Logger, key []byte) chi.Router {
	// Default handler config
	hcDef := handler.NewDefHandlerConfig(storage, logger)
	// API handler config
	hcAPI := handler.NewAPIHandlerConfig(storage, logger)
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(mwCookie.NewMiddlewareCookieMan(key))
	r.Use(mwCompress.NewGzipMiddleware())
	r.Use(mwLogger.NewMiddlewareLogger(logger))
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Get("/{id}", handler.RedirectToFullURL(hcDef))
	r.Get("/ping", handler.PingPgDB(hcDef))
	r.Post("/", handler.CreateShortURL(hcDef))
	r.Post("/api/shorten", handler.CreateShortURL(hcAPI))
	r.Post("/api/shorten/batch", handler.CreateBatchShortURL(hcAPI))
	r.Get("/api/user/urls", handler.ReturnUserURLs(hcAPI))
	return r
}
