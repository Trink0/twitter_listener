package source

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/fiorix/go-redis/redis"
)

// NewConfigSource creates a new instance of Redis-based ConfigSource.
// The URL should not spcify any db number, e.g. "127.0.0.1:6379".
// Internally, connection URL is then constructed from dbURL and a database number.
func NewConfigSource(dbURL string, appDB, userDB int) ConfigSource {
	return &redisConfigSource{dbURL, appDB, userDB}
}

// redisConfigSource is Redis-based implementation of ConfigSource
type redisConfigSource struct {
	dbURL  string
	appDB  int
	userDB int
}

func (s *redisConfigSource) ListAppNames() ([]string, error) {
	return s.newClient(s.appDB).Keys("*")
}

func (s *redisConfigSource) GetApp(name string) (app *Application, getErr error) {
	jsonApp, err := s.newClient(s.appDB).Get(name)
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

func (s *redisConfigSource) ListTwitterIDs(name string) ([]string, error) {
	users := s.newClient(s.userDB)

	userIDs, err := users.SMembers("customer:" + name)
	if err != nil {
		return nil, err
	}

	twitterIDs := make([]string, 0, len(userIDs))
	for _, userID := range userIDs {
		jsonUser, getErr := users.Get(userID)
		if getErr != nil {
			log.Printf("ERROR getting user %q of app %q", userID, name)
			continue
		}

		user := &User{}
		if json.Unmarshal([]byte(jsonUser), user) != nil {
			log.Printf("ERROR parsing JSON user data of %q (%s)", userID, name)
			continue
		}
		if user.Metadata == nil {
			log.Printf("WARNING: no metadata for user %q (%s)", userID, name)
			continue
		}

		// TODO: find a better way to get twitter ID,
		// e.g. store it in user.Metadata["twitter.user.id"]
		if _, ok := user.Metadata["twitter.user.screenName"]; ok {
			twitterIDs = append(twitterIDs, user.Username)
		}
	}

	return twitterIDs, nil
}

func (s *redisConfigSource) Subscribe(topic string, msg chan redis.PubSubMessage, stop chan bool) error {
	return s.newClient(s.appDB).Subscribe(topic, msg, stop)
}

func (s *redisConfigSource) connectionURL(db int) string {
	return fmt.Sprintf("%s db=%d", s.dbURL, db)
}

func (s *redisConfigSource) newClient(db int) *redis.Client {
	return redis.New(s.connectionURL(db))
}
