package tenki

import (
	"image"
	"net/http"
	"time"
)

const (
	// TenkiStaticURL ...
	TenkiStaticURL = "https://static.tenki.jp/static-images"
	// TenkiDynamicTimestampPath ...
	TenkiDynamicTimestampPath = "/radar/2006/01/02/15/04/00"
)

// Location ...
type Location interface {
	GetEntry() Entry
}

// Entry ...
type Entry interface {
	ReferenceURL() string
	Image(client ...*http.Client) (image.Image, error)
}

var supported = map[string]Location{
	"japan": Japan{},
	"osaka": Osaka{},
}

// GetLocation ...
func GetLocation(name string) Location {
	return supported[name]
}

func truncateTime(t time.Time, unit time.Duration) time.Time {
	return t.Add(-1 * unit).Round(1 * unit)
}
