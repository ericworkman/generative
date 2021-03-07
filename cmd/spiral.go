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

// spiralCmd represents the spiral command
var spiralCmd = &cobra.Command{
	Use:   "spiral",
	Short: "Create a logarithmic spiral",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("spiral called")
		params := sketch.SpiralParams{
			DestWidth:  width,
			DestHeight: height,
			Iterations: limitByIterations,
			Beta:       beta,
			Mu:         mu,
		}

		csketch := sketch.NewSpiralSketch(params)

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
	rootCmd.AddCommand(spiralCmd)

	spiralCmd.Flags().StringVarP(&outputImgName, "out", "o", "out.png", "Output image name")
	spiralCmd.Flags().IntVarP(&limitByIterations, "iterations", "i", 3, "Number of iterations")
	spiralCmd.Flags().IntVarP(&width, "width", "", 1920, "Width of output")
	spiralCmd.Flags().IntVarP(&height, "height", "", 1080, "Height of output")
	spiralCmd.Flags().Float64VarP(&beta, "beta", "", 0.100, "Tweakable")
	spiralCmd.Flags().Float64VarP(&mu, "mu", "", 0.100, "Tweakable")
}
