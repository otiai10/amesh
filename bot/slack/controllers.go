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

	message := createResponseMessage(context.Background(), payload)
	if err := postMessage(message); err != nil {
		log.Fatalln(err)
	}

	return
}

func createResponseMessage(ctx context.Context, payload *Payload) Message {

	words := spell.Parse(payload.Event.Text)[1:]
	if len(words) == 0 {
		return ame(ctx, payload)
	}

	command := words[0]

	if command == "予報" {
		return forecast(ctx, payload)
	}

	return Message{
		Channel: payload.Event.Channel,
		Text:    fmt.Sprintf("ちょっと何言ってるかわからない\n> %v", words),
	}
}
