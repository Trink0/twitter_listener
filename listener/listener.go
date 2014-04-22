package listener

import (
	"log"
)

// APP_TOPIC is Redis channel name
const APP_TOPIC = "appchanges"

// Listener is a Twitter Streaming API client.
type Listener interface {
	// Start initiates a connection and reads from it indefinitely in a goroutine
	Start()
	Stop()
	Restart()
	IsActive() bool
	Name() string
	UpdateUsers(userIds []string)
	UpdateApp(app *Application)
}

// NewListener creates a new listener with credentials provided by the app.
func NewListener(app *Application, userIds []string, qc chan *Tweet, errc chan int) Listener {
	return listenerFactory(app, userIds, qc, errc)
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
	listener := NewListener(app, userIDs, qc, errc)
	listener.Start()

	aw := NewAppWatcher(APP_TOPIC, []Listener{listener}, s)
	return aw.Watch(qc, errc)
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

	errc := make(chan int, len(appNames))
	qc := make(chan *Tweet, len(appNames)*100)
	queue.Start(qc)

	allListeners := make([]Listener, 0, len(appNames))
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

		listener := NewListener(storedApp, userIDs, qc, errc)
		allListeners = append(allListeners, listener)
		go listener.Start()
	}
	aw := NewAppWatcher(APP_TOPIC, allListeners, s)
	return aw.Watch(qc, errc)
}

// listenerFactory is by NewListener to create a new listener struct.
// Meant overwritten in tests.
var listenerFactory = func(a *Application, userIds []string, qc chan *Tweet, errc chan int) Listener {
	return &httpStreamer{app: a, users: userIds, queue: qc, errc: errc}
}
