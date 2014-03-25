package listener

import (
	"log"
)

// Listener is a Twitter Streaming API client.
type Listener interface {
	// Start initiates a connection and reads from it indefinitely in a goroutine
	// c is a control channel used by the listener to communicate an exit error.
	Start(c chan int)
}

// NewListener creates a new listener with credentials provided by the app.
func NewListener(app *Application, userIds []string) Listener {
	return listenerFactory(app, userIds)
}

// httpStreamer is a default implementation of the Listener over HTTP
// using Public Stream API.
type httpStreamer struct {
	app   *Application
	users []string
}

func (s *httpStreamer) Start(c chan int) {
	log.Printf("Starting listener %q", s.app.Name)
	go s.stream(c)
}

// StartOne creates and starts one listener for the specified application.
func StartOne(s AppStore, u AppStore, appName string) error {
	app, err := s.GetApp(appName)
	if err != nil {
		return err
	}
	userIds, userErr := u.ListAppUserIds(appName)
	if userErr != nil {
		return userErr
	}

	c := make(chan int, 1)
	listener := NewListener(app, userIds)
	listener.Start(c)

	waitAll(c, 1)
	return nil
}

// StartAll creates and starts a new listener for each application
// registered in the store.
func StartAll(s AppStore, u AppStore) error {
	appNames, err := s.ListAppNames()
	if err != nil {
		return err
	}
	if len(appNames) == 0 {
		log.Print("No applications found. Exiting.")
		return nil
	}

	c := make(chan int, len(appNames))
	for _, name := range appNames {
		storedApp, getErr := s.GetApp(name)
		if getErr != nil {
			log.Printf("ERROR fetching app %q: %v", name, getErr)
			continue
		}

		userIds, userErr := u.ListAppUserIds(name)
		if userErr != nil {
			return userErr
		}

		listener := NewListener(storedApp, userIds)
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
var listenerFactory = func(a *Application, userIds []string) Listener {
	return &httpStreamer{app: a, users: userIds}
}
