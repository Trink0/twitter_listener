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

type AppStore struct {
	connUrl string
}

func NewAppStore(connUrl string) *AppStore {
	return &AppStore{connUrl}
}
