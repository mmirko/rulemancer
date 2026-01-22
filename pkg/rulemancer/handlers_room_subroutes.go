/*
Copyright © 2026 Mirko Mariotti mirko@mirkomariotti.it
*/
package rulemancer

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (e *Engine) roomSubRoutes(r chi.Router) {
	r.Route("/", func(r chi.Router) {
		r.Get("/assert", e.apiAssert)
		r.Get("/query", e.apiQuery)
	})
}

func (e *Engine) apiAssert(w http.ResponseWriter, r *http.Request) {
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
