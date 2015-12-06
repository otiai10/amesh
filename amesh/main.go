// main
package main

import (
	"flag"
	"fmt"
	"image"
	"net/http"
	"os"

	"github.com/otiai10/amesh"
	"github.com/otiai10/gat"

	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

var (
	geo, mesh       bool
	daemon          bool
	notifierservice string
	twitter         struct {
		consumer struct {
			key, secret string
		}
		account struct {
			accessToken, accessTokenSecret string
		}
		target string
	}
)

func onerror(err error) {
	if err == nil {
		return
	}
	fmt.Println(err)
	os.Exit(1)
}

func init() {
	flag.BoolVar(&geo, "g", false, "地形を描画")
	flag.BoolVar(&mesh, "m", false, "県境を描画")
	flag.BoolVar(&daemon, "d", false, "daemonモード起動")
	flag.StringVar(&notifierservice, "n", "", "daemonモードの時の通知サービス [twitter|slack]")
	flag.StringVar(&twitter.consumer.key, "tw_consumer_key", "", "-n=twitterの時必要. Twitterのコンシューマキー")
	flag.StringVar(&twitter.consumer.secret, "tw_consumer_secret", "", "-n=twitterの時必要. Twitterのコンシューマシークレット")
	flag.StringVar(&twitter.account.accessToken, "tw_access_token", "", "-n=twitterの時必要. Twitterアカウントのアクセストークン")
	flag.StringVar(&twitter.account.accessTokenSecret, "tw_access_token_secret", "", "-n=twitterの時必要. Twitterアカウントのアクセストークンシークレット")
	flag.StringVar(&twitter.target, "tw_target_acount", "", "-n=twitterの時必要. 雨がふってた時にメンションする@アカウント名")
	flag.Parse()
}

func main() {

	if daemon {
		startDaemon()
		return
	}

	entry := amesh.GetEntry()

	meshLayer, err := getImage(entry.Mesh)
	onerror(err)

	base := image.NewRGBA(meshLayer.Bounds())

	if geo {
		geoLayer, err := getImage(entry.Map)
		onerror(err)
		draw.Draw(base, base.Bounds(), geoLayer, image.Point{0, 0}, 0)
	}

	draw.Draw(base, meshLayer.Bounds(), meshLayer, image.Point{0, 0}, 0)

	if mesh {
		mapLayer, err := getImage(entry.Mask)
		onerror(err)
		draw.Draw(base, base.Bounds(), mapLayer, image.Point{0, 0}, 0)
	}

	gat.NewClient(gat.GetTerminal()).Set(gat.SimpleBorder{}).PrintImage(base)
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

	var notifier amesh.Notifier
	switch notifierservice {
	case "twitter":
		notifier = amesh.NewTwitterNotifier(
			twitter.consumer.key, twitter.consumer.secret,
			twitter.account.accessToken, twitter.account.accessTokenSecret,
		)
	}
	// debug
	observer.IsRaining = func(ev amesh.Event) bool {
		return true
	}
	observer.On(amesh.Rain, func(ev amesh.Event) {
		if notifier != nil && twitter.target != "" {
			// とりあえず俺
			notifier.Notify(fmt.Sprintf("@otiai10 雨がふってるよ！\n%s", ev.Timestamp.String()))
		}
	})
	observer.Start()
}
