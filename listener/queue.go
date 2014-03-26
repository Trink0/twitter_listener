package listener

import (
// "time"
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
	return &kestrelQueue{endPoint}
}

type kestrelQueue struct {
	endPoint string
}

func (k *kestrelQueue) Start(qc chan *Tweet) {
	//...
}
