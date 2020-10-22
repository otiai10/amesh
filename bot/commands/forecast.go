package commands

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/otiai10/amesh/bot/slack"
	"github.com/otiai10/goapis/openweathermap"
)

// ForecastCommand ...
type ForecastCommand struct{}

// Match ...
func (cmd ForecastCommand) Match(payload *slack.Payload) bool {
	if len(payload.Ext.Words) == 0 {
		return false
	}
	return payload.Ext.Words[0] == "予報" || payload.Ext.Words[0] == "forecast"
}

// Handle ...
func (cmd ForecastCommand) Handle(ctx context.Context, payload *slack.Payload) slack.Message {
	client := openweathermap.New(os.Getenv("OPENWEATHERMAP_API_KEY"))
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return slack.Message{Channel: payload.Event.Channel, Text: err.Error()}
	}
	res, err := client.ByCityName("Tokyo", nil)
	if err != nil {
		return slack.Message{Channel: payload.Event.Channel, Text: err.Error()}
	}
	if len(res.Forecasts) == 0 || len(res.Forecasts[0].Weather) == 0 {
		return slack.Message{Channel: payload.Event.Channel, Text: "Not enough forecast entries."}
	}
	message := slack.Message{
		Channel: payload.Event.Channel,
		Blocks:  []slack.Block{},
	}

	// {{{ 日付で分けて、Blockを作っていく
	var blockdate int
	var block *slack.Block = nil
	for _, forecast := range res.Forecasts {
		_, month, date := time.Unix(forecast.Timestamp, 0).In(loc).Date()
		// 新しい日付であればBlockを初期化
		if date != blockdate {
			if block != nil {
				message.Blocks = append(message.Blocks, *block)
			}
			block = &slack.Block{Type: "context", Elements: []slack.Element{{Type: "plain_text", Text: fmt.Sprintf("%d/%d", month, date)}}}
			blockdate = date
		}
		w := forecast.Weather[0]
		block.Elements = append(
			block.Elements,
			// Element{Type: "plain_text", Text: time.Unix(forecast.Timestamp, 0).In(loc).Format("15:04")},
			slack.Element{Type: "image", ImageURL: w.IconURL("@2x"), AltText: w.Description},
			// Element{Type: "plain_text", Text: "|"},
		)
	}
	if block != nil {
		message.Blocks = append(message.Blocks, *block)
	}
	// }}}

	return message
}

// Help ...
func (cmd ForecastCommand) Help(payload *slack.Payload) slack.Message {
	return slack.Message{
		Channel: payload.Event.Channel,
		Text:    "天気予報コマンド",
	}
}
