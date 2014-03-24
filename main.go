package main

import (
	"log"

	"github.com/xpeppers/twitter_listener/listener"
)

func main() {
	store := listener.NewAppStore("192.168.10.10:6379 db=1")
	err := listener.StartAll(store)
	if err != nil {
		log.Fatal(err)
	}
}
