package config

import (
	"github.com/joho/godotenv"
	"os"
)

const (
	defStorageName = "./.storage"
	defSrvAddr     = "localhost:8080"
	defBaseAddr    = "http://localhost:8080"
)

type AppEnv struct {
	StorageName string
	SrvAddr     string
	BaseAddr    string
}

func NewAppEnv() (*AppEnv, error) {
	aEnv := &AppEnv{
		defStorageName,
		defSrvAddr,
		defBaseAddr,
	}
	dir, err := os.Getwd()
	if err != nil {
		// return default env
		return aEnv, nil
		//return nil, err
	}
	err = godotenv.Load(dir + "/.env")
	if err != nil {
		// return default env
		return aEnv, nil
		//return nil, err
	}
	aEnv = &AppEnv{
		loadEnv(dir, "URL_STORAGE"),
		loadEnv(dir, "HTTP_SERVER"),
		loadEnv(dir, "BASE_ADDR"),
	}
	return aEnv, nil
}

func loadEnv(dir, key string) string {
	// load environments
	lParam := os.Getenv(key)
	// dir structure
	switch key {
	case "URL_STORAGE":
		{
			if lParam == "" {
				lParam = dir + "/.storage"
			}
			break
		}
	case "HTTP_SERVER":
		{
			if lParam == "" {
				lParam = "localhost:8080"
			}
			break
		}
	case "BASE_ADDR":
		{
			if lParam == "" {
				lParam = "http://localhost:8080"
			}
			break
		}
	}
	return lParam
}
