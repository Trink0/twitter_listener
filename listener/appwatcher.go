package listener

import (
	"log"

	"github.com/fiorix/go-redis/redis"
)

type AppWatcher struct {
	topic     string
	listeners []Listener
	store     Store
}

func NewAppWatcher(topic string, listeners []Listener, store Store) *AppWatcher {
	return &AppWatcher{topic, listeners, store}
}

func (a *AppWatcher) Watch(qc chan *Tweet, errc chan int) error {
	msgc := make(chan redis.PubSubMessage)
	stopc := make(chan bool, 1)
	if err := a.store.Subscribe(a.topic, msgc, stopc); err != nil {
		log.Printf("Error subscribing to channel: %q", err)
		return err
	}

	go a.loop(msgc, qc, errc)
	a.waitAll(errc)

	stopc <- true
	return nil
}

func (a *AppWatcher) loop(msgc chan redis.PubSubMessage, qc chan *Tweet, errc chan int) {
	for {
		msg := <-msgc
		if msg.Error != nil {
			log.Printf("Message channel error %s", msg.Error)
			continue
		}
		appName := msg.Value
		log.Printf("Received message for app: %s", appName)
		var listener Listener
		for _, l := range a.listeners {
			if l.Name() == appName {
				listener = l
				break
			}
		}

		app, err := a.store.GetApp(appName)
		if err != nil {
			log.Printf("Message channel error %s", err)
			continue
		}
		userIDs, userErr := a.store.ListTwitterIDs(appName)
		if userErr != nil {
			log.Printf("Message channel error %s", userErr)
			continue
		}
		if len(userIDs) == 0 {
			log.Printf("No users found for app %q. Exiting.", appName)
			continue
		}
		if listener == nil {
			listener = NewListener(app, userIDs, qc, errc)
			a.listeners = append(a.listeners, listener)
			log.Printf("Staring new listener for application %q.", appName)
			go listener.Start()
		} else {
			listener.UpdateApp(app)
			listener.UpdateUsers(userIDs)
			log.Printf("Restarting listener for application %q.", appName)
			go listener.Restart()
		}

	}

}

// waitAll reads from the
func (a *AppWatcher) waitAll(errc chan int) {
	for {
		status := <-errc
		log.Printf("Listener exited with status %d", status)
		if !a.hasActiveListeners() {
			log.Printf("Quit application")
			return
		}
	}
}

func (a *AppWatcher) hasActiveListeners() bool {
	for _, l := range a.listeners {
		if l.IsActive() {
			return true
		}
	}
	return false
}
