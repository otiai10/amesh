package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"time"

	"google.golang.org/appengine/taskqueue"

	"github.com/otiai10/amesh"
	"github.com/otiai10/amesh/server/middlewares"
	m "github.com/otiai10/marmoset"
)

// Slack ...
type Slack struct {
	BotAccessToken string
	Channels       string
	Verification   string
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
		"channels": {payload.Event.Channel},
	})

	t, err := taskqueue.Add(ctx, t, "")
	if err != nil {
		render.JSON(http.StatusBadRequest, m.P{
			"message": err.Error(),
		})
		return
	}

	// FIXME: とりあえず寿司
	err = slack.postMessage(ctx, ":sushi:", payload.Event.Channel)
	render.JSON(http.StatusOK, map[string]interface{}{
		"queue_name": t.Name,
		"error":      err,
	})

}

// QueueURL ...
func (slack *Slack) QueueURL() string {
	return "/queue/slack"
}

// HandleQueue ...
func (slack *Slack) HandleQueue(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithDeadline(middlewares.Context(r), time.Now().Add(60*time.Second))
	defer cancel()
	client := middlewares.HTTPClient(ctx)
	channels := r.FormValue("channels")

	render := m.Render(w, true)

	entry := amesh.GetEntry()
	img, err := entry.Image(true, true, client)
	if err != nil {
		render.JSON(http.StatusOK, m.P{"error": slack.postMessage(ctx, err.Error(), channels)})
		return
	}

	buf := new(bytes.Buffer)
	if err := png.Encode(buf, img); err != nil {
		render.JSON(http.StatusOK, m.P{"error": slack.postMessage(ctx, err.Error(), channels)})
		return
	}

	postbody := new(bytes.Buffer)
	writer := multipart.NewWriter(postbody)

	f, err := writer.CreateFormFile("file", "amesh.png")
	if err != nil {
		render.JSON(http.StatusOK, m.P{"error": slack.postMessage(ctx, err.Error(), channels)})
		return
	}

	if _, err := io.Copy(f, buf); err != nil {
		render.JSON(http.StatusOK, m.P{"error": slack.postMessage(ctx, err.Error(), channels)})
		return
	}

	if err := writer.WriteField("token", slack.BotAccessToken); err != nil {
		render.JSON(http.StatusOK, m.P{"error": slack.postMessage(ctx, err.Error(), channels)})
		return
	}

	if err := writer.WriteField("channels", channels); err != nil {
		render.JSON(http.StatusOK, m.P{"error": slack.postMessage(ctx, err.Error(), channels)})
		return
	}

	if err := writer.Close(); err != nil {
		render.JSON(http.StatusOK, m.P{"error": slack.postMessage(ctx, err.Error(), channels)})
		return
	}

	req, err := http.NewRequest("POST", "https://slack.com/api/files.upload", postbody)
	if err != nil {
		render.JSON(http.StatusOK, m.P{"error": slack.postMessage(ctx, err.Error(), channels)})
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := client.Do(req)
	if err != nil {
		render.JSON(http.StatusOK, m.P{"error": slack.postMessage(ctx, err.Error(), channels)})
		return
	}
	if res.StatusCode != http.StatusOK {
		render.JSON(http.StatusOK, m.P{"error": slack.postMessage(ctx, err.Error(), channels)})
		return
	}
	defer res.Body.Close()

	response := new(APIResponse)
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		render.JSON(http.StatusOK, m.P{"error": slack.postMessage(ctx, err.Error(), channels)})
		return
	}

	render.JSON(http.StatusOK, response)
}

func (slack *Slack) postMessage(ctx context.Context, text, channel string) error {
	log := middlewares.Log(ctx)
	client := middlewares.HTTPClient(ctx)
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(map[string]interface{}{
		"as_user": true,
		"channel": channel,
		"text":    text,
	})
	req, err := http.NewRequest("POST", "https://slack.com/api/chat.postMessage", body)
	if err != nil {
		log.Errorf("Failed to construct http request: %v", err)
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", slack.BotAccessToken))
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Errorf("Failed to post message `%s`: %v", text, err)
		return err
	}
	response := new(APIResponse)
	json.NewDecoder(res.Body).Decode(response)
	if !response.OK {
		log.Errorf("Response from Slack is not ok: %s", response.Error)
		return fmt.Errorf(response.Error)
	}
	return nil
}

// APIResponse は、APIのレスポンス、のはしょったの
type APIResponse struct {
	OK    bool   `json:"ok"`
	Error string `json:"error"`
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
