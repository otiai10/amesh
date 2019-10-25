package cli

import "time"

const defaultLocation = "Asia/Tokyo"

// FIXME: ameshは東京のやつなのでとりあえずAsia/Tokyoだけにします。
// ゆうて日本国内なら同じtimezoneなのであんまり問題無いと思うけど、
// 今後世界都市対応するなら、これは問題になりますね。
func getNow() (time.Time, error) {
	location, err := time.LoadLocation(defaultLocation)
	if err != nil {
		return time.Time{}, err
	}
	return time.Now().In(location), nil
}
