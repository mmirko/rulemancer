/*
Copyright © 2026 Mirko Mariotti mirko@mirkomariotti.it
*/
package rulemancer

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mmirko/rulemancer/pkg/game"
)

func (e *Engine) roomSubRoutes(r chi.Router) {
	r.Route("/", func(r chi.Router) {
		r.Post("/assert", e.apiAssert)
		r.Get("/query", e.apiQuery)
		r.Get("/facts", e.apiGetFacts)

	})
}

func (e *Engine) apiAssert(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if room, err := e.searchRoom(id); err != nil {
		Error(w, http.StatusNotFound, "room not found")
		return
	} else {

		result, err := game.GenericHandler(&game.Config{Debug: e.Debug}, w, r)
		if err != nil {
			Error(w, http.StatusBadRequest, err.Error())
			return
		}
		if e.Debug {
			fmt.Printf("Asserting fact in room %s: %s\n", id, result)
		}
		if err := room.clipsInstance.AssertFact(result); err != nil {
			Error(w, http.StatusInternalServerError, "failed to assert")
			return
		}
		if e.Debug {
			fmt.Printf("Launching in room %s: run\n", id)
		}
		if err := room.clipsInstance.Run(); err != nil {
			Error(w, http.StatusInternalServerError, "failed to run")
			return
		}
		JSON(w, http.StatusOK, map[string]any{
			"status": "asserted",
		})
	}
}

func (e *Engine) apiQuery(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if room, err := e.searchRoom(id); err != nil {
		Error(w, http.StatusNotFound, "room not found")
		return
	} else {
		JSON(w, http.StatusOK, map[string]any{
			"clips_instance": room.clipsInstance.Info(),
		})
	}
}

func (e *Engine) apiGetFacts(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if room, err := e.searchRoom(id); err != nil {
		Error(w, http.StatusNotFound, "room not found")
		return
	} else {
		facts, err := room.clipsInstance.QueryFacts("*")
		if err != nil {
			Error(w, http.StatusInternalServerError, "failed to get facts")
			return
		}
		JSON(w, http.StatusOK, map[string]any{
			"facts": facts,
		})
	}
}
