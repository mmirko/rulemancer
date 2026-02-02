/*
Copyright © 2025 Mirko Mariotti mirko@mirkomariotti.it
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var rebuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the engine extra elements (experimental)",
	Long:  `Write extra elements for games, things like call from curl, JSON interfaces and other stuff.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := e.RebuildEngine(); err != nil {
			cmd.Println("Failed to build tools:")
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
