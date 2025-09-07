package handler

import (
	"fmt"
	"github.com/AVitaminOz-z-z/urlshortener.git/internal/storage"
	"io"
	"net/http"
	"regexp"
)

const (
	AvailableContentTypeRgx = `^text/plain(|.+)$`
	DefContentType          = "text/plain; charset=utf-8"
)

func checkContentType(r *http.Request) bool {
	return (regexp.MustCompile(AvailableContentTypeRgx)).MatchString(r.Header.Get("Content-Type"))
}

func writeDefaultHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", DefContentType)
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

func writeFullURL(fullURL string, w http.ResponseWriter) {
	writeDefaultHeader(w)
	w.Header().Set("Location", fullURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
	_, _ = w.Write(nil)
}

func ManagePOST(urlStorage *storage.URLStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// unsupported content type
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
		//WriteShortURL(r.Host+"/"+urlStorage.ReturnShortURL(string(body)), w)
		writeShortURL(urlStorage.BaseURL+"/"+urlStorage.ReturnShortURL(string(body)), w)
	}
}

func ManageGET(urlStorage *storage.URLStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// unsupported content type
		/*if !checkContentType(r) {
			writeBadRequest(w, "unsupported content type")
			return
		}
		id := chi.URLParam(r, "id")*/
		id := r.URL.Path[1:]
		fullURL := urlStorage.ReturnFullURL(id)
		if fullURL == "" {
			writeBadRequest(w, fmt.Sprintf("unmanaged short url by {id} = %s", id))
			return
		}
		writeFullURL(fullURL, w)
	}
}
