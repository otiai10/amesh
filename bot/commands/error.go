package commands

import (
	"fmt"
	"log"

	"github.com/otiai10/amesh/bot/slack"
	"google.golang.org/api/googleapi"
)

func wrapError(err error, channel ...string) slack.Message {
	channel = append(channel, "bot-dev")
	switch e := err.(type) {
	case *googleapi.Error:
		return slack.Message{Channel: channel[0], Text: fmt.Sprintf("```\n%s\n```", e.Body)}
	default:
		log.Printf("[DEBUG][writer.Write] UNCAUGHT TYPED ERROR: %T", err)
		return slack.Message{Channel: channel[0], Text: fmt.Sprintf("```\n%s\n```", err.Error())}
	}
}
