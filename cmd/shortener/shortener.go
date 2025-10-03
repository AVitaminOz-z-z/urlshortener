package main

import (
	"fmt"
	"github.com/AVitaminOz-z-z/urlshortener.git/internal/application"
	"github.com/AVitaminOz-z-z/urlshortener.git/internal/common"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// new empty application
	app := application.NewApp()

	// prepare application
	if err := app.PrepareApp(); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	// save url-storage to file
	defer func() {
		_ = app.Close()
	}()

	// notify os signals
	doneOS := make(chan os.Signal, 1)
	signal.Notify(doneOS, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// run application
	log := app.GetLogger()
	go func() {
		log.Info(fmt.Sprintf("%s started on %v {DB-Engine: %v}", common.AppSrvVersion, app.GetServerAddr(), app.AppStorage.UseDBEngine))
		if err := app.OnAir(); err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
	}()

	// finish application
	s := <-doneOS
	log.Info(fmt.Sprintf("%s (%s) was finished via <%v> signal", common.AppSrvVersion, app.GetServerAddr(), s))
}
