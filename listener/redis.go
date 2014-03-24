package listener

import (
	"encoding/json"
	"fmt"

	"github.com/fiorix/go-redis/redis"
)

// redisAppStore is Redis-based implementation of AppStore
type redisAppStore struct {
	connUrl string
}

// NewAppStore creates a new instance of Redis-based AppStore.
// Connection URL should also specify db number, e.g. "127.0.0.1:6379 db=1".
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
