package sketch

import (
	"fmt"
	"image"
	"image/color"

	"github.com/fogleman/gg"
	"gitlab.com/ericworkman/generative/util"
)

// StackParams contains user input
type StackParams struct {
	DestWidth  int
	DestHeight int
}

// StackSketch is the canvas and grid wrapper
type StackSketch struct {
	StackParams
	source       image.Image
	DC           *gg.Context
	sourceWidth  int
	sourceHeight int
}

// NewStackSketch creates a stack sketch
func NewStackSketch(source image.Image, params StackParams) *StackSketch {
	fmt.Println("Starting Sketch")

	s := &StackSketch{StackParams: params}
	bounds := source.Bounds()
	s.sourceWidth, s.sourceHeight = bounds.Max.X, bounds.Max.Y
	s.source = source

	// canvas is a gg image context and contains what gets drawn to the screen
	canvas := gg.NewContext(s.DestWidth, s.DestHeight)
	canvas.SetLineWidth(0.0)
	canvas.SetColor(color.White)
	canvas.DrawRectangle(0, 0, float64(s.DestWidth), float64(s.DestHeight))
	canvas.FillPreserve()
	canvas.Stroke()
	s.DC = canvas
	return s
}

// Output returns the canvas as an image
func (s *StackSketch) Output() image.Image {
	return s.DC.Image()
}

// Update performs a single iteration
func (s *StackSketch) Update(i int) {
	limitX := float64(s.DestWidth) / float64(i)
	limitY := float64(s.DestHeight) / float64(i)
	//fmt.Println("limitX limitY", limitX, limitY)

	for x := 0.0; x < float64(s.DestWidth); x += limitX {
		for y := 0.0; y < float64(s.DestHeight); y += limitY {
			r, g, b := util.Rgb255(s.source.At(int(x+limitX/2), int(y+limitY/2)))
			s.DC.SetRGBA255(r, g, b, 255/i)
			s.DC.DrawEllipse(x+limitX/2, y+limitY/2, limitX/2, limitY/2)
			s.DC.FillPreserve()
			s.DC.Stroke()
		}
	}

}
