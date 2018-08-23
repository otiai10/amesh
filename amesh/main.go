// main
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/otiai10/amesh"
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

	entry := amesh.GetEntry()
	merged, err := entry.Image(geo, mask)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	switch {
	case !usepix && render.ITermImageSupported():
		iterm := &render.ITerm{Scale: 0.5}
		iterm.Render(os.Stdout, merged)
	case !usepix && render.SixelSupported():
		sixel := &render.Sixel{Scale: 0.5}
		sixel.Render(os.Stdout, merged)
	default:
		cellgrid := &render.CellGrid{}
		cellgrid.Render(os.Stdout, merged)
	}
	fmt.Println("#amesh", entry.URL)
}
