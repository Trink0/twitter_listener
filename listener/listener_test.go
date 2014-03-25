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
	listenerFactory = func(app *Application, userIds []string) Listener {
		return dummy
	}

	if err := StartOne(store, "chumhum"); err != nil {
		t.Fatal(err)
	}
	if !dummy.startCalled {
		t.Error("Didn't call listener.Start()")
	}
}

func TestStartAllEmpty(t *testing.T) {
	store := &dummyStore{}
	store.listAppNames = func() ([]string, error) {
		return []string{}, nil
	}

	err := StartAll(store)
	if err != nil {
		t.Fatal(err)
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
	listenerFactory = func(app *Application, userIds []string) Listener {
		dummy := &dummyListener{name: app.Name, users: userIds}
		listeners = append(listeners, dummy)
		return dummy
	}

	err := StartAll(store)
	if err != nil {
		t.Fatal(err)
	}

	if listLen := len(listeners); listLen != 2 {
		t.Fatalf("Have %d listeners, want 2", listLen)
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
