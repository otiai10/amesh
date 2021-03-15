package cli

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

// Timelapse タイムラプス表示
func Timelapse(r render.Renderer, minutes, delay int, loop bool) error {

	fmt.Printf("直近%d分間の降雨画像を取得中", minutes)

	now, err := getNow()
	if err != nil {
		return err
	}

	start := now.Add(time.Duration(-1*minutes) * time.Minute)
	entries := amesh.GetEntries(start, now)

	progress := func(i int) { fmt.Print(".") }
	_, err = entries.GetImages(progress)
	if err != nil {
		return err
	}

	// まずクリアする
	fmt.Printf("\033c")

	var moveCursorToTop = func() {
		fmt.Print("\033[s\033[H\033[1;32m")
	}

	for _, entry := range entries {
		moveCursorToTop()
		r.Render(os.Stdout, entry.Image)
		fmt.Fprintln(os.Stdout, entry.Time.String())
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}

	fmt.Print("\033[0m")

	return nil
}
