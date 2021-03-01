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

var (
	crawlCount = 3
	crawlStart = "center"
)

// crawlCmd represents the crawl command
var crawlCmd = &cobra.Command{
	Use:   "crawl",
	Short: "Create crawling lines from a center point",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("crawl called")
		params := sketch.CrawlParams{
			DestWidth:  width,
			DestHeight: height,
			Iterations: limitByIterations,
			Count:      crawlCount,
			Start:      crawlStart,
		}

		csketch := sketch.NewCrawlSketch(params)

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
			if i%(util.MaxInt(limitByIterations/10, 1)) == 0 {
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
	rootCmd.AddCommand(crawlCmd)
	crawlCmd.Flags().StringVarP(&outputImgName, "out", "o", "out.png", "Output image name")
	crawlCmd.Flags().IntVarP(&limitByIterations, "iterations", "i", 1, "Number of iterations")
	crawlCmd.Flags().IntVarP(&width, "width", "", 1920, "Width of output")
	crawlCmd.Flags().IntVarP(&height, "height", "", 1080, "Height of output")
	crawlCmd.Flags().IntVarP(&crawlCount, "count", "", 3, "Number of crawlers")
	crawlCmd.Flags().StringVarP(&crawlStart, "start", "", "center", "center or corner starting location")
}
