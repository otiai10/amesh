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

	message := createResponseMessage(context.Background(), payload)
	if err := postMessage(message); err != nil {
		log.Fatalln(err)
	}

	return
}

func createResponseMessage(ctx context.Context, payload *Payload) Message {
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
