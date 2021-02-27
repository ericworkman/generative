package sketch

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"

	"github.com/fogleman/gg"
	"github.com/teacat/noire"
)

var (
	andersonColors = [...][3]uint8{
		{224, 105, 99},
		{119, 194, 169},
		{45, 225, 100},
		{135, 174, 99},
		{232, 223, 104},
		{56, 125, 179},
	}
	sky   = [3]int{46, 59, 75}
	water = [3]int{52, 56, 57}
	dark  = [3]int{36, 47, 62}
)

// AndersonParams contains externally-provided parameters
type AndersonParams struct {
	// tweakable parameters for the cli
	DestWidth  int
	DestHeight int
	Iterations int
}

// AndersonSketch wraps all the components needed to draw the spiral sketch
type AndersonSketch struct {
	AndersonParams
	DC          *gg.Context
	currentR    float64
	horizon     int // also max height of each slot
	slot        float64
	slotOffsets [6]int
}

// NewAndersonSketch initializes the canvas and AndersonSketch
// The sketch has a top sky zone and a bottom water zone made of a blue-ish color on top and a gray-blue color on bottom
// Each color fits into a "slot" of equal possible width.
// The edges on the left and right can be different sizes than the rest of the slots.
// The water rectangle mirrors the sky rectangle but gets muted out.
// The sky rectangle has progressively more opaque and vibrant colored squares with paint marks and fades.
// The progression of colors can be offset to either side of the slot.
// The water rectangle contains aberations that look vaguely like waves.
func NewAndersonSketch(params AndersonParams) *AndersonSketch {
	fmt.Println("Starting Sketch")

	s := &AndersonSketch{AndersonParams: params, currentR: 2.0}
	s.horizon = randIntRangeFrom(s.DestHeight/5, s.DestHeight*4/5)
	//s.horizon = 300
	s.slot = float64(s.DestWidth) / 8.0

	// canvas is a gg image context and contains what gets drawn to the screen
	canvas := gg.NewContext(s.DestWidth, s.DestHeight)
	canvas.SetLineWidth(0.0)
	// set sky rectangle
	canvas.SetRGBA255(sky[0], sky[1], sky[2], 255)
	canvas.DrawRectangle(0, 0, float64(s.DestWidth), float64(s.horizon))
	canvas.FillPreserve()
	canvas.Stroke()
	// set water rectangle
	canvas.SetRGBA255(water[0], water[1], water[2], 255)
	canvas.DrawRectangle(0, float64(s.horizon), float64(s.DestWidth), float64(s.DestHeight))
	canvas.FillPreserve()
	canvas.Stroke()
	s.DC = canvas

	rand.Shuffle(len(andersonColors), func(i, j int) {
		andersonColors[i], andersonColors[j] = andersonColors[j], andersonColors[i]
	})

	// slot offsets, 1 for left and -1 for right
	slotOffsets := [6]int{}
	for j := 0; j < len(andersonColors); j++ {
		if rand.Intn(100) > 50 {
			slotOffsets[j] = 1
		} else {
			slotOffsets[j] = -1
		}
	}
	s.slotOffsets = slotOffsets

	return s
}

// Output produces an image output of the current state of the sketch
func (s *AndersonSketch) Output() image.Image {
	return s.DC.Image()
}

