package slack

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/otiai10/marmoset"
	"github.com/otiai10/spell"
)

// HandleIndex ...
func HandleIndex(w http.ResponseWriter, r *http.Request) {
	render := marmoset.Render(w, true)
	render.JSON(http.StatusOK, marmoset.P{"message": "hello"})
}

// HandleWebhook ...
func HandleWebhook(w http.ResponseWriter, r *http.Request) {
	render := marmoset.Render(w, true)
	payload := &Payload{}
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		render.JSON(http.StatusBadRequest, marmoset.P{"message": err.Error()})
		return
	}

	if payload.Token != os.Getenv("SLACK_VERIFICATION_TOKEN") {
		render.JSON(http.StatusBadRequest, marmoset.P{"message": "invalid verification"})
		return
	}

	if payload.Type == "url_verification" {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(payload.Challenge))
		return
	}

	// https://api.slack.com/events-api#the-events-api__responding-to-events
	render.JSON(http.StatusAccepted, marmoset.P{"message": "ok"})

	go handle(context.Background(), payload)

	return
}

func handle(ctx context.Context, payload *Payload) {
	message := createResponseMessage(context.Background(), payload)
	if err := postMessage(message); err != nil {
		log.Fatalln(err)
	}
}

func createResponseMessage(ctx context.Context, payload *Payload) Message {

	payload.Ext.Words = spell.Parse(payload.Event.Text)[1:]

	commands := []Command{
		AmeshCommand{},
		ImageCommand{},
		ForecastCommand{},
	}

	for _, cmd := range commands {
		if cmd.Match(payload) {
			return cmd.Handle(ctx, payload)
		}
	}

	return Message{
		Channel: payload.Event.Channel,
		Text:    fmt.Sprintf("ちょっと何言ってるかわからない\n> %v", payload.Ext.Words),
	}
}
