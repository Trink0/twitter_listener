package listener

import (
	"encoding/json"
	"fmt"

	"github.com/fiorix/go-redis/redis"
)

type redisAppStore struct {
	connUrl string
}

func NewAppStore(connUrl string) AppStore {
	return &redisAppStore{connUrl}
}

func (s *redisAppStore) ListAppNames() ([]string, error) {
	rc := redis.New(s.connUrl)
	return rc.Keys("*")
}

func (s *redisAppStore) GetApp(name string) (app *Application, getErr error) {
	rc := redis.New(s.connUrl)

	jsonApp, err := rc.Get(name)
	if err != nil {
		return nil, err
	}
	if jsonApp == "" {
		return nil, fmt.Errorf("App %q not found", name)
	}

	app = &Application{}
	getErr = json.Unmarshal([]byte(jsonApp), app)
	return
}
