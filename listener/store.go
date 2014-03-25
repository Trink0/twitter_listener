package listener

// Application is a customer/merchant
type Application struct {
	ApiKey           string `json:"apiKey"`
	Name             string `json:"name"`
	TwConsumerKey    string `json:"twitterConsumerKey"`
	TwConsumerSecret string `json:"twitterConsumerSecret"`
	TwAccessToken    string `json:"twitterAccessToken"`
	TwTokenSecret    string `json:"twitterTokenSecret"`
}

// Store is the client for a storage backend where all apps data are located at.
type Store interface {
	// ListAppNames fetches a list of names of all currently registered apps.
	ListAppNames() ([]string, error)
	// GetApp fetches a single application data identified by its name.
	GetApp(name string) (*Application, error)
	// GetApp fetches a single application data identified by its name.
	ListAppUserIds(name string) ([]string, error)
}
