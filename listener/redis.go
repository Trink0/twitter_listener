// +build !test

package listener

import (
	"encoding/json"
	"fmt"

	"github.com/fiorix/go-redis/redis"
)

func (s *AppStore) ListAppNames() ([]string, error) {
	rc := redis.New(s.connUrl)
	return rc.Keys("*")
}

func (s *AppStore) GetApp(name string) (*Application, error) {
	rc := redis.New(s.connUrl)

	jsonApp, err := rc.Get(name)
	if err != nil {
		return nil, err
	}
	if jsonApp == "" {
		return nil, fmt.Errorf("App %q not found", name)
	}

	storedApp := &Application{}
	if err := json.Unmarshal([]byte(jsonApp), storedApp); err != nil {
		return nil, err
	}
	return storedApp, nil
}
