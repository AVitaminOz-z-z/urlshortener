package handler

import (
	"encoding/json"
	"fmt"
	"github.com/AVitaminOz-z-z/urlshortener.git/internal/storage"
	"io"
	"net/http"
	"regexp"
)

const (
	AvailableContentTypeRgx    = `^text/plain(|.+)$`
	AvailableAPIContentTypeRgx = `^application/json(|.+)$`
	DefContentType             = "text/plain; charset=utf-8"
	DefAPIContentType          = "application/json"
)

type (
	APIRequest struct {
		URL string `json:"url,omitempty"`
	}

	APIResult struct {
		Result string `json:"result,omitempty"`
	}
)

func checkContentType(r *http.Request) bool {
	return (regexp.MustCompile(AvailableContentTypeRgx)).MatchString(r.Header.Get("Content-Type"))
}

func checkAPIContentType(r *http.Request) bool {
	return (regexp.MustCompile(AvailableAPIContentTypeRgx)).MatchString(r.Header.Get("Content-Type"))
}

func writeDefaultHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", DefContentType)
}

func writeAPIDefaultHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", DefAPIContentType)
}

func writeBadRequest(w http.ResponseWriter, ext string) {
	writeDefaultHeader(w)
	http.Error(w, fmt.Sprintf("Bad Request (%s)", ext), http.StatusBadRequest)
}

func writeShortURL(shortURL string, w http.ResponseWriter) {
	writeDefaultHeader(w)
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(shortURL))
}

func writeAPIShortURL(shortURL string, w http.ResponseWriter) {
	writeAPIDefaultHeader(w)
	w.WriteHeader(http.StatusCreated)
	// make & write result
	resp := APIResult{
		Result: shortURL,
	}
	enc := json.NewEncoder(w)
	_ = enc.Encode(resp)
}

func writeFullURL(fullURL string, w http.ResponseWriter) {
	writeDefaultHeader(w)
	w.Header().Set("Location", fullURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
	_, _ = w.Write(nil)
}

func CreateShortURL(storage *storage.URLStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// check content type
		if !checkContentType(r) {
			writeBadRequest(w, "unsupported content type")
			return
		}
		// body is empty or read error
		body, err := io.ReadAll(r.Body)
		if len(body) == 0 || err != nil {
			writeBadRequest(w, "body is empty or read error")
			return
		}
		// finding and send short url
		writeShortURL(fmt.Sprintf("%s/%s", storage.BaseURL, storage.ReturnShortURL(string(body))), w)
	}
}

func CreateAPIShortURL(storage *storage.URLStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// check content type
		if !checkAPIContentType(r) {
			writeBadRequest(w, "unsupported content type")
			return
		}
		// decode error
		var req APIRequest
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&req); err != nil {
			writeBadRequest(w, fmt.Sprintf("can't unmarshal request to struct: %v", err))
			return
		}
		// req URL is empty
		if req.URL == "" {
			writeBadRequest(w, "request URL is empty")
			return
		}

		// finding and send short url
		writeAPIShortURL(fmt.Sprintf("%s/%s", storage.BaseURL, storage.ReturnShortURL(req.URL)), w)
	}
}

func RedirectToFullURL(storage *storage.URLStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// unsupported content type
		/*if !checkContentType(r) {
			writeBadRequest(w, "unsupported content type")
			return
		}
		id := chi.URLParam(r, "id")*/
		id := r.URL.Path[1:]
		redirectURL := storage.ReturnFullURL(id)
		if redirectURL == "" {
			writeBadRequest(w, fmt.Sprintf("unmanaged short url by {id} = %s", id))
			return
		}
		writeFullURL(redirectURL, w)
	}
}
