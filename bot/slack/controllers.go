package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image/png"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/otiai10/amesh/lib/amesh"
	"github.com/otiai10/marmoset"
)

// HandleIndex ...
func HandleIndex(w http.ResponseWriter, r *http.Request) {
	render := marmoset.Render(w, true)
	render.JSON(http.StatusOK, marmoset.P{"message": "hello"})
}

// HandleWebhook ...
func HandleWebhook(w http.ResponseWriter, r *http.Request) {
	render := marmoset.Render(w, true)
	payload := &Payload{}
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		render.JSON(http.StatusBadRequest, marmoset.P{"message": err.Error()})
		return
	}

	if payload.Token != os.Getenv("SLACK_VERIFICATION_TOKEN") {
		render.JSON(http.StatusBadRequest, marmoset.P{"message": "invalid verification"})
		return
	}

	if payload.Type == "url_verification" {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(payload.Challenge))
		return
	}

	// https://api.slack.com/events-api#the-events-api__responding-to-events
	render.JSON(http.StatusAccepted, marmoset.P{"message": "ok"})

	// {{{
	tokyo, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatalln(err)
		return
	}
	entry := amesh.GetEntry(time.Now().In(tokyo))
	img, err := entry.Image(true, true)
	if err != nil {
		log.Fatalln(err)
		return
	}
	buf := bytes.NewBuffer(nil)
	if err := png.Encode(buf, img); err != nil {
		log.Fatalln(err)
		return
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalln(err)
		return
	}
	bname := os.Getenv("GOOGLE_STORAGE_BUCKET_NAME")
	bucket := client.Bucket(bname)
	datetime := entry.Time.Format("2006-0102-1504")
	fname := fmt.Sprintf("%s.png", datetime)
	obj := bucket.Object(fname)
	writer := obj.NewWriter(ctx)
	_, err = writer.Write(buf.Bytes())
	if err != nil {
		log.Fatalln(err)
		return
	}
	writer.Close()
	publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bname, fname)
	// }}}

	body := bytes.NewBuffer(nil)
	err = json.NewEncoder(body).Encode(map[string]interface{}{
		"channel": payload.Event.Channel,
		// https://api.slack.com/messaging/composing/layouts#sending-messages
		"blocks": []Block{{Type: "image", ImageURL: publicURL, AltText: datetime}},
	})
	if err != nil {
		log.Fatalln(err)
		return
	}

	// https://api.slack.com/methods/chat.postMessage
	req, err := http.NewRequest("POST", "https://slack.com/api/chat.postMessage", body)
	if err != nil {
		log.Fatalln(err)
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("SLACK_BOT_USER_OAUTH_ACCESS_TOKEN")))
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
		return
	}
	if res.StatusCode >= 400 {
		log.Fatalln(res.Status)
		return
	}
	log.Println(res.Status)
}
