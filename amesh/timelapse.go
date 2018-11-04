package main

import (
	"fmt"
	"image"
	"os"
	"time"

	"github.com/otiai10/amesh/lib/amesh"
	"github.com/otiai10/gat/render"
)

type snapshot struct {
	Image *image.RGBA
	Time  time.Time
}

// タイムラプス表示
func timelapse(r render.Renderer, minutes int) error {

	fmt.Printf("直近%d分間の降雨画像を取得中...", minutes)

	snapshots, err := getSnapshots(time.Duration(minutes) * time.Minute)
	if err != nil {
		return err
	}

	// まずクリアする
	fmt.Printf("\033c")

	// FIXME: カーソルを一番上に持っていく作業
	var moveCursorToTop = func() {
		height := 1000
		fmt.Printf("\033[s\033[%dA\033[1;32m", height)
	}

	for _, s := range snapshots {
		moveCursorToTop()
		r.Render(os.Stdout, s.Image)
		fmt.Println(s.Time.String())
		time.Sleep(500 * time.Millisecond)
	}

	// TODO: カラーリングをリセットする

	return nil
}

func getSnapshots(dur time.Duration) (snapshots []snapshot, err error) {

	sheets := int((int64(dur) / int64(5*time.Minute))) + 1
	for i := 0; i < sheets; i++ {
		t := time.Now().Add(time.Duration(-5*(sheets-i)) * time.Minute)
		entry := amesh.GetEntry(t)
		img, err := entry.Image(true, true)
		if err != nil {
			return nil, err
		}
		snapshots = append(snapshots, snapshot{img, entry.Time})
	}
	return snapshots, nil
}
