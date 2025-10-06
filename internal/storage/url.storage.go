package storage

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/AVitaminOz-z-z/urlshortener.git/internal/common"
	"net/http"
	"os"
)

/*
func (us *URLStorage) ReturnShortURL(ctx context.Context, url string) string {
	if us.UseDBEngine {
		return us.returnPgDBShortURL(ctx, url)
	}
	return us.updateFileStorage(url)
}

func (us *URLStorage) updateFileStorage(url string) string {
	if val, ok := us.Storage.POSTStorage[url]; !ok {
		rs := us.randomString()
		sha256 := us.sha256Sum(url + rs)
		us.Storage.POSTStorage[url] = []string{sha256, rs}
		us.Storage.GETStorage[sha256] = []string{url, rs}
		return sha256
	} else {
		return val[0]
	}
}
*/

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

type APIShorURL struct {
	HTTPCode int    `json:"http_code,omitempty"`
	ShortURL string `json:"short_url,omitempty"`
	BaseURL  string `json:"base_url,omitempty"`
}

type APIBatchShorURLs struct {
	HTTPCode int `json:"http_code,omitempty"`
	APIBatchResponseA
}

type Storage struct {
	POSTStorage map[string][]string `json:"POSTStorage,omitempty"`
	GETStorage  map[string][]string `json:"GETStorage,omitempty"`
}

type URLStorage struct {
	StorageName string `json:"StorageName,omitempty"`
	BaseURL     string `json:"BaseAddr,omitempty"`
	Storage     `json:"Storage,omitempty"`
	*PgDB
	UseDBEngine bool
}

func (us *URLStorage) ReturnShortURL(ctx context.Context, url string) (*APIShorURL, error) {
	if us.UseDBEngine {
		b := us.returnPgDBShortURL(ctx, url)
		apiSU, err := us.byteAToAPIShorURL(b)
		if err != nil {
			return nil, err
		}
		apiSU.BaseURL = us.BaseURL
		return apiSU, nil
	}
	return us.updateFileStorage(url), nil
}

func (us *URLStorage) ReturnBatchShortURL(ctx context.Context, batch APIBatchRequestA) (*APIBatchShorURLs, error) {
	if us.UseDBEngine {
		batchBytes, err := json.Marshal(batch)
		if err != nil {
			return nil, err
		}

		respBytes := us.returnPgDBBatchShortURLs(ctx, batchBytes, us.BaseURL+"/")
		batchResponse := APIBatchResponseA{}
		err = json.Unmarshal(respBytes, &batchResponse)
		if err != nil {
			return nil, err
		}

		return &APIBatchShorURLs{
			http.StatusCreated,
			batchResponse,
		}, nil
	}
	return nil, errors.New("unsupported storage engine")
}

func (us *URLStorage) ReturnFullURL(ctx context.Context, short string) string {
	if us.UseDBEngine {
		return us.returnPgDBFullURL(ctx, short)
	}
	return us.returnFileStorageFullURL(short)
}

func (us *URLStorage) returnFileStorageFullURL(short string) string {
	if val, ok := us.Storage.GETStorage[short]; !ok {
		return ""
	} else {
		return val[0]
	}
}

func (us *URLStorage) updateFileStorage(url string) *APIShorURL {
	// http.StatusCreated
	if val, ok := us.Storage.POSTStorage[url]; !ok {
		rs := us.randomString()
		sha256 := us.sha256Sum(url + rs)
		us.Storage.POSTStorage[url] = []string{sha256, rs}
		us.Storage.GETStorage[sha256] = []string{url, rs}
		return &APIShorURL{
			HTTPCode: http.StatusCreated,
			ShortURL: sha256,
			BaseURL:  us.BaseURL,
		}
	} else {
		// http.StatusConflict
		return &APIShorURL{
			HTTPCode: http.StatusConflict,
			ShortURL: val[0],
			BaseURL:  us.BaseURL,
		}
	}
}

func (us *URLStorage) fileStorageToByteA() ([]byte, error) {
	if b, err := json.Marshal(us.Storage); err != nil {
		return nil, err
	} else {
		return b, nil
	}
}

func (us *URLStorage) SaveStorage(storageName string) error {
	if !us.UseDBEngine {
		if err := us.prepareStorageFile(storageName); err != nil {
			return err
		}
		sData, err := us.fileStorageToByteA()
		if err != nil {
			return err
		}
		if err = os.WriteFile(storageName, sData, os.FileMode(0644)); err != nil {
			return err
		}
		return nil
	}
	return us.PgDB.Close()
}

func (us *URLStorage) prepareStorageFile(storageName string) error {
	_, err := os.Stat(storageName)
	if err != nil {
		if os.IsNotExist(err) {
			var f *os.File
			f, err = os.OpenFile(storageName, os.O_RDWR|os.O_CREATE, os.FileMode(0644))
			if err != nil {
				return err
			}
			if err = f.Close(); err != nil {
				return err
			}
			return nil
		} else {
			return err
		}
	}
	return nil
}

func (us *URLStorage) ResetFileStorage(storageName, addr string) error {
	us.StorageName = storageName
	us.BaseURL = addr
	return us.loadFileStorage(storageName)
}

func (us *URLStorage) loadFileStorage(storageName string) error {
	err := us.prepareStorageFile(storageName)
	if err != nil {
		return err
	}
	fBytes, err := os.ReadFile(storageName)
	if err != nil {
		return err
	}
	if len(fBytes) != 0 {
		err = json.Unmarshal(fBytes, &us.Storage)
		if err != nil {
			return err
		}
	}
	return nil
}

func (us *URLStorage) sha256Sum(s string) string {
	return common.GetSHA256Sum(s)
}

func (us *URLStorage) randomString() string {
	return common.GetRandomString(common.MinRndStrLen)
}

func NewURLFileStorage(storageName string) (*URLStorage, error) {
	us := &URLStorage{
		StorageName: storageName,
		Storage:     Storage{POSTStorage: make(map[string][]string), GETStorage: make(map[string][]string)},
	}
	err := us.loadFileStorage(storageName)
	return us, err
}

func (us *URLStorage) SetBaseURL(addr string) {
	us.BaseURL = addr
}
