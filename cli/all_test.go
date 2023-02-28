package cli

import (
	"image"
	"io"
	"testing"

	. "github.com/otiai10/mint"
)

type PseudoRenderer struct{}

func (r PseudoRenderer) Render(w io.Writer, img image.Image) error {
	return nil
}

func (r PseudoRenderer) SetScale(s float64) error {
	return nil
}

func TestAmesh(t *testing.T) {
	err := Amesh(PseudoRenderer{}, true, true)
	Expect(t, err).ToBe(nil)
}
