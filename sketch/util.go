package sketch

import (
  "errors"
  "fmt"
  "image"
  "image/color"
  "image/png"
  "net/http"
  "os"
  "math/rand"
)

func LoadUnsplashImage(width, height int, url string) (image.Image,error) {
  if url == "" {
    url = fmt.Sprintf("https://source.unsplash.com/random/%dx%d", width, height)
  }

  req, _ := http.NewRequest("GET", url, nil)
  var lastUrlQuery string

  client := new(http.Client)
  client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
    if len(via) > 10 {
      return errors.New("too many redirects")
    }
    lastUrlQuery = req.URL.String()
    return nil
  }

  res, err := client.Do(req)

  if err != nil {
    return nil, err
  }
  fmt.Println(lastUrlQuery)

  img, _, err := image.Decode(res.Body)
  return img, err
}

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

func rgb255(c color.Color) (r, g, b int) {
  r0, g0, b0, _ := c.RGBA()
  return int(r0 / 257), int(g0 / 257), int(b0/ 257)
}

func randRange(max int) int {
  return -max + rand.Intn(2 * max)
}

func randFloat64Range(max float64) float64 {
  return -max + rand.Float64() * 2 * max
}

func randFloat64RangeFrom(min, max float64) float64 {
  return min + rand.Float64() * max
}
