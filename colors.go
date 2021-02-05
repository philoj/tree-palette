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
	"fmt"
)

// Color express A color as A n-dimensional point in the RGBA space for usage in the kd-tree search algorithm.
// This supports both RGBA and RGB(no alpha) spaces since latter would reduce computing for cases where transparency is not important.
type Color interface {

	// Dimensions returns the total number of dimensions(3 for RGB, 4 for RGBA).
	Dimensions() int

	// Dimension returns the value of the i-th dimension, say R,G,B and/or A.
	Dimension(i int) uint32
}

// PaletteColor is A Color inside an indexed color palette.
type PaletteColor interface {
	Color

	// Index returns palette index of the color
	Index() int
}

// ColorRGBA Example Color implementation.
type ColorRGBA struct {
	R, G, B, A   uint32 // R,G,B, and A are considered as dimensions 0,1,2 and 3 respectively.
	AlphaChannel bool   // If false, alpha values are ignored.
}

// color.Color implementation
func (c ColorRGBA) RGBA() (uint32, uint32, uint32, uint32) {
	if c.AlphaChannel {
		return c.R, c.G, c.B, c.A
	}
	return c.R, c.G, c.B, 0xffff
}

func (c ColorRGBA) Dimensions() int {
	if c.AlphaChannel {
		return 4
	} else {
		return 3
	}
}

func (c ColorRGBA) Dimension(i int) uint32 {
	switch i {
	case 0:
		return c.R
	case 1:
		return c.G
	case 2:
		return c.B
	case 3:
		if c.AlphaChannel {
			return c.A
		}
		fallthrough
	default:
		panic(fmt.Errorf("invalid dimension %d: expected [0-%d]", i, c.Dimensions()))
	}
}

func (c ColorRGBA) String() string {
	if c.AlphaChannel {
		return fmt.Sprintf("{R:%f, G:%f, B:%d, A:%f}",
			float64(c.R)/float64(0xffff)*255,
			float64(c.G)/float64(0xffff)*255,
			float64(c.B)/float64(0xffff)*255,
			float64(c.A)/float64(0xffff),
		)
	} else {
		return fmt.Sprintf("{R:%f, G:%f, B:%f}",
			float64(c.R)/float64(0xffff)*255,
			float64(c.G)/float64(0xffff)*255,
			float64(c.B)/float64(0xffff)*255,
		)
	}
}

// NewTransparentColor creates a transparent color. R,G,B values are in range [0-255], A in range [0-1]
func NewTransparentColor(R, G, B int, A float64) ColorRGBA {
	r, g, b, a := float64(R)/255*float64(0xffff),
		float64(G)/255*float64(0xffff),
		float64(B)/255*float64(0xffff),
		A*float64(0xffff)
	return ColorRGBA{
		R:            uint32(r),
		G:            uint32(g),
		B:            uint32(b),
		A:            uint32(a),
		AlphaChannel: true,
	}
}

// NewOpaqueColor creates a opaque color. R,G,B values are in range [0-255]
func NewOpaqueColor(R, G, B int) ColorRGBA {
	r, g, b := float64(R)/255*float64(0xffff),
		float64(G)/255*float64(0xffff),
		float64(B)/255*float64(0xffff)
	return ColorRGBA{
		R:            uint32(r),
		G:            uint32(g),
		B:            uint32(b),
		AlphaChannel: false,
	}
}

// IndexedColorRGBA Example PaletteColor implementation.
type IndexedColorRGBA struct {
	ColorRGBA
	Id   int    // Id is the color's unique index
	Name string // A human readable name. Used in stringer
}

func (ic IndexedColorRGBA) Index() int {
	return ic.Id
}

func (ic IndexedColorRGBA) String() string {
	return fmt.Sprintf("%s(%d)", ic.Name, ic.Id)
}

// NewTransparentColor creates a transparent palette color.
// R,G,B values are in range [0-255], A in range [0-1].
// id is the unique id for the color across the palette.
// name is just any human readable identifier, no need to be unique.
func NewTransparentPaletteColor(R, G, B int, A float64, id int, name string) IndexedColorRGBA {
	return IndexedColorRGBA{
		Id:        id,
		Name:      name,
		ColorRGBA: NewTransparentColor(R, G, B, A),
	}
}

// NewOpaquePaletteColor creates an opaque palette color.
// R,G,B values are in range [0-255].
// id is the unique id for the color across the palette.
// name is just any human readable identifier, no need to be unique.
func NewOpaquePaletteColor(R, G, B int, id int, name string) IndexedColorRGBA {
	return IndexedColorRGBA{
		Id:        id,
		Name:      name,
		ColorRGBA: NewOpaqueColor(R, G, B),
	}
}
