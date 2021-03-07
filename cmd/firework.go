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

// fireworkCmd represents the firework command
var fireworkCmd = &cobra.Command{
	Use:   "firework",
	Short: "Generate a firework",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("firework called")
		params := sketch.FireworkParams{
			DestWidth:  width,
			DestHeight: height,
			Iterations: limitByIterations,
		}

		csketch := sketch.NewFireworkSketch(params)

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
			if i%(util.MaxInt(1, limitByIterations/10)) == 0 {
				fmt.Println("Iteration", i)
			}

			csketch.Update(i)
			if save == true {
				util.SaveOutput(csketch.Output(), outputImgName)
			}
		}

		util.SaveOutput(csketch.Output(), outputImgName)
	},
}

func init() {
	rootCmd.AddCommand(fireworkCmd)
	fireworkCmd.Flags().StringVarP(&outputImgName, "out", "o", "out.png", "Output image name")
	fireworkCmd.Flags().IntVarP(&limitByIterations, "iterations", "i", 3, "Number of iterations")
	fireworkCmd.Flags().IntVarP(&width, "width", "", 1920, "Width of output")
	fireworkCmd.Flags().IntVarP(&height, "height", "", 1080, "Height of output")
}
