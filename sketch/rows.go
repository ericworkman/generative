package sketch

import (
	"fmt"
	"image"
	"image/color"

	"github.com/fogleman/gg"
	"gitlab.com/ericworkman/generative/util"
)

// RowsParams contains externally-provided parameters
type RowsParams struct {
	// tweakable parameters for the cli
	DestWidth  int
	DestHeight int
	Vignette   bool
	Size       float64
}

// RowsSketch wraps all the components needed to draw the sketch
type RowsSketch struct {
	RowsParams
	DC           *gg.Context
	source       image.Image
	sourceWidth  int
	sourceHeight int
}

// NewRowsSketch initializes the canvas and RowsSketch
func NewRowsSketch(source image.Image, params RowsParams) *RowsSketch {
	fmt.Println("Starting Sketch")

	s := &RowsSketch{RowsParams: params}
	bounds := source.Bounds()
	s.sourceWidth, s.sourceHeight = bounds.Max.X, bounds.Max.Y
	s.source = source

	// canvas is a gg image context and contains what gets drawn to the screen
	canvas := gg.NewContext(s.DestWidth, s.DestHeight)
	canvas.SetLineWidth(0.0)
	canvas.SetColor(color.Black)
	canvas.DrawRectangle(0, 0, float64(s.DestWidth), float64(s.DestHeight))
	canvas.FillPreserve()
	canvas.Stroke()
	canvas.SetLineCapRound()
	s.DC = canvas

	s.DC.SetLineWidth(0.0)
	return s
}

// Draw completes the drawing
func (s *RowsSketch) Draw() {
	alpha := 200.0

	spacing := s.Size
	iteration := 0.0
	for x := spacing; x < (float64(s.sourceWidth)-spacing)/2; x += spacing {
		iteration += 2.0
		endx := float64(s.DestWidth) - iteration*spacing
		for y := spacing; y < float64(s.sourceHeight)-spacing; y += spacing {
			r, g, b := util.Rgb255(s.source.At(int(x), int(y)))
			s.DC.SetRGBA255(r, g, b, int(alpha))
			s.DC.DrawRoundedRectangle(x, y, endx, spacing, spacing/4)
			//s.DC.DrawRectangle(x, y, endx, spacing)
			s.DC.FillPreserve()
			s.DC.Stroke()
		}
	}

}

// Output produces an image output of the current state of the sketch
func (s *RowsSketch) Output() image.Image {
	return s.DC.Image()
}
