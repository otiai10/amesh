package amesh

import (
	"fmt"

	"github.com/mrjones/oauth"
)

// Notifier ...
type Notifier interface {
	Notify(string) error
}

// SlackNotifier ...
type SlackNotifier struct {
}

// Notify ...
func (notifier *SlackNotifier) Notify(msg string) error {
	return nil
}

// NewTwitterNotifier ...
func NewTwitterNotifier(consumerKey, consumerSecret, axToken, axTokenSecret string) *TwitterNotifier {
	return &TwitterNotifier{
		Consumer: oauth.NewConsumer(
			consumerKey,
			consumerSecret,
			oauth.ServiceProvider{
				AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
				RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
				AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
			},
		),
		AccessToken: &oauth.AccessToken{
			Token:          axToken,
			Secret:         axTokenSecret,
			AdditionalData: make(map[string]string),
		},
	}
}

// TwitterNotifier ...
type TwitterNotifier struct {
	Consumer    *oauth.Consumer
	AccessToken *oauth.AccessToken
}

// Notify ...
func (notifier *TwitterNotifier) Notify(msg string) error {
	if notifier.Consumer == nil || notifier.AccessToken == nil {
		return fmt.Errorf("cannot use twitter API without auth")
	}
	res, err := notifier.Consumer.Post(
		"https://api.twitter.com/1.1/statuses/update.json",
		map[string]string{
			"status": msg,
		},
		notifier.AccessToken,
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || 300 <= res.StatusCode {
		return fmt.Errorf(res.Status)
	}

	return nil
}
