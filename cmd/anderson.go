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

// andersonCmd represents the anderson command
var andersonCmd = &cobra.Command{
	Use:   "anderson",
	Short: "Create art based on Jason Anderson's work",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("anderson called")
		params := sketch.AndersonParams{
			DestWidth:  width,
			DestHeight: height,
			Iterations: limitByIterations,
		}

		csketch := sketch.NewAndersonSketch(params)

		// catch the sigterm signal for ctrl-c quitting mostly
		// save the output at this point
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			util.SaveOutput(csketch.Output(), outputImgName)
			os.Exit(1)
		}()

		for i := 0; i <= limitByIterations; i++ {
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
	rootCmd.AddCommand(andersonCmd)

	andersonCmd.Flags().StringVarP(&outputImgName, "out", "o", "out.png", "Output image name")
	andersonCmd.Flags().IntVarP(&limitByIterations, "iterations", "i", 3, "Number of iterations")
	andersonCmd.Flags().IntVarP(&width, "width", "", 1920, "Width of output")
	andersonCmd.Flags().IntVarP(&height, "height", "", 1080, "Height of output")
}
