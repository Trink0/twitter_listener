package listener

import (
	"encoding/json"
	"fmt"

	"github.com/fiorix/go-redis/redis"
)

// redisStore is Redis-based implementation of Store
type redisStore struct {
	connUrl string
}

// NewStore creates a new instance of Redis-based Store.
// Connection URL should also specify db number, e.g. "127.0.0.1:6379 db=1".
func NewStore(connUrl string) Store {
	return &redisStore{connUrl}
}

func (s *redisStore) ListAppNames() ([]string, error) {
	rc := redis.New(s.connUrl)
	return rc.Keys("*")
}

func (s *redisStore) GetApp(name string) (app *Application, getErr error) {
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

func (s *redisStore) ListAppUserIds(name string) ([]string, error) {
	rc := redis.New(s.connUrl)

	userIds, err := rc.SMembers("customer:" + name)
	if err != nil {
		return nil, err
	}
	if len(userIds) == 0 {
		return nil, fmt.Errorf("User for App %q not found", name)
	}

	for i := 0; i < len(userIds); i++ {
		//jsonUser, err := rc.Get(userIds[i])
	}
	return []string{"15170239", "1585341620"}, nil
}
