package listener

import (
	"encoding/json"
	"log"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

// Activity is a Beancounter Activity object being pushed down the queue.
type Activity struct {
	Id              string          `json:"id"`
	Verb            string          `json:"verb"`
	Object          ActivityObject  `json:"object"`
	Context         ActivityContext `json:"context"`
	ApplicationName string          `json:"applicationName"`
}

type ActivityObject struct {
	Type string   `json:"type"`
	URL  string   `json:"url"`
	Text string   `json:"text"`
	Urls []string `json:"urls"`
}

type ActivityContext struct {
	Service  string `json:"service"`
	Username string `json:"username"`
	Date     int64  `json:"date"`
}

type Queue interface {
	Start(qc chan *Tweet)
}

func NewQueue(endPoint string, queueName string) Queue {
	client := memcache.New(endPoint)
	putFn := func(payload []byte) {
		putErr := client.Set(&memcache.Item{Key: queueName, Value: payload})
		if putErr != nil {
			log.Printf("ERROR queue: %v", putErr)
		}
	}
	return &kestrelQueue{endPoint, queueName, putFn}
}

type kestrelQueue struct {
	endPoint  string
	queueName string
	putFn     func([]byte)
}

func (k *kestrelQueue) Start(qc chan *Tweet) {
	log.Printf("Connecting to Kestrel on %s", k.endPoint)
	go k.loop(qc)
}

func (k *kestrelQueue) loop(qc chan *Tweet) {
	for {
		tweet := <-qc
		activity := tweetToActivity(tweet)
		log.Printf("ENQUEUE: %+v", activity)
		payload, err := json.Marshal(activity)
		if err != nil {
			log.Printf("ERROR encoding activity: %v", err)
			continue
		}

		k.putFn(payload)
	}
}

func tweetToActivity(tweet *Tweet) *Activity {
	activity := &Activity{
		Id:   uuid(),
		Verb: "TWEET",
		Object: ActivityObject{
			Type: "TWEET",
			URL:  tweet.statusUrl(),
			Text: tweet.Text,
		},
		Context: ActivityContext{
			Service:  "twitter",
			Username: tweet.User.ID,
		},
		ApplicationName: tweet.AppName,
	}
	createdAt, _ := time.Parse(time.RubyDate, tweet.CreatedAt)
	activity.Context.Date = createdAt.UnixNano() / 1E6

	activity.Object.Urls = make([]string, len(tweet.Entities.URLs))
	for i, entityUrl := range tweet.Entities.URLs {
		activity.Object.Urls[i] = entityUrl.ExpandedURL
	}

	return activity
}
