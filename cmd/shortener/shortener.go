package main

import (
	"flag"
	"fmt"
	"github.com/AVitaminOz-z-z/urlshortener.git/internal/config"
	"github.com/AVitaminOz-z-z/urlshortener.git/internal/router"
	"github.com/AVitaminOz-z-z/urlshortener.git/internal/storage"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var srvAddr, baseAddr *string

func parseArgs(aEnv *config.AppEnv) {
	// parse cmd args
	flag.CommandLine.SetOutput(os.Stdout)
	srvAddr = flag.String("a", aEnv.SrvAddr, "[-a\t| --a]\t->\tserver:port")
	baseAddr = flag.String("b", aEnv.BaseAddr, "[-b\t| --b]\t->\thttp(s)://server:port")
	if len(os.Args[1:]) > 0 {
		flag.Parse()
	}

	if *srvAddr != "" {
		aEnv.SrvAddr = *srvAddr
	}

	if *baseAddr != "" {
		aEnv.BaseAddr = *baseAddr
	}
}

func main() {
	// loading environment
	appEnv, err := config.NewAppEnv()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	// creating url-storage
	urlStorage, err := storage.NewURLStorage(appEnv.StorageName)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	defer func() {
		_ = urlStorage.SaveStorage(appEnv.StorageName)
	}()

	// parse cmd args
	parseArgs(appEnv)
	urlStorage.SetBaseAddr(appEnv.BaseAddr)

	// creating router
	R := router.NewURLRouter(urlStorage)

	// notify os signals
	doneOS := make(chan os.Signal, 1)
	signal.Notify(doneOS, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// listen http-server
	go func() {
		fmt.Printf("AVitaminOz-z-z HTTP-Server v0.1 started on %v\n", appEnv.SrvAddr)
		err = http.ListenAndServe(appEnv.SrvAddr, R)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	s := <-doneOS
	fmt.Printf("AVitaminOz-z-z HTTP-Server v0.1 (%s) was finished. Got <%v> signal.\n", appEnv.SrvAddr, s)
}
