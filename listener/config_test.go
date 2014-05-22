package listener

import (
	"github.com/Trink0/twitter_listener/source"
	// "github.com/fiorix/go-redis/redis"
)

// dummyConfigSource is a fake implementation of Store.
type dummyConfigSource struct {
	listAppNames   func() ([]string, error)
	getApp         func(string) (*source.Application, error)
	listTwitterIDs func(string) ([]string, error)
	subscribe      func(topic string, msg chan source.Notification, stop chan bool) error
}

func (s *dummyConfigSource) ListAppNames() ([]string, error) {
	return s.listAppNames()
}

func (s *dummyConfigSource) GetApp(name string) (*source.Application, error) {
	return s.getApp(name)
}

func (s *dummyConfigSource) ListTwitterIDs(name string) ([]string, error) {
	return s.listTwitterIDs(name)
}

func (s *dummyConfigSource) Subscribe(topic string, msg chan source.Notification, stop chan bool) error {
	return s.subscribe(topic, msg, stop)
}
