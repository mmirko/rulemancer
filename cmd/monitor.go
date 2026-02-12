/*
Copyright © 2026 Mirko Mariotti mirko@mirkomariotti.it
*/
package cmd

import (
	"log"

	"github.com/mmirko/rulemancer/pkg/rulemancer"
	"github.com/spf13/cobra"
)

var baseURL string

// monitorCmd represents the monitor command
var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Monitor events of the engine",
	Long:  `Monitor the engine for games, system or other events.`,
	Run: func(cmd *cobra.Command, args []string) {

		// Initialize the engine with the secret
		e = rulemancer.NewEngine("")

		if cfgFile != "" {
			err := e.LoadConfig(cfgFile)
			if err != nil {
				log.Fatalf("Error loading config file: %v", err)
			}
		}
		if err := e.Monitor(baseURL); err != nil {
			log.Fatalf("Error monitoring engine: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(monitorCmd)
	monitorCmd.Flags().StringVarP(&baseURL, "baseurl", "b", "wss://localhost:3000/api/v1/system/ws", "Base URL for the engine websocket API")
}
