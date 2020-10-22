package commands

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/otiai10/amesh/bot/slack"
	"github.com/otiai10/goapis/google"
)

// GoogleCommand ...
type GoogleCommand struct{}

// GoogleCommandError ...
type GoogleCommandError error

// ErrorGoogleNotFound ...
var ErrorGoogleNotFound GoogleCommandError = errors.New("検索結果ゼロ件でした")

// ErrorGoogleNoQueryGiven ...
var ErrorGoogleNoQueryGiven GoogleCommandError = errors.New("検索語が指定されていません")

// Match ...
func (cmd GoogleCommand) Match(payload *slack.Payload) bool {
	if len(payload.Ext.Words) == 0 {
		return false
	}
	return payload.Ext.Words[0] == "google" || payload.Ext.Words[0] == "ggl"
}

// Handle ...
func (cmd GoogleCommand) Handle(ctx context.Context, payload *slack.Payload) slack.Message {
	client := google.Client{
		APIKey:               os.Getenv("GOOGLE_CUSTOMSEARCH_API_KEY"),
		CustomSearchEngineID: os.Getenv("GOOGLE_CUSTOMSEARCH_ENGINE_ID"),
	}
	words := payload.Ext.Words[1:]
	if len(words) == 0 {
		return wrapError(payload, ErrorGoogleNoQueryGiven)
	}
	query := strings.Join(words, "+")
	res, err := client.CustomSearch(url.Values{"q": {query}, "hl": {"ja"}})
	if err != nil {
		return wrapError(payload, err)
	}
	if len(res.Items) == 0 {
		return wrapError(payload, ErrorGoogleNotFound)
	}
	item := res.Items[0]
	return slack.Message{
		Channel: payload.Event.Channel,
		Text:    fmt.Sprintf("> %v\n%s\n", words, item.Link),
	}
}

// Help ...
func (cmd GoogleCommand) Help(payload *slack.Payload) slack.Message {
	return slack.Message{
		Channel: payload.Event.Channel,
		Text:    "Google検索コマンド",
	}
}
