package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"gitlab.com/ericworkman/generative/sketch"
	"gitlab.com/ericworkman/generative/util"
)

var (
	divisions = 12
)

// flipCmd represents the stack command
var flipCmd = &cobra.Command{
	Use:   "flip",
	Short: "Flip and style an image using diamonds",
	Long:  `Single pass only`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("flip called")
		img, _ := util.LoadUnsplashImage(width, height, url)

		params := sketch.FlipParams{
			DestWidth:  width,
			DestHeight: height,
		}

		csketch := sketch.NewFlipSketch(img, params)

		// catch the sigterm signal for ctrl-c quitting mostly
		// save the output at this point
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			util.SaveOutput(csketch.Output(), outputImgName)
			os.Exit(1)
		}()

		csketch.Draw(divisions)
		if save == true {
			util.SaveOutput(csketch.Output(), outputImgName)
		}

		util.SaveOutput(csketch.Output(), outputImgName)
	},
}

func init() {
	rootCmd.AddCommand(flipCmd)

	flipCmd.Flags().StringVarP(&url, "url", "u", "", "A url to an image")
	flipCmd.Flags().StringVarP(&outputImgName, "out", "o", "out.png", "Output image name")
	flipCmd.Flags().IntVarP(&width, "width", "", 1920, "Width of output")
	flipCmd.Flags().IntVarP(&height, "height", "", 1080, "Height of output")
	flipCmd.Flags().BoolVarP(&save, "save", "s", false, "Save output regularly")
	flipCmd.Flags().IntVarP(&divisions, "divisions", "d", 12, "Divisions of height")
}
