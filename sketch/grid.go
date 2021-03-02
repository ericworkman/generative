package sketch

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/fogleman/gg"
	"gitlab.com/ericworkman/generative/util"
)

// GridParams contains externally-provided parameters
type GridParams struct {
	// tweakable parameters for the cli
	DestWidth  int
	DestHeight int
	Vignette   bool
	Size       float64
}

// GridSketch wraps all the components needed to draw the sketch
type GridSketch struct {
	GridParams
	DC           *gg.Context
	source       image.Image
	sourceWidth  int
	sourceHeight int
}

// NewGridSketch initializes the canvas and GridSketch
func NewGridSketch(source image.Image, params GridParams) *GridSketch {
	fmt.Println("Starting Sketch")

	s := &GridSketch{GridParams: params}
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
	s.DC = canvas

	s.DC.SetLineWidth(0.0)
	return s
}

// Draw completes the drawing
func (s *GridSketch) Draw() {
	spacing := s.Size
	for x := spacing; x < float64(s.sourceWidth); x += spacing {
		alpha := 255.0
		if s.Vignette {
			alpha = 255 - math.Abs(255.0*(float64(s.sourceWidth/2)-x)/float64(s.sourceWidth/2))
		}
		for y := spacing; y < float64(s.sourceHeight); y += spacing {
			r, g, b := util.Rgb255(s.source.At(int(x), int(y)))
			s.DC.SetRGBA255(r, g, b, int(alpha))
			s.DC.DrawCircle(float64(x), float64(y), float64(spacing/2))
			s.DC.FillPreserve()
			s.DC.Stroke()
		}
	}

}

// Output produces an image output of the current state of the sketch
func (s *GridSketch) Output() image.Image {
	return s.DC.Image()
}
