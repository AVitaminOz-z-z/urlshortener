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
	defPgDSN       = "postgres://user:pass@localhost:5432/dbname?sslmode=disable"
	defSrvKeyMsg   = "37Eq93iXuIKDqSGd"
)

type AppArgs struct {
	ArgsLen     int
	EnvFile     string
	SrvAddress  string
	BaseURL     string
	StorageName string
	PgDSN       string
	SrvKeyMsg   string
}

type AppEnv struct {
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
			defPgDSN,
			defSrvKeyMsg,
		},
	}
}

func (env *AppEnv) SrvKey() []byte {
	return []byte(env.SrvKeyMsg)
}

func (env *AppEnv) LoadEnv(envFile string) error {
	if err := godotenv.Load(envFile); err != nil {
		return err
	}
	env.StorageName = loadEnv("FILE_STORAGE_PATH")
	env.SrvAddress = loadEnv("SERVER_ADDRESS")
	env.BaseURL = loadEnv("BASE_URL")
	env.PgDSN = loadEnv("DATABASE_DSN")
	env.PgDSN = loadEnv("DATABASE_DSN")
	env.SrvKeyMsg = loadEnv("SRV_KEY")
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
	if lParam == "" {
		switch key {
		case "FILE_STORAGE_PATH":
			{
				lParam = defStorageName
				break
			}
		case "SERVER_ADDRESS":
			{
				lParam = defSrvAddress
				break
			}
		case "BASE_URL":
			{
				lParam = defBaseURL
				break
			}
		case "DATABASE_DSN":
			{
				lParam = defPgDSN
				break
			}
		case "SRV_KEY":
			{
				lParam = defSrvKeyMsg
				break
			}
		}
	}
	return lParam
}
