package amesh

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/mrjones/oauth"
)

// Notifier ...
type Notifier interface {
	Notify(string) error
}

// SlackNotifier ...
type SlackNotifier struct {
	Token   string
	Channel string
}

// NewSlackNotifier ...
func NewSlackNotifier(token, channel string) *SlackNotifier {
	return &SlackNotifier{
		Token:   token,
		Channel: channel,
	}
}

// Notify ...
func (notifier *SlackNotifier) Notify(msg string) error {

	query := url.Values{}
	query.Add("token", notifier.Token)
	query.Add("channel", notifier.Channel)
	query.Add("as_user", "true")
	query.Add("text", msg)

	res, err := http.Get("https://slack.com/api/chat.postMessage?" + query.Encode())
	if err != nil {
		return err
	}
	defer res.Body.Close()

	resp := struct {
		OK    bool   `json:"ok"`
		Error string `json:"error"`
	}{}

	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return err
	}

	if !resp.OK {
		return fmt.Errorf("slack said: %s", resp.Error)
	}

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
