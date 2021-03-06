package main

import (
	"flag"
	"log"

	"github.com/Trink0/twitter_listener/listener"
	"github.com/Trink0/twitter_listener/source"
)

var (
	// Apps & users database URL
	dbURL string
	// Users DB
	userDB int
	// Apps DB
	appDB     int
	queueURL  string
	queueName string
)

func main() {
	flag.StringVar(&dbURL, "dburl", "127.0.0.1:6379", "Redis database connection URL.")
	flag.IntVar(&userDB, "userdb", 0, "Redis Users DB number.")
	flag.IntVar(&appDB, "appdb", 1, "Redis Applications DB number.")
	flag.StringVar(&queueURL, "qurl", "127.0.0.1:22133", "Kestrel queue connection URL.")
	flag.StringVar(&queueName, "qname", "social-web-activities", "Kestrel queue name.")
	flag.Parse()

	if len(flag.Args()) > 1 {
		flag.Usage()
		return
	}

	source := source.NewRedisSource(dbURL, appDB, userDB)
	dest := listener.NewQueue(queueURL, queueName)

	var startErr error
	if appName := flag.Arg(0); appName != "" {
		startErr = listener.StartOne(appName, source, dest)
	} else {
		startErr = listener.StartAll(source, dest)
	}

	if startErr != nil {
		log.Fatal(startErr)
	}
}
