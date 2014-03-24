package listener

import (
	"log"
	"time"
)

type Listener interface {
	Start(c chan int)
}

type httpListener struct {
	app *Application
}

func (l *httpListener) Start(c chan int) {
	log.Printf("Starting Listner: %s", l.app.Name)
	go l.stream(c)
}

func (l *httpListener) stream(c chan int) {
	time.Sleep(time.Second * 5)
	c <- 1
}

func StartAll(s AppStore) error {
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

		listener := listenerFactory(storedApp)
		listener.Start(c)
	}

	count := 0
	for {
		status := <-c
		log.Printf("Listener Exiting with status: %d", status)
		if count += 1; count == len(appNames) {
			log.Printf("Exiting application")
			break
		}
	}

	return nil
}

var listenerFactory = func(a *Application) Listener {
	return &httpListener{app: a}
}
