package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.com/ericworkman/generative/sketch"
	"gitlab.com/ericworkman/generative/util"
)

var (
	vignette = false
	size     = 20.0
)

// gridCmd represents the grid command
var gridCmd = &cobra.Command{
	Use:   "grid",
	Short: "Create a grided image",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("grid called")

		img, _ := util.LoadUnsplashImage(width, height, url)

		params := sketch.GridParams{
			DestWidth:  width,
			DestHeight: height,
			Vignette:   vignette,
			Size:       size,
		}

		csketch := sketch.NewGridSketch(img, params)
		csketch.Draw()

		util.SaveOutput(csketch.Output(), outputImgName)
	},
}

func init() {
	rootCmd.AddCommand(gridCmd)
	gridCmd.Flags().StringVarP(&url, "url", "u", "", "A url to an image")
	gridCmd.Flags().StringVarP(&outputImgName, "out", "o", "out.png", "Output image name")
	gridCmd.Flags().IntVarP(&width, "width", "", 1920, "Width of output")
	gridCmd.Flags().IntVarP(&height, "height", "", 1080, "Height of output")
	gridCmd.Flags().Float64VarP(&size, "size", "s", 20.0, "Size of grid")
	gridCmd.Flags().BoolVarP(&vignette, "vignette", "", false, "Vignette on the x-axis")
}
