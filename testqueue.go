package main

import (
	// "log"

	"github.com/xpeppers/twitter_listener/listener"
)

func main() {
	c := make(chan *listener.Tweet)
	queue := listener.NewQueue("127.0.0.1:2229")
	queue.Start(c)
	c <- &listener.Tweet{
		ID:   "123123123",
		Text: "Hello there",
		User: listener.TweetUser{ID: "4545"},
	}
}
