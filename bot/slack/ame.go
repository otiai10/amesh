package slack

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/otiai10/amesh/lib/amesh"
)

// デフォルトの振る舞い。アメッシュを表示
func ame(ctx context.Context, payload *Payload) Message {

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
