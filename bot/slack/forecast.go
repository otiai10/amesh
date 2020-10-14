package slack

import (
	"context"
	"fmt"
	"os"

	"github.com/otiai10/openweathermap"
)

func forecast(ctx context.Context, payload *Payload) Message {
	client := openweathermap.New(os.Getenv("OPENWEATHERMAP_API_KEY"))
	res, err := client.ByCityName("Tokyo", nil)
	if err != nil {
		return Message{
			Channel: payload.Event.Channel,
			Text:    err.Error(),
		}
	}
	message := Message{Channel: payload.Event.Channel}

	for _, forecast := range res.Forecasts {
		if len(forecast.Weather) == 0 {
			continue
		}
		w := forecast.Weather[0]
		message.Blocks = append(message.Blocks, Block{
			Type: "context",
			Elements: []Element{
				{Type: "image", ImageURL: w.IconURL(), AltText: w.Description},
				{Type: "plain_text", Text: fmt.Sprintf("%s | %s", w.Main, w.Description)},
			},
		})
	}

	return message
}
