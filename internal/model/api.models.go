package model

import (
	"github.com/AVitaminOz-z-z/urlshortener.git/internal/storage"
	"log/slog"
)

const (
	DefAvailableContentTypeRgx = `^text/plain(|.+)$`
	APIAvailableContentTypeRgx = `^application/json(|.+)$`
	DefContentType             = "text/plain; charset=utf-8"
	APIContentType             = "application/json"
)

// API request & response

type APIRequest struct {
	URL string `json:"url,omitempty"`
}

type APIResult struct {
	Result string `json:"result,omitempty"`
}

// handler config

type HandlerDefaults struct {
	ContentType             string
	AvailableContentTypeRgx string
}

type HandlerConfig struct {
	storage *storage.URLStorage
	logger  *slog.Logger
	*HandlerDefaults
}

func (hc *HandlerConfig) GetStorage() *storage.URLStorage {
	return hc.storage
}

func (hc *HandlerConfig) GetPgDB() *storage.PgDB {
	return hc.storage.PgDB
}

func (hc *HandlerConfig) GetLogger() *slog.Logger {
	return hc.logger
}

func (hc *HandlerConfig) ResetLogger(l *slog.Logger) {
	hc.logger = l
}

func NewDefHandlerConfig(storage *storage.URLStorage, logger *slog.Logger) *HandlerConfig {
	return &HandlerConfig{
		storage: storage,
		logger:  logger,
		HandlerDefaults: &HandlerDefaults{
			ContentType:             DefContentType,
			AvailableContentTypeRgx: DefAvailableContentTypeRgx,
		},
	}
}

func NewAPIHandlerConfig(storage *storage.URLStorage, logger *slog.Logger) *HandlerConfig {
	hc := NewDefHandlerConfig(storage, logger)
	hc.HandlerDefaults = &HandlerDefaults{
		ContentType:             APIContentType,
		AvailableContentTypeRgx: APIAvailableContentTypeRgx,
	}
	return hc
}
