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

// CrawlParams contains externally-provided parameters
type CrawlParams struct {
	// tweakable parameters for the cli
	DestWidth  int
	DestHeight int
	Iterations int
	Count      int
	Start      string
}

// CrawlSketch wraps all the components needed to draw the sketch
type CrawlSketch struct {
	CrawlParams
	DC       *gg.Context
	crawlers []crawler
}

type point struct {
	x float64
	y float64
}

type crawler struct {
	start      point
	current    point
	history    []point
	theta      float64
	thetaRange float64
	r          float64
	c          noire.Color
	light      noire.Color
}

func (c *crawler) crawl(s *CrawlSketch) {
	if c.current.x >= 0 && c.current.x < float64(s.DestWidth) && c.current.y >= 0 && c.current.y < float64(s.DestHeight) {
		awayAngle := c.theta + util.RandFloat64Range(c.thetaRange)

		xx1 := c.r * math.Cos(awayAngle)
		yy1 := c.r * math.Sin(awayAngle)

		current := point{c.current.x + xx1, c.current.y + yy1}
		c.current = current
		c.history = append(c.history, current)
	}
}

func (s *CrawlSketch) addCrawler(start string) {
	x := float64(s.DestWidth / 2)
	y := float64(s.DestHeight / 2)
	theta := util.RandFloat64RangeFrom(0, 2*math.Pi)
	thetaRange := 2 * math.Pi / 3
	if start == "corner" {
		x = 10.0
		y = 10.0
		theta = util.RandFloat64RangeFrom(0.05, math.Pi/2.05)
		thetaRange = math.Pi / 2
	}

	r := 5.0

	xx := x + r*math.Cos(theta)
	yy := y + r*math.Sin(theta)

	c := noire.NewRGBA(rand.Float64()*128, rand.Float64()*128, 128+rand.Float64()*127, 1)
	lightc := c.Lighten(.35)

	crawly := crawler{start: point{xx, yy}, current: point{xx, yy}, theta: theta, thetaRange: thetaRange, r: r, history: []point{{x: xx, y: yy}}, c: c, light: lightc}
	s.crawlers = append(s.crawlers, crawly)
}

// NewCrawlSketch initializes the canvas and CrawlSketch
func NewCrawlSketch(params CrawlParams) *CrawlSketch {
	fmt.Println("Starting Sketch")

	s := &CrawlSketch{CrawlParams: params}

	// canvas is a gg image context and contains what gets drawn to the screen
	canvas := gg.NewContext(s.DestWidth, s.DestHeight)
	canvas.SetLineWidth(0.0)
	canvas.SetColor(color.White)
	canvas.DrawRectangle(0, 0, float64(s.DestWidth), float64(s.DestHeight))
	canvas.FillPreserve()
	canvas.Stroke()
	s.DC = canvas

	s.DC.SetLineWidth(1.0)

	for i := 0; i < s.Count; i++ {
		s.addCrawler(s.Start)
	}

	return s
}

// Output produces an image output of the current state of the sketch
func (s *CrawlSketch) Output() image.Image {
	//draw all background lines, then all foreground lines
	for j := 0; j < len(s.crawlers); j++ {
		crawly := s.crawlers[j]
		r, g, b := crawly.light.RGB()
		s.DC.SetRGBA255(int(r), int(g), int(b), 15)
		for k := 0; k < len(crawly.history); k++ {
			p := crawly.history[k]
			s.DC.DrawLine(crawly.start.x, crawly.start.y, p.x, p.y)
			s.DC.Stroke()
		}
	}

	for j := 0; j < len(s.crawlers); j++ {
		crawly := s.crawlers[j]
		prevX := crawly.start.x
		prevY := crawly.start.y
		r, g, b := crawly.c.RGB()
		s.DC.SetRGBA255(int(r), int(g), int(b), 255)
		for k := 0; k < len(crawly.history); k++ {
			p := crawly.history[k]
			s.DC.DrawLine(prevX, prevY, p.x, p.y)
			s.DC.Stroke()
			prevX = p.x
			prevY = p.y
		}
	}

	return s.DC.Image()
}

// Update makes a logical step into generation
func (s *CrawlSketch) Update(i int) {

	for j := 0; j < len(s.crawlers); j++ {
		crawly := s.crawlers[j]
		crawly.crawl(s)
		s.crawlers[j] = crawly
	}

}
