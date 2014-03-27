package listener

import (
	"encoding/json"
	"log"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

// Tweet is what Twitter sends back as stream items.
type Tweet struct {
	ID   string `json:"id_str"`
	Text string
	// Indicates whether the tweet was truncated, i.e. > 140 chars?
	Truncated bool
	Entities  struct {
		URLs []TweetEntityUrl `json:"urls"`
	}
	// Source is a Twitter client name used to post the tweet.
	Source    string
	CreatedAt string/*time.Time*/ `json:"created_at"`
	User      TweetUser

	AppName string
}

func (t *Tweet) statusUrl() string {
	return "https://twitter.com/" + t.User.ScreenName + "/statuses/" + t.ID
}

type TweetUser struct {
	ID         string `json:"id_str"`
	ScreenName string `json:"screen_name"`
	Photo      string `json:"profile_image_url_https"`
	// Country code, e.g. "en", "it" specified by the user.
	Lang string
}

// TweetEntityUrl is Tweet.Entities.URLs slice item.
type TweetEntityUrl struct {
	URL         string `json:"url"`
	ExpandedURL string `json:"expanded_url"`
	DisplayURL  string `json:"display_url"`
}

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
