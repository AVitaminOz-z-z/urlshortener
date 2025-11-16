package handler

import (
	"bytes"
	"github.com/AVitaminOz-z-z/urlshortener.git/internal/logger"
	"github.com/AVitaminOz-z-z/urlshortener.git/internal/storage"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

/*
https://google.com
http://localhost:8080/cf41b6052ebb502e1ca4d72544333018475c73d740c8f74d027e1f1256e3c7be

https://yandex.ru
http://localhost:8080/7f6845db5740b940d7bcb2bd8f52f3b34501544416f499a5c27c55ae11ab0336

https://github.com
http://localhost:8080/82d9082ba2edaab19779147cd25ffe1e1483a498571e98b0d89eb44847baddf4
*/

func TestCreateShortURL(t *testing.T) {
	const (
		WantShortURL = "http://localhost:8080/7f6845db5740b940d7bcb2bd8f52f3b34501544416f499a5c27c55ae11ab0336"
		TestURL      = "https://yandex.ru"
	)

	us, _ := storage.NewURLFileStorage("./.test_storage")
	us.SetBaseURL("http://localhost:8080")
	var fn http.HandlerFunc

	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(TestURL)))
	r.Header.Set("Content-Type", "text/plaint; charset=utf-8")

	w := httptest.NewRecorder()
	l := logger.NewDiscardLogger()

	fn = CreateShortURL(NewDefHandlerConfig(us, l))
	fn(w, r)

	res := w.Result()

	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusConflict {
		t.Errorf("Code (%d && %d) was expected, but %d was received", http.StatusCreated, http.StatusConflict, res.StatusCode)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Body read error %v", err)
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if len(resBody) == 0 {
		t.Errorf("Body must be present in response")
	}

	if string(resBody) != WantShortURL {
		t.Errorf("Body %v was expected, but %v was received", WantShortURL, string(resBody))
	}
}

func TestAPICreateShortURL(t *testing.T) {
	const (
		WantShortURL = "{\"result\":\"http://localhost:8080/7f6845db5740b940d7bcb2bd8f52f3b34501544416f499a5c27c55ae11ab0336\"}\x0a"
		TestURL      = "{\"url\":\"https://yandex.ru\"}"
	)

	us, _ := storage.NewURLFileStorage("./.test_storage")
	us.SetBaseURL("http://localhost:8080")
	var fn http.HandlerFunc

	r := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader([]byte(TestURL)))
	r.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	l := logger.NewDiscardLogger()

	fn = CreateShortURL(NewAPIHandlerConfig(us, l))
	fn(w, r)

	res := w.Result()

	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusConflict {
		t.Errorf("Code (%d && %d) was expected, but %d was received", http.StatusCreated, http.StatusConflict, res.StatusCode)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Body read error %v", err)
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if len(resBody) == 0 {
		t.Errorf("Body must be present in response")
	}

	if string(resBody) != WantShortURL {
		t.Errorf("Body %v was expected, but %v was received", WantShortURL, string(resBody))
	}
}

func TestAPICreateBatchShortURLs(t *testing.T) {
	const (
		WantShortURL = "[{\"correlation_id\":\"cc0aa201-a5f9-40ea-8e74-b92273f337b5\",\"short_url\":\"http://localhost:8080/82d9082ba2edaab19779147cd25ffe1e1483a498571e98b0d89eb44847baddf4\"},{\"correlation_id\":\"996b268f-2763-4d1e-b73a-5d66e6a2e3bc\",\"short_url\":\"http://localhost:8080/cf41b6052ebb502e1ca4d72544333018475c73d740c8f74d027e1f1256e3c7be\"}]\x0a"
		TestURL      = "[{\"correlation_id\":\"cc0aa201-a5f9-40ea-8e74-b92273f337b5\",\"original_url\":\"https://github.com\"},{\"correlation_id\":\"996b268f-2763-4d1e-b73a-5d66e6a2e3bc\",\"original_url\":\"https://google.com\"}]"
	)

	us, _ := storage.NewURLFileStorage("./.test_storage")
	us.SetBaseURL("http://localhost:8080")
	var fn http.HandlerFunc

	r := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewReader([]byte(TestURL)))
	r.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	l := logger.NewDiscardLogger()

	fn = CreateBatchShortURL(NewAPIHandlerConfig(us, l))
	fn(w, r)

	res := w.Result()

	if res.StatusCode != http.StatusCreated {
		t.Errorf("Code %d was expected, but %d was received", http.StatusCreated, res.StatusCode)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Body read error %v", err)
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if len(resBody) == 0 {
		t.Errorf("Body must be present in response")
	}

	if string(resBody) != WantShortURL {
		t.Errorf("Body %v was expected, but %v was received", WantShortURL, string(resBody))
	}
}

func TestRedirectToFullURL(t *testing.T) {
	const (
		WantLocation = "https://github.com"
		TargetPath   = "/82d9082ba2edaab19779147cd25ffe1e1483a498571e98b0d89eb44847baddf4"
	)

	us, _ := storage.NewURLFileStorage("./.test_storage")
	us.SetBaseURL("http://localhost:8080")
	var fn http.HandlerFunc

	r := httptest.NewRequest(http.MethodGet, TargetPath, nil)
	r.Context().Value(middleware.URLFormatCtxKey)
	r.Header.Set("Content-Type", "text/plaint; charset=utf-8")

	w := httptest.NewRecorder()

	fn = RedirectToFullURL(NewDefHandlerConfig(us, nil))
	fn(w, r)

	res := w.Result()

	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusTemporaryRedirect {
		t.Errorf("Code %d was expected, but %d was received", http.StatusTemporaryRedirect, res.StatusCode)
	}

	if res.Header.Get("Location") != WantLocation {
		t.Errorf("Location %v was expected, but %v was received", res.Header.Get("Location"), WantLocation)
	}
}
