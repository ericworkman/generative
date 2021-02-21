package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"gitlab.com/ericworkman/generative/sketch"
)

// stackCmd represents the stack command
var stackCmd = &cobra.Command{
	Use:   "stack",
	Short: "Create an image of a stack of transparent shapes on a finer and finer grid",
	Long:  `Use iterations < 150 for regular usage`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("stack called")
		img, _ := sketch.LoadUnsplashImage(width, height, url)

		params := sketch.StackParams{
			DestWidth:  width,
			DestHeight: height,
		}

		csketch := sketch.NewStackSketch(img, params)

		// catch the sigterm signal for ctrl-c quitting mostly
		// save the output at this point
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			sketch.SaveOutput(csketch.Output(), outputImgName)
			os.Exit(1)
		}()

		for i := 1; i <= limitByIterations; i++ {
			fmt.Println("Iteration", i)

			csketch.Update(i)
			if save == true {
				sketch.SaveOutput(csketch.Output(), outputImgName)
			}
		}

		sketch.SaveOutput(csketch.Output(), outputImgName)
	},
}

func init() {
	rootCmd.AddCommand(stackCmd)

	stackCmd.Flags().StringVarP(&url, "url", "u", "", "A url to an image")
	stackCmd.Flags().StringVarP(&outputImgName, "out", "o", "out.png", "Output image name")
	stackCmd.Flags().IntVarP(&limitByIterations, "iterations", "i", 3, "Number of iterations")
	stackCmd.Flags().IntVarP(&width, "width", "", 1920, "Width of output")
	stackCmd.Flags().IntVarP(&height, "height", "", 1080, "Height of output")
	stackCmd.Flags().BoolVarP(&save, "save", "s", false, "Save output regularly")
}
