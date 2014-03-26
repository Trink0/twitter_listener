package listener

import (
	"reflect"
	"testing"
)

type dummyListener struct {
	name        string
	users       []string
	startCalled bool
}

func (l *dummyListener) Start(c chan int) {
	l.startCalled = true
	c <- 1
}

type dummyQueue struct {
	qc chan *Tweet
}

func (d *dummyQueue) Start(c chan *Tweet) {
	d.qc = c
}

func TestStartOne(t *testing.T) {
	store := &dummyStore{}
	store.getApp = func(name string) (*Application, error) {
		if name != "chumhum" {
			t.Errorf("Have app name %q, want chumhum", name)
		}
		return &Application{Name: name}, nil
	}
	store.listTwitterIDs = func(name string) ([]string, error) {
		return []string{"15170239", "1585341620"}, nil
	}

	dummy := &dummyListener{}
	listenerFactory = func(app *Application, userIds []string, qc chan *Tweet) Listener {
		return dummy
	}

	queue := &dummyQueue{}
	if err := StartOne("chumhum", store, queue); err != nil {
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
	store := &dummyStore{}
	store.listAppNames = func() ([]string, error) {
		return []string{}, nil
	}

	queue := &dummyQueue{}
	err := StartAll(store, queue)
	if err != nil {
		t.Fatal(err)
	}
	if queue.qc != nil {
		t.Error("Queue channel should be nil")
	}
}

func TestStartAll(t *testing.T) {
	twitterIDs := []string{"15170239", "1585341620"}
	store := &dummyStore{}
	store.listAppNames = func() ([]string, error) {
		return []string{"chumhum", "xpeppers"}, nil
	}
	store.getApp = func(name string) (*Application, error) {
		return &Application{Name: name}, nil
	}
	store.listTwitterIDs = func(name string) ([]string, error) {
		return twitterIDs, nil
	}

	listeners := make([]*dummyListener, 0)
	listenerFactory = func(app *Application, userIds []string, qc chan *Tweet) Listener {
		dummy := &dummyListener{name: app.Name, users: userIds}
		listeners = append(listeners, dummy)
		return dummy
	}

	queue := &dummyQueue{}
	err := StartAll(store, queue)
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
