// main
package main

import (
	"fmt"
	"image"
	"net/http"

	"github.com/otiai10/amesh"
	"github.com/otiai10/gcat"

	_ "image/gif"
)

func main() {
	entry := amesh.GetEntry()
	res, err := http.Get(entry.Mesh)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	img, _, err := image.Decode(res.Body)
	if err != nil {
		panic(err)
	}
	gcat.OfTerminal().PrintImage(img)
	fmt.Println("#amesh", entry.URL)
}
