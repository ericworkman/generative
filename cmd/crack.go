package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"image"
	"image/gif"

	"github.com/andybons/gogif"
	"github.com/spf13/cobra"

	"gitlab.com/ericworkman/generative/sketch"
	"gitlab.com/ericworkman/generative/util"
)

var crackCmd = &cobra.Command{
	Use:   "crack",
	Short: "Create sketches in the style of Jared Tarbell",
	Long: `Create a sketch of growing cracks that "crystalize"
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("crack called")

		params := sketch.CrackParams{
			DestWidth:      width,
			DestHeight:     height,
			CrackLimit:     25,
			Seeds:          width/10 + height/10,
			StartingCracks: 5,
		}

		csketch := sketch.NewCrackSketch(params)

		f, err := os.Create("my-image.gif")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		outGif := &gif.GIF{}

		// catch the sigterm signal for ctrl-c quitting mostly
		// save the output at this point
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			util.SaveOutput(csketch.Output(), outputImgName)
			if gifme == true {
				simage := csketch.Output()
				bounds := simage.Bounds()
				palettedImage := image.NewPaletted(bounds, nil)
				quantizer := gogif.MedianCutQuantizer{NumColor: 64}
				quantizer.Quantize(palettedImage, bounds, simage, image.ZP)

				// Add new frame to animated GIF
				outGif.Image = append(outGif.Image, palettedImage)
				outGif.Delay = append(outGif.Delay, 100)
				gif.EncodeAll(f, outGif)
			}
			os.Exit(1)
		}()

		for i := 0; i < limitByIterations; i++ {
			fmt.Println("Iteration", i)

			csketch.Update()

			if (gifme == true) && (i%10 == 0) {
				simage := csketch.Output()
				bounds := simage.Bounds()
				palettedImage := image.NewPaletted(bounds, nil)
				quantizer := gogif.MedianCutQuantizer{NumColor: 64}
				quantizer.Quantize(palettedImage, bounds, simage, image.ZP)

				// Add new frame to animated GIF
				outGif.Image = append(outGif.Image, palettedImage)
				outGif.Delay = append(outGif.Delay, 0)
			}
		}

		util.SaveOutput(csketch.Output(), outputImgName)
		if gifme == true {
			simage := csketch.Output()
			bounds := simage.Bounds()
			palettedImage := image.NewPaletted(bounds, nil)
			quantizer := gogif.MedianCutQuantizer{NumColor: 64}
			quantizer.Quantize(palettedImage, bounds, simage, image.ZP)

			// Add new frame to animated GIF
			outGif.Image = append(outGif.Image, palettedImage)
			outGif.Delay = append(outGif.Delay, 100)
			gif.EncodeAll(f, outGif)
		}

	},
}

func init() {
	rootCmd.AddCommand(crackCmd)

	crackCmd.Flags().StringVarP(&outputImgName, "out", "o", "out.png", "Output image name")
	crackCmd.Flags().IntVarP(&limitByIterations, "iterations", "i", 0, "Number of iterations")
	crackCmd.Flags().IntVarP(&width, "width", "", 1920, "Width of output")
	crackCmd.Flags().IntVarP(&height, "height", "", 1080, "Height of output")
	crackCmd.Flags().BoolVarP(&save, "save", "s", false, "Save output regularly")
	crackCmd.Flags().BoolVarP(&gifme, "gif", "", false, "Create a gif of the results")

}
