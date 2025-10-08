package storage

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/AVitaminOz-z-z/urlshortener.git/internal/common"
	m "github.com/AVitaminOz-z-z/urlshortener.git/internal/model"
	"net/http"
	"os"
)

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

func (us *URLStorage) ReturnShortURL(ctx context.Context, url string) (*m.APIShorURL, error) {
	prefix := us.BaseURL + "/"
	if us.UseDBEngine {
		httpCode := http.StatusCreated
		// handle PgDBStorage Errors
		short, err := us.returnPgDBShortURL(ctx, url, prefix)
		if err != nil {
			if !errors.Is(err, common.ErrStatusConflict) {
				return nil, err
			} else {
				httpCode = http.StatusConflict
			}
		}
		return &m.APIShorURL{
			HTTPCode: httpCode,
			ShortURL: short,
		}, nil
	}
	return us.returnFileStorageShortURL(url, prefix), nil
}

func (us *URLStorage) ReturnBatchShortURL(ctx context.Context, batch m.APIBatchRequestA) (*m.APIBatchShorURLs, error) {
	prefix := us.BaseURL + "/"
	if us.UseDBEngine {
		batchBytes, err := json.Marshal(batch)
		if err != nil {
			return nil, err
		}
		respBytes, err := us.returnPgDBBatchShortURLs(ctx, batchBytes, prefix)
		if err != nil {
			return nil, err
		}
		batchResponse := m.APIBatchResponseA{}
		err = json.Unmarshal(respBytes, &batchResponse)
		if err != nil {
			return nil, err
		}
		return &m.APIBatchShorURLs{
			HTTPCode:          http.StatusCreated,
			APIBatchResponseA: batchResponse,
		}, nil
	}
	return us.returnFileStorageBatchShortURLs(batch, prefix)
}

func (us *URLStorage) returnFileStorageURL(url string, prefix string) (string, bool) {
	if val, ok := us.Storage.POSTStorage[url]; !ok {
		rs := us.randomString()
		sha256 := us.sha256Sum(url + rs)
		us.Storage.POSTStorage[url] = []string{sha256, rs}
		us.Storage.GETStorage[sha256] = []string{url, rs}
		return prefix + sha256, ok
	} else {
		return prefix + val[0], ok
	}
}

func (us *URLStorage) returnFileStorageBatchShortURLs(batch m.APIBatchRequestA, prefix string) (*m.APIBatchShorURLs, error) {
	batchResponseA := make(m.APIBatchResponseA, len(batch))
	batchShorURLs := &m.APIBatchShorURLs{
		HTTPCode:          http.StatusCreated,
		APIBatchResponseA: batchResponseA,
	}
	for i, v := range batch {
		batchResponseA[i].CorrelationID = v.CorrelationID
		sURL, _ := us.returnFileStorageURL(v.OriginalURL, prefix)
		batchResponseA[i].ShortURL = sURL

	}
	return batchShorURLs, nil
}

func (us *URLStorage) ReturnFullURL(ctx context.Context, short string) (string, error) {
	if us.UseDBEngine {
		return us.returnPgDBFullURL(ctx, short)
	}
	return us.returnFileStorageFullURL(short)
}

func (us *URLStorage) returnFileStorageFullURL(short string) (string, error) {
	if val, ok := us.Storage.GETStorage[short]; !ok {
		return "", nil
	} else {
		return val[0], nil
	}
}

func (us *URLStorage) returnFileStorageShortURL(url string, prefix string) *m.APIShorURL {
	short, exists := us.returnFileStorageURL(url, prefix)
	httpCode := http.StatusCreated
	if exists {
		httpCode = http.StatusConflict
	}
	return &m.APIShorURL{
		HTTPCode: httpCode,
		ShortURL: short,
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
