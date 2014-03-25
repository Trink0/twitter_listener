package main

import (
	"flag"
	"log"

	"github.com/xpeppers/twitter_listener/listener"
)

var (
	// Apps & users database URL
	dbURL string
	// Users DB
	userDB int
	// Apps DB
	appDB int
)

func main() {
	flag.StringVar(&dbURL, "dburl", "127.0.0.1:6379", "Redis database connection URL.")
	flag.IntVar(&userDB, "userdb", 0, "Redis Users DB number.")
	flag.IntVar(&appDB, "appdb", 1, "Redis Applications DB number.")
	flag.Parse()

	if len(flag.Args()) > 1 {
		flag.Usage()
		return
	}

	store := listener.NewStore(dbURL, appDB, userDB)

	var startErr error
	if appName := flag.Arg(0); appName != "" {
		startErr = listener.StartOne(store, appName)
	} else {
		startErr = listener.StartAll(store)
	}

	if startErr != nil {
		log.Fatal(startErr)
	}
}
