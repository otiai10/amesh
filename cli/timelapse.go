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

	snapshots, err := getSnapshots(time.Duration(minutes) * time.Minute)
	if err != nil {
		return err
	}

	// まずクリアする
	fmt.Printf("\033c")

	var moveCursorToTop = func() {
		fmt.Print("\033[s\033[H\033[1;32m")
	}

	length := len(snapshots)
	for i := 0; true; i++ {
		if i == length && !loop {
			break
		}
		s := snapshots[i%length]
		moveCursorToTop()
		r.Render(os.Stdout, s.Image)
		fmt.Println(s.Time.String())
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}

	fmt.Print("\033[0m")

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
		fmt.Print(".")
	}
	return snapshots, nil
}
