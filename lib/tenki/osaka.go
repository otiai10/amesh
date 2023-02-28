package tenki

import (
	"image"
	"log"
	"net/http"
	"time"
)

const (
	// TenkiStaticURL ...
	OsakaEntryPath = "/pref-30-large.jpg"
)

// Osaka ...
type Osaka struct{}

// OsakaEntry ...
type OsakaEntry struct {
	URL string
}

// GetEntry ...
func (osaka Osaka) GetEntry() Entry {
	area := "Asia/Tokyo"
	loc, err := time.LoadLocation(area)
	if err != nil {
		log.Fatalf("Failed to load location `%s`", area)
	}
	now := truncateTime(time.Now().In(loc), 5*time.Minute)
	return OsakaEntry{
		URL: TenkiStaticURL + now.Format(TenkiDynamicTimestampPath) + OsakaEntryPath,
	}
}

// Image ...
func (osaka OsakaEntry) Image(client ...*http.Client) (image.Image, error) {
	if len(client) == 0 {
		client = append(client, http.DefaultClient)
	}
	res, err := client[0].Get(osaka.URL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	img, _, err := image.Decode(res.Body)
	return img, err
}

// ReferenceURL ...
func (osaka OsakaEntry) ReferenceURL() string {
	return "https://tenki.jp/"
}
