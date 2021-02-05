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

package treepalette_test

import (
	"fmt"
	"github.com/philoj/tree-palette"
	"github.com/stretchr/testify/assert"
	"math"
	"math/rand"
	"testing"
)

func TestTreePalette_ConvertColor(t *testing.T) {
	tests := []struct {
		name    string
		color   treepalette.Color
		alpha   bool
		palette []treepalette.PaletteColor
		output  treepalette.Color
	}{
		{
			name:  "nil",
			color: nil,
			alpha: true,
			palette: []treepalette.PaletteColor{
				treepalette.NewOpaquePaletteColor(199, 44, 58, 7, "BRICK RED"),
			},
			output: nil,
		},
		{
			name:    "empty palette",
			color:   treepalette.NewOpaquePaletteColor(2, 181, 160, 12, "PERSIAN GREEN"),
			alpha:   true,
			palette: []treepalette.PaletteColor{},
			output:  nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := treepalette.NewPalette(test.palette, test.alpha)
			assert.Equal(t, test.output, p.ConvertColor(test.color))
		})
	}
}
func TestTreePalette_ConvertColorRandomAlpha(t *testing.T) {
	for i := int64(1); i <= 10; i++ {
		rand.Seed(i)
		c := randomColor(true)
		p := randomPalette(12, true)
		r := closestIndex(c, p)
		t.Run(fmt.Sprintf("random transparent palette: seed %d", i), func(t *testing.T) {
			p := treepalette.NewPalette(p, true)
			cc := p.ConvertColor(c)
			assert.Equal(t, r, cc.Index())
		})
	}
}
func TestTreePalette_ConvertColorRandomNoAlpha(t *testing.T) {
	for i := int64(1); i <= 10; i++ {
		rand.Seed(i)
		c := randomColor(false)
		p := randomPalette(12, false)
		r := closestIndex(c, p)
		t.Run(fmt.Sprintf("random opaque palette: seed %d", i), func(t *testing.T) {
			p := treepalette.NewPalette(p, false)
			cc := p.ConvertColor(c)
			assert.Equal(t, r, cc.Index())
		})
	}
}

// Helper functions

func randomColor(alpha bool) treepalette.ColorRGBA {
	if alpha {
		return treepalette.NewTransparentColor(rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Float64())
	} else {
		return treepalette.NewOpaqueColor(rand.Intn(255), rand.Intn(255), rand.Intn(255))
	}
}

func randomPalette(n int, alpha bool) []treepalette.PaletteColor {
	p := make([]treepalette.PaletteColor, n)
	for i := range p {
		p[i] = &treepalette.IndexedColorRGBA{
			ColorRGBA: randomColor(alpha),
			Id:        i,
		}
	}
	return p
}

func closestIndex(c treepalette.Color, p []treepalette.PaletteColor) int {
	var minD int64 = math.MaxInt64
	var result int
	for _, cc := range p {
		var d int64 = 0
		for i := 0; i < cc.Dimensions(); i++ {
			v2, v1 := int64(cc.Dimension(i)), int64(c.Dimension(i))
			d += (v2 - v1) * (v2 - v1)
		}
		if d < minD {
			minD, result = d, cc.Index()
		}
	}
	return result
}
