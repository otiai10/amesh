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

// Bot ...
type Bot struct {
	Commands []Command
}

// Command ...
type Command interface {
	Match(*Payload) bool
	Handle(context.Context, *Payload) Message
	Help(*Payload) Message
}

func (bot Bot) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	go bot.handle(context.Background(), payload)

	return
}

func (bot Bot) handle(ctx context.Context, payload *Payload) {
	message := bot.createResponseMessage(context.Background(), payload)
	if err := postMessage(message); err != nil {
		log.Println(err)
	}
}

func (bot Bot) createResponseMessage(ctx context.Context, payload *Payload) (message Message) {

	defer func() {
		if r := recover(); r != nil {
			message = Message{
				Channel: payload.Event.Channel,
				Text:    fmt.Sprintf("🤪\n> %v\n```\n%v\n```", payload.Ext.Words, r),
			}
		}
	}()

	payload.Ext.Words = spell.Parse(payload.Event.Text)[1:]

	for _, cmd := range bot.Commands {
		if cmd.Match(payload) {
			return cmd.Handle(ctx, payload)
		}
	}

	return Message{
		Channel: payload.Event.Channel,
		Text:    fmt.Sprintf("ちょっと何言ってるかわからない\n> %v", payload.Ext.Words),
	}
}
