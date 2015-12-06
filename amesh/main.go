// main
package main

import (
	"flag"
	"fmt"
	"image"
	"log"
	"net/http"
	"os"
	"strings"

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

	switch os.Getenv("AMESH_NOTIFICATION_SERVICE") {
	case "twitter":
		observer.Notifier = amesh.NewTwitterNotifier(
			os.Getenv("AMESH_TWITTER_CONSUMER_KEY"),
			os.Getenv("AMESH_TWITTER_CONSUMER_SECRET"),
			os.Getenv("AMESH_TWITTER_ACCESS_TOKEN"),
			os.Getenv("AMESH_TWITTER_ACCESS_TOKEN_SECRET"),
		)
	}
	users := strings.Split(os.Getenv("AMESH_NOTIFICATION_USERS"), ",")

	observer.On(amesh.Rain, func(ev amesh.Event) {
		msg := fmt.Sprintf("%s 雨がふってるよ！\n%s\n%s",
			strings.Join(users, " "), ev.Timestamp.Format("15:04:05"), amesh.AmeshURL,
		)
		if observer.Notifier != nil {
			log.Println("[RAIN NOTIFICATION]", observer.Notifier.Notify(msg))
		} else {
			log.Println("[RAIN]", msg)
		}
	})
	observer.Start()
}
