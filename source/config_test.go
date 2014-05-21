package source

import (
	"github.com/fiorix/go-redis/redis"
)

// dummyConfigSource is a fake implementation of Store.
type dummyConfigSource struct {
	listAppNames   func() ([]string, error)
	getApp         func(string) (*Application, error)
	listTwitterIDs func(string) ([]string, error)
	subscribe      func(topic string, msg chan redis.PubSubMessage, stop chan bool) error
}

func (s *dummyConfigSource) ListAppNames() ([]string, error) {
	return s.listAppNames()
}

func (s *dummyConfigSource) GetApp(name string) (*Application, error) {
	return s.getApp(name)
}

func (s *dummyConfigSource) ListTwitterIDs(name string) ([]string, error) {
	return s.listTwitterIDs(name)
}

func (s *dummyConfigSource) Subscribe(topic string, msg chan redis.PubSubMessage, stop chan bool) error {
	return s.subscribe(topic, msg, stop)
}