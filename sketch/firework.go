package sketch

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"

	"github.com/fogleman/gg"
	"gitlab.com/ericworkman/generative/util"
)

// FireworkParams contains externally-provided parameters
type FireworkParams struct {
	// tweakable parameters for the cli
	DestWidth  int
	DestHeight int
	Iterations int
}

// FireworkSketch wraps all the components needed to draw the firework sketch
type FireworkSketch struct {
	FireworkParams
	DC    *gg.Context
	slope float64
	x1    int
}

// NewFireworkSketch initializes the canvas and FireworkSketch
func NewFireworkSketch(params FireworkParams) *FireworkSketch {
	fmt.Println("Starting Sketch")

	s := &FireworkSketch{FireworkParams: params}

	// Draw a line from some middle point on the left to the inverse point on the right
	// all bursts will be below this line, save for the offsets
	startX := int(util.RandFloat64RangeFrom(0.33*float64(params.DestHeight), 0.67*float64(params.DestHeight)))
	s.slope = float64(s.DestHeight-2*startX) / float64(s.DestWidth)
	s.x1 = startX
	//fmt.Println("y = (", s.slope, ") * x + ", s.x1)

	// canvas is a gg image context and contains what gets drawn to the screen
	canvas := gg.NewContext(s.DestWidth, s.DestHeight)
	canvas.SetLineWidth(0.0)
	canvas.SetColor(color.Black)
	canvas.DrawRectangle(0, 0, float64(s.DestWidth), float64(s.DestHeight))
	canvas.FillPreserve()
	canvas.Stroke()
	s.DC = canvas
	return s
}

// Output produces an image output of the current state of the sketch
func (s *FireworkSketch) Output() image.Image {
	return s.DC.Image()
}

// Update makes a logical step into generation
func (s *FireworkSketch) Update(i int) {
	rndX := rand.Float64() * float64(s.DestWidth)
	rndY := util.RandFloat64RangeFrom(s.slope*rndX+float64(s.x1), float64(s.DestHeight))

	// burst
	color := [3]int{253, 255, 240}
	offsets := [...][2]float64{
		{0, 0},
		{-0.75, -0.5},
		{-0.75, 0.5},
		{0.75, 0.5},
		{0.75, -0.5},

		{-0.25, 1.5},
		{0.25, 1.5},

		{0, 2},
		{0, 3},
		{0, 5},
		{0, 8},
	}
	scale := 255 * i / s.Iterations
	alpha := util.MaxInt(10, scale-50)
	radius := util.MaxFloat64(7, float64(60-3*scale/4))
	//fmt.Println(scale, alpha, radius)

	for _, offset := range offsets {
		s.DC.SetRGBA255(color[0], color[1], color[2], alpha)
		s.DC.DrawCircle(rndX+offset[0]*radius, rndY+offset[1]*radius, radius)
		s.DC.FillPreserve()
		s.DC.Stroke()
	}
}
