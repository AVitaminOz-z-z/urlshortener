package handler

import (
	"encoding/json"
	"fmt"
	m "github.com/AVitaminOz-z-z/urlshortener.git/internal/model"
	"io"
	"log/slog"
	"net/http"
)

/*func checkContentType(r *http.Request, hc *m.HandlerConfig) bool {
	return (regexp.MustCompile(hc.AvailableContentTypeRgx)).MatchString(r.Header.Get("Content-Type"))
}*/

func writeDefaultHeader(w http.ResponseWriter, hc *m.HandlerConfig) {
	w.Header().Set("Content-Type", hc.ContentType)
}

func writeBadRequest(w http.ResponseWriter, hc *m.HandlerConfig, ext string) {
	writeDefaultHeader(w, hc)
	hc.GetLogger().Warn("Bad Request", slog.String("Message", ext))
	http.Error(w, fmt.Sprintf("Bad Request (%s)", ext), http.StatusBadRequest)
}

func writeShortURL(w http.ResponseWriter, hc *m.HandlerConfig, shortURL string) {
	writeDefaultHeader(w, hc)
	w.WriteHeader(http.StatusCreated)
	// make & write result
	if hc.ContentType != m.APIContentType {
		_, _ = w.Write([]byte(shortURL))
	} else {
		resp := m.APIResult{
			Result: shortURL,
		}
		enc := json.NewEncoder(w)
		_ = enc.Encode(resp)
	}
}

func writeFullURL(w http.ResponseWriter, hc *m.HandlerConfig, fullURL string) {
	writeDefaultHeader(w, hc)
	w.Header().Set("Location", fullURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
	_, _ = w.Write(nil)
}

func CreateShortURL(hc *m.HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// check content type
		/*if !checkContentType(r, hc) {
			writeBadRequest(w, hc, "unsupported content type")
			return
		}*/
		// parse body
		var url string
		if hc.ContentType != m.APIContentType {
			// body is empty or read error
			body, err := io.ReadAll(r.Body)
			if len(body) == 0 || err != nil {
				writeBadRequest(w, hc, fmt.Sprintf("body is empty or read error <err = %v>", err))
				return
			}
			url = string(body)
		} else {
			req := m.APIRequest{}
			dec := json.NewDecoder(r.Body)
			if err := dec.Decode(&req); err != nil {
				writeBadRequest(w, hc, fmt.Sprintf("can't unmarshal request to struct: %v", err))
				return
			}
			// req URL is empty
			if req.URL == "" {
				writeBadRequest(w, hc, "request URL is empty")
				return
			}
			url = req.URL
		}
		// finding and send short url
		writeShortURL(w, hc, fmt.Sprintf("%s/%s", hc.GetStorage().BaseURL, hc.GetStorage().ReturnShortURL(url)))
	}
}

func RedirectToFullURL(hc *m.HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// check content type
		/*if !checkContentType(r, hc) {
			writeBadRequest(w, hc, "unsupported content type")
			return
		}*/
		//id := chi.URLParam(r, "id")
		id := r.URL.Path[1:]
		redirectURL := hc.GetStorage().ReturnFullURL(id)
		if redirectURL == "" {
			writeBadRequest(w, hc, fmt.Sprintf("unmanaged short url by {id} = %s", id))
			return
		}
		writeFullURL(w, hc, redirectURL)
	}
}
