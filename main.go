package main

import (
	"flag"
	"log"

	"github.com/xpeppers/twitter_listener/listener"
)

var (
	StoreUrl      string
	userStoreUrl  string
	singleAppName string
)

func main() {
	flag.StringVar(&StoreUrl, "db", "127.0.0.1:6379 db=1",
		"Redis database connection URL.")
	flag.StringVar(&userStoreUrl, "dbuser", "127.0.0.1:6379 db=0",
		"Redis database connection URL.")
	flag.StringVar(&singleAppName, "app", "",
		"Launches single listener only for a specific app if not empty.")
	flag.Parse()

	if len(flag.Args()) > 0 {
		flag.Usage()
		return
	}

	store := listener.NewStore(StoreUrl)
	userStore := listener.NewStore(userStoreUrl)

	var startErr error
	if singleAppName != "" {
		startErr = listener.StartOne(store, userStore, singleAppName)
	} else {
		startErr = listener.StartAll(store, userStore)
	}

	if startErr != nil {
		log.Fatal(startErr)
	}
}
