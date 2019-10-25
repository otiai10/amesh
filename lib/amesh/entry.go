package amesh

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"time"
)

const (
	// AmeshURL は「東京アメッシュ」のURLです
	AmeshURL = "http://tokyo-ame.jwa.or.jp"
	// 地形図
	map000 = "/map/map000.jpg"
	// 地名図
	msk000 = "/map/msk000.png"
	// 雨分布画像の時刻対応フォーマット
	mesh = "/mesh/000/200601021504.gif"
	// 雨分布画像が更新される間隔
	unit time.Duration = 5 * time.Minute
	// アメッシュは東京だけなので
	defaultLocation = "Asia/Tokyo"
)

// Entry ...
type Entry struct {
	URL  string    `json:"url"`
	Map  string    `json:"map"`
	Mesh string    `json:"mesh"`
	Mask string    `json:"mask"`
	Time time.Time `json:"time"`

	IsRainingFunc func(image.Image) (bool, error)
}

// GetEntry ...
func GetEntry(t time.Time) Entry {
	t = truncateTime(t)
	return Entry{
		URL:  AmeshURL,
		Map:  getMap(),
		Mesh: getMesh(t),
		Mask: getMask(),
		Time: t,
	}
}

func truncateTime(t time.Time) time.Time {
	return t.Add(-1 * unit).Round(1 * unit)
}

func getMap() string {
	return AmeshURL + map000
}

func getMask() string {
	return AmeshURL + msk000
}

func getMesh(t time.Time) string {
	return AmeshURL + t.Format(mesh)
}

// Image fetches image data from URL and merge them if needed.
func (entry Entry) Image(geo, mask bool, client ...*http.Client) (*image.RGBA, error) {

	// If client not specified, use default HTTP client.
	// This is because, for example, Google App Engine requires HTTP client with context.
	if len(client) == 0 {
		client = append(client, http.DefaultClient)
	}

	meshlayer, err := entry.getImageFor(entry.Mesh, client[0])
	if err != nil {
		return nil, fmt.Errorf("failed to get image for mesh: %v", err)
	}

	merged := image.NewRGBA(meshlayer.Bounds())

	if geo {
		geolayer, err := entry.getImageFor(entry.Map, client[0])
		if err != nil {
			return nil, fmt.Errorf("failed to get image for geo: %v", err)
		}
		draw.Draw(merged, geolayer.Bounds(), geolayer, image.Point{0, 0}, 0)
	}

	draw.Draw(merged, meshlayer.Bounds(), meshlayer, image.Point{0, 0}, 0)

	if mask {
		masklayer, err := entry.getImageFor(entry.Mask, client[0])
		if err != nil {
			return nil, fmt.Errorf("failed to get image for mask: %v", err)
		}
		draw.Draw(merged, masklayer.Bounds(), masklayer, image.Point{0, 0}, 0)
	}

	return merged, nil
}

func (entry Entry) getImageFor(imgurl string, client *http.Client) (image.Image, error) {
	res, err := client.Get(imgurl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode >= 400 {
		return nil, fmt.Errorf(res.Status)
	}
	img, _, err := image.Decode(res.Body)
	return img, err
}

// IsRaining ...
func (entry *Entry) IsRaining(cliet *http.Client) (bool, error) {

	res, err := cliet.Get(entry.Mesh)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	img, _, err := image.Decode(res.Body)
	if err != nil {
		return false, err
	}

	if entry.IsRainingFunc != nil {
		return entry.IsRainingFunc(img)
	}

	max := img.Bounds().Max
	var hit, all float64 = 0, float64(max.X) * float64(max.Y)

	for y := 1; y < max.Y-1; y++ {
		for x := 1; x < max.X-1; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			if r+g+b+a > 100 {
				hit++
			}
		}
	}
	var threshold float64 = 30
	if (hit*100)/all > threshold {
		return true, nil
	}

	return false, nil
}
