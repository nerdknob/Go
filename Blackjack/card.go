//go:generate fyne bundle --prefix=img -o bundled.go ./images

package main

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type Card struct {
	Rank  string
	Suit  string
	Value int
	Image canvas.Image
}

func cardFace(r string, s string) fyne.Resource {
	card := fmt.Sprintf("%v%v", r, s)
	for k, v := range images {
		if k == card {
			return v
		}
	}
	mgs := fmt.Sprintf("Could not find card image for %v of %v", r, s)
	log.Fatal(mgs)
	return nil
}
