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

// Each User belongs to an Application.
type User struct {
	Id       string
	Username string
	Metadata map[string]string
}

// Store is the client for a storage backend where all apps data are located at.
type Store interface {
	// ListAppNames fetches a list of names of all currently registered apps.
	ListAppNames() ([]string, error)
	// GetApp fetches a single application data identified by its name.
	GetApp(name string) (*Application, error)
	// ListTwitterIDs returns Twitter IDs of all users that belong to an app
	// identified by the given name.
	ListTwitterIDs(name string) ([]string, error)
}
