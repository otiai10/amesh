package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"

	"google.golang.org/appengine/taskqueue"

	"github.com/otiai10/amesh/server/middlewares"
	m "github.com/otiai10/marmoset"
)

const (
	slackMethodClean = "clean"
	slackMethodShow  = "show"
)

var (
	slackDirectMentionTextFormat = regexp.MustCompile("^<@([0-9A-Z]+)>[ 　]*(.+)?$")
)

type (
	// Slack Handler
	Slack struct {
		UserAccessToken string
		BotAccessToken  string
		Verification    string
	}

	// SlackAPIResponse は、APIのレスポンス、のはしょったの
	SlackAPIResponse struct {
		OK    bool   `json:"ok"`
		Error string `json:"error"`

		Message struct {
			Text      string `json:"text"`
			Timestamp string `json:"ts"`
		} `json:"message"`

		// files.list のレスポンス
		Files []*SlackFile `json:"files,omitempty"`
	}

	// SlackPayload は、Events API でくるやつ、のはしょったの
	SlackPayload struct {
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

	// SlackFile files.list のレスポンス参照
	SlackFile struct {
		ID    string `json:"id"`
		Name  string `json:"string"`
		Title string `json:"title"`
	}
)

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

	slack.UserAccessToken = os.Getenv("SLACK_USER_ACCESS_TOKEN")

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

	payload := new(SlackPayload)
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

	middlewares.Log(ctx).Debugf("%+v\n", payload.Event)
	matches := slackDirectMentionTextFormat.FindStringSubmatch(payload.Event.Text)
	if len(matches) == 0 {
		render.JSON(http.StatusOK, m.P{"message": fmt.Sprintf("ignore this text `%s`", payload.Event.Text)})
		return
	}
	bot := matches[1]
	text := matches[2]

	// メンションの内容から、TaskQueueの種類を変える
	var t *taskqueue.Task
	switch text {
	case slackMethodClean: // このチャンネルに、このbotが投稿したファイルを全消しするタスク
		t = taskqueue.NewPOSTTask(slack.QueueURL(), url.Values{"channel": {payload.Event.Channel}, "method": {slackMethodClean}, "bot": {bot}})
	case "": // アメッシュ画像のアップロードをするタスク
		t = taskqueue.NewPOSTTask(slack.QueueURL(), url.Values{"channel": {payload.Event.Channel}, "method": {slackMethodShow}})
	default:
		render.JSON(http.StatusOK, m.P{"accepted": false})
		return
	}

	if _, err := taskqueue.Add(ctx, t, ""); err != nil {
		slack.onError(ctx, w, err, payload.Event.Channel)
		return
	}
	render.JSON(http.StatusOK, m.P{"accepted": true, "text": text})
}

// QueueURL ...
func (slack *Slack) QueueURL() string {
	return "/queue/slack"
}

// HandleQueue ...
func (slack *Slack) HandleQueue(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithDeadline(middlewares.Context(r), time.Now().Add(60*time.Second))
	defer cancel()

	channel := r.FormValue("channel")
	method := r.FormValue("method")
	bot := r.FormValue("bot")

	switch method {
	case slackMethodClean:
		if err := slack.methodClean(ctx, channel, bot); err != nil {
			slack.onError(ctx, w, err, channel)
			return
		}
	case slackMethodShow:
		if err := slack.methodShow(ctx, channel); err != nil {
			slack.onError(ctx, w, err, channel)
			return
		}
	}

	m.Render(w, true).JSON(http.StatusOK, m.P{
		"method":  method,
		"channel": channel,
		"bot":     bot,
	})

}

func (slack *Slack) onError(ctx context.Context, w http.ResponseWriter, err error, channel string) {
	res, err := slack.postMessage(ctx, err.Error(), channel)
	m.Render(w, true).JSON(http.StatusOK, m.P{
		"response": res,
		"error":    err,
	})
}

func (slack *Slack) postMessage(ctx context.Context, text, channel string) (*SlackAPIResponse, error) {
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
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", slack.BotAccessToken))
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Errorf("Failed to post message `%s`: %v", text, err)
		return nil, err
	}
	response := new(SlackAPIResponse)
	json.NewDecoder(res.Body).Decode(response)
	if !response.OK {
		log.Errorf("Response from Slack is not ok: %s", response.Error)
		return nil, fmt.Errorf(response.Error)
	}
	return response, nil
}
