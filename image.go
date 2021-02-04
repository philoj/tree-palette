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
package treePalette

import (
	"image"
	"image/color"
	"sort"
)

// paletted wraps A source image into A 'paletted' image.
type paletted struct {
	src image.Image // src original image
	p   *Palette
}

func (i *paletted) ColorModel() color.Model {
	return i.p
}
func (i *paletted) Bounds() image.Rectangle {
	return i.src.Bounds()
}
func (i *paletted) At(x, y int) color.Color {
	c := ColorRGBA{AlphaChannel: i.p.alpha}
	c.R, c.G, c.B, c.A = i.src.At(x, y).RGBA()
	return i.p.Convert(c)
}

func (i *paletted) ColorIndexAt(x, y int) int {
	c := ColorRGBA{AlphaChannel: i.p.alpha}
	c.R, c.G, c.B, c.A = i.src.At(x, y).RGBA()
	return i.p.ConvertColor(c).Index()
}

// ApplyPalette applies the palette onto A given image and returns new image with Palette as color.Model.
func (t *Palette) ApplyPalette(img image.Image) image.Image {
	return &paletted{
		src: img,
		p:   t,
	}
}

// Rank ranks the colors in the Palette based on counts of pixels of each PaletteColor in the given image.
// Returns A rank list of colors(most occurrences first) and A map with count of pixels for each color index.
func (t *Palette) Rank(img image.Image) ([]PaletteColor, map[int]int) {
	count := make(map[int]int)
	var colors []PaletteColor
	pImg := &paletted{
		src: img,
		p:   t,
	}
	b := img.Bounds()
	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx(); x++ {
			index := pImg.ColorIndexAt(x, y)
			_, ok := count[index]
			if !ok {
				colors = append(colors, t.lookup[index])
			}
			count[index]++
		}
	}
	sort.Slice(colors, func(i, j int) bool {
		return count[colors[i].Index()] > count[colors[j].Index()]
	})
	return colors, count
}

// RankByIndex ranks the colors in the Palette based on counts of pixels of each PaletteColor in the given image.
// Returns A rank list of color indexes(most occurrences first) and A map with count of pixels for each color index.
func (t *Palette) RankByIndex(img image.Image) ([]int, map[int]int) {
	count := make(map[int]int)
	var colors []int
	pImg := &paletted{
		src: img,
		p:   t,
	}
	b := img.Bounds()
	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx(); x++ {
			index := pImg.ColorIndexAt(x, y)
			_, ok := count[index]
			if !ok {
				colors = append(colors, index)
			}
			count[index]++
		}
	}
	sort.Slice(colors, func(i, j int) bool {
		return count[colors[i]] > count[colors[j]]
	})
	return colors, count
}
