package application

import (
	"flag"
	"github.com/AVitaminOz-z-z/urlshortener.git/internal/config"
	"github.com/AVitaminOz-z-z/urlshortener.git/internal/logger"
	"github.com/AVitaminOz-z-z/urlshortener.git/internal/router"
	"github.com/AVitaminOz-z-z/urlshortener.git/internal/storage"
	"github.com/go-chi/chi/v5"
	"io"
	"log/slog"
	"net/http"
	"os"
)

type AppServer struct {
	AppEnv     *config.AppEnv
	AppStorage *storage.URLStorage
	AppRouter  chi.Router
	AppLogger  *slog.Logger
}

func NewApp() *AppServer {
	return &AppServer{
		nil,
		nil,
		nil,
		nil,
	}
}

func (a *AppServer) setEnv() error {
	if env, err := config.NewAppEnv(); err != nil {
		return err
	} else {
		a.AppEnv = env
		return nil
	}
}

func (a *AppServer) setStorage(name string) error {
	if uStorage, err := storage.NewURLStorage(name); err != nil {
		return err
	} else {
		a.AppStorage = uStorage
		a.getURLStorage().SetBaseURL(a.getAppEnv().BaseURL)
		return nil
	}
}

func (a *AppServer) setLogger(w io.Writer) {
	a.AppLogger = logger.NewLogger(w)
}

func (a *AppServer) setRouter() {
	a.AppRouter = router.NewURLRouter(a.getURLStorage(), a.GetLogger())
}

func (a *AppServer) getRouter() chi.Router {
	return a.AppRouter
}

func (a *AppServer) GetLogger() *slog.Logger {
	return a.AppLogger
}

func (a *AppServer) GetServerAddr() string {
	return a.getAppEnv().SrvAddress
}

func (a *AppServer) getURLStorageName() string {
	return a.getAppEnv().StorageName
}

func (a *AppServer) getURLStorage() *storage.URLStorage {
	return a.AppStorage
}

func (a *AppServer) SaveURLStorage() error {
	return a.getURLStorage().SaveStorage(a.getURLStorageName())
}

func (a *AppServer) getAppEnv() *config.AppEnv {
	return a.AppEnv
}

func (a *AppServer) getOSArgs(env *config.AppEnv) *config.AppArgs {
	pAppArgs := &config.AppArgs{}

	flag.CommandLine.SetOutput(os.Stdout)
	flag.StringVar(&pAppArgs.EnvFile, "env", env.EnvFile, "[-env\t| --env]\t->\t/path/to/env/file")
	flag.StringVar(&pAppArgs.SrvAddress, "a", env.SrvAddress, "[-a\t| --a]\t->\t[server]:port")
	flag.StringVar(&pAppArgs.BaseURL, "b", env.BaseURL, "[-b\t| --b]\t->\thttp(s)://server:port")
	flag.StringVar(&pAppArgs.StorageName, "f", env.StorageName, "[-f\t| --f]\t->\tpath/to/storage/file")

	if len(os.Args[1:]) > 0 {
		flag.Parse()
		pAppArgs.ArgsLen = len(os.Args[1:])
	}

	return pAppArgs
}

func (a *AppServer) resetAppArgs(args *config.AppArgs) error {
	var err error
	if args.ArgsLen > 0 {
		a.AppEnv.AppArgs = *args
		err = a.getURLStorage().ResetStorage(a.getAppEnv().StorageName, a.getAppEnv().BaseURL)
		/*a.getURLStorage().SetBaseURL(a.getAppEnv().BaseURL)*/
	}
	return err
}

func (a *AppServer) PrepareApp() error {
	// loading environment
	if err := a.setEnv(); err != nil {
		return err
	}

	// creating url-storage
	if err := a.setStorage(a.getURLStorageName()); err != nil {
		return err
	}

	// creating logger
	a.setLogger(os.Stdout)

	// reset App-args via OS-args
	if err := a.resetAppArgs(a.getOSArgs(a.getAppEnv())); err != nil {
		return err
	}

	// creating router
	a.setRouter()

	return nil
}

func (a *AppServer) OnAir() error {
	return http.ListenAndServe(a.GetServerAddr(), a.getRouter())
}
