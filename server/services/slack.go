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
	"os"
	"time"

	"github.com/otiai10/amesh"
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

// ServeHTTP ...
func (slack *Slack) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithDeadline(middlewares.Context(r), time.Now().Add(60*time.Second))
	defer cancel()
	client := middlewares.HTTPClient(ctx)
	log := middlewares.Log(ctx)
	render := m.Render(w, true)

	payload := new(Payload)
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		render.JSON(http.StatusBadRequest, m.P{
			"message": err.Error(),
		})
		return
	}
	log.Debugf("%+v\n", payload)

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

	if payload.ID == slack.lastEventID {
		log.Errorf("duplicated event id: %v", payload.ID)
		w.WriteHeader(200)
		return
	}

	// 以下の処理は時間がかかるのでもうHTTPレスポンス返しちゃいます
	render.JSON(http.StatusOK, map[string]interface{}{})

	entry := amesh.GetEntry()
	img, err := entry.Image(true, true, client)
	if err != nil {
		log.Errorf("E01: %v", err)
		return
	}

	buf := new(bytes.Buffer)
	if err := png.Encode(buf, img); err != nil {
		log.Errorf("E02: %v", err)
		return
	}

	postbody := new(bytes.Buffer)
	writer := multipart.NewWriter(postbody)

	f, err := writer.CreateFormFile("file", "amesh.png")
	if err != nil {
		log.Errorf("E03: %v", err)
		return
	}

	if _, err := io.Copy(f, buf); err != nil {
		log.Errorf("E04: %v", err)
		return
	}

	if err := writer.WriteField("token", slack.BotAccessToken); err != nil {
		log.Errorf("E04: %v", err)
		return
	}

	if err := writer.WriteField("channels", payload.Event.Channel); err != nil {
		log.Errorf("E05: %v", err)
		return
	}

	if err := writer.Close(); err != nil {
		log.Errorf("E06: %v", err)
		return
	}

	req, err := http.NewRequest("POST", "https://slack.com/api/files.upload", postbody)
	if err != nil {
		log.Errorf("E07: %v", err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := client.Do(req)
	if err != nil {
		log.Errorf("E08: %v", err)
		return
	}
	if res.StatusCode != http.StatusOK {
		log.Errorf("E09: %v", err)
		return
	}

	response := map[string]interface{}{}
	json.NewDecoder(res.Body).Decode(&response)
	res.Body.Close()
	log.Debugf("%+v\n", response)
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
