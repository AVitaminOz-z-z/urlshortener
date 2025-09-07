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
	// new empty application.go
	app := application.NewApp()

	// prepare application.go
	if err := app.PrepareApp(); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	// save url-storage to file
	defer func() {
		_ = app.SaveURLStorage()
	}()

	// notify os signals
	doneOS := make(chan os.Signal, 1)
	signal.Notify(doneOS, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// run application
	go func() {
		fmt.Printf("%s started on %v\n", common.AppSrvVersion, app.GetServerAddr())
		if err := app.OnAir(); err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
	}()

	// finish application
	s := <-doneOS
	fmt.Printf("%s (%s) was finished via <%v> signal\n", common.AppSrvVersion, app.GetServerAddr(), s)
}
