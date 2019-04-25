package services

// import (
// 	"context"
// 	"fmt"
// 	"testing"

// 	. "github.com/otiai10/mint"
// )

// func TestSlack_PostMessage(t *testing.T) {
// 	ctx := context.Background()
// 	slack := &Slack{
// 		BotAccessToken: "xoxb-???????",
// 	}
// 	message := Message{
// 		Text:   "This is test",
// 		Blocks: []Block{},
// 	}
// 	for i := 0; i < 20; i++ {
// 		message.Blocks = append(message.Blocks, Block{
// 			Type: "context",
// 			Elements: []BlockElement{
// 				{
// 					Type: "image", AltText: "Sunny",
// 					ImageURL: "https://openweathermap.org/img/w/01d.png",
// 				},
// 				{
// 					Type: "mrkdwn", Text: fmt.Sprintf("%02d:00", i),
// 				},
// 				{
// 					Type: "mrkdwn", Text: "Sunny",
// 				},
// 			},
// 		})
// 	}
// 	res, err := slack.PostMessage(ctx, "otiai10-dev", message)
// 	Expect(t, err).ToBe(nil)
// 	Expect(t, res).Not().ToBe(nil)
// }
