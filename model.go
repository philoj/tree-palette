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
package treepalette

import "image/color"

//
// color.Model implementation for Palette.
//

// Convert converts the given color into one of the palette colors.
func (t *Palette) Convert(p color.Color) color.Color {
	c := ColorRGBA{AlphaChannel: t.alpha}
	c.R, c.G, c.B, c.A = p.RGBA()
	c.AlphaChannel = t.alpha
	res := t.ConvertColor(c)
	cc := ColorRGBA{AlphaChannel: t.alpha}
	cc.R, cc.G, cc.B = res.Dimension(0), res.Dimension(1), res.Dimension(2)
	if t.alpha {
		cc.A = res.Dimension(3)
	} else {
		cc.A = 0xffff
	}
	return cc
}

// Create Palette as color.Model from A list of color.Color
// alpha if false, ignores transparency(A) values
func NewPalettedColorModel(colors []color.Color, alpha bool) color.Model {
	var nodes []PaletteColor
	for i, p := range colors {
		r, g, b, a := p.RGBA()
		nodes = append(nodes, IndexedColorRGBA{
			Id: i,
			ColorRGBA: ColorRGBA{
				R:            r,
				G:            g,
				B:            b,
				A:            a,
				AlphaChannel: alpha},
		})
	}
	return NewPalette(nodes, alpha)
}
