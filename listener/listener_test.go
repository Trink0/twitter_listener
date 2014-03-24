package listener

import "testing"

type dummyListener struct {
	name        string
	startCalled bool
}

func (l *dummyListener) Start(c chan int) {
	l.startCalled = true
	c <- 1
}

func TestStartAllEmpty(t *testing.T) {
	store := &dummyAppStore{}
	store.listAppNames = func() ([]string, error) {
		return []string{}, nil
	}

	err := StartAll(store)
	if err != nil {
		t.Fatal(err)
	}
}

func TestStartAll(t *testing.T) {
	store := &dummyAppStore{}
	store.listAppNames = func() ([]string, error) {
		return []string{"chumhum", "xpeppers"}, nil
	}
	store.getApp = func(name string) (*Application, error) {
		return &Application{Name: name}, nil
	}

	listeners := make([]*dummyListener, 0)
	listenerFactory = func(app *Application) Listener {
		dummy := &dummyListener{name: app.Name}
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
	}
}
