package sketch

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"

	"github.com/fogleman/gg"
	"github.com/teacat/noire"
	"gitlab.com/ericworkman/generative/util"
)

var (
	sunColors = [...][3]int{
		{233, 168, 6},
	}
)

// SunParams contains externally-provided parameters
type SunParams struct {
	// tweakable parameters for the cli
	DestWidth  int
	DestHeight int
	SunRadius  float64
	LineWidth  float64
}

// SunSketch wraps all the components needed to draw the sketch
type SunSketch struct {
	SunParams
	DC *gg.Context
}

// NewSunSketch initializes the canvas and SunSketch
func NewSunSketch(params SunParams) *SunSketch {
	fmt.Println("Starting Sketch")

	s := &SunSketch{SunParams: params}

	// canvas is a gg image context and contains what gets drawn to the screen
	canvas := gg.NewContext(s.DestWidth, s.DestHeight)
	canvas.SetLineWidth(0.0)
	canvas.SetColor(color.White)
	canvas.DrawRectangle(0, 0, float64(s.DestWidth), float64(s.DestHeight))
	canvas.FillPreserve()
	canvas.Stroke()
	s.DC = canvas

	s.DC.SetLineWidth(s.LineWidth)

	return s
}

// Output produces an image output of the current state of the sketch
func (s *SunSketch) Output() image.Image {
	return s.DC.Image()
}

// Draw completes the drawing
func (s *SunSketch) Draw() {
	x := float64(s.DestWidth / 2)
	y := float64(s.DestHeight / 2)

	sun := noire.NewRGB(233, 168, 6)
	sun = sun.Tint(util.RandFloat64RangeFrom(-0.2, 0.2))
	sr, sg, sb := sun.RGB()
	s.DC.SetRGB255(int(sr), int(sg), int(sb))
	s.DC.DrawCircle(x, y, s.SunRadius)
	s.DC.Fill()
	s.DC.Stroke()

	skyColor := noire.NewRGB(29, 103, 131)

	for r := s.SunRadius + 1.5*s.LineWidth; r <= math.Sqrt((x*x)+(y*y)); r += (s.LineWidth * 2) {
		offset := util.RandFloat64RangeFrom(0, 1.0)
		start := 0.0
		distance := 0.0
		end := 0.0
		gap := util.MaxFloat64(s.LineWidth/5, 1.5) / r

		for i := offset; i < 1.0+offset; i = end {
			start = i
			distance = util.RandFloat64RangeFrom(i, util.MinFloat64(1.0+offset-i, 0.45))
			end = util.MinFloat64(start+distance, 1.0+offset)

			chosen := skyColor
			chance := rand.Intn(100)
			if chance < 50 {
				chosen = chosen.Tint(util.RandFloat64RangeFrom(0, 0.4))
			} else {
				chosen = chosen.Shade(util.RandFloat64RangeFrom(0, 0.25))
			}
			re, g, b := chosen.RGB()

			s.DC.SetRGB255(int(re), int(g), int(b))

			s.DC.DrawArc(x, y, r, (start+gap)*2*math.Pi, end*2*math.Pi)
			s.DC.Stroke()
		}

	}

}
