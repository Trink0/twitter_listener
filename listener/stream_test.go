package listener

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/Trink0/twitter_listener/source"
)

const (
	DEFAULT_APP_NAME = "TestAppName"
	DEFAULT_USER_ID  = "011100"
	TEST_FILTER_URL  = "Test Filter Url"
)

func TestLoopOneTweet(t *testing.T) {
	filterURL = TEST_FILTER_URL
	tweet := &Tweet{ID: "123456", Text: "tweet text", User: TweetUser{ID: DEFAULT_USER_ID}}

	b, _ := json.Marshal(tweet)
	var buffer bytes.Buffer
	buffer.Write(b)
	buffer.Write(endOfTweet)
	fakeStream := bytes.NewReader(buffer.Bytes())

	q := make(chan *Tweet, 1)
	streamer := &httpStreamer{
		app:   &source.Application{Name: DEFAULT_APP_NAME},
		users: []string{DEFAULT_USER_ID},
		queue: q,
		stopc: make(chan bool, 1),
		errc:  make(chan int, 1),
	}
	streamer.loop(fakeStream)

	select {
	case tweeted := <-q:
		tweet.AppName = DEFAULT_APP_NAME
		if !reflect.DeepEqual(tweeted, tweet) {
			t.Errorf("Have %+v expected %+v", tweeted, tweet)
		}
	case <-time.Tick(time.Millisecond * 5):
		t.Fatal("Exepected message but received nothing")
	}

}
func TestLoopGarbageTweet(t *testing.T) {
	filterURL = TEST_FILTER_URL
	fakeStream := bytes.NewReader([]byte("garbage tweet"))

	q := make(chan *Tweet, 1)
	streamer := &httpStreamer{
		app:   &source.Application{},
		users: []string{DEFAULT_USER_ID},
		queue: q,
		stopc: make(chan bool, 1),
		errc:  make(chan int, 1),
	}
	streamer.loop(fakeStream)

	select {
	case tweeted := <-q:
		t.Errorf("Didn't expect %+v", tweeted)
	case <-time.Tick(time.Millisecond * 5):
		// test pass
	}
}
func TestLoopTweetWithNotFollowedUser(t *testing.T) {
	filterURL = TEST_FILTER_URL
	tweet := &Tweet{ID: "123456", Text: "tweet text", User: TweetUser{ID: "4444444"}}

	b, _ := json.Marshal(tweet)
	var buffer bytes.Buffer
	buffer.Write(b)
	buffer.Write(endOfTweet)
	fakeStream := bytes.NewReader(buffer.Bytes())

	q := make(chan *Tweet, 1)
	streamer := &httpStreamer{
		app:   &source.Application{},
		users: []string{DEFAULT_USER_ID},
		queue: q,
		stopc: make(chan bool, 1),
		errc:  make(chan int, 1),
	}
	streamer.loop(fakeStream)

	select {
	case tweeted := <-q:
		t.Errorf("Didn't expect %+v", tweeted)
	case <-time.Tick(time.Millisecond * 5):
		// test pass
	}
}

func TestUpdateUsers(t *testing.T) {
	q := make(chan *Tweet, 1)
	streamer := &httpStreamer{
		app:   &source.Application{},
		users: []string{DEFAULT_USER_ID},
		queue: q,
		stopc: make(chan bool, 1),
		errc:  make(chan int, 1),
	}
	users_upd := []string{"aaa-bbb", "ccc-ddd"}
	streamer.UpdateUsers(users_upd)

	if len(streamer.users) != 2 {
		t.Fatal("Expected 2 elements in users")
	}
	if streamer.users[0] != "aaa-bbb" || streamer.users[1] != "ccc-ddd" {
		t.Fatalf("unexpected values: %+v %+v", streamer.users[0], streamer.users[1])
	}
}
