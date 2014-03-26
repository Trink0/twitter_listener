package listener

import (
	"encoding/json"
	"testing"
)

func TestQueueLoop(t *testing.T) {

	qk := kestrelQueue{putFn: func(payload []byte) {
		activity := &Activity{}
		if err := json.Unmarshal(payload, activity); err != nil {
			t.Fail()
		}
	},
	}

	qc := make(chan *Tweet, 1)
	go qk.loop(qc)
	qc <- &Tweet{}
}
