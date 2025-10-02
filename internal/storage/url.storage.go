package storage

import (
	"encoding/json"
	"github.com/AVitaminOz-z-z/urlshortener.git/internal/common"
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

func (us *URLStorage) ReturnShortURL(url string) string {
	if us.UseDBEngine {
		return us.returnPgDBShortURL(url)
	}
	return us.updateFileStorage(url)
}

func (us *URLStorage) ReturnFullURL(short string) string {
	if us.UseDBEngine {
		return us.returnPgDBFullURL(short)
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

func (us *URLStorage) fileStorageToByteA() ([]byte, error) {
	if b, err := json.Marshal(us.Storage); err != nil {
		return nil, err
	} else {
		return b, nil
	}
}

func (us *URLStorage) SaveFileStorage(storageName string) error {
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
	return nil
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
