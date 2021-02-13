package cmd

import (
	"github.com/spf13/cobra"

  "gitlab.com/ericworkman/generative/sketch"
)

var (
  outputImgName = ""
  url = ""
  limitByIterations = 0
  limitBySize = float64(5)
  reduction = float64(0)
  ratio = float64(0)
  alpha = float64(0)
  alphaIncrease = float64(0)
  jitter = float64(0)
  edge = false
  edgeMin = 0
  edgeMax = 0
  width = 1920
  height = 1080
  inversionThreshold = float64(5)
)

var layerCmd = &cobra.Command{
	Use:   "layer",
	Short: "Create sketches in the style of Preslav Rachev",
	Long: `Create a sketch of overlapping shapes with various drawing options
`,
	Run: func(cmd *cobra.Command, args []string) {

    img, _ := sketch.LoadUnsplashImage(width, height, url)

    if edgeMin > edgeMax {
      edgeMax = edgeMin
    }

    params := sketch.LayerParams{
      DestWidth: width,
      DestHeight: height,
      PathRatio: ratio,
      PathReduction: reduction,
      PathMin: limitBySize,
      PathJitter: int(jitter * float64(width)),
      InitialAlpha: alpha,
      AlphaIncrease: alphaIncrease,
      MinEdgeCount: edgeMin,
      MaxEdgeCount: edgeMax,
      Edge: edge,
      PathInversionThreshold: inversionThreshold,
    }

    lsketch := sketch.NewLayerSketch(img, params)

    if limitByIterations == 0 {
      for i := 0; lsketch.PathSize >= params.PathMin; i++ {
        lsketch.Update()
      }
    } else {
      for i := 0; i < limitByIterations; i++ {
        lsketch.Update()
      }
    }

    sketch.SaveOutput(lsketch.Output(), outputImgName)
	},
}

func init() {
	rootCmd.AddCommand(layerCmd)

	layerCmd.Flags().StringVarP(&url, "url", "u", "", "A url to an image")
	layerCmd.Flags().StringVarP(&outputImgName, "out", "o", "out.png", "Output image name")
	layerCmd.Flags().IntVarP(&width, "width", "", 1920, "Width of output")
	layerCmd.Flags().IntVarP(&height, "height", "", 1080, "Height of output")

	layerCmd.Flags().IntVarP(&limitByIterations, "iterations", "i", 0, "Number of iterations")
	layerCmd.Flags().Float64VarP(&reduction, "reduction", "", 0.001, "Reduction per iteration")
	layerCmd.Flags().Float64VarP(&limitBySize, "minsize", "", 5.0, "Minimun size of paths")
	layerCmd.Flags().Float64VarP(&ratio, "ratio", "", 0.50, "Starting path size as a ratio of image width")
	layerCmd.Flags().Float64VarP(&alpha, "alpha", "a", 0.1, "Starting alpha")
	layerCmd.Flags().Float64VarP(&alphaIncrease, "alphaIncrease", "", 0.006, "Increase of alpha per iteration")
	layerCmd.Flags().IntVarP(&edgeMin, "edgeMinimum", "", 0, "Minimum number of edges of path")
	layerCmd.Flags().IntVarP(&edgeMax, "edgeMaximum", "", 0, "Maximum number of edges of path")
	layerCmd.Flags().Float64VarP(&jitter, "jitter", "", 0.007, "Jitter multiplier")
	layerCmd.Flags().BoolVarP(&edge, "edge", "", false, "Paint edges with inversion")
	layerCmd.Flags().Float64VarP(&inversionThreshold, "inversion", "", 0.05, "Size at which to invert the color")
}
