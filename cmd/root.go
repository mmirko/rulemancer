/*
Copyright © 2025 Mirko Mariotti mirko@mirkomariotti.it

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"log"
	"os"

	"github.com/mmirko/rulemancer/pkg/rulemancer"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rulemancer",
	Short: "CLIPS-based Engine for rules-based games",
	Long:  `rulemancer is a CLIPS-based go application to manage and serve rules-based games over an HTTP API.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

var cfgFile string     // config file
var rulePool string    // Knowledge base pool, used to store the knowledge base in form of a Prolog files
var testPool string    // Test pool, used to store test files
var TLSCertFile string // TLS Certificate file
var TLSKeyFile string  // TLS Key file

var e *rulemancer.Engine // Engine object

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "rulemancer.json", "config file (default rulemancer.json)")
	rootCmd.PersistentFlags().StringVarP(&rulePool, "rulepool", "k", "rulepool", "Knowledge base pool directory (default rulepool)")
	rootCmd.PersistentFlags().StringVarP(&testPool, "testpool", "t", "testpool", "Test pool directory (default testpool)")
	rootCmd.PersistentFlags().StringVarP(&TLSCertFile, "tlscert", "", "server.crt", "TLS Certificate file (default server.crt)")
	rootCmd.PersistentFlags().StringVarP(&TLSKeyFile, "tlskey", "", "server.key", "TLS Key file (default server.key)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	e = rulemancer.NewEngine()

	if cfgFile != "" {
		err := e.LoadConfig(cfgFile)
		if err != nil {
			log.Fatalf("Error loading config file: %v", err)
		}
	}

	// Override rule pool if specified
	if rulePool != "" {
		e.RulePool = rulePool
	}
	// Override TLS cert and key if specified
	if TLSCertFile != "" {
		e.TLSCertFile = TLSCertFile
	}
	if TLSKeyFile != "" {
		e.TLSKeyFile = TLSKeyFile
	}
}
