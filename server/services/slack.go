package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"

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

	slack.Channels = os.Getenv("SLACK_CHANNELS")

	return nil
}

// ServeHTTP ...
func (slack *Slack) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	render := m.Render(w, true)

	// TODO: ちゃんとstructにしたほうがいい
	body := map[string]interface{}{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		render.JSON(http.StatusBadRequest, m.P{
			"message": err.Error(),
		})
		return
	}

	if body["token"] != slack.Verification {
		render.JSON(http.StatusBadRequest, m.P{
			"message": "invalid verification",
		})
		return
	}

	if body["type"] == "url_verification" {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(body["challenge"].(string)))
		return
	}

	if body["type"] != "app_mention" {
		render.JSON(http.StatusOK, m.P{"message": fmt.Sprintf("ignore this type of events: %v", body["type"])})
		return
	}

	entry := amesh.GetEntry()
	ctx := middlewares.Context(r)
	client := middlewares.HTTPClient(ctx)

	img, err := entry.Image(true, true, client)
	if err != nil {
		render.JSON(http.StatusBadRequest, m.P{"message": err.Error()})
		return
	}

	buf := new(bytes.Buffer)
	if err := png.Encode(buf, img); err != nil {
		render.JSON(http.StatusBadRequest, m.P{"message": err.Error()})
		return
	}

	postbody := new(bytes.Buffer)
	writer := multipart.NewWriter(postbody)

	f, err := writer.CreateFormFile("file", "amesh.png")
	if err != nil {
		render.JSON(http.StatusBadRequest, m.P{"message": err.Error()})
		return
	}

	if _, err := io.Copy(f, buf); err != nil {
		render.JSON(http.StatusBadRequest, m.P{"message": err.Error()})
		return
	}

	if err := writer.WriteField("token", slack.BotAccessToken); err != nil {
		render.JSON(http.StatusBadRequest, m.P{"message": err.Error()})
		return
	}

	if err := writer.WriteField("channels", slack.Channels); err != nil {
		render.JSON(http.StatusBadRequest, m.P{"message": err.Error()})
		return
	}

	if err := writer.Close(); err != nil {
		render.JSON(http.StatusBadRequest, m.P{"message": err.Error()})
		return
	}

	req, err := http.NewRequest("POST", "https://slack.com/api/files.upload", postbody)
	if err != nil {
		render.JSON(http.StatusBadRequest, m.P{"message": err.Error()})
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := client.Do(req)
	if err != nil {
		render.JSON(http.StatusBadRequest, m.P{"message": err.Error()})
		return
	}
	if res.StatusCode != http.StatusOK {
		render.JSON(http.StatusBadRequest, m.P{"message": fmt.Sprintf("http status is not OK: (%d) %s", res.StatusCode, res.Status)})
		return
	}
	defer res.Body.Close()

	response := map[string]interface{}{}
	json.NewDecoder(res.Body).Decode(&response)

	render.JSON(http.StatusOK, response)
}
