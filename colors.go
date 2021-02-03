/*
 * Copyright 2021 Philoj Johny
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
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
	"fmt"
)

// color.Color express a color as a n-dimensional point in the RGBA space for usage in the kd-tree search algorithm.
// This supports both RGBA and RGB(no alpha) spaces since latter would reduce computing for cases where transparency is not important.
type Color interface {

	// Dimensions returns the total number of dimensions(3 for RGB, 4 for RGBA).
	Dimensions() int

	// Dimension returns the value of the i-th dimension, say r,g,b and/or a.
	Dimension(i int) uint32
}

// PaletteColor is a Color inside an indexed color palette.
type PaletteColor interface {
	Color

	// Index returns palette index of the color
	Index() int
}


// rgba Example Color implementation.
type rgba struct {
	r, g, b, a   uint32 // r,g,b, and a are considered as dimensions 0,1,2 and 3 respectively.
	alphaChannel bool   // If false, alpha values are ignored.
}

// color.Color implementation
func (c rgba) RGBA() (uint32, uint32, uint32, uint32) {
	if c.alphaChannel {
		return c.r, c.g, c.b, c.a
	}
	return c.r, c.g, c.b, 0xffff
}

func (c rgba) Dimensions() int {
	if c.alphaChannel {
		return 4
	} else {
		return 3
	}
}

func (c rgba) Dimension(i int) uint32 {
	switch i {
	case 0:
		return c.r
	case 1:
		return c.g
	case 2:
		return c.b
	case 3:
		if c.alphaChannel {
			return c.a
		}
		fallthrough
	default:
		panic(fmt.Errorf("invalid dimension %d: expected [0-%d]", i, c.Dimensions()))
	}
}

// indexedColor Example PaletteColor implementation.
type indexedColor struct {
	rgba
	id int // id is the color's unique index
}

func (ic indexedColor) Index() int {
	return ic.id
}
