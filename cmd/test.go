/*
Copyright © 2025 Mirko Mariotti mirko@mirkomariotti.it
*/
package cmd

/*
#cgo CFLAGS: -I core
#cgo LDFLAGS: -L core -lclips -lm
#include <stdlib.h>

void* clips_create();
void clips_destroy(void*);
void clips_load(void*, const char*);
void clips_reset(void*);
void clips_run(void*);
void clips_assert(void*, const char*);
*/
import "C"
import (
	"os"
	"unsafe"

	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test the rules engine",
	Long:  `A command to open the rules database and run pre-defined tests. It does not rely on the configuration file, but uses the rulepool and testpool directories as given in the flags.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if rulePool directory exists
		if _, err := os.Stat(rulePool); os.IsNotExist(err) {
			cmd.Println("Rule pool directory does not exist:", rulePool)
			return
		}
		if rulesFiles, err := os.ReadDir(rulePool); err != nil {
			cmd.Println("Failed to read rule pool directory:", err.Error())
			return
		} else {
			// Create a CLIPS environment
			env := C.clips_create()
			defer C.clips_destroy(env)

			for _, file := range rulesFiles {
				if !file.IsDir() {
					cfile := C.CString(rulePool + "/" + file.Name())
					defer C.free(unsafe.Pointer(cfile))
					C.clips_load(env, cfile)
				}
			}
			C.clips_reset(env)
			C.clips_run(env)

			if testFiles, err := os.ReadDir(testPool); err != nil {
				cmd.Println("Failed to read test pool directory:", err.Error())
				return
			} else {
				for _, file := range testFiles {
					if !file.IsDir() {
						testFilePath := testPool + "/" + file.Name()
						testBytes, err := os.ReadFile(testFilePath)
						if err != nil {
							cmd.Println("Failed to read test file:", testFilePath, err.Error())
							continue
						}
						fact := C.CString(string(testBytes))
						C.clips_assert(env, fact)
						C.clips_run(env)
						C.free(unsafe.Pointer(fact))
					}
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
