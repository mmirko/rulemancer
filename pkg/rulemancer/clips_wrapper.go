package rulemancer

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
	"fmt"
	"os"
	"unsafe"
)

type ClipsInstance struct {
	e     *Engine
	cl    unsafe.Pointer
	sChan chan struct{}
	qChan chan struct{}
}

func (e *Engine) NewClipsInstance() *ClipsInstance {
	return &ClipsInstance{
		e:     e,
		sChan: make(chan struct{}),
		qChan: make(chan struct{}),
	}
}

// InitClips initializes the CLIPS environment for the instance, it is ment to be called once per instance
func (ci *ClipsInstance) InitClips() error {
	// Initialize CLIPS environment
	env := C.clips_create()
	ci.cl = env
	ci.spawnSerializer()

	return nil
}

// LoadKnowledgeBase loads the knowledge base from the rule pool directory, it is ment to be called once per instance after InitClips
func (ci *ClipsInstance) LoadKnowledgeBase() error {
	// Load knowledge base from the pool directory
	rulePool := ci.e.RulePool
	if _, err := os.Stat(rulePool); os.IsNotExist(err) {
		return fmt.Errorf("rule pool directory does not exist: %s", rulePool)
	}
	if rulesFiles, err := os.ReadDir(rulePool); err != nil {
		return fmt.Errorf("failed to read rule pool directory: %w", err)
	} else {
		// Load each rule file into CLIPS
		for _, file := range rulesFiles {
			if !file.IsDir() {
				cfile := C.CString(rulePool + "/" + file.Name())
				defer C.free(unsafe.Pointer(cfile))
				C.clips_load(ci.cl, cfile)
			}
		}
		C.clips_reset(ci.cl)
		C.clips_run(ci.cl)
	}
	return nil
}

func (ci *ClipsInstance) spawnSerializer() {
	go func() {
		for {
			select {
			case ci.sChan <- struct{}{}:
				// Every time a request is made, this channel will be read to ensure serialized access
			case ci.qChan <- struct{}{}:
				// On dispose, exit the goroutine
				return
			}
		}
	}()
}

func (ci *ClipsInstance) Info() map[string]string {
	// Return basic info about the CLIPS instance
	if ci.cl == nil {
		return map[string]string{
			"status": "uninitialized",
		}
	}
	<-ci.sChan
	return map[string]string{
		"status":        "running",
		"engine":        "CLIPS",
		"version":       "6.40",
		"rule_pool":     ci.e.RulePool,
		"instance_addr": fmt.Sprintf("%p", ci.cl),
	}
}

func (ci *ClipsInstance) AssertFact(fact string) error {
	// Assert a fact into the CLIPS environment
	if ci.cl == nil {
		return fmt.Errorf("CLIPS instance not initialized")
	}
	<-ci.sChan
	return nil
}

func (ci *ClipsInstance) Run() error {
	// Run the CLIPS engine
	if ci.cl == nil {
		return fmt.Errorf("CLIPS instance not initialized")
	}
	<-ci.sChan
	C.clips_run(ci.cl)
	return nil
}

func (ci *ClipsInstance) QueryFacts(pattern string) ([]string, error) {
	// Query facts matching the pattern
	if ci.cl == nil {
		return nil, fmt.Errorf("CLIPS instance not initialized")
	}
	<-ci.sChan
	return nil, nil
}

func (ci *ClipsInstance) Dispose() {
	if ci.cl == nil {
		return
	}
	<-ci.qChan
	// Destroy the CLIPS environment
	C.clips_destroy(ci.cl)
}
