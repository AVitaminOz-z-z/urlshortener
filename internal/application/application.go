package application

import (
	"context"
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
	return &AppServer{}
}

func (a *AppServer) setEnv() error {
	env, err := config.NewAppEnv()
	if err != nil {
		return err
	}
	a.AppEnv = env
	return nil
}

func (a *AppServer) Close() error {
	return a.SaveURLStorage()
}

func (a *AppServer) setPgDB(ctx context.Context, dsn string) error {
	return a.AppStorage.SetPgDB(ctx, dsn)
}

func (a *AppServer) setFileStorage(name string) error {
	uStorage, err := storage.NewURLFileStorage(name)
	if err != nil {
		return err
	}
	uStorage.SetBaseURL(a.AppEnv.BaseURL)
	a.AppStorage = uStorage
	return nil
}

func (a *AppServer) setLogger(w io.Writer) {
	a.AppLogger = logger.NewLogger(w)
}

func (a *AppServer) setRouter() {
	a.AppRouter = router.NewURLRouter(a.getURLStorage(), a.GetLogger(), a.getSrvKey())
}

func (a *AppServer) getRouter() chi.Router {
	return a.AppRouter
}

func (a *AppServer) getPgDB() *storage.PgDB {
	return a.AppStorage.PgDB
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

func (a *AppServer) getPgDSN() string {
	return a.getAppEnv().PgDSN
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

func (a *AppServer) getSrvKey() []byte {
	return a.AppEnv.SrvKey()
}

func (a *AppServer) getOSArgs(env *config.AppEnv) *config.AppArgs {
	pAppArgs := &config.AppArgs{}

	flag.CommandLine.SetOutput(os.Stdout)
	flag.StringVar(&pAppArgs.EnvFile, "env", env.EnvFile, "[-env\t| --env]\t->\t/path/to/env/file")
	flag.StringVar(&pAppArgs.SrvAddress, "a", env.SrvAddress, "[-a\t| --a]\t->\t[server]:port")
	flag.StringVar(&pAppArgs.BaseURL, "b", env.BaseURL, "[-b\t| --b]\t->\thttp(s)://server:port")
	flag.StringVar(&pAppArgs.StorageName, "f", env.StorageName, "[-f\t| --f]\t->\tpath/to/storage/file")
	flag.StringVar(&pAppArgs.PgDSN, "d", env.PgDSN, "[-d\t| --d]\t->\tPgSQL DSN conn string")

	if len(os.Args[1:]) > 0 {
		flag.Parse()
		pAppArgs.ArgsLen = len(os.Args[1:])
	}

	return pAppArgs
}

func (a *AppServer) resetAppArgs(ctx context.Context, args *config.AppArgs) error {
	// no args
	if args.ArgsLen == 0 {
		return nil
	}
	// new args
	a.AppEnv.AppArgs = *args

	// reset file storage
	if err := a.getURLStorage().ResetFileStorage(a.getAppEnv().StorageName, a.getAppEnv().BaseURL); err != nil {
		return err
	}

	// reset DB-engine
	_ = a.setPgDB(ctx, a.getPgDSN())

	return nil
}

func (a *AppServer) PrepareApp() error {
	// main context
	ctx := context.Background()

	// loading environment
	if err := a.setEnv(); err != nil {
		return err
	}

	// creating url-storage
	if err := a.setFileStorage(a.getURLStorageName()); err != nil {
		return err
	}

	// creating database engine
	_ = a.setPgDB(ctx, a.getPgDSN())

	// creating logger
	a.setLogger(os.Stdout)

	// reset App-args via OS-args
	if err := a.resetAppArgs(ctx, a.getOSArgs(a.getAppEnv())); err != nil {
		return err
	}

	// creating router
	a.setRouter()

	return nil
}

func (a *AppServer) OnAir() error {
	return http.ListenAndServe(a.GetServerAddr(), a.getRouter())
}
