/*
Copyright © 2026 Mirko Mariotti mirko@mirkomariotti.it
*/
package rulemancer

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func (e *Engine) querySubRoutes(r chi.Router) {
	r.Route("/", func(r chi.Router) {
		r.Post("/{relation}", e.apiQuery)
	})
}

func (e *Engine) apiQuery(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	statusItem := chi.URLParam(r, "relation")

	if room, err := e.searchRoom(id); err != nil {
		Error(w, http.StatusNotFound, "room not found")
		return
	} else {

		if !elementsInSlice(e.Config.Querables, statusItem) {
			if e.Debug {
				l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiQuery]")+" ", 0)
				l.Printf("Relation not querable in room %s: %s", id, statusItem)
			}
			Error(w, http.StatusBadRequest, "relation not querable")
			return
		}

		response := make(map[string][]map[string]string)

		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/apiQuery]")+" ", 0)
			l.Printf("Querying status in room %s: %s", id, statusItem)
		}
		if factList, err := room.clipsInstance.QueryFacts(statusItem); err != nil {
			if e.Debug {
				l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiQuery]")+" ", 0)
				l.Printf("Error querying status in room %s - %s: %v", id, statusItem, err)
			}
			Error(w, http.StatusInternalServerError, "failed to query status")
			return
		} else {
			if e.Debug {
				l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/apiQuery]")+" ", 0)
				l.Printf("Status in room %s - %s: %+v", id, statusItem, factList)
			}
			if factMap, err := genericFactToMap(e.Config, statusItem, factList); err != nil {
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiQuery]")+" ", 0)
					l.Printf("Error converting fact to struct in room %s - %s: %v", id, statusItem, err)
				}
				Error(w, http.StatusInternalServerError, "failed to convert fact to struct")
				return
			} else {
				response[statusItem] = factMap
			}
		}
		JSON(w, http.StatusOK, map[string]any{
			"response": response,
		})

	}
}
