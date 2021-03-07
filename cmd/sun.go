package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.com/ericworkman/generative/sketch"
	"gitlab.com/ericworkman/generative/util"
)

var sunCmd = &cobra.Command{
	Use:   "sun",
	Short: "Create a stylized sun and sky inspired by https://www.reddit.com/user/yum_paste",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("sun called")
		params := sketch.SunParams{
			DestWidth:  width,
			DestHeight: height,
			SunRadius:  beta,
			LineWidth:  mu,
		}

		ssketch := sketch.NewSunSketch(params)

		ssketch.Draw()

		util.SaveOutput(ssketch.Output(), outputImgName)
	},
}

func init() {
	rootCmd.AddCommand(sunCmd)
	sunCmd.Flags().StringVarP(&outputImgName, "out", "o", "out.png", "Output image name")
	sunCmd.Flags().IntVarP(&width, "width", "", 1920, "Width of output")
	sunCmd.Flags().IntVarP(&height, "height", "", 1080, "Height of output")
	sunCmd.Flags().Float64VarP(&beta, "beta", "", 50, "Radius of sun")
	sunCmd.Flags().Float64VarP(&mu, "mu", "", 5.0, "Thickness of lines")
}
