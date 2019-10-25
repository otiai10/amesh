package cli

import (
	"fmt"
	"os"

	"github.com/otiai10/amesh/lib/amesh"
	"github.com/otiai10/gat/render"
)

// Amesh デフォルトのアメッシュを表示
func Amesh(r render.Renderer, geo, mask bool) error {

	now, err := getNow()
	if err != nil {
		return err
	}

	entry := amesh.GetEntry(now)
	merged, err := entry.Image(geo, mask)
	if err != nil {
		return err
	}
	if err := r.Render(os.Stdout, merged); err != nil {
		return err
	}
	fmt.Println("#amesh", entry.URL)
	return nil
}
