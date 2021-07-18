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

// mondrianCmd represents the stack command
var mondrianCmd = &cobra.Command{
	Use:   "mondrian",
	Short: "Create rectangles with border from a sample image",
	Long:  `Use iterations < 150 for regular usage`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("mondrian called")
		img, _ := util.LoadUnsplashImage(width, height, url)

		params := sketch.MondrianParams{
			DestWidth:  width,
			DestHeight: height,
		}

		csketch := sketch.NewMondrianSketch(img, params)

		// catch the sigterm signal for ctrl-c quitting mostly
		// save the output at this point
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			util.SaveOutput(csketch.Output(), outputImgName)
			os.Exit(1)
		}()

		for i := 1; i <= limitByIterations; i++ {
			fmt.Println("Iteration", i)

			csketch.Update(i)
			if save == true {
				util.SaveOutput(csketch.Output(), outputImgName)
			}
		}

		util.SaveOutput(csketch.Output(), outputImgName)
	},
}

func init() {
	rootCmd.AddCommand(mondrianCmd)

	mondrianCmd.Flags().StringVarP(&url, "url", "u", "", "A url to an image")
	mondrianCmd.Flags().StringVarP(&outputImgName, "out", "o", "out.png", "Output image name")
	mondrianCmd.Flags().IntVarP(&limitByIterations, "iterations", "i", 3, "Number of iterations")
	mondrianCmd.Flags().IntVarP(&width, "width", "", 1920, "Width of output")
	mondrianCmd.Flags().IntVarP(&height, "height", "", 1080, "Height of output")
	mondrianCmd.Flags().BoolVarP(&save, "save", "s", false, "Save output regularly")
}
