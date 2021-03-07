package sketch

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"

	"github.com/fogleman/gg"
	"github.com/teacat/noire"
)

var (
	crystalColors = [...]noire.Color{
		noire.NewRGB(63, 132, 229),
		noire.NewRGB(73, 65, 109),
		noire.NewRGB(178, 13, 48),
		noire.NewRGB(129, 120, 23),
		noire.NewRGB(224, 239, 222),
		noire.NewRGB(208, 188, 213),
		noire.NewRGB(234, 186, 107),
		noire.NewRGB(131, 144, 115),
		noire.NewRGB(191, 168, 158),
		noire.NewRGB(254, 94, 645),
		noire.NewRGB(164, 194, 168),
		noire.NewRGB(205, 237, 246),
		noire.NewRGB(239, 123, 69),
		noire.NewRGB(216, 71, 39),
		noire.NewRGB(65, 69, 53),
		noire.NewRGB(242, 227, 188),
		noire.NewRGB(193, 152, 117),
		noire.NewRGB(105, 143, 63),
		noire.NewRGB(163, 154, 146),
	}
)

// GrowthParams contains externally-provided parameters
type GrowthParams struct {
	// tweakable parameters for the cli
	DestWidth     int
	DestHeight    int
	StartingSeeds int
}

// GrowthSketch wraps all the components needed to draw the sketch
type GrowthSketch struct {
	GrowthParams
	DC    *gg.Context
	Seeds []seed
}

type seed struct {
	x      int
	y      int
	r      int
	c      noire.Color
	colorR int
	colorG int
	colorB int
	grew   bool
}

// NewGrowthSketch initializes the canvas and GrowthSketch
func NewGrowthSketch(params GrowthParams) *GrowthSketch {
	fmt.Println("Starting Sketch")

	s := &GrowthSketch{GrowthParams: params}

	// canvas is a gg image context and contains what gets drawn to the screen
	canvas := gg.NewContext(s.DestWidth, s.DestHeight)
	canvas.SetLineWidth(0.0)
	canvas.SetColor(color.White)
	canvas.DrawRectangle(0, 0, float64(s.DestWidth), float64(s.DestHeight))
	canvas.FillPreserve()
	canvas.Stroke()
	s.DC = canvas

	for i := 0; i < s.StartingSeeds; i++ {
		c := crystalColors[rand.Intn(len(crystalColors))]
		r, g, b := c.RGB()
		s.Seeds = append(s.Seeds, seed{x: rand.Intn(s.DestWidth), y: rand.Intn(s.DestHeight), r: 0, c: c, colorR: int(r), colorG: int(g), colorB: int(b)})
	}

	return s
}

// Output produces an image output of the current state of the sketch
func (s *GrowthSketch) Output() image.Image {
	return s.DC.Image()
}

// Draw completes the drawing
func (s *GrowthSketch) Draw() {
	// grow the seeds until there is no room left
	// expand all "borders" unless a color already exists there
	// track which seeds grow in an iteration
	// stopping point for iterations is set as higher that it should ever take, shortcut later when all seeds stop
	for k := 0; k < s.DestHeight*s.DestWidth; k++ {
		for i := 0; i < len(s.Seeds); i++ {
			seed := s.Seeds[i]
			seed.grew = false

			s.DC.SetRGB255(seed.colorR, seed.colorG, seed.colorB)

			for m := -seed.r; m <= seed.r; m++ {
				for n := -seed.r; n <= seed.r; n++ {
					x := seed.x + m
					y := seed.y + n

					if x >= 0 && x <= int(s.DestWidth) && y >= 0 && y <= int(s.DestHeight) {
						cr, cg, cb, _ := s.DC.Image().At(x, y).RGBA()
						// image is white, which is 65535 in all channels
						if cr == 65535 && cg == 65535 && cb == 65535 {
							s.DC.SetPixel(x, y)
							seed.grew = true
						}
					}
				}
			}

			seed.r++
			s.Seeds[i] = seed
		}
		// check if all seeds have stopped growing
		someSeedsContinuing := false
		for _, seed := range s.Seeds {
			if seed.grew == true {
				someSeedsContinuing = true
			}
		}
		if someSeedsContinuing != true {
			k = s.DestHeight * s.DestWidth
		}
	}
}