// Update makes a logical step into generation
func (s *AndersonSketch) Update(i int) {
	for j := 0; j < len(andersonColors); j++ {
		acolor := andersonColors[j]

		x := s.slot + float64(j)*s.slot
		y := float64(s.horizon)

		maxWidth := float64(s.slot) / float64(i+1)
		nextStepWidth := float64(s.slot) / float64(i+2)
		w := maxWidth
		if i != 0 {
			w = randFloat64RangeFrom(nextStepWidth+(maxWidth-nextStepWidth)/2, maxWidth)
		}
		// push the starting place to the right if selected at initilization
		if s.slotOffsets[j] == -1 {
			x += s.slot - w
		}

		maxHeight := y / float64(i+1)
		nextStepHeight := y / float64(i+2)
		h := randFloat64RangeFrom(nextStepHeight, maxHeight)

		// gradient is two circles: first is the solid color and is the smaller of the two
		// second is the transparent color and is larger
		// solid color is 100% within the first circle and transitions between the smaller and larger circles
		// transparent color is 100% outside the second circle.
		// Ensure the first circle is entirely below the horizon, so that the base is a solid color.
		// Jitter left and right and radius of larger circle for some variation
		grad := gg.NewRadialGradient(x+w/2, y+5, 5, x+w/2+randFloat64Range(5), y+5, h+randFloat64Range(5))

		alpha := minFloat64(0.2+0.2*float64(i), 1.0)
		solid := color.RGBA{}
		solid.R = uint8(alpha * float64(acolor[0]))
		solid.G = uint8(alpha * float64(acolor[1]))
		solid.B = uint8(alpha * float64(acolor[2]))
		solid.A = uint8(alpha * 255)

		transparent := color.RGBA{}
		transparent.R = 0
		transparent.G = 0
		transparent.B = 0
		transparent.A = 0

		grad.AddColorStop(0, solid)
		grad.AddColorStop(1, transparent)

		s.DC.SetFillStyle(grad)
		s.DC.DrawRectangle(x, y, w, -h)
		s.DC.Fill()
		s.DC.Stroke()

		// water mirror
		wgrad := gg.NewRadialGradient(x+w/2, y-5, 5, x+w/2, y-5, h)
		wcolor := noire.NewRGB(float64(acolor[0]), float64(acolor[1]), float64(acolor[2]))
		r, g, b := wcolor.Darken(0.33).RGB()

		walpha := minFloat64(0.2+0.2*float64(i), 1.0)
		wsolid := color.RGBA{}
		wsolid.R = uint8(walpha * float64(r))
		wsolid.G = uint8(walpha * float64(g))
		wsolid.B = uint8(walpha * float64(b))
		wsolid.A = uint8(walpha * 255)

		wtransparent := color.RGBA{}
		wtransparent.R = 0
		wtransparent.G = 0
		wtransparent.B = 0
		wtransparent.A = 0

		wgrad.AddColorStop(0, wsolid)
		wgrad.AddColorStop(1, wtransparent)

		waterScale := 1.4
		s.DC.SetFillStyle(wgrad)
		s.DC.DrawRectangle(x, y, w, h*waterScale)
		s.DC.Fill()
		s.DC.Stroke()

		// if this is the last iteration, we are going to put some full strength color and near black blocks in the front
		// and draw the sides
		if i == s.Iterations {
			// shuffle, we'll care about the first 3 only
			// 0 = darken
			// 1 = brighten
			// 2 = near-black
			// 3 = same
			options := [...]int{0, 1, 2, 2, 2, 2, 3}
			rand.Shuffle(len(options), func(i, j int) {
				options[i], options[j] = options[j], options[i]
			})

			lcolor := noire.NewRGB(float64(acolor[0]), float64(acolor[1]), float64(acolor[2]))
			// reset
			x = s.slot + float64(j)*s.slot
			h := nextStepHeight
			down := y
			wi := math.Round(s.slot / 3)
			he := -h / 2
			hj := 0.0

			for t := 0; t < 3; t++ {
				left := math.Round(x + float64(t)*s.slot/3)
				if t == 1 {
					hj = randFloat64Range(he / 5)
				} else {
					hj = 0.0
				}
				switch options[t] {
				case 0:
					dr, dg, db := lcolor.Darken(0.15).RGB()
					s.DC.SetRGBA255(int(dr), int(dg), int(db), 255)
					s.DC.DrawRectangle(left, down, wi, he-hj)
					s.DC.Fill()
					s.DC.Stroke()
				case 1:
					br, bg, bb := lcolor.Brighten(0.15).RGB()
					s.DC.SetRGBA255(int(br), int(bg), int(bb), 255)
					s.DC.DrawRectangle(left, down, wi, he-hj)
					s.DC.Fill()
					s.DC.Stroke()
				case 2:
					s.DC.SetRGBA255(dark[0], dark[1], dark[2], 255)
					s.DC.DrawRectangle(left, down, wi, he-hj)
					s.DC.Fill()
					s.DC.Stroke()
				case 3:
					nr, ng, nb := lcolor.RGB()
					s.DC.SetRGBA255(int(nr), int(ng), int(nb), 255)
					s.DC.DrawRectangle(left, down, wi, he-hj)
					s.DC.Fill()
					s.DC.Stroke()
				}
			}

			// draw the shapes on the far edges of the sketch
			// left
			s.DC.SetRGBA255(dark[0], dark[1], dark[2], 255)
			s.DC.MoveTo(0, float64(s.horizon))
			s.DC.LineTo(s.slot, float64(s.horizon))
			s.DC.LineTo(s.slot, float64(s.horizon)+he)
			s.DC.LineTo(0, float64(s.horizon)+2*he)
			s.DC.ClosePath()
			s.DC.Fill()
			s.DC.Stroke()
			// right
			s.DC.SetRGBA255(dark[0], dark[1], dark[2], 255)
			s.DC.MoveTo(float64(s.DestWidth)-s.slot, float64(s.horizon))
			s.DC.LineTo(float64(s.DestWidth), float64(s.horizon))
			s.DC.LineTo(float64(s.DestWidth), float64(s.horizon)+2*he)
			s.DC.LineTo(float64(s.DestWidth)-s.slot, float64(s.horizon)+he)
			s.DC.ClosePath()
			s.DC.Fill()
			s.DC.Stroke()

			// add ripples
			ripples := (s.DestHeight - s.horizon) / 100
			for k := 0; k < ripples; k++ {
				s.DC.SetRGBA255(dark[0], dark[1], dark[2], 10)
				s.DC.DrawRectangle(0, float64(s.horizon+100*k)+randFloat64Range(5.0), float64(s.DestWidth), 10+randFloat64Range(3))
				s.DC.Fill()
				s.DC.Stroke()
			}
		}
	}
}
