package cli

import (
	"fmt"
	"os"

	"github.com/otiai10/amesh/lib/tenki"
	"github.com/otiai10/gat/render"
)

// Tenki
func Tenki(r render.Renderer, location tenki.Location) error {
	entry := location.GetEntry()
	img, err := entry.Image()
	if err != nil {
		return err
	}
	if err := r.Render(os.Stdout, img); err != nil {
		return err
	}
	fmt.Printf("Powered by %s\n", entry.ReferenceURL())
	return nil
}
