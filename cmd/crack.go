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
  save = false
)

var crackCmd = &cobra.Command{
	Use:   "crack",
	Short: "Create sketches in the stle of Jared Tarbell",
	Long: `Create a sketch of growing cracks that "crystalize"
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("crack called")

    params := sketch.CrackParams{
      DestWidth: width,
      DestHeight: height,
      CrackLimit: 10,
      Seeds: width / 10 + height / 10,
      StartingCracks: 2,
    }

    csketch := sketch.NewCrackSketch(params)

    // catch the sigterm signal for ctrl-c quitting mostly
    // save the output at this point
    c := make(chan os.Signal)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        sketch.SaveOutput(csketch.Output(), outputImgName)
        os.Exit(1)
    }()

    for i := 0; i < limitByIterations; i++ {
      fmt.Println("Iteration", i)

      csketch.Update()

      // save the output every so often so that we don't just lose a lot of work
      // this isn't a nice way of handling it, but we'll live
      if (save == true) && (i % 100 == 0) {
        sketch.SaveOutput(csketch.Output(), outputImgName)
      }
    }

    sketch.SaveOutput(csketch.Output(), outputImgName)
	},
}

func init() {
	rootCmd.AddCommand(crackCmd)

	crackCmd.Flags().StringVarP(&outputImgName, "out", "o", "out.png", "Output image name")
	crackCmd.Flags().IntVarP(&limitByIterations, "iterations", "i", 0, "Number of iterations")
	crackCmd.Flags().IntVarP(&width, "width", "", 1920, "Width of output")
	crackCmd.Flags().IntVarP(&height, "height", "", 1080, "Height of output")
	crackCmd.Flags().BoolVarP(&save, "save", "s", false, "Save output regularly")
}

