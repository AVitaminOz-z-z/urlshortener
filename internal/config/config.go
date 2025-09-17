package config

import (
	"github.com/joho/godotenv"
	"os"
)

const (
	defEnvFile     = "./.env"
	defStorageName = "./.storage"
	defSrvAddress  = "localhost:8080"
	defBaseURL     = "http://localhost:8080"
)

type AppArgs struct {
	ArgsLen     int
	EnvFile     string
	SrvAddress  string
	BaseURL     string
	StorageName string
}

type AppEnv struct {
	//StorageName string
	AppArgs
}

func newDefAppEnv() *AppEnv {
	return &AppEnv{
		//defStorageName,
		AppArgs{0,
			defEnvFile,
			defSrvAddress,
			defSrvAddress,
			defStorageName,
		},
	}
}

func (env *AppEnv) LoadEnv(envFile string) error {
	if err := godotenv.Load(envFile); err != nil {
		return err
	}
	//env.StorageName = loadEnv("URL_STORAGE")
	//env.StorageName = loadEnv("FILE_STORAGE_PATH")
	env.StorageName = loadEnv("FILE_STORAGE_PATH")
	env.SrvAddress = loadEnv("SERVER_ADDRESS")
	env.BaseURL = loadEnv("BASE_URL")
	return nil
}

func NewAppEnv() (*AppEnv, error) {
	aEnv := newDefAppEnv()
	if err := aEnv.LoadEnv(aEnv.EnvFile); err != nil {
		return nil, err
	}
	return aEnv, nil
}

func loadEnv(key string) string {
	// load environments
	lParam := os.Getenv(key)
	switch key {
	/*case "URL_STORAGE":
	{
		if lParam == "" {
			lParam = defStorageName
		}
		break
	}*/
	case "FILE_STORAGE_PATH":
		{
			if lParam == "" {
				lParam = defStorageName
			}
			break
		}
	case "SERVER_ADDRESS":
		{
			if lParam == "" {
				lParam = defSrvAddress
			}
			break
		}
	case "BASE_URL":
		{
			if lParam == "" {
				lParam = defBaseURL
			}
			break
		}
	}
	return lParam
}
