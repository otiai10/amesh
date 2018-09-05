package services

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/otiai10/amesh/server/plugins"
)

// Service は、Slackなど、webhookを受けたり返したりするサービスのインターフェースです。
type Service interface {
	WebhookURL() string
	HandleWebhook(http.ResponseWriter, *http.Request)
	QueueURL() string
	HandleQueue(http.ResponseWriter, *http.Request)
}

// Handler ...
func Handler(p ...plugins.Plugin) Service {

	name := strings.ToUpper(os.Getenv("SERVICE"))
	switch name {
	case "SLACK":
		slack := &Slack{Plugins: p}
		if err := slack.Init(); err != nil {
			log.Fatalln(err)
		}
		return slack
	default:
		log.Fatalf("chat service name '%s' undefined or unknown", name)
	}

	return nil
}
