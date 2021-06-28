package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.com/ericworkman/generative/sketch"
	"gitlab.com/ericworkman/generative/util"
)

var ()

// rowsCmd represents the rows command
var rowsCmd = &cobra.Command{
	Use:   "rows",
	Short: "Create a row-based image",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("rows called")

		img, _ := util.LoadUnsplashImage(width, height, url)

		params := sketch.RowsParams{
			DestWidth:  width,
			DestHeight: height,
			Vignette:   vignette,
			Size:       size,
		}

		csketch := sketch.NewRowsSketch(img, params)
		csketch.Draw()

		util.SaveOutput(csketch.Output(), outputImgName)
	},
}

func init() {
	rootCmd.AddCommand(rowsCmd)
	rowsCmd.Flags().StringVarP(&url, "url", "u", "", "A url to an image")
	rowsCmd.Flags().StringVarP(&outputImgName, "out", "o", "out.png", "Output image name")
	rowsCmd.Flags().IntVarP(&width, "width", "", 1920, "Width of output")
	rowsCmd.Flags().IntVarP(&height, "height", "", 1080, "Height of output")
	rowsCmd.Flags().Float64VarP(&size, "size", "s", 20.0, "Size of grid")
	rowsCmd.Flags().BoolVarP(&vignette, "vignette", "", false, "Vignette on the x-axis")
}
