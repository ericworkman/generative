package sketch

import (
	"fmt"
	"image"
	"image/color"

	"github.com/fogleman/gg"
)

var (
	stack_colors = [...][3]int{
		{172, 68, 6},
		{201, 148, 89},
		{128, 44, 8},
		{154, 135, 109},
		{215, 207, 185},
		{79, 68, 59},
		{244, 172, 68},
		{234, 204, 147},
		{59, 44, 28},
		{61, 62, 68},
		{221, 89, 64},
		{252, 180, 140},
		{96, 40, 28},
		{160, 92, 92},
	}
)

type StackParams struct {
	// tweakable parameters for the cli
	DestWidth  int
	DestHeight int
}

type StackSketch struct {
	// canvas and grid wrapper
	StackParams
	source       image.Image
	DC           *gg.Context
	sourceWidth  int
	sourceHeight int
}

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
	s.DC = canvas
	return s
}

func (s *StackSketch) Output() image.Image {
	return s.DC.Image()
}

func (s *StackSketch) Update(i int) {
	limitX := float64(s.DestWidth) / float64(i)
	limitY := float64(s.DestHeight) / float64(i)
	//fmt.Println("limitX limitY", limitX, limitY)

	for x := 0.0; x < float64(s.DestWidth); x += limitX {
		for y := 0.0; y < float64(s.DestHeight); y += limitY {
			r, g, b := rgb255(s.source.At(int(x+limitX/2), int(y+limitY/2)))
			s.DC.SetRGBA255(r, g, b, 255/i)
			s.DC.DrawEllipse(x+limitX/2, y+limitY/2, limitX/2, limitY/2)
			s.DC.FillPreserve()
			s.DC.Stroke()
		}
	}

}
