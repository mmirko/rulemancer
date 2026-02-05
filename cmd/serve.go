/*
Copyright © 2025 Mirko Mariotti mirko@mirkomariotti.it
*/
package cmd

import (
	"log"
	"os"

	"github.com/mmirko/rulemancer/pkg/rulemancer"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the engine",
	Long:  `Spawn the engine and serve it over an HTTP API.`,
	Run: func(cmd *cobra.Command, args []string) {

		// Search for JWT secret in environment variable if not provided as flag
		if secret == "" {
			secret = os.Getenv("RULEMANCER_JWT_SECRET")
			if secret == "" {
				log.Fatal("JWT secret must be provided via --secret flag or RULEMANCER_JWT_SECRET environment variable")
			}
		}

		// Initialize the engine with the secret
		e = rulemancer.NewEngine(secret)

		if cfgFile != "" {
			err := e.LoadConfig(cfgFile)
			if err != nil {
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[cmd/serve]")+" ", 0)
					l.Printf("Error loading config file: %v", err)
				}
				log.Fatalf("Error loading config file: %v", err)
			}
		}

		// Override TLS cert and key if specified
		if TLSCertFile != "" {
			if e.Debug {
				l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[cmd/root]")+" ", 0)
				l.Printf("Overriding TLS cert file to %s", TLSCertFile)
			}
			e.TLSCertFile = TLSCertFile
		}
		if TLSKeyFile != "" {
			if e.Debug {
				l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[cmd/root]")+" ", 0)
				l.Printf("Overriding TLS key file to %s", TLSKeyFile)
			}
			e.TLSKeyFile = TLSKeyFile
		}

		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[cmd/serve]")+" ", 0)
			l.Println("Starting engine in debug mode...")
		}
		if err := e.SpawnEngine(); err != nil {
			if e.Debug {
				l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[cmd/serve]")+" ", 0)
				l.Printf("Error spawning engine: %v", err)
			}
			log.Fatalf("Error spawning engine: %v", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
