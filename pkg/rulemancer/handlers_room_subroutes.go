/*
Copyright © 2026 Mirko Mariotti mirko@mirkomariotti.it
*/
package rulemancer

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/mmirko/rulemancer/pkg/game"
)

func (e *Engine) roomSubRoutes(r chi.Router) {
	r.Route("/", func(r chi.Router) {
		r.Post("/assert", e.apiAssert)
		r.Route("/query", e.querySubRoutes)

		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/roomSubRoutes]")+" ", 0)
			l.Printf("Debug mode enabled: adding /facts endpoints")
			r.Get("/facts", e.apiGetFacts)
		}
	})
}

func (e *Engine) apiAssert(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if room, err := e.searchRoom(id); err != nil {
		Error(w, http.StatusNotFound, "room not found")
		return
	} else {

		assertType, assertItem, err := game.GenericAssertHandler(&game.Config{Debug: e.Debug}, w, r)
		if err != nil {
			Error(w, http.StatusBadRequest, err.Error())
			return
		}
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/apiAssert]")+" ", 0)
			l.Printf("Asserting fact in room %s: %s", id, assertItem)
		}
		if err := room.clipsInstance.AssertFact(assertItem); err != nil {
			Error(w, http.StatusInternalServerError, "failed to assert")
			return
		}
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/apiAssert]")+" ", 0)
			l.Printf("Launching in room %s: run", id)
		}
		if err := room.clipsInstance.Run(); err != nil {
			Error(w, http.StatusInternalServerError, "failed to run")
			return
		}
		statusList, err := e.responseForType(assertType)
		if err != nil {
			Error(w, http.StatusInternalServerError, "failed to get status")
			return
		}

		response := make(map[string][]map[string]string)

		for _, statusItem := range statusList {
			if e.Debug {
				l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/apiAssert]")+" ", 0)
				l.Printf("Querying status in room %s: %s", id, statusItem)
			}
			if factList, err := room.clipsInstance.QueryFacts(statusItem); err != nil {
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiAssert]")+" ", 0)
					l.Printf("Error querying status in room %s - %s: %v", id, statusItem, err)
				}
				Error(w, http.StatusInternalServerError, "failed to query status")
				return
			} else {
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/apiAssert]")+" ", 0)
					l.Printf("Status in room %s - %s: %+v", id, statusItem, factList)
				}
				if factMap, err := genericFactToMap(e.Config, statusItem, factList); err != nil {
					if e.Debug {
						l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiAssert]")+" ", 0)
						l.Printf("Error converting fact to struct in room %s - %s: %v", id, statusItem, err)
					}
					Error(w, http.StatusInternalServerError, "failed to convert fact to struct")
					return
				} else {
					response[statusItem] = factMap
				}
			}
		}
		JSON(w, http.StatusOK, map[string]any{
			"status":   "asserted",
			"response": response,
		})
	}
}

func (e *Engine) apiGetFacts(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if room, err := e.searchRoom(id); err != nil {
		Error(w, http.StatusNotFound, "room not found")
		return
	} else {
		facts, err := room.clipsInstance.QueryFactsAllFacts()
		if err != nil {
			Error(w, http.StatusInternalServerError, "failed to get facts")
			return
		}
		JSON(w, http.StatusOK, map[string]any{
			"facts": facts,
		})
	}
}
