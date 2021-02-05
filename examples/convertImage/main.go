/*
 * Copyright 2021 Philoj Johny
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain A copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */
package main

import (
	"fmt"
	treePalette "github.com/philoj/tree-palette"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"
)

// Convert the given image into one using the given palette and serve both images so that one can preview the results.
func main() {
	paletteColors := []treePalette.PaletteColor{
		treePalette.NewOpaquePaletteColor(255, 211, 92, 1, "DANDELION"),
		treePalette.NewOpaquePaletteColor(255, 130, 1, 2, "DARK ORANGE"),
		treePalette.NewOpaquePaletteColor(243, 114, 82, 3, "CRUSTA"),
		treePalette.NewOpaquePaletteColor(199, 44, 58, 7, "BRICK RED"),
		treePalette.NewOpaquePaletteColor(234, 62, 112, 8, "DARK PINK"),
		treePalette.NewOpaquePaletteColor(149, 69, 103, 9, "CADILLAC"),
		treePalette.NewOpaquePaletteColor(75, 196, 213, 10, "MEDIUM TURQUOISE"),
		treePalette.NewOpaquePaletteColor(1, 128, 181, 11, "PACIFIC BLUE"),
		treePalette.NewOpaquePaletteColor(2, 181, 160, 12, "PERSIAN GREEN"),
		treePalette.NewOpaquePaletteColor(138, 151, 71, 13, "OLD OLIVE"),
	}
	palette := treePalette.NewPalette(paletteColors, false)

	// open and decode image
	f, err := os.Open("image.jpg")
	if err != nil {
		panic(err)
	}
	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	palettedImage := palette.ApplyPalette(img)

	http.HandleFunc("/", ServeHtml(`
		<html><body>
			<a href="/originalImage.jpg">Original Image</a>
			<a href="/palettedImage.jpg">Paletted Image</a>
		</body></html>
	`))
	http.HandleFunc("/palettedImage.jpg", ServeImage(palettedImage))
	http.HandleFunc("/originalImage.jpg", ServeImage(img))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func ServeImage(img image.Image) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		headers := w.Header()
		headers["Content-Type"] = []string{"image/jpeg"}
		err := jpeg.Encode(w, img, nil)
		if err != nil {
			panic(err)
		}
	}
}
func ServeHtml(html string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprint(w, html)
		if err != nil {
			panic(err)
		}
	}
}
