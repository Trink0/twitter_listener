package listener

import (
	"log"
	"time"
)

type Listener struct {
	app *Application
}

func (l *Listener) start(c chan int) {
	log.Printf("Starting Listner: %s", l.app.Name)
	time.Sleep(time.Second * 5)
	c <- 1
}

func StartAll(s *AppStore) error {
	appNames, err := s.ListAppNames()
	if err != nil {
		return err
	}

	c := make(chan int)
	for _, name := range appNames {
		storedApp, getErr := s.GetApp(name)
		if getErr != nil {
			log.Printf("ERROR fetching app %q: %v", name, getErr)
			continue
		}

		listener := &Listener{app: storedApp}
		go listener.start(c)
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
