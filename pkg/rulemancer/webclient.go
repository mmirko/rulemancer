package rulemancer

import (
	"fmt"
	"log"
	"net/http"
	"os"

	chi "github.com/go-chi/chi/v5"
)

func (e *Engine) templateGenerator(baseDir, templateDir string) (map[string]string, error) {
	templateMap := make(map[string]string)

	if templateFiles, err := os.ReadDir(baseDir + "/" + templateDir); err == nil {
		// Process each template file
		for _, file := range templateFiles {
			if !file.IsDir() {
				content, err := os.ReadFile(baseDir + "/" + templateDir + "/" + file.Name())
				if err != nil {
					return nil, err
				}
				templateMap[templateDir+"/"+file.Name()] = string(content)
			} else {
				subTemplateMap, err := e.templateGenerator(baseDir, templateDir+"/"+file.Name())
				if err != nil {
					return nil, err
				}
				for k, v := range subTemplateMap {
					templateMap[k] = v
				}
			}
		}
	}

	return templateMap, nil
}

func (e *Engine) handlerGenerator(path string, templateMap map[string]string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(templateMap[path]))
	}
}

func (e *Engine) webClientRoutes(r chi.Router) {

	templateDir := "pkg/rulemancer/templates/webclient"

	templateMap, err := e.templateGenerator(templateDir, "")
	if err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/SpawnWebClient]")+" ", 0)
			l.Printf("Failed to generate templates: %v", err)
		}
	}
	fmt.Println(templateMap)
	r.Route("/", func(r chi.Router) {
		for path := range templateMap {
			r.Get(path, e.handlerGenerator(path, templateMap))
		}
	})
}
