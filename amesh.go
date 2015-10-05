package amesh

import "time"

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
)

// Entry ...
type Entry struct {
	URL  string `json:"url"`
	Map  string `json:"map"`
	Mesh string `json:"mesh"`
	Mask string `json:"mask"`
}

// GetEntry ...
func GetEntry() Entry {
	return Entry{
		URL:  AmeshURL,
		Map:  getMap(),
		Mesh: getMesh(),
		Mask: getMask(),
	}
}

func getMap() string {
	return AmeshURL + map000
}

func getMask() string {
	return AmeshURL + msk000
}

func getMesh() string {
	return AmeshURL + time.Now().Add(-1*unit).Round(1*unit).Format(mesh)
}
