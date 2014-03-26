package listener

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestLoopOneTweet(t *testing.T) {
	tweet := &Tweet{ID: "123456", Text: "tweet text", User: TweetUser{ID: "011100"}}

	b, _ := json.Marshal(tweet)
	var buffer bytes.Buffer
	buffer.Write(b)
	buffer.Write(endOfTweet)
	fakeStream := bytes.NewReader(buffer.Bytes())

	q := make(chan *Tweet, 1)
	streamer := &httpStreamer{app: &Application{}, users: []string{"011100"}, queue: q}
	streamer.loop(fakeStream)

	select {
	case tweeted := <-q:
		if !reflect.DeepEqual(tweeted, tweet) {
			t.Errorf("Have %+v expected %+v", tweeted, tweet)
		}
	case <-time.Tick(time.Millisecond):
		t.Fatal("Exepected message but received nothing")
	}

}
func TestLoopGarbageTweet(t *testing.T) {
	fakeStream := bytes.NewReader([]byte("garbage tweet"))

	q := make(chan *Tweet, 1)
	streamer := &httpStreamer{app: &Application{}, users: []string{"011100"}, queue: q}
	streamer.loop(fakeStream)

	select {
	case tweeted := <-q:
		t.Errorf("Didn't expect %+v", tweeted)
	case <-time.Tick(time.Millisecond):
		// test pass
	}
}
func TestLoopTweetWithNotFollowedUser(t *testing.T) {
	tweet := &Tweet{ID: "123456", Text: "tweet text", User: TweetUser{ID: "4444444"}}

	b, _ := json.Marshal(tweet)
	var buffer bytes.Buffer
	buffer.Write(b)
	buffer.Write(endOfTweet)
	fakeStream := bytes.NewReader(buffer.Bytes())

	q := make(chan *Tweet, 1)
	streamer := &httpStreamer{app: &Application{}, users: []string{"011100"}, queue: q}
	streamer.loop(fakeStream)

	select {
	case tweeted := <-q:
		t.Errorf("Didn't expect %+v", tweeted)
	case <-time.Tick(time.Millisecond):
		// test pass
	}
}
