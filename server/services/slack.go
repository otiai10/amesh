package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"google.golang.org/appengine/taskqueue"

	"github.com/otiai10/amesh/server/middlewares"
	m "github.com/otiai10/marmoset"
)

// Slack ...
type Slack struct {
	BotAccessToken string
	Channels       string
	Verification   string

	// https://stackoverflow.com/questions/50715387/slack-events-api-triggers-multiple-times-by-one-message
	lastEventID string
}

// Init サービスの初期化
func (slack *Slack) Init() error {

	slack.BotAccessToken = os.Getenv("SLACK_BOT_ACCESS_TOKEN")
	if slack.BotAccessToken == "" {
		return fmt.Errorf("SLACK_BOT_ACCESS_TOKEN is not specified")
	}

	slack.Verification = os.Getenv("SLACK_VERIFICATION")
	if slack.Verification == "" {
		return fmt.Errorf("SLACK_VERIFICATION is not specified")
	}

	return nil
}

// WebhookURL ...
func (slack *Slack) WebhookURL() string {
	return "/webhook/slack"
}

// HandleWebhook ...
func (slack *Slack) HandleWebhook(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithDeadline(middlewares.Context(r), time.Now().Add(30*time.Second))
	defer cancel()
	render := m.Render(w, true)

	payload := new(Payload)
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		render.JSON(http.StatusBadRequest, m.P{
			"message": err.Error(),
		})
		return
	}

	if payload.Token != slack.Verification {
		render.JSON(http.StatusBadRequest, m.P{
			"message": "invalid verification",
		})
		return
	}

	if payload.Type == "url_verification" {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(payload.Challenge))
		return
	}

	if payload.Event.Type != "app_mention" {
		render.JSON(http.StatusOK, m.P{"message": fmt.Sprintf("ignore this type of events: %v", payload.Event.Type)})
		return
	}

	t := taskqueue.NewPOSTTask(slack.QueueURL(), url.Values{
		"message": []string{"Hello, otiai10"},
	})

	t, err := taskqueue.Add(ctx, t, "")
	if err != nil {
		render.JSON(http.StatusBadRequest, m.P{
			"message": err.Error(),
		})
		return
	}
	render.JSON(http.StatusOK, map[string]interface{}{
		"queue_name": t.Name,
	})

}

// QueueURL ...
func (slack *Slack) QueueURL() string {
	return "/queue/slack"
}

// HandleQueue ...
func (slack *Slack) HandleQueue(w http.ResponseWriter, r *http.Request) {

	render := m.Render(w, true)
	render.JSON(http.StatusOK, m.P{
		"message": "HandleQueue",
	})

	// entry := amesh.GetEntry()
	// img, err := entry.Image(true, true, client)
	// if err != nil {
	// 	log.Errorf("E01: %v", err)
	// 	return
	// }

	// buf := new(bytes.Buffer)
	// if err := png.Encode(buf, img); err != nil {
	// 	log.Errorf("E02: %v", err)
	// 	return
	// }

	// postbody := new(bytes.Buffer)
	// writer := multipart.NewWriter(postbody)

	// f, err := writer.CreateFormFile("file", "amesh.png")
	// if err != nil {
	// 	log.Errorf("E03: %v", err)
	// 	return
	// }

	// if _, err := io.Copy(f, buf); err != nil {
	// 	log.Errorf("E04: %v", err)
	// 	return
	// }

	// if err := writer.WriteField("token", slack.BotAccessToken); err != nil {
	// 	log.Errorf("E04: %v", err)
	// 	return
	// }

	// if err := writer.WriteField("channels", payload.Event.Channel); err != nil {
	// 	log.Errorf("E05: %v", err)
	// 	return
	// }

	// if err := writer.Close(); err != nil {
	// 	log.Errorf("E06: %v", err)
	// 	return
	// }

	// req, err := http.NewRequest("POST", "https://slack.com/api/files.upload", postbody)
	// if err != nil {
	// 	log.Errorf("E07: %v", err)
	// 	return
	// }
	// req.Header.Set("Content-Type", writer.FormDataContentType())

	// res, err := client.Do(req)
	// if err != nil {
	// 	log.Errorf("E08: %v", err)
	// 	return
	// }
	// if res.StatusCode != http.StatusOK {
	// 	log.Errorf("E09: %v", err)
	// 	return
	// }

	// response := map[string]interface{}{}
	// json.NewDecoder(res.Body).Decode(&response)
	// res.Body.Close()
	// log.Debugf("%+v\n", response)
}

// Payload は、Events API でくるやつ、のはしょったの
type Payload struct {
	Token       string   `json:"token"`
	TeamID      string   `json:"team_id"`
	APIAppID    string   `json:"api_app_id"`
	Type        string   `json:"type"`
	ID          string   `json:"event_id"`
	AuthedUsers []string `json:"authed_users"`
	Timestamp   int64    `json:"event_time"`
	Event       struct {
		Type      string `json:"type"`
		User      string `json:"user"`
		Text      string `json:"text"`
		Channel   string `json:"channel"`
		Timestamp string `json:"event_ts"`
	} `json:"event"`
	Challenge string `json:"challenge"`
}
