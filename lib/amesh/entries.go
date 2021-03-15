package amesh

import (
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"net/http"
	"time"
)

// Entries ...
type Entries []*Entry

// GetEntries ...
func GetEntries(start, end time.Time) (entries Entries) {
	t := truncateTime(start)
	entries = append(entries, GetEntry(t))
	for t := t.Add(unit); t.Before(end); t = t.Add(unit) {
		entries = append(entries, GetEntry(t))
	}
	return
}

// ToImages ...
func (entries Entries) GetImages(progress func(int), client ...*http.Client) ([]*image.RGBA, error) {
	images := make([]*image.RGBA, len(entries), len(entries))
	for i, entry := range entries {
		img, err := entry.GetImage(true, true, client...)
		if err != nil {
			return images, err
		}
		images[i] = img
		if progress != nil {
			progress(i)
		}
	}
	return images, nil
}

// ToGif delay == msec
func (entries Entries) ToGif(delay int, loop bool) (*gif.GIF, error) {
	dest := &gif.GIF{LoopCount: 5}
	if loop {
		dest.LoopCount = 0
	}
	images, err := entries.GetImages(nil)
	if err != nil {
		return nil, err
	}
	for _, img := range images {
		paletted := image.NewPaletted(img.Bounds(), palette.Plan9)
		draw.Draw(paletted, paletted.Bounds(), img, img.Bounds().Min, draw.Over)
		dest.Image = append(dest.Image, paletted)
		dest.Delay = append(dest.Delay, delay/10) // Because it's 100ths
	}
	return dest, err
}
