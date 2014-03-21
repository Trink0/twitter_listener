package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/fiorix/go-redis/redis"
)

// Application is a customer/merchant
type Application struct {
	ApiKey           string `json:"apiKey"`
	Name             string `json:"name"`
	TwConsumerKey    string `json:"twitterConsumerKey"`
	TwConsumerSecret string `json:"twitterConsumerSecret"`
	TwAccessToken    string `json:"twitterAccessToken"`
	TwTokenSecret    string `json:"twitterTokenSecret"`
}

type Listener struct {
	app *Application
}

func (l *Listener) start(c chan int) {
	log.Printf("Starting Listner: %s", l.app.Name)
	time.Sleep(time.Second * 5)
	c <- 1
}

func main() {
	rc := redis.New("192.168.10.10:6379 db=1")
	keys, err := rc.Keys("*")
	if err != nil {
		log.Fatal(err)
	}
	c := make(chan int)
	for _, appName := range keys {
		jsonApp, err := rc.Get(appName)
		if err != nil {
			log.Fatal(err)
		}

		if jsonApp == "" {
			log.Fatal("Application not found")
		}
		storedApp := &Application{}
		if err := json.Unmarshal([]byte(jsonApp), storedApp); err != nil {
			log.Fatal(err)
		}

		listener := &Listener{app: storedApp}
		go listener.start(c)
	}
	count := 0
	for {
		status := <-c
		log.Printf("Listener Exiting with status: %d", status)
		if count += 1; count == len(keys) {
			log.Printf("Exiting application")
			break
		}
	}
}
