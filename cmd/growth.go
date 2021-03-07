package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.com/ericworkman/generative/sketch"
	"gitlab.com/ericworkman/generative/util"
)

var growthCmd = &cobra.Command{
	Use:   "growth",
	Short: "Grow crystals from seeds",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("growth called")
		params := sketch.GrowthParams{
			DestWidth:     width,
			DestHeight:    height,
			StartingSeeds: seeds,
		}

		ssketch := sketch.NewGrowthSketch(params)

		ssketch.Draw()

		util.SaveOutput(ssketch.Output(), outputImgName)
	},
}

func init() {
	rootCmd.AddCommand(growthCmd)
	growthCmd.Flags().StringVarP(&outputImgName, "out", "o", "out.png", "Output image name")
	growthCmd.Flags().IntVarP(&width, "width", "", 1920, "Width of output")
	growthCmd.Flags().IntVarP(&height, "height", "", 1080, "Height of output")
	growthCmd.Flags().IntVarP(&seeds, "seeds", "", 5, "Number of starting seeds")
}
