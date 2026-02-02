/*
Copyright © 2025 Mirko Mariotti mirko@mirkomariotti.it
*/
package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the engine",
	Long:  `Spawn the engine and serve it over an HTTP API.`,
	Run: func(cmd *cobra.Command, args []string) {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[cmd/serve]")+" ", 0)
			l.Println("Starting engine in debug mode...")
		}
		if err := e.SpawnEngine(); err != nil {
			cmd.Println("Failed to spawn engine:")
			cmd.Println(err.Error())
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
