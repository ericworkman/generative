package sketch

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"

	"github.com/fogleman/gg"
)

var (
	spiralColors = [...][3]int{
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

// SpiralParams contains externally-provided parameters
type SpiralParams struct {
	// tweakable parameters for the cli
	DestWidth  int
	DestHeight int
	Iterations int
	// made up parameter names for a logarithmic spiral parameters
	// see https://www.wolframalpha.com/input/?i=parametric+plot+%281%2Be%5E%280.1+t%29sin+t%2C+1%2Be%5E%280.1t%29cos+t%29+for+t%3D-20+to+10
	Beta float64
	Mu   float64
}

// SpiralSketch wraps all the components needed to draw the spiral sketch
type SpiralSketch struct {
	SpiralParams
	DC       *gg.Context
	currentR float64
}

// NewSpiralSketch initializes the canvas and SpiralSketch
func NewSpiralSketch(params SpiralParams) *SpiralSketch {
	fmt.Println("Starting Sketch")

	s := &SpiralSketch{SpiralParams: params, currentR: 2.0}

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

// Output produces an image output of the current state of the sketch
func (s *SpiralSketch) Output() image.Image {
	return s.DC.Image()
}

// Update makes a logical step into generation
func (s *SpiralSketch) Update(i int) {
	centerX := float64(s.DestWidth) / 2.0
	centerY := float64(s.DestHeight) / 2.0

	// offset controls how many get skipped in the tight center
	// growth controls how close the pattern stays
	j := math.Pi*15 + float64(i)/math.Pi
	x := centerX + s.Beta*math.Exp(j*s.Mu)*math.Cos(j)
	y := centerY + s.Beta*math.Exp(j*s.Mu)*math.Sin(j)

	color := spiralColors[rand.Intn(len(spiralColors))]
	s.DC.SetRGBA255(color[0], color[1], color[2], i)
	// logistic growth of radius, barely noticable in practice I think
	s.currentR += 0.006 * float64(i) * float64(s.Iterations-i) / float64(s.Iterations)
	s.DC.DrawCircle(x, y, s.currentR)
	s.DC.FillPreserve()
	s.DC.Stroke()
}
