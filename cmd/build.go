/*
Copyright © 2025 Mirko Mariotti mirko@mirkomariotti.it
*/
package cmd

import (
	"log"

	"github.com/mmirko/rulemancer/pkg/rulemancer"
	"github.com/spf13/cobra"
)

var interfaceDir string

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the engine extra elements (experimental)",
	Long:  `Write extra elements for games, things like call from curl, JSON interfaces and other stuff.`,
	Run: func(cmd *cobra.Command, args []string) {

		// Initialize the engine with the secret
		e = rulemancer.NewEngine("")

		if cfgFile != "" {
			err := e.LoadConfig(cfgFile)
			if err != nil {
				log.Fatalf("Error loading config file: %v", err)
			}
		}

		if err := e.BuildEngineExtras(interfaceDir); err != nil {
			log.Fatalf("Error building engine extras: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
	buildCmd.Flags().StringVarP(&interfaceDir, "interface", "i", "interface", "Output directory for engine extras")
}
