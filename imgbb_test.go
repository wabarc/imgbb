// Copyright 2021 Wayback Archiver. All rights reserved.
// Use of this source code is governed by the GNU GPL v3
// license that can be found in the LICENSE file.

package imgbb // import "github.com/wabarc/imgbb"

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
)

func genImage(height int) *os.File {
	width := 1000

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// Colors are defined by Red, Green, Blue, Alpha uint8 values.
	cyan := color.RGBA{100, 200, 200, 0xff}

	// Set color for each pixel.
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			switch {
			case x < width/2 && y < height/2: // upper left quadrant
				img.Set(x, y, cyan)
			case x >= width/2 && y >= height/2: // lower right quadrant
				img.Set(x, y, color.White)
			default:
				// Use zero value.
			}
		}
	}

	// Encode as PNG.
	f, _ := os.Create(os.TempDir() + "/image.png")
	png.Encode(f, img)

	return f
}

func TestUploadByURI(t *testing.T) {
	f := genImage(1024)
	defer os.Remove(f.Name())

	i := NewImgBB(nil, "")
	if url, err := i.Upload(f.Name()); err != nil {
		t.Fatal(err)
	} else {
		t.Log(url)
	}
}
