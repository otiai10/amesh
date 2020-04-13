// main
package main

import (
	"flag"
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	"os"

	"github.com/otiai10/amesh/cli"
	"github.com/otiai10/gat/render"
)

var (
	geo, mask bool
	usepix    bool
	scale     float64

	// 以下、タイムラプスでのみ有効
	lapse   bool
	minutes int
	delay   int
	loop    bool
)

func init() {
	flag.BoolVar(&lapse, "a", false, "タイムラプス表示")
	flag.IntVar(&minutes, "m", 40, "タイムラプスの取得直近時間（分）")
	flag.IntVar(&delay, "d", 200, "タイムラプスアニメーションのfps（msec）")
	flag.BoolVar(&loop, "l", false, "タイムラプスアニメーションをループ表示")
	flag.BoolVar(&geo, "g", true, "地形を描画")
	flag.BoolVar(&mask, "b", true, "県境を描画")
	flag.BoolVar(&usepix, "p", false, "iTermであってもピクセル画で表示")
	flag.Float64Var(&scale, "s", 0.8, "表示拡大倍率")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "東京アメッシュをCLIに表示するコマンドです。\n利用可能なオプション:\n")
		flag.PrintDefaults()
	}
	flag.Parse()
}

func main() {
	renderer := render.GetDefaultRenderer()
	renderer.SetScale(scale)
	subcommand := flag.Arg(0)
	switch {
	case subcommand == "typhoon":
		onerror(cli.Typhoon(renderer))
	case lapse:
		onerror(cli.Timelapse(renderer, minutes, delay, loop))
	default:
		onerror(cli.Amesh(renderer, geo, mask))
	}
}

func onerror(err error) {
	if err == nil {
		return
	}
	fmt.Println(err)
	os.Exit(1)
}
