package sketch

import (
	"image"
	"image/color"
	"math"
	"math/rand"

	"github.com/fogleman/gg"
)

type LayerParams struct {
	DestWidth              int
	DestHeight             int
	PathRatio              float64
	PathReduction          float64
	PathMin                float64
	PathJitter             int
	InitialAlpha           float64
	AlphaIncrease          float64
	MinEdgeCount           int
	MaxEdgeCount           int
	Edge                   bool
	PathInversionThreshold float64
}

type LayerSketch struct {
	LayerParams
	source          image.Image
	DC              *gg.Context
	sourceWidth     int
	sourceHeight    int
	InitialPathSize float64
	PathSize        float64
}

func NewLayerSketch(source image.Image, layerParams LayerParams) *LayerSketch {
	s := &LayerSketch{LayerParams: layerParams}
	bounds := source.Bounds()
	s.sourceWidth, s.sourceHeight = bounds.Max.X, bounds.Max.Y
	s.PathSize = s.PathRatio * float64(s.DestWidth)
	s.InitialPathSize = s.PathSize

	canvas := gg.NewContext(s.DestWidth, s.DestHeight)
	canvas.SetColor(color.Black)
	canvas.DrawRectangle(0, 0, float64(s.sourceWidth), float64(s.sourceHeight))
	canvas.FillPreserve()

	s.source = source
	s.DC = canvas
	return s
}

func (s *LayerSketch) Output() image.Image {
	return s.DC.Image()
}

func (s *LayerSketch) Update() {
	rndX := rand.Float64() * float64(s.sourceWidth)
	rndY := rand.Float64() * float64(s.sourceHeight)
	r, g, b := rgb255(s.source.At(int(rndX), int(rndY)))

	destX := rndX * float64(s.DestWidth) / float64(s.sourceWidth)
	destX += float64(randRange(s.PathJitter))
	destY := rndY * float64(s.DestHeight) / float64(s.sourceHeight)
	destY += float64(randRange(s.PathJitter))

	s.DC.SetRGBA255(r, g, b, int(s.InitialAlpha))
	edges := s.MinEdgeCount + rand.Intn(s.MaxEdgeCount-s.MinEdgeCount+1)
	if edges < 2 {
		s.DC.DrawCircle(destX, destY, s.PathSize)
		s.DC.FillPreserve()
	} else if edges == 2 {
		s.DC.SetLineWidth(10.00)
		randAngle := rand.Float64() * float64(360)
		s.DC.DrawLine(destX, destY, destX+s.PathSize*math.Cos(randAngle), destY+s.PathSize*math.Sin(randAngle))
		s.DC.StrokePreserve()
	} else {
		s.DC.DrawRegularPolygon(edges, destX, destY, s.PathSize, rand.Float64())
		s.DC.FillPreserve()
	}

	if s.Edge && s.PathSize <= s.PathInversionThreshold*s.InitialPathSize {
		if (r+g+b)/3 < 128 {
			s.DC.SetRGBA255(255, 255, 255, int(s.InitialAlpha*2))
		} else {
			s.DC.SetRGBA255(0, 0, 0, int(s.InitialAlpha*2))
		}
	}

	s.DC.Stroke()

	s.PathSize -= s.PathReduction * s.PathSize
	s.InitialAlpha += s.AlphaIncrease
}
