package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"gitlab.com/ericworkman/generative/sketch"
)

var (
	beta = 0.101
	mu   = 0.10
)

// spiralCmd represents the spiral command
var spiralCmd = &cobra.Command{
	Use:   "spiral",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
	rootCmd.AddCommand(spiralCmd)

	spiralCmd.Flags().StringVarP(&outputImgName, "out", "o", "out.png", "Output image name")
	spiralCmd.Flags().IntVarP(&limitByIterations, "iterations", "i", 3, "Number of iterations")
	spiralCmd.Flags().IntVarP(&width, "width", "", 1920, "Width of output")
	spiralCmd.Flags().IntVarP(&height, "height", "", 1080, "Height of output")
	layerCmd.Flags().Float64VarP(&beta, "beta", "", 0.100, "Tweakable")
	layerCmd.Flags().Float64VarP(&mu, "mu", "", 0.100, "Tweakable")
}
