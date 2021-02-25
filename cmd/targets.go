package cmd

import (
	"fmt"

	"fuzznav/util"
	"github.com/spf13/cobra"
)

// targetsCmd represents the targets command
var targetsCmd = &cobra.Command{
	Use:     "targets",
	Aliases: []string{"t", "targ", "targs", "target"},
	Short:   "Show fuzz targets",
	Long:    `Show fuzz targets and wordlists used to fuzz the target.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("targets called")
		util.TargetsMap(util.ReadStdinAndParse())
	},
}

func init() {
	rootCmd.AddCommand(targetsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// targetsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// targetsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
