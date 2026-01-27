/*
Copyright © 2025 Mirko Mariotti mirko@mirkomariotti.it
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// rebuildCmd represents the rebuild command
var rebuildCmd = &cobra.Command{
	Use:   "rebuild",
	Short: "Rebuild the engine (experimental)",
	Long:  `Rewrite some parts of the engine's internal state and data structures parsing the rules database. This is an experimental feature and may not work as expected.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := e.RebuildEngine(); err != nil {
			cmd.Println("Failed to rebuild engine:")
			cmd.Println(err.Error())
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(rebuildCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
