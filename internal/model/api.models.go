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
	pgDB    *storage.PgDB
	*HandlerDefaults
}

func (hc *HandlerConfig) GetStorage() *storage.URLStorage {
	return hc.storage
}

func (hc *HandlerConfig) GetPgDB() *storage.PgDB {
	return hc.pgDB
}

func (hc *HandlerConfig) GetLogger() *slog.Logger {
	return hc.logger
}

func NewDefHandlerConfig(storage *storage.URLStorage, logger *slog.Logger, pgDB *storage.PgDB) *HandlerConfig {
	return &HandlerConfig{
		storage: storage,
		logger:  logger,
		pgDB:    pgDB,
		HandlerDefaults: &HandlerDefaults{
			ContentType:             DefContentType,
			AvailableContentTypeRgx: DefAvailableContentTypeRgx,
		},
	}
}

func NewAPIHandlerConfig(storage *storage.URLStorage, logger *slog.Logger, pgDB *storage.PgDB) *HandlerConfig {
	hc := NewDefHandlerConfig(storage, logger, pgDB)
	hc.HandlerDefaults = &HandlerDefaults{
		ContentType:             APIContentType,
		AvailableContentTypeRgx: APIAvailableContentTypeRgx,
	}
	return hc
}
