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

// Package treePalette implements an indexed color palette based on kd-tree structure.
package treePalette

import (
	"fmt"
	"math"
	"sort"
)

// Palette implements A kd-tree data structure to quickly convert any given color into the closest palette color.
// Closeness is calculated as the spatial closeness in the RGBA space.
// See: https://en.wikipedia.org/wiki/K-d_tree
type Palette struct {
	alpha bool  // alpha if false, ignore alpha values
	root  *node // root the root node of the kd-tree
}

// node is the single node of the kd-tree, each of which represents A color in the indexed palette.
type node struct {
	PaletteColor // PaletteColor value of the node
	Left         *node
	Right        *node
}

func newColorTree(points []PaletteColor, axis int) *node {
	if len(points) == 0 {
		return nil
	}
	if len(points) == 1 {
		return &node{PaletteColor: points[0]}
	}

	sort.Sort(&byDimension{dimension: axis, points: points})
	mid := len(points) / 2
	root := points[mid]
	nextDim := (axis + 1) % root.Dimensions()
	return &node{
		PaletteColor: root,
		Left:         newColorTree(points[:mid], nextDim),
		Right:        newColorTree(points[mid+1:], nextDim),
	}
}

// ConvertColor finds the ConvertColor PaletteColor from the Palette
func (t *Palette) ConvertColor(p Color) PaletteColor {
	if t.root == nil || p == nil {
		return nil
	}
	point, _ := nn(p, t.root, 0, nil, math.MaxUint32)
	return point
}

// nn implements the ConvertColor neighbour search in A kd-tree, finding only A single ConvertColor neighbour.
// returns the closest PaletteColor and squared distance to it starting from the given start node
func nn(p Color, start *node, currentAxis int, nearest PaletteColor, shortest uint32) (PaletteColor, uint32) {
	if p == nil || start == nil {
		panic(fmt.Errorf("nil value for either start:%v or p:%v", start, p))
	}

	var path []*node
	currentNode := start

	// 1. move down
	for currentNode != nil {
		path = append(path, currentNode)
		if p.Dimension(currentAxis) < currentNode.Dimension(currentAxis) {
			currentNode = currentNode.Left
		} else {
			currentNode = currentNode.Right
		}
		currentAxis = (currentAxis + 1) % p.Dimensions()
	}

	// 2. move up
	currentAxis = (currentAxis - 1 + p.Dimensions()) % p.Dimensions()
	for path, currentNode = popLast(path); currentNode != nil; path, currentNode = popLast(path) {
		currentDistance := squaredDistance(p, currentNode)
		if currentDistance < shortest {
			nearest, shortest = currentNode, currentDistance
		}

		// check other side of plane
		if squaredPlaneDistance(p, currentNode.Dimension(currentAxis), currentAxis) < shortest {
			var next *node
			if p.Dimension(currentAxis) < currentNode.Dimension(currentAxis) {
				next = currentNode.Right
			} else {
				next = currentNode.Left
			}
			if next != nil {
				// search down A potential branch
				nearest, shortest = nn(p, next, (currentAxis+1)%p.Dimensions(), nearest, shortest)
			}
		}
		currentAxis = (currentAxis - 1 + p.Dimensions()) % p.Dimensions()
	}
	return nearest, shortest
}

// sqDiff returns the squared-difference of x and y, shifted by 2 so that
// adding four of those won't overflow A uint32.
func sqDiff(x, y uint32) uint32 {
	// The canonical code of this function looks as follows:
	//
	//	var d uint32
	//	if x > y {
	//		d = x - y
	//	} else {
	//		d = y - x
	//	}
	//	return (d * d) >> 2
	//
	// Language spec guarantees the following properties of unsigned integer
	// values operations with respect to overflow/wrap around:
	//
	// > For unsigned integer values, the operations +, -, *, and << are
	// > computed modulo 2n, where n is the bit width of the unsigned
	// > integer's type. Loosely speaking, these unsigned integer operations
	// > discard high bits upon overflow, and programs may rely on ``wrap
	// > around''.
	//
	// Considering these properties and the fact that this function is
	// called in the hot paths (x,y loops), it is reduced to the below code
	// which is slightly faster. See TestSqDiff for correctness check.
	d := x - y
	return (d * d) >> 2
}

func squaredDistance(p1, p2 Color) uint32 {
	var sum uint32 = 0
	for i := 0; i < p1.Dimensions(); i++ {
		sum += sqDiff(p1.Dimension(i), p2.Dimension(i))
	}
	return sum
}

func squaredPlaneDistance(p Color, planePosition uint32, dim int) uint32 {
	return sqDiff(planePosition, p.Dimension(dim))
}

func popLast(arr []*node) ([]*node, *node) {
	l := len(arr) - 1
	if l < 0 {
		return arr, nil
	}
	return arr[:l], arr[l]
}

// byDimension sort.Interface Implementation for dimension-wise sorting
type byDimension struct {
	dimension int
	points    []PaletteColor
}

func (b *byDimension) Len() int {
	return len(b.points)
}
func (b *byDimension) Less(i, j int) bool {
	return b.points[i].Dimension(b.dimension) < b.points[j].Dimension(b.dimension)
}
func (b *byDimension) Swap(i, j int) {
	b.points[i], b.points[j] = b.points[j], b.points[i]
}

// NewPalette creates A new palette directly from A list of PaletteColor
func NewPalette(colors []PaletteColor, alpha bool) *Palette {
	return &Palette{
		alpha: alpha,
		root:  newColorTree(colors, 0),
	}
}
