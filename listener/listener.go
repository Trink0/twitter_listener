package listener

import (
	"log"
)

// Listener is a Twitter Streaming API client.
type Listener interface {
	// Start initiates a connection and reads from it indefinitely in a goroutine
	// c is a control channel used by the listener to communicate an exit error.
	Start(c chan int)
	Restart(c chan int)
	Name() string
}

// NewListener creates a new listener with credentials provided by the app.
func NewListener(app *Application, userIds []string, qc chan *Tweet) Listener {
	return listenerFactory(app, userIds, qc)
}

// StartOne creates and starts one listener for the specified application.
func StartOne(appName string, s Store, queue Queue) error {
	app, err := s.GetApp(appName)
	if err != nil {
		return err
	}
	userIDs, userErr := s.ListTwitterIDs(appName)
	if userErr != nil {
		return userErr
	}
	if len(userIDs) == 0 {
		log.Printf("No users found for app %q. Exiting.", appName)
		return nil
	}

	qc := make(chan *Tweet, 100)
	queue.Start(qc)

	errc := make(chan int, 1)
	listener := NewListener(app, userIDs, qc)
	listener.Start(errc)

	aw := NewAppWatcher("topic", []Listener{listener}, s)
	return aw.Watch(qc, errc)

	// waitAll(c, 1)
	// return nil
}

// StartAll creates and starts a new listener for each application
// registered in the store.
func StartAll(s Store, queue Queue) error {
	appNames, err := s.ListAppNames()
	if err != nil {
		return err
	}
	if len(appNames) == 0 {
		log.Print("No applications found. Exiting.")
		return nil
	}

	c := make(chan int, len(appNames))
	qc := make(chan *Tweet, len(appNames)*100)
	queue.Start(qc)

	for _, name := range appNames {
		storedApp, getErr := s.GetApp(name)
		if getErr != nil {
			log.Printf("ERROR fetching app %q: %v", name, getErr)
			continue
		}

		userIDs, userErr := s.ListTwitterIDs(name)
		if userErr != nil {
			return userErr
		}

		listener := NewListener(storedApp, userIDs, qc)
		listener.Start(c)
	}

	waitAll(c, len(appNames))
	return nil
}

// waitAll reads from the
func waitAll(c chan int, n int) {
	count := 0
	for {
		status := <-c
		log.Printf("Listener exited with status %d", status)
		if count += 1; count == n {
			log.Printf("Quit application")
			break
		}
	}
}

// listenerFactory is by NewListener to create a new listener struct.
// Meant overwritten in tests.
var listenerFactory = func(a *Application, userIds []string, qc chan *Tweet) Listener {
	return &httpStreamer{app: a, users: userIds, queue: qc}
}
