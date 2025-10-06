package handler

import (
	"encoding/json"
	"fmt"
	m "github.com/AVitaminOz-z-z/urlshortener.git/internal/model"

	//m "github.com/AVitaminOz-z-z/urlshortener.git/internal/model"
	"github.com/AVitaminOz-z-z/urlshortener.git/internal/storage"
	"io"
	"log/slog"
	"net/http"
)

// API request & response

/*type APIRequest struct {
	URL string `json:"url,omitempty"`
}

type APIBatchRequest struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	OriginalURL   string `json:"original_url,omitempty"`
}
type APIBatchRequestA []APIBatchRequest

type APIBatchResponse struct {
	CorrelationId string `json:"correlation_id,omitempty"`
	ShortURL      string `json:"short_url,omitempty"`
}
type APIBatchResponseA []APIBatchResponse

type APIResult struct {
	Result string `json:"result,omitempty"`
}*/

/*
func checkContentType(r *http.Request, hc *m.HandlerConfig) bool {
	return (regexp.MustCompile(hc.AvailableContentTypeRgx)).MatchString(r.Header.Get("Content-Type"))
}

func writeBadRequest(w http.ResponseWriter, hc *m.HandlerConfig, ext string) {
	writeDefaultHeader(w, hc)
	hc.GetLogger().Warn("Bad Request", slog.String("Message", ext))
	http.Error(w, fmt.Sprintf("Bad Request (%s)", ext), http.StatusBadRequest)
}

func write5xx(w http.ResponseWriter, hc *m.HandlerConfig, ext string) {
	writeDefaultHeader(w, hc)
	hc.GetLogger().Error("Server Error", slog.String("Message", ext))
	http.Error(w, fmt.Sprintf("Server Error (%s)", ext), http.StatusInternalServerError)
}

func CreateShortURL(hc *m.HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// parse body
		var url string
		if hc.ContentType != m.APIContentType {
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
				// writeBadRequest(w, hc, "request URL is empty")
				writeError(w, hc, "request URL is empty", http.StatusBadRequest)
				return
			}
			url = req.URL
		}
		// finding and send short url
		writeShortURL(w, hc, fmt.Sprintf("%s/%s", hc.GetStorage().BaseURL, hc.GetStorage().ReturnShortURL(r.Context(), url)))
	}
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
*/

func writeDefaultHeader(w http.ResponseWriter, hc *m.HandlerConfig) {
	w.Header().Set("Content-Type", hc.ContentType)
}

func writeError(w http.ResponseWriter, hc *m.HandlerConfig, ext string, code int) {
	writeDefaultHeader(w, hc)
	if code >= 500 {
		hc.GetLogger().Error(http.StatusText(code), slog.String("Message", ext))
	} else {
		hc.GetLogger().Warn(http.StatusText(code), slog.String("Message", ext))
	}
	http.Error(w, fmt.Sprintf("%s (%s)", http.StatusText(code), ext), code)
}

func writeOK(w http.ResponseWriter, hc *m.HandlerConfig) {
	writeDefaultHeader(w, hc)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(nil)
}

func writeShortURL(w http.ResponseWriter, hc *m.HandlerConfig, APIShortURL *storage.APIShorURL) {
	writeDefaultHeader(w, hc)
	w.WriteHeader(APIShortURL.HTTPCode)
	shortURL := fmt.Sprintf("%s/%s", APIShortURL.BaseURL, APIShortURL.ShortURL)
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

func writeBatchShortURLs(w http.ResponseWriter, hc *m.HandlerConfig, APIBatchShortURLs *storage.APIBatchShorURLs) {
	writeDefaultHeader(w, hc)
	w.WriteHeader(APIBatchShortURLs.HTTPCode)
	resp := APIBatchShortURLs.APIBatchResponseA
	enc := json.NewEncoder(w)
	_ = enc.Encode(resp)
}

func writeFullURL(w http.ResponseWriter, hc *m.HandlerConfig, fullURL string) {
	writeDefaultHeader(w, hc)
	w.Header().Set("Location", fullURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
	_, _ = w.Write(nil)
}

func CreateBatchShortURL(hc *m.HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// parse body
		//var url string

		req := storage.APIBatchRequestA{}
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

		//writeError(w, hc, fmt.Sprintf("request struct is not empty <len(req) = %+v>", req), http.StatusBadRequest)

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
				// writeBadRequest(w, hc, "request URL is empty")
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

func RedirectToFullURL(hc *m.HandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// check content type
		/*if !checkContentType(r, hc) {
			writeBadRequest(w, hc, "unsupported content type")
			return
		}*/
		//id := chi.URLParam(r, "id")
		id := r.URL.Path[1:]
		redirectURL := hc.GetStorage().ReturnFullURL(r.Context(), id)
		if redirectURL == "" {
			writeError(w, hc, fmt.Sprintf("unmanaged short url by {id} = %s", id), http.StatusBadRequest)
			return
		}
		writeFullURL(w, hc, redirectURL)
	}
}

func PingPgDB(hc *m.HandlerConfig) http.HandlerFunc {
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
