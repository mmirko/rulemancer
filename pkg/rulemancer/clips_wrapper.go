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
char* find_facts_as_string(void*, const char*);
char* find_all_facts_as_string(void*);
*/
import "C"
import (
	"fmt"
	"log"
	"os"
	"unsafe"
)

type ClipsInstance struct {
	e     *Engine
	cl    unsafe.Pointer
	sChan chan struct{} // serialize channel
	rChan chan struct{} // response channel
	qChan chan struct{} // quit channel
}

func (e *Engine) NewClipsInstance() *ClipsInstance {
	return &ClipsInstance{
		e:     e,
		sChan: make(chan struct{}),
		rChan: make(chan struct{}),
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
				<-ci.rChan
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
	response := map[string]string{
		"status":        "running",
		"engine":        "CLIPS",
		"version":       "6.40",
		"rule_pool":     ci.e.RulePool,
		"instance_addr": fmt.Sprintf("%p", ci.cl),
	}
	ci.rChan <- struct{}{}
	return response
}

func (ci *ClipsInstance) AssertFact(fact string) error {
	// Assert a fact into the CLIPS environment
	if ci.cl == nil {
		return fmt.Errorf("CLIPS instance not initialized")
	}
	<-ci.sChan
	cFact := C.CString(fact)
	defer C.free(unsafe.Pointer(cFact))
	C.clips_assert(ci.cl, cFact)
	ci.rChan <- struct{}{}
	return nil
}

func (ci *ClipsInstance) Run() error {
	// Run the CLIPS engine
	if ci.cl == nil {
		return fmt.Errorf("CLIPS instance not initialized")
	}
	<-ci.sChan
	C.clips_run(ci.cl)
	ci.rChan <- struct{}{}
	return nil
}

func (ci *ClipsInstance) QueryFacts(relation string) (string, error) {
	// Query facts matching the pattern
	if ci.cl == nil {
		return "", fmt.Errorf("CLIPS instance not initialized")
	}
	<-ci.sChan
	cRelation := C.CString(relation)
	defer C.free(unsafe.Pointer(cRelation))
	facts := C.find_facts_as_string(ci.cl, cRelation)
	defer C.free(unsafe.Pointer(facts))
	goFacts := sanitizeFacts(C.GoString(facts))
	if ci.e.Debug {
		l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/QueryFacts]")+" ", 0)
		l.Println("Queried facts raw:", goFacts)
	}
	ci.rChan <- struct{}{}
	return goFacts, nil
}

func (ci *ClipsInstance) QueryFactsAllFacts() (string, error) {
	// Query all facts
	if ci.cl == nil {
		return "", fmt.Errorf("CLIPS instance not initialized")
	}
	<-ci.sChan
	facts := C.find_all_facts_as_string(ci.cl)
	defer C.free(unsafe.Pointer(facts))
	goFacts := sanitizeFacts(C.GoString(facts))
	if ci.e.Debug {
		l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/QueryFactsAllFacts]")+" ", 0)
		l.Println("Queried all facts raw:", goFacts)
	}
	ci.rChan <- struct{}{}
	return goFacts, nil
}

func (ci *ClipsInstance) Dispose() {
	if ci.cl == nil {
		return
	}
	<-ci.qChan
	// Destroy the CLIPS environment
	C.clips_destroy(ci.cl)
}
