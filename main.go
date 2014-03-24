package main

import (
	"flag"
	"log"

	"github.com/xpeppers/twitter_listener/listener"
)

var (
	appStoreUrl   string
	singleAppName string
)

func main() {
	flag.StringVar(&appStoreUrl, "db", "127.0.0.1:6379 db=1",
		"Redis database connection URL.")
	flag.StringVar(&singleAppName, "a", "",
		"Launch listener only for a single app. Otherwise listen for all apps.")
	flag.Parse()

	store := listener.NewAppStore(appStoreUrl)

	var startErr error
	if singleAppName != "" {
		startErr = listener.StartOne(store, singleAppName)
	} else {
		startErr = listener.StartAll(store)
	}

	if startErr != nil {
		log.Fatal(startErr)
	}
}
