package listener

import (
	"encoding/json"
	"testing"
	"time"
)

const (
	TEST_TEXT_TWEET = "This is my test Tweet"
)

func TestQueueLoop(t *testing.T) {
	var isCalled bool
	kq := kestrelQueue{
		putFn: func(payload []byte) {
			activity := &Activity{}
			if err := json.Unmarshal(payload, activity); err != nil {
				t.Fail()
			}
			if activity.Object.Text != TEST_TEXT_TWEET {
				t.Fatalf("Activity text is not correct. Got %s, Expected %s", activity.Object.Text, TEST_TEXT_TWEET)
			}
			isCalled = true
		},
	}

	qc := make(chan *Tweet, 1)
	go kq.loop(qc)

	tweet := &Tweet{
		Text: TEST_TEXT_TWEET,
	}
	qc <- tweet
	time.Sleep(time.Millisecond * 10)
	if !isCalled {
		t.Fatal("Loop func not called")
	}
}
