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
func timelapse(r render.Renderer) error {

	snapshots, err := getSnapshots(6)
	if err != nil {
		return err
	}

	// まずクリアする
	fmt.Printf("\033c")

	for _, s := range snapshots {
		// カーソル位置を **十分** 戻す
		// FIXME: 十分？
		fmt.Printf("\033[s\033[32A\033[1;32m")
		r.Render(os.Stdout, s.Image)
		fmt.Println(s.Time.String())
		time.Sleep(500 * time.Millisecond)
	}

	// TODO: カラーリングをリセットする

	return nil
}

func getSnapshots(length int) (snapshots []snapshot, err error) {
	for i := 0; i < length; i++ {
		t := time.Now().Add(time.Duration(-5*(length-i)) * time.Minute)
		entry := amesh.GetEntry(t)
		img, err := entry.Image(true, true)
		if err != nil {
			return nil, err
		}
		snapshots = append(snapshots, snapshot{img, entry.Time})
	}
	return snapshots, nil
}
