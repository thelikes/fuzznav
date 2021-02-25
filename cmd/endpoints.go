package cmd

import (
	"fmt"

	"fuzznav/util"
	"github.com/spf13/cobra"
)

// endpointsCmd represents the endpoints command
var endpointsCmd = &cobra.Command{
	Use:     "endpoints",
	Aliases: []string{"e", "ep", "eps"},
	Short:   "Show discovered endpoints",
	Long:    `Show discovered endpoints, status, size, words, and lines.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("[+] endpoints called")
		util.EndpointsMap(util.ReadStdinAndParse())
	},
}

func init() {
	rootCmd.AddCommand(endpointsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// endpointsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// endpointsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
