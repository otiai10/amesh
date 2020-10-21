package slack

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"math/rand"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/otiai10/amesh/lib/amesh"
	"github.com/otiai10/goapis/google"
	"github.com/otiai10/goapis/openweathermap"
)

// Command ...
type Command interface {
	Match(*Payload) bool
	Handle(context.Context, *Payload) Message
	Help(*Payload) Message
}

// AmeshCommand ...
type AmeshCommand struct{}

// Match ...
func (cmd AmeshCommand) Match(payload *Payload) bool {
	return len(payload.Ext.Words) == 0
}

// Handle ...
func (cmd AmeshCommand) Handle(ctx context.Context, payload *Payload) Message {
	tokyo, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return Message{Channel: payload.Event.Channel, Text: err.Error()}
	}
	entry := amesh.GetEntry(time.Now().In(tokyo))
	img, err := entry.Image(true, true)
	if err != nil {
		return Message{Channel: payload.Event.Channel, Text: err.Error()}
	}
	buf := bytes.NewBuffer(nil)
	if err := png.Encode(buf, img); err != nil {
		return Message{Channel: payload.Event.Channel, Text: err.Error()}
	}
	client, err := storage.NewClient(ctx)
	if err != nil {
		return Message{Channel: payload.Event.Channel, Text: err.Error()}
	}

	bname := os.Getenv("GOOGLE_STORAGE_BUCKET_NAME")
	bucket := client.Bucket(bname)
	datetime := entry.Time.Format("2006-0102-1504")
	fname := fmt.Sprintf("%s.png", datetime)
	obj := bucket.Object(fname)
	writer := obj.NewWriter(ctx)
	_, err = writer.Write(buf.Bytes())
	if err != nil {
		return Message{Channel: payload.Event.Channel, Text: err.Error()}
	}

	if err := writer.Close(); err != nil {
		return Message{Channel: payload.Event.Channel, Text: err.Error()}
	}

	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bname, fname)
	return Message{
		Channel: payload.Event.Channel,
		Blocks:  []Block{{Type: "image", ImageURL: url, AltText: datetime}},
	}
}

// Help ...
func (cmd AmeshCommand) Help(payload *Payload) Message {
	return Message{
		Channel: payload.Event.Channel,
		Text:    "デフォルトのアメッシュコマンド",
	}
}

// ForecastCommand ...
type ForecastCommand struct{}

// Match ...
func (cmd ForecastCommand) Match(payload *Payload) bool {
	if len(payload.Ext.Words) == 0 {
		return false
	}
	return payload.Ext.Words[0] == "予報" || payload.Ext.Words[0] == "forecast"
}

// Handle ...
func (cmd ForecastCommand) Handle(ctx context.Context, payload *Payload) Message {
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

// Help ...
func (cmd ForecastCommand) Help(payload *Payload) Message {
	return Message{
		Channel: payload.Event.Channel,
		Text:    "天気予報コマンド",
	}
}

// ImageCommand ...
type ImageCommand struct{}

// Match ...
func (cmd ImageCommand) Match(payload *Payload) bool {
	if len(payload.Ext.Words) == 0 {
		return false
	}
	return payload.Ext.Words[0] == "img" || payload.Ext.Words[0] == "image"
}

// Handle ...
func (cmd ImageCommand) Handle(ctx context.Context, payload *Payload) Message {
	client := google.Client{
		APIKey:               os.Getenv("GOOGLE_CUSTOMSEARCH_API_KEY"),
		CustomSearchEngineID: os.Getenv("GOOGLE_CUSTOMSEARCH_ENGINE_ID"),
	}
	words := payload.Ext.Words[1:]
	if len(words) == 0 {
		return Message{Channel: payload.Event.Channel, Text: "検索語が無いです"}
	}
	query := strings.Join(words, "+")
	rand.Seed(time.Now().Unix())
	res, err := client.SearchImage(query, 1+rand.Intn(10))
	if err != nil {
		return Message{Channel: payload.Event.Channel, Text: fmt.Sprintf("%v\n> %v", err.Error(), words)}
	}
	if len(res.Items) == 0 {
		return Message{Channel: payload.Event.Channel, Text: fmt.Sprintf("結果はゼロでした\n> %v", words)}
	}
	// TODO: ランダムにひとつ選ぶ
	item := res.Items[0]
	return Message{
		Channel: payload.Event.Channel,
		Blocks:  []Block{{Type: "image", ImageURL: item.Link, AltText: query}},
	}
}

// Help ...
func (cmd ImageCommand) Help(payload *Payload) Message {
	return Message{
		Channel: payload.Event.Channel,
		Text:    "画像検索コマンド",
	}
}
