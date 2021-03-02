package util

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"net/http"
	"os"
)

// LoadUnsplashImage either fetches a random image from Unsplash and prints the url of it,
// or it takes the provided url and loads it
func LoadUnsplashImage(width, height int, url string) (image.Image, error) {
	if url == "" {
		url = fmt.Sprintf("https://source.unsplash.com/random/%dx%d", width, height)
	}

	req, _ := http.NewRequest("GET", url, nil)
	var lastURLQuery string

	client := new(http.Client)
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) > 10 {
			return errors.New("too many redirects")
		}
		lastURLQuery = req.URL.String()
		return nil
	}

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	fmt.Println(lastURLQuery)

	img, _, err := image.Decode(res.Body)
	return img, err
}

// SaveOutput writes an image to a file
func SaveOutput(img image.Image, filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	//Encode and Save
	err = png.Encode(f, img)
	if err != nil {
		return err
	}

	return nil
}

// Rgb255 converts a color.Color to r, g, b 0-255
func Rgb255(c color.Color) (r, g, b int) {
	r0, g0, b0, _ := c.RGBA()
	return int(r0 / 257), int(g0 / 257), int(b0 / 257)
}

// RandRange returns an int between -max and max
func RandRange(max int) int {
	return -max + rand.Intn(2*max)
}

// RandFloat64Range returns a float64 between -max and max
func RandFloat64Range(max float64) float64 {
	return -max + rand.Float64()*2*max
}

// RandFloat64RangeFrom returns a float64 between min and max
func RandFloat64RangeFrom(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// RandIntRangeFrom returns an int between min and max
func RandIntRangeFrom(min, max int) int {
	return min + rand.Intn(max-min)
}

// MaxInt returns the larger of two ints
func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// MinInt returns the smaller of two ints
func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// MaxFloat64 returns the larger of two float64s
func MaxFloat64(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// MinFloat64 returns the smaller of two float64s
func MinFloat64(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
