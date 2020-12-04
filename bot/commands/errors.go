package commands

import (
	"fmt"

	"github.com/otiai10/amesh/bot/slack"
	"google.golang.org/api/googleapi"
)

func wrapError(payload *slack.Payload, err error) *slack.Message {
	switch e := err.(type) {
	case *googleapi.Error:
		return &slack.Message{
			Channel: payload.Event.Channel,
			Text:    fmt.Sprintf("```\n%s\n```", e.Body),
		}
	case GoogleCommandError, ForecastCommandError:
		return &slack.Message{
			Channel: payload.Event.Channel,
			Text:    fmt.Sprintf("```\n%s\n```", err.Error()),
		}
	default:
		typemessage := fmt.Sprintf("UNCAUGHT TYPED ERROR: %T", err)
		return &slack.Message{
			Channel: payload.Event.Channel,
			Text:    fmt.Sprintf("```\n%s\n\n%s\n```", typemessage, err.Error()),
		}
	}
}
