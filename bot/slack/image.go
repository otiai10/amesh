package slack

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/otiai10/goapis/google"
	"github.com/otiai10/spell"
)

// 画像検索
// https://programmablesearchengine.google.com/cse/all
func searchImage(ctx context.Context, payload *Payload) Message {
	client := google.Client{
		APIKey:               os.Getenv("GOOGLE_CUSTOMSEARCH_API_KEY"),
		CustomSearchEngineID: os.Getenv("GOOGLE_CUSTOMSEARCH_ENGINE_ID"),
		// Referer:           "http://localhost:8080",
	}
	words := spell.Parse(payload.Event.Text)[2:]
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
