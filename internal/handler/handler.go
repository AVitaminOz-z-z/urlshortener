package handler

import (
	"encoding/json"
	"fmt"
	"github.com/AVitaminOz-z-z/urlshortener.git/internal/common"
	m "github.com/AVitaminOz-z-z/urlshortener.git/internal/model"
	"github.com/AVitaminOz-z-z/urlshortener.git/internal/storage"
	"io"
	"log/slog"
	"net/http"
)

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
			ContentType:             common.DefContentType,
			AvailableContentTypeRgx: common.DefAvailableContentTypeRgx,
		},
	}
}

func NewAPIHandlerConfig(storage *storage.URLStorage, logger *slog.Logger) *HandlerConfig {
	hc := NewDefHandlerConfig(storage, logger)
	hc.HandlerDefaults = &HandlerDefaults{
		ContentType:             common.APIContentType,
		AvailableContentTypeRgx: common.APIAvailableContentTypeRgx,
	}
	return hc
}

func writeDefaultHeader(w http.ResponseWriter, hc *HandlerConfig) {
	w.Header().Set("Content-Type", hc.ContentType)
}

func writeLog(hc *HandlerConfig, code int, message string) {
	switch {
	case code >= 500:
		hc.GetLogger().Error(http.StatusText(code), slog.String("Message", message))
	case code >= 400 && code < 500:
		hc.GetLogger().Warn(http.StatusText(code), slog.String("Message", message))
	}
}

func writeError(w http.ResponseWriter, hc *HandlerConfig, message string, code int) {
	writeDefaultHeader(w, hc)
	writeLog(hc, code, message)
	http.Error(w, fmt.Sprintf("%s (%s)", http.StatusText(code), message), code)
}

func writeOK(w http.ResponseWriter, hc *HandlerConfig) {
	writeDefaultHeader(w, hc)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(nil)
}

func writeShortURL(w http.ResponseWriter, hc *HandlerConfig, APIShortURL *m.APIShorURL) {
	writeDefaultHeader(w, hc)
	writeLog(hc, APIShortURL.HTTPCode, http.StatusText(APIShortURL.HTTPCode))
	w.WriteHeader(APIShortURL.HTTPCode)
	// make & write result
	if hc.ContentType != common.APIContentType {
		_, _ = w.Write([]byte(APIShortURL.ShortURL))
	} else {
		resp := m.APIResult{
			Result: APIShortURL.ShortURL,
		}
		enc := json.NewEncoder(w)
		_ = enc.Encode(resp)
	}
}

func writeBatchShortURLs(w http.ResponseWriter, hc *HandlerConfig, APIBatchShortURLs *m.APIBatchShorURLs) {
	writeDefaultHeader(w, hc)
	w.WriteHeader(APIBatchShortURLs.HTTPCode)
	resp := APIBatchShortURLs.APIBatchResponseA
	enc := json.NewEncoder(w)
	_ = enc.Encode(resp)
}

func writeFullURL(w http.ResponseWriter, hc *HandlerConfig, fullURL string) {
	writeDefaultHeader(w, hc)
	w.Header().Set("Location", fullURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
	_, _ = w.Write(nil)
}

func CreateBatchShortURL(hc *HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// parse body
		req := m.APIBatchRequestA{}
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&req); err != nil {
			writeError(w, hc, fmt.Sprintf("can't unmarshal request to struct: %v", err), http.StatusBadRequest)
			return
		}
		// struct is empty
		if len(req) == 0 {
			writeError(w, hc, fmt.Sprintf("request struct is empty <len(req) = %d>", len(req)), http.StatusBadRequest)
			return
		}
		// finding and send short url
		us := hc.GetStorage()
		batchSU, err := us.ReturnBatchShortURL(r.Context(), req)
		// can't get batch short urls
		if err != nil {
			writeError(w, hc, err.Error(), http.StatusInternalServerError)
			return
		}
		writeBatchShortURLs(w, hc, batchSU)
	}
}

func CreateShortURL(hc *HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// parse body
		var url string
		if hc.ContentType != common.APIContentType {
			// body is empty or read error
			body, err := io.ReadAll(r.Body)
			if len(body) == 0 || err != nil {
				writeError(w, hc, fmt.Sprintf("body is empty or read error <err = %v>", err), http.StatusBadRequest)
				return
			}
			url = string(body)
		} else {
			req := m.APIRequest{}
			dec := json.NewDecoder(r.Body)
			if err := dec.Decode(&req); err != nil {
				writeError(w, hc, fmt.Sprintf("can't unmarshal request to struct: %v", err), http.StatusBadRequest)
				return
			}
			// req URL is empty
			if req.URL == "" {
				writeError(w, hc, "request URL is empty", http.StatusBadRequest)
				return
			}
			url = req.URL
		}
		// finding and send short url
		us := hc.GetStorage()
		apiSU, err := us.ReturnShortURL(r.Context(), url)
		// can't get short url
		if err != nil {
			writeError(w, hc, err.Error(), http.StatusInternalServerError)
			return
		}
		writeShortURL(w, hc, apiSU)
	}
}

func RedirectToFullURL(hc *HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//id := chi.URLParam(r, "id")
		id := r.URL.Path[1:]
		redirectURL, err := hc.GetStorage().ReturnFullURL(r.Context(), id)
		if err != nil {
			writeError(w, hc, err.Error(), http.StatusInternalServerError)
			return
		}
		if redirectURL == "" {
			writeError(w, hc, fmt.Sprintf("unmanaged short url by {id} = %s", id), http.StatusBadRequest)
			return
		}
		writeFullURL(w, hc, redirectURL)
	}
}

func PingPgDB(hc *HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if hc.GetPgDB() == nil {
			writeError(w, hc, "PgDB-engine not defined", http.StatusInternalServerError)
			return
		}
		if err := hc.GetPgDB().Ping(); err != nil {
			writeError(w, hc, err.Error(), http.StatusInternalServerError)
			return
		}
		writeOK(w, hc)
	}
}
