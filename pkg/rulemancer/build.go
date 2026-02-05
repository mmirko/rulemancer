package rulemancer

import (
	"fmt"
	"log"
	"os"
)

type relevantFact struct {
	relation   string
	slots      []string
	multislots []string
}

func (e *Engine) BuildEngineExtras(shellOutdir string) error {
	// The rebuild engine reads the rules games directories and the assertables,results and querables from there.
	// Then uses re2c to write the various artifacts needed to interact with the engine.

	e.loadGames()

	// Define the templates directory
	templateDir := "pkg/rulemancer/templates"

	templateMap := make(map[string]string)

	// Load the shell templates
	if _, err := os.Stat(templateDir + "/shell"); os.IsNotExist(err) {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/BuildEngineExtras]")+" ", 0)
			l.Printf("Shell template directory does not exist: %s", templateDir+"/shell")
		}
		return fmt.Errorf("template location does not exist: %s", templateDir+"/shell")
	}
	if templateFiles, err := os.ReadDir(templateDir + "/shell"); err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/BuildEngineExtras]")+" ", 0)
			l.Printf("Error reading shell template directory: %v", err)
		}
		return fmt.Errorf("failed to read rules location: %w", err)
	} else {
		// Process each template file
		for _, file := range templateFiles {
			if !file.IsDir() {
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/BuildEngineExtras]")+" ", 0)
					l.Printf("Processing shell template file: %s", file.Name())
				}
				content, err := os.ReadFile(templateDir + "/shell/" + file.Name())
				if err != nil {
					if e.Debug {
						l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/BuildEngineExtras]")+" ", 0)
						l.Printf("Error reading template file %s: %v", file.Name(), err)
					}
					return fmt.Errorf("failed to read template file %s: %w", file.Name(), err)
				}
				templateMap[file.Name()] = string(content)
			}
		}
	}

	// Create the output directory if it doesn't exist
	if _, err := os.Stat(shellOutdir); os.IsNotExist(err) {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/BuildEngineExtras]")+" ", 0)
			l.Printf("Output directory does not exist, creating: %s", shellOutdir)
		}
		if err := os.MkdirAll(shellOutdir, 0755); err != nil {
			if e.Debug {
				l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/BuildEngineExtras]")+" ", 0)
				l.Printf("Error creating output directory %s: %v", shellOutdir, err)
			}
			return fmt.Errorf("failed to create output directory %s: %w", shellOutdir, err)
		}
	}

	for _, game := range e.games {
		gameName := game.name
		rulesLocation := game.rulesLocation
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/BuildEngineExtras]")+" ", 0)
			l.Printf("Building engine extras for game: %s from rules location: %s", gameName, rulesLocation)
		}
		// Load a game from the specified rules location
		if _, err := os.Stat(rulesLocation); os.IsNotExist(err) {
			return fmt.Errorf("rules location does not exist: %s", rulesLocation)
		}
		if rulesFiles, err := os.ReadDir(rulesLocation); err != nil {
			return fmt.Errorf("failed to read rules location: %w", err)
		} else {
			// Load each rule file into CLIPS
			for _, file := range rulesFiles {
				// load the file
				if !file.IsDir() {
					if e.Debug {
						l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/BuildEngineExtras]")+" ", 0)
						l.Printf("Processing rule file: %s", file.Name())
					}
					fileContent, err := os.ReadFile(rulesLocation + "/" + file.Name())
					if err != nil {
						if e.Debug {
							l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/BuildEngineExtras]")+" ", 0)
							l.Printf("Error reading rule file %s: %v", file.Name(), err)
						}
						return fmt.Errorf("failed to read rule file %s: %w", file.Name(), err)
					}

					pd := e.newProtocolData()
					if err := pd.Compile(string(fileContent) + "\x00"); err != nil {
						if e.Debug {
							l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/BuildEngineExtras]")+" ", 0)
							l.Printf("Error compiling rule file %s: %v", file.Name(), err)
						}
						return fmt.Errorf("failed to compile rule file %s: %w", file.Name(), err)
					}

				}
			}
		}
	}

	// TODO
	return nil
}
