# treePalette

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/kyroy/kdtree/blob/master/LICENSE)

An indexed color palette implementation in Go on top of a [k-d tree](https://en.wikipedia.org/wiki/K-d_tree) for fast color lookups. Also rank a palette against an image to identify prominent colors.

- Transparent(RGBA) and opaque(RGB) palettes
- Direct image conversion
- Image pixel counting and color ranking, for prominent color analysis

kd-tree implementation adapted from: [kyroy/kdtree](https://github.com/kyroy/kdtree)

## Usage

```bash
go get github.com/philoj/tree-palette
```

```go
import "github.com/philoj/tree-palette"
````

### Create a color model for color lookups.
```go
m := NewPalettedColorModel([]color.Color{
        // list of colors in the palette
        }, false // ignore alpha values
    )
equivalentColor := m.Convert(someColor)
```


### Color ranking and image color analysis

Start by implementing `treePalette.PaletteColor` and `treePalette.Color` interfaces:

```go
// Color express A color as A n-dimensional point in the RGBA space for usage in the kd-tree search algorithm.
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
```

Or use included implementations `treePalette.ColorRGBA` and `treePalette.IndexedColorRGBA` respectively:
```go
// Unknown color
c := treePalette.NewOpaqueColor(121,201,10)

// Palette colors
p1 := treePalette.NewOpaquePaletteColor(255, 130, 1, 2, "DARK ORANGE") // R,G,B, unique-id, name
p2 := treePalette.NewOpaquePaletteColor(1, 128, 181, 11, "PACIFIC BLUE")

// Create palette
palette := treePalette.NewPalette([]treePalette.PaletteColor{p1,p2}, false)

// Equivalent color
equivalent := palette.Convert(c)

// Convert an image.Image
palettedImage := palette.ApplyPalette(img)

// Rank the palette against all the pixels in an image.Image
colors, colorCount := palette.Rank(img)
fmt.Printf("Most frequent color is %s. It appears %d times.", colors[0], colorCount[colors[0].Index()])
```

