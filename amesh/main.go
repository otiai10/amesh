// main
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/otiai10/amesh"
	"github.com/otiai10/amesh/plugins/typhoon"
	"github.com/otiai10/gat/render"

	_ "image/gif"
	_ "image/jpeg"
)

var (
	geo, mask bool
	usepix    bool
)

func init() {
	flag.BoolVar(&geo, "g", true, "地形を描画")
	flag.BoolVar(&mask, "b", true, "県境を描画")
	flag.BoolVar(&usepix, "p", false, "iTermであってもピクセル画で表示")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "東京アメッシュをCLIに表示するコマンドです。\n利用可能なオプション:\n")
		flag.PrintDefaults()
	}
	flag.Parse()
}

func main() {

	var renderer render.Renderer
	switch {
	case !usepix && render.ITermImageSupported():
		renderer = &render.ITerm{Scale: 0.5}
	case !usepix && render.SixelSupported():
		renderer = &render.Sixel{Scale: 0.5}
	default:
		renderer = &render.CellGrid{}
	}

	onerror := func(err error) {
		if err == nil {
			return
		}
		fmt.Println(err)
		os.Exit(1)
	}

	// とりあえず
	switch flag.Arg(0) {
	case "typhoon":
		entry, err := typhoon.GetEntry(http.DefaultClient)
		onerror(err)
		img, err := entry.Image(http.DefaultClient)
		onerror(err)
		renderer.Render(os.Stdout, img)
		fmt.Println("#tenkijp", entry.Reference)
		return
	default:
		entry := amesh.GetEntry()
		merged, err := entry.Image(geo, mask)
		onerror(err)
		renderer.Render(os.Stdout, merged)
		fmt.Println("#amesh", entry.URL)
	}

}
