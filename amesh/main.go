// main
package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"image"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/otiai10/amesh"
	"github.com/otiai10/gat"

	_ "image/gif"
	_ "image/jpeg"
	"image/png"
)

var (
	geo, mask bool
	usepix    bool
	daemon    bool
)

func onerror(err error) {
	if err == nil {
		return
	}
	fmt.Println(err)
	os.Exit(1)
}

func init() {
	flag.BoolVar(&geo, "g", true, "地形を描画")
	flag.BoolVar(&mask, "b", true, "県境を描画")
	flag.BoolVar(&usepix, "p", false, "iTermであってもピクセル画で表示")
	flag.BoolVar(&daemon, "d", false, "daemonモード起動")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "東京アメッシュをCLIに表示するコマンドです。\n利用可能なオプション:\n")
		flag.PrintDefaults()
	}
	flag.Parse()
}

func main() {

	if daemon {
		startDaemon()
		return
	}

	entry := amesh.GetEntry()

	merged, err := entry.Image(geo, mask)
	onerror(err)

	switch {
	case !usepix && os.Getenv("TERM_PROGRAM") == "iTerm.app":
		buf := bytes.NewBuffer(nil)
		err = png.Encode(buf, merged)
		onerror(err)
		size := merged.Rect.Size()
		r := 2
		encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
		fmt.Fprintf(os.Stdout, "\033]1337;File=;width=%dpx;height=%dpx;inline=1:%s\a\n", size.X/r, size.Y/r, encoded)
	default:
		gat.NewClient(gat.GetTerminal()).Set(gat.SimpleBorder{}).PrintImage(merged)
	}
	fmt.Println("#amesh", entry.URL)

}

func getImage(url string) (image.Image, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	img, _, err := image.Decode(res.Body)
	return img, err
}

func startDaemon() {

	observer := amesh.NewObserver()

	switch os.Getenv("AMESH_NOTIFICATION_SERVICE") {
	case "slack":
		observer.Notifier = amesh.NewSlackNotifier(
			os.Getenv("AMESH_SLACK_TOKEN"),
			os.Getenv("AMESH_SLACK_CHANNEL"),
		)
	case "twitter":
		observer.Notifier = amesh.NewTwitterNotifier(
			os.Getenv("AMESH_TWITTER_CONSUMER_KEY"),
			os.Getenv("AMESH_TWITTER_CONSUMER_SECRET"),
			os.Getenv("AMESH_TWITTER_ACCESS_TOKEN"),
			os.Getenv("AMESH_TWITTER_ACCESS_TOKEN_SECRET"),
		)
	}
	users := strings.Split(os.Getenv("AMESH_NOTIFICATION_USERS"), ",")

	observer.On(amesh.Rain, func(ev amesh.Event) error {
		msg := fmt.Sprintf("%s 雨がふってるよ！\n%s %s",
			strings.Join(users, " "), amesh.AmeshURL, ev.Timestamp.Format("15:04:05"),
		)
		log.Println("[RAIN]", msg)
		if observer.LastRain.IsZero() && observer.Notifier != nil {
			return observer.Notifier.Notify(msg)
		}
		if observer.LastRain.IsZero() {
			observer.LastRain = ev.Timestamp // to throttle notification
		}
		if ev.Timestamp.After(observer.LastRain.Add(observer.NotificationInterval)) {
			observer.LastRain = time.Time{} // reset to notify again
		}
		return nil
	})
	observer.Start()
}
