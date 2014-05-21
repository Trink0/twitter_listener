package listener

import (
	"reflect"
	"testing"

	"github.com/Trink0/twitter_listener/source"
	"github.com/fiorix/go-redis/redis"
)

type dummyListener struct {
	name        string
	users       []string
	startCalled bool
	stopCalled  bool
	c           chan int
}

func (l *dummyListener) Start() {
	l.startCalled = true
	l.c <- 1
}

func (l *dummyListener) Stop() {
	l.stopCalled = true
}

func (l *dummyListener) Restart() {
	l.Stop()
	l.Start()
}

func (l *dummyListener) IsActive() bool {
	return false
}

func (l *dummyListener) Name() string {
	return l.name
}

func (l *dummyListener) UpdateApp(app *source.Application) {}

func (l *dummyListener) UpdateUsers(users []string) {}

type dummyQueue struct {
	qc chan *Tweet
}

func (d *dummyQueue) Start(c chan *Tweet) {
	d.qc = c
}

func TestStartOne(t *testing.T) {
	config := &dummyConfigSource{}
	config.getApp = func(name string) (*source.Application, error) {
		if name != "chumhum" {
			t.Errorf("Have app name %q, want chumhum", name)
		}
		return &source.Application{Name: name}, nil
	}
	config.listTwitterIDs = func(name string) ([]string, error) {
		return []string{"15170239", "1585341620"}, nil
	}
	config.subscribe = func(topic string, msg chan redis.PubSubMessage, stop chan bool) error {
		return nil
	}

	dummy := &dummyListener{}
	listenerFactory = func(app *source.Application, userIds []string, qc chan *Tweet, c chan int) Listener {
		dummy.c = c
		return dummy
	}

	queue := &dummyQueue{}
	if err := StartOne("chumhum", config, queue); err != nil {
		t.Fatal(err)
	}
	if !dummy.startCalled {
		t.Error("Didn't call listener.Start()")
	}
	if queue.qc == nil {
		t.Error("Expected queue channel is nil")
	}
}

func TestStartAllEmpty(t *testing.T) {
	config := &dummyConfigSource{}
	config.listAppNames = func() ([]string, error) {
		return []string{}, nil
	}

	queue := &dummyQueue{}
	err := StartAll(config, queue)
	if err != nil {
		t.Fatal(err)
	}
	if queue.qc != nil {
		t.Error("Queue channel should be nil")
	}
}

func TestStartAll(t *testing.T) {
	twitterIDs := []string{"15170239", "1585341620"}
	config := &dummyConfigSource{}
	config.listAppNames = func() ([]string, error) {
		return []string{"chumhum", "xpeppers"}, nil
	}
	config.getApp = func(name string) (*source.Application, error) {
		return &source.Application{Name: name}, nil
	}
	config.listTwitterIDs = func(name string) ([]string, error) {
		return twitterIDs, nil
	}
	config.subscribe = func(topic string, msg chan redis.PubSubMessage, stop chan bool) error {
		return nil
	}

	listeners := make([]*dummyListener, 0)
	listenerFactory = func(app *source.Application, userIds []string, qc chan *Tweet, c chan int) Listener {
		dummy := &dummyListener{name: app.Name, users: userIds, c: c}
		listeners = append(listeners, dummy)
		return dummy
	}

	queue := &dummyQueue{}
	err := StartAll(config, queue)
	if err != nil {
		t.Fatal(err)
	}

	if listLen := len(listeners); listLen != 2 {
		t.Fatalf("Have %d listeners, want 2", listLen)
	}
	if queue.qc == nil {
		t.Error("Expected queue channel is nil")
	}

	for _, l := range listeners {
		if !l.startCalled {
			t.Errorf("Didn't start listener %q", l.name)
		}
		if l.name != "xpeppers" && l.name != "chumhum" {
			t.Errorf("Have app name %q, want either xpeppers or chumhum", l.name)
		}
		if !reflect.DeepEqual(l.users, twitterIDs) {
			t.Errorf("Have twitter IDs: %v, want %v", l.users, twitterIDs)
		}
	}
}
