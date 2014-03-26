package main

import (
	"github.com/xpeppers/twitter_listener/listener"
	"time"
)

func main() {
	c := make(chan *listener.Tweet)
	queue := listener.NewQueue("192.168.10.10:22133")
	queue.Start(c)
	c <- &listener.Tweet{
		ID:   "123123123",
		Text: "Hello there",
		User: listener.TweetUser{ID: "4545"},
	}
	time.Sleep(time.Second * 10)
}
