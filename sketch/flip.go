package sketch

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
	"math/rand"

	"github.com/fogleman/gg"
)

// FlipParams contains user input
type FlipParams struct {
	DestWidth  int
	DestHeight int
}

// FlipSketch is the canvas and grid wrapper
type FlipSketch struct {
	FlipParams
	source       image.Image
	DC           *gg.Context
	sourceWidth  int
	sourceHeight int
	xOffset      float64
	yOffset      float64
}

// NewFlipSketch creates a stack sketch
func NewFlipSketch(source image.Image, params FlipParams) *FlipSketch {
	fmt.Println("Starting Sketch")

	s := &FlipSketch{FlipParams: params}
	bounds := source.Bounds()
	s.sourceWidth, s.sourceHeight = bounds.Max.X, bounds.Max.Y
	s.source = source
	s.xOffset = float64(s.sourceWidth) / 2
	s.yOffset = float64(s.sourceHeight) / 2

	// canvas is a gg image context and contains what gets drawn to the screen
	canvas := gg.NewContext(s.DestWidth*2, s.DestHeight*2)
	canvas.SetLineWidth(0.0)
	canvas.SetColor(color.White)
	canvas.DrawRectangle(s.xOffset, s.yOffset, float64(s.DestWidth), float64(s.DestHeight))
	canvas.FillPreserve()
	canvas.Stroke()
	s.DC = canvas
	return s
}

// Output trims and returns the canvas as an image
func (s *FlipSketch) Output() image.Image {
	src := s.DC.Image()
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)
	draw.Draw(dst, bounds, src, bounds.Min, draw.Src)

	return dst.SubImage(image.Rect(int(s.xOffset), int(s.yOffset), int(s.xOffset)+s.DestWidth, int(s.yOffset)+s.DestHeight))
}

// Draw performs the algorithm on the image
func (s *FlipSketch) Draw(divisions int) {
	r := float64(s.DestHeight) / float64(divisions)
	maxRows := int(math.Ceil(float64(s.DestHeight) / (r / 1.4142)))
	maxCols := int(math.Ceil(float64(s.DestWidth) / (r * 1.4142)))
	//fmt.Println(r, maxRows, maxCols)
	rot := math.Pi / 4
	row := 0

	// draw on a bigger workspace than necessary and over draw the bottom and right most edges to make things fit
	for y := s.yOffset; y < float64(s.DestHeight)+s.yOffset+r/1.4142; y += r / 1.4142 {
		fmt.Println("row", row)
		col := 0

		// for debugging, draw markers for the row
		if false {
			s.DC.SetLineWidth(10.0)
			s.DC.DrawLine(s.xOffset, y, s.xOffset+float64(s.DestWidth), y)
			s.DC.Stroke()
			s.DC.SetLineWidth(0.0)
		}

		for x := s.xOffset + 0.7071*r*float64(row%2); x < float64(s.DestWidth)+s.xOffset+r; x += r * 1.4142 {
			// For each cell, draw the shape and clip it so we only paint the image onto that cell.
			// Some cells get rotated 180 degrees, some get a white border, but never the edge cells.
			// This whole thing is a little complicated and pretty slow. I'm probably doing this really inefficiently,
			// but I can't quite grasp a better way of doing this.
			// This uses subcontexts from gg to restrict operations to smaller sections of the canvas.
			s.DC.Push()
			s.DC.DrawRegularPolygon(4, x, y, r/1.4142, rot)
			s.DC.Clip()
			s.DC.Push()

			if (row > 2) && (row < maxRows-1) && (col > 2) && (col < maxCols-1) && (rand.Intn(100) < 25) {
				s.DC.RotateAbout(math.Pi, x, y)
			}

			s.DC.DrawImage(s.source, int(s.xOffset), int(s.yOffset))
			s.DC.Fill()
			s.DC.Stroke()

			if (row > 2) && (row < maxRows-1) && (col > 2) && (col < maxCols-1) && (rand.Intn(100) < 3) {
				s.DC.Push()
				s.DC.SetLineWidth(10.0)
				//s.DC.DrawRegularPolygon(3, x, y, r, rot)
				s.DC.DrawRegularPolygon(4, x, y, r/1.4142, rot)
				//s.DC.DrawCircle(x, y, 10)
				s.DC.Stroke()
				s.DC.SetLineWidth(0.0)
				s.DC.Pop()
			}

			// for debugging, mark the center of each cell
			if false {
				s.DC.Push()
				s.DC.SetLineWidth(10.0)
				s.DC.DrawCircle(x, y, 10)
				s.DC.Stroke()
				s.DC.SetLineWidth(0.0)
				s.DC.Pop()
			}

			s.DC.Pop()

			s.DC.ResetClip()
			s.DC.Pop()
			col++
		}
		row++
	}

}
