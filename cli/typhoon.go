package cli

import (
	"fmt"
	"net/http"
	"os"

	"github.com/otiai10/amesh/lib/typhoon"
	"github.com/otiai10/gat/render"
)

// Typhoon 台風情報を表示
func Typhoon(r render.Renderer) error {
	entry, err := typhoon.GetEntry(http.DefaultClient)
	if err != nil {
		return err
	}
	img, err := entry.Image(http.DefaultClient)
	if err != nil {
		return err
	}
	r.Render(os.Stdout, img)
	fmt.Println("#tenkijp", entry.Reference)
	return nil
}
