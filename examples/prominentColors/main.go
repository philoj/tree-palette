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
	"github.com/philoj/tree-palette"
	"image"
	_ "image/jpeg"
	"os"
)

// Find prominent palette colors in the given image.
func main() {

	paletteColors := []treepalette.PaletteColor{
		treepalette.NewOpaquePaletteColor(255, 211, 92, 1, "DANDELION"),
		treepalette.NewOpaquePaletteColor(255, 130, 1, 2, "DARK ORANGE"),
		treepalette.NewOpaquePaletteColor(243, 114, 82, 3, "CRUSTA"),
		treepalette.NewOpaquePaletteColor(199, 44, 58, 7, "BRICK RED"),
		treepalette.NewOpaquePaletteColor(234, 62, 112, 8, "DARK PINK"),
		treepalette.NewOpaquePaletteColor(149, 69, 103, 9, "CADILLAC"),
		treepalette.NewOpaquePaletteColor(75, 196, 213, 10, "MEDIUM TURQUOISE"),
		treepalette.NewOpaquePaletteColor(1, 128, 181, 11, "PACIFIC BLUE"),
		treepalette.NewOpaquePaletteColor(2, 181, 160, 12, "PERSIAN GREEN"),
		treepalette.NewOpaquePaletteColor(138, 151, 71, 13, "OLD OLIVE"),
	}
	palette := treepalette.NewPalette(paletteColors, false)

	// open and decode image
	f, err := os.Open("image.jpg")
	if err != nil {
		panic(err)
	}
	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	// Find prominent colors from image
	colors, count := palette.Rank(img)
	for i, c := range colors {
		fmt.Printf("%d. %s - %d Occurances\n", i+1, c, count[c.Index()])
	}
}
