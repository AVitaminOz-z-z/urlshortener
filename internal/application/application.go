package application

import (
	"flag"
	"github.com/AVitaminOz-z-z/urlshortener.git/internal/config"
	"github.com/AVitaminOz-z-z/urlshortener.git/internal/router"
	"github.com/AVitaminOz-z-z/urlshortener.git/internal/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
	"os"
)

type AppServer struct {
	AppEnv     *config.AppEnv
	AppStorage *storage.URLStorage
	AppRouter  chi.Router
}

func NewApp() *AppServer {
	return &AppServer{
		nil,
		nil,
		nil,
	}
}

func (a *AppServer) SetEnv() error {
	if env, err := config.NewAppEnv(); err != nil {
		return err
	} else {
		a.AppEnv = env
		return nil
	}
}

func (a *AppServer) SetStorage(name string) error {
	if uStorage, err := storage.NewURLStorage(name); err != nil {
		return err
	} else {
		a.AppStorage = uStorage
		a.GetURLStorage().SetBaseURL(a.GetAppEnv().BaseURL)
		return nil
	}
}

func (a *AppServer) SetRouter(storage *storage.URLStorage) {
	a.AppRouter = router.NewURLRouter(storage)
}

func (a *AppServer) GetRouter() chi.Router {
	return a.AppRouter
}

func (a *AppServer) GetServerAddr() string {
	return a.GetAppEnv().SrvAddress
}

func (a *AppServer) GetURLStorageName() string {
	return a.GetAppEnv().StorageName
}

func (a *AppServer) GetURLStorage() *storage.URLStorage {
	return a.AppStorage
}

func (a *AppServer) SaveURLStorage() error {
	return a.GetURLStorage().SaveStorage(a.GetURLStorageName())
}

func (a *AppServer) GetAppEnv() *config.AppEnv {
	return a.AppEnv
}

func (a *AppServer) GetOSArgs(env *config.AppEnv) *config.AppArgs {
	type ptrAppArgs struct {
		ArgsLen    int
		EnvFile    *string
		SrvAddress *string
		BaseURL    *string
	}

	pAppArgs := ptrAppArgs{
		0,
		nil,
		nil,
		nil,
	}

	flag.CommandLine.SetOutput(os.Stdout)
	pAppArgs.EnvFile = flag.String("env", env.EnvFile, "[-env\t| --env]\t->\t/path/to/env/file")
	pAppArgs.SrvAddress = flag.String("a", env.SrvAddress, "[-a\t| --a]\t->\tserver:port")
	pAppArgs.BaseURL = flag.String("b", env.BaseURL, "[-b\t| --b]\t->\thttp(s)://server:port")

	osArgsLen := len(os.Args[1:])
	if osArgsLen > 0 {
		flag.Parse()
	}

	return &config.AppArgs{
		ArgsLen:    osArgsLen,
		EnvFile:    *pAppArgs.EnvFile,
		SrvAddress: *pAppArgs.SrvAddress,
		BaseURL:    *pAppArgs.BaseURL,
	}
}

func (a *AppServer) ResetAppArgs(args *config.AppArgs) {
	if args.ArgsLen > 0 {
		a.AppEnv.AppArgs = *args
		a.GetURLStorage().SetBaseURL(a.GetAppEnv().BaseURL)
	}
}

func (a *AppServer) PrepareApp() error {
	// loading environment
	if err := a.SetEnv(); err != nil {
		return err
	}

	// creating url-storage
	if err := a.SetStorage(a.GetURLStorageName()); err != nil {
		return err
	}

	// reset App Args via OS Args
	a.ResetAppArgs(a.GetOSArgs(a.GetAppEnv()))

	// creating router
	a.SetRouter(a.GetURLStorage())

	return nil
}

func (a *AppServer) OnAir() error {
	return http.ListenAndServe(a.GetServerAddr(), a.GetRouter())
}
