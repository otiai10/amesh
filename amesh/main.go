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
	geo, mesh bool
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
	flag.Parse()
}

func main() {
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
