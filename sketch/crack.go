package sketch

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"

	"github.com/fogleman/gg"
	"gitlab.com/ericworkman/generative/util"
)

const (
	blankAngle = 10001
)

var (
	crackColors = [...][3]int{
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

// CrackParams contains user input options
type CrackParams struct {
	// tweakable parameters for the cli
	DestWidth      int
	DestHeight     int
	CrackLimit     int
	Seeds          int
	StartingCracks int
}

// CrackSketch contains a canvas, a grid, a set of cracks, and some other information
type CrackSketch struct {
	// This is ported almost directly from Jared Tarbell's Processing sketch "Substrate".
	// See details at http://www.complexification.net/gallery/machines/substrate/
	// There are a few spots with a couple changes mostly to fit into golang and gg.
	// This hasn't been optimized and very likely has bugs, but it does produce nice results.
	CrackParams
	DC       *gg.Context
	GridSize int
	Grid     []int
	cracks   []crack
}

type crack struct {
	// the dark lines that look like cracks
	X  float64
	Y  float64
	T  float64 // direction in degrees
	SP sandPainter
}

// Grow a crack it its direction t. Color to the side of it some distance using the sandPainter.
func (c *crack) Move(sketch *CrackSketch) {
	c.X += 0.42 * math.Cos(c.T*math.Pi/180)
	c.Y += 0.42 * math.Sin(c.T*math.Pi/180)

	// bound check
	z := 0.25
	cx := int(c.X + util.RandFloat64Range(z))
	cy := int(c.Y + util.RandFloat64Range(z))

	// draw sand painter
	c.RegionColor(sketch)

	// draw black crack
	sketch.DC.SetRGBA255(0, 0, 0, 180)

	// TODO: replace jitter
	x := int(c.X + util.RandFloat64Range(z))
	y := int(c.Y + util.RandFloat64Range(z))
	sketch.DC.SetPixel(x, y)

	sketch.DC.Stroke()

	if (cx >= 0) && (cy >= 0) && (cx < sketch.DestWidth) && (cy < sketch.DestHeight) {
		// within bounds of canvas

		if (sketch.Grid[cy*sketch.DestWidth+cx] > 10000) || (math.Abs(float64(sketch.Grid[cy*sketch.DestWidth+cx])-c.T) < 5.0) {
			// continue growing
			sketch.Grid[cy*sketch.DestWidth+cx] = int(c.T)

		} else if math.Abs(float64(sketch.Grid[cy*sketch.DestWidth+cx])-c.T) > 2.0 {
			// found a different crack, so this crack ends
			c.findStart(sketch)
			makecrack(sketch)
		}

	} else {
		// out of bounds, stop cracking
		c.findStart(sketch)
		makecrack(sketch)
	}
}

func (c *crack) findStart(sketch *CrackSketch) {
	// pick random spots on the canvas until a crack is found
	// a crack is any cell on the grid with a degree value between -360 and 360, or really less than the blank value
	// limit the number of times this attempts to find a crack with the timeout
	var px, py int
	found := false
	timeout := 0
	for ok := true; ok; ok = ((found == false) && (timeout <= 10000)) {
		timeout++
		px = rand.Intn(sketch.DestWidth)
		py = rand.Intn(sketch.DestHeight)
		if sketch.Grid[py*sketch.DestWidth+px] < 10000 {
			found = true
		}
	}

	if found == true {
		// found a starting point, so now pick a perpendicular angle to the existing crack angle
		// we add some angle jitter here too for interest
		a := sketch.Grid[py*sketch.DestWidth+px]
		if rand.Intn(100) < 50 {
			a -= 90 + util.RandRange(3)
		} else {
			a += 90 + util.RandRange(3)
		}
		c.T = float64(a)
		c.X = float64(px) // + 0.61 * math.Cos(crack.T * math.Pi / 180)
		c.Y = float64(py) // + 0.61 * math.Sin(crack.T * math.Pi / 180)
		c.SP = newsandPainter()
	}
}

func makecrack(sketch *CrackSketch) crack {
	crack := crack{}
	// only make a new crack if there's a slot for one
	if len(sketch.cracks) < sketch.CrackLimit {
		crack.findStart(sketch)
		sketch.cracks = append(sketch.cracks, crack)
	}
	return crack
}

// NewCrackSketch sets up the wrapper components
func NewCrackSketch(crackParams CrackParams) *CrackSketch {
	fmt.Println("Starting Sketch")

	s := &CrackSketch{CrackParams: crackParams}

	// the grid is dimensionally the same as the canvas, but contains angles in degrees or a blank value
	cgrid := make([]int, s.DestWidth*s.DestHeight)
	s.GridSize = len(cgrid)
	for i := 0; i < s.GridSize; i++ {
		cgrid[i] = blankAngle
	}

	// preseed some spots in the grid with real angles
	for k := 0; k < s.Seeds; k++ {
		i := rand.Intn(s.DestWidth*s.DestHeight - 1)
		cgrid[i] = rand.Intn(360)
	}

	s.Grid = cgrid

	// start the cracks
	for k := 0; k < s.StartingCracks; k++ {
		makecrack(s)
	}

	// canvas is a gg image context and contains what gets drawn to the screen
	canvas := gg.NewContext(s.DestWidth, s.DestHeight)
	canvas.SetColor(color.White)
	canvas.DrawRectangle(0, 0, float64(s.DestWidth), float64(s.DestHeight))
	canvas.FillPreserve()
	s.DC = canvas
	return s
}

// Output creates the image from the canvas
func (s *CrackSketch) Output() image.Image {
	return s.DC.Image()
}

// Update iterates through all of the cracks
func (s *CrackSketch) Update() {
	// allow the cracks to grow a step
	for i := 0; i < len(s.cracks); i++ {
		s.cracks[i].Move(s)
	}
}

type sandPainter struct {
	// creates transparent "grains of sands" perpendicular to the crack with a lot of variation
	// contains color components and a grain size
	R         int
	G         int
	B         int
	GrainSize float64
}

func newsandPainter() sandPainter {
	// aim for desert colors, a slight departure from Tarbell's
	// Tarbell's version takes colors from an image, while this one selects from a predefined list of colors
	color := crackColors[rand.Intn(len(crackColors))]
	sp := sandPainter{R: color[0], G: color[1], B: color[2], GrainSize: util.RandFloat64RangeFrom(0.01, 0.01)}
	return sp
}

func (sp *sandPainter) render(s *CrackSketch, x, y, ox, oy float64) {
	// modulate gain, clamping it between 0 and 1.0
	sp.GrainSize += util.RandFloat64Range(0.050)
	maxg := 1.0
	if sp.GrainSize < 0 {
		sp.GrainSize = 0
	}
	if sp.GrainSize > maxg {
		sp.GrainSize = maxg
	}

	// proportion grain count for smoothness
	grains := int(math.Sqrt(float64((ox-x)*(ox-x) + (oy-y)*(oy-y))))

	// draw the sand grains
	w := sp.GrainSize / float64(grains-1)
	for i := 0; i < grains; i++ {
		a := 7
		x := ox + (x-ox)*math.Sin(math.Sin(float64(i)*w))
		y := oy + (y-oy)*math.Sin(math.Sin(float64(i)*w))
		s.DC.SetRGBA255(sp.R, sp.G, sp.B, a)
		s.DC.DrawPoint(x, y, 0.6)
		s.DC.Stroke()
	}
}

func (c *crack) RegionColor(s *CrackSketch) {
	// find the open region that can be colored that's perpendicular to the crack at the new pixel
	// we use the boundary of this open space to determine how to draw the sand
	rx := c.X
	ry := c.Y

	openspace := true
	for ok := true; ok; ok = openspace {
		// move perpendicular to crack
		rx += 0.81 * math.Sin(c.T*math.Pi/180)
		ry -= 0.81 * math.Cos(c.T*math.Pi/180)
		cx := int(rx)
		cy := int(ry)

		// limit the maximum size of the region to be within the bounds of the canvas and only a percent of the dimensions
		if (cx >= 0) && (cy >= 0) && (cx < s.DestWidth) && (cy < s.DestHeight) && (int(math.Abs(float64(cx)-c.X)) < s.DestWidth/10) && (int(math.Abs(float64(cy)-c.Y)) < s.DestHeight/10) {
			if s.Grid[cy*s.DestWidth+cx] <= 10000 {
				openspace = false
			}
		} else {
			openspace = false
		}
	}

	c.SP.render(s, rx, ry, c.X, c.Y)
}
