package listener

import (
	"encoding/json"
	// "errors"
	"log"
	// "strconv"
	// "strings"

	"github.com/alindeman/go-kestrel"
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
}

type Queue interface {
	Start(qc chan *Tweet)
}

func NewQueue(endPoint string) Queue {
	// ep := strings.SplitN(endPoint, ":", 2)
	// if ep[0] == "" {
	// 	return nil, errors.New("Invalid queue host")
	// }
	// port := 2229
	// if len(ep) == 2 {
	// 	port, portErr := strconv.Atoi(ep[1])
	// 	if portErr != nil {
	// 		return nil, portErr
	// 	}
	// }
	// return &kestrelQueue{ep[0], port}, nil
	return &kestrelQueue{"127.0.0.1", 2229}
}

type kestrelQueue struct {
	host string
	port int
}

func (k *kestrelQueue) Start(qc chan *Tweet) {
	log.Printf("Connecting to Kestrel on %s:%d", k.host, k.port)
	client := kestrel.NewClient(k.host, k.port)
	go k.loop(client, qc)
}

func (k *kestrelQueue) loop(client *kestrel.Client, qc chan *Tweet) {
	for {
		tweet := <-qc
		log.Printf("ENQUEUE: %+v", tweet)
		payload, err := json.Marshal(tweet)
		if err != nil {
			log.Printf("ERROR encoding activity: %v", err)
			continue
		}

		n, putErr := client.Put("social-web-activities", [][]byte{payload})
		if putErr != nil {
			log.Printf("ERROR queue: %v", putErr)
			continue
		}
		if n < 1 {
			log.Println("ERROR: unable to put")
		}
	}
}
