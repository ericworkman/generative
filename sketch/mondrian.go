package sketch

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"

	"github.com/fogleman/gg"
	"gitlab.com/ericworkman/generative/util"
)

// MondrianParams contains user input
type MondrianParams struct {
	DestWidth  int
	DestHeight int
}

// MondrianSketch is the canvas and grid wrapper
type MondrianSketch struct {
	MondrianParams
	source       image.Image
	DC           *gg.Context
	sourceWidth  int
	sourceHeight int
}

// NewMondrianSketch creates a stack sketch
func NewMondrianSketch(source image.Image, params MondrianParams) *MondrianSketch {
	fmt.Println("Starting Sketch")

	s := &MondrianSketch{MondrianParams: params}
	bounds := source.Bounds()
	s.sourceWidth, s.sourceHeight = bounds.Max.X, bounds.Max.Y
	s.source = source

	// canvas is a gg image context and contains what gets drawn to the screen
	canvas := gg.NewContext(s.DestWidth, s.DestHeight)
	canvas.SetLineWidth(0.0)
	canvas.SetColor(color.White)
	canvas.DrawRectangle(0, 0, float64(s.DestWidth), float64(s.DestHeight))
	canvas.FillPreserve()
	canvas.DrawImage(source, 0, 0)
	canvas.Stroke()
	s.DC = canvas
	return s
}

// Output returns the canvas as an image
func (s *MondrianSketch) Output() image.Image {
	return s.DC.Image()
}

// Update performs a single iteration
func (s *MondrianSketch) Update(i int) {
	rndX := rand.Float64() * float64(s.sourceWidth)
	rndY := rand.Float64() * float64(s.sourceHeight)
	r, g, b := util.Rgb255(s.source.At(int(rndX), int(rndY)))

	destX := rndX * float64(s.DestWidth) / float64(s.sourceWidth)
	destY := rndY * float64(s.DestHeight) / float64(s.sourceHeight)

	size := 0.01*float64(s.sourceWidth) + rand.Float64()*0.15*float64(s.sourceWidth)

	// black border
	s.DC.SetRGBA255(0, 0, 0, 255)
	s.DC.DrawRegularPolygon(4, destX, destY, size+10, 0)
	s.DC.FillPreserve()
	s.DC.Stroke()

	s.DC.SetRGBA255(r, g, b, 255)
	s.DC.DrawRegularPolygon(4, destX, destY, size, 0)
	s.DC.FillPreserve()
	s.DC.Stroke()
}
