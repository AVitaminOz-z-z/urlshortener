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
	BaseAddr    string `json:"BaseAddr,omitempty"`
	Storage     `json:"Storage,omitempty"`
}

func (us *URLStorage) ReturnShortURL(url string) string {
	return us.updateStorage(url)
}

func (us *URLStorage) ReturnFullURL(short string) string {
	if val, ok := us.Storage.GETStorage[short]; !ok {
		return ""
	} else {
		return val[0]
	}
}

func (us *URLStorage) updateStorage(url string) string {
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

func (us *URLStorage) storageToByteA() ([]byte, error) {
	if b, err := json.Marshal(us.Storage); err != nil {
		return nil, err
	} else {
		return b, nil
	}
}

func (us *URLStorage) SaveStorage(StorageName string) error {
	err := us.prepareStorageFile(StorageName)
	if err != nil {
		return err
	}
	sData, err := us.storageToByteA()
	if err != nil {
		return err
	}
	err = os.WriteFile(StorageName, sData, os.FileMode(0644))
	if err != nil {
		return err
	}
	return nil
}

func (us *URLStorage) prepareStorageFile(StorageName string) error {
	_, err := os.Stat(StorageName)
	if err != nil {
		if os.IsNotExist(err) {
			var f *os.File
			f, err = os.OpenFile(StorageName, os.O_RDWR|os.O_CREATE, os.FileMode(0644))
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

func (us *URLStorage) loadStorage(StorageName string) error {
	err := us.prepareStorageFile(StorageName)
	if err != nil {
		return err
	}
	fBytes, err := os.ReadFile(StorageName)
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

func NewURLStorage(StorageName string) (*URLStorage, error) {
	us := &URLStorage{
		StorageName,
		"",
		Storage{make(map[string][]string), make(map[string][]string)},
	}
	var err error
	err = us.loadStorage(StorageName)
	return us, err
}

func (us *URLStorage) SetBaseAddr(addr string) {
	us.BaseAddr = addr
}
