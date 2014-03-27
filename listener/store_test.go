package listener

import (
	"github.com/fiorix/go-redis/redis"
)

// dummyStore is a fake implementation of Store.
type dummyStore struct {
	listAppNames   func() ([]string, error)
	getApp         func(string) (*Application, error)
	listTwitterIDs func(string) ([]string, error)
	subscribe      func(topic string, msg chan redis.PubSubMessage, stop chan bool) error
}

func (s *dummyStore) ListAppNames() ([]string, error) {
	return s.listAppNames()
}

func (s *dummyStore) GetApp(name string) (*Application, error) {
	return s.getApp(name)
}

func (s *dummyStore) ListTwitterIDs(name string) ([]string, error) {
	return s.listTwitterIDs(name)
}

func (s *dummyStore) Subscribe(topic string, msg chan redis.PubSubMessage, stop chan bool) error {
	return s.subscribe(topic, msg, stop)
}
