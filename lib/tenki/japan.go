package tenki

import (
	"image"
	"log"
	"net/http"
	"time"
)

const (
	// JapanEntryPath ...
	JapanEntryPath = "/japan-detail-large.jpg"
)

// Japan ...
type Japan struct{}

// JapanEntry ...
type JapanEntry struct {
	URL string
}

// GetEntry ...
func (japan Japan) GetEntry() Entry {
	area := "Asia/Tokyo"
	loc, err := time.LoadLocation(area)
	if err != nil {
		log.Fatalf("Failed to load location `%s`", area)
	}
	now := truncateTime(time.Now().In(loc), 5*time.Minute)
	return JapanEntry{
		URL: TenkiStaticURL + now.Format(TenkiDynamicTimestampPath) + JapanEntryPath,
	}
}

// Image ...
func (entry JapanEntry) Image(client ...*http.Client) (image.Image, error) {
	if len(client) == 0 {
		client = append(client, http.DefaultClient)
	}
	res, err := client[0].Get(entry.URL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	img, _, err := image.Decode(res.Body)
	return img, err
}

// ReferenceURL ...
func (entry JapanEntry) ReferenceURL() string {
	return "https://tenki.jp/"
}
