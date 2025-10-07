package handler

import (
	"bytes"
	"github.com/AVitaminOz-z-z/urlshortener.git/internal/storage"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateShortURL(t *testing.T) {
	const (
		WantShortURL = "http://localhost:8080/746b8a4841e54de93f485104e6f0020b5c4c38f0e03845ff5578ef835ea8ddcc"
		TestURL      = "https://yandex.ru"
	)

	us, _ := storage.NewURLFileStorage("./.test_storage")
	us.SetBaseURL("http://localhost:8080")
	var fn http.HandlerFunc

	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(TestURL)))
	r.Header.Set("Content-Type", "text/plaint; charset=utf-8")

	w := httptest.NewRecorder()

	fn = CreateShortURL(NewDefHandlerConfig(us, nil))
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
		WantShortURL = "{\"result\":\"http://localhost:8080/746b8a4841e54de93f485104e6f0020b5c4c38f0e03845ff5578ef835ea8ddcc\"}\x0a"
		TestURL      = "{\"url\":\"https://yandex.ru\"}"
	)

	us, _ := storage.NewURLFileStorage("./.test_storage")
	us.SetBaseURL("http://localhost:8080")
	var fn http.HandlerFunc

	r := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader([]byte(TestURL)))
	r.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	fn = CreateShortURL(NewAPIHandlerConfig(us, nil))
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
		WantShortURL = "[{\"correlation_id\":\"cc0aa201-a5f9-40ea-8e74-b92273f337b5\",\"short_url\":\"http://localhost:8080/249110f758ff4d188b124b94787b758a93539865f89609bfdf681fea588b6fba\"},{\"correlation_id\":\"996b268f-2763-4d1e-b73a-5d66e6a2e3bc\",\"short_url\":\"http://localhost:8080/aebf6688de772690fe9e3f3839a2063833da43f31dbe3d6b677a16dd943a67b0\"}]\x0a"
		TestURL      = "[{\"correlation_id\":\"cc0aa201-a5f9-40ea-8e74-b92273f337b5\",\"original_url\":\"https://github.com\"},{\"correlation_id\":\"996b268f-2763-4d1e-b73a-5d66e6a2e3bc\",\"original_url\":\"https://google.com\"}]"
	)

	us, _ := storage.NewURLFileStorage("./.test_storage")
	us.SetBaseURL("http://localhost:8080")
	var fn http.HandlerFunc

	r := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewReader([]byte(TestURL)))
	r.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	fn = CreateBatchShortURL(NewAPIHandlerConfig(us, nil))
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
		TargetPath   = "/249110f758ff4d188b124b94787b758a93539865f89609bfdf681fea588b6fba"
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
