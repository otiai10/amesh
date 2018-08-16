package services

import (
	"log"
	"net/http"
	"os"
	"strings"
)

// Handler ...
func Handler() http.Handler {

	name := strings.ToUpper(os.Getenv("SERVICE"))
	switch name {
	case "SLACK":
		slack := new(Slack)
		if err := slack.Init(); err != nil {
			log.Fatalln(err)
		}
		return slack
	default:
		log.Fatalf("chat service name '%s' undefined or unknown", name)
	}

	return nil
}
