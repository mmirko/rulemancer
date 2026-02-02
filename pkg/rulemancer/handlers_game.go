/*
Copyright © 2026 Mirko Mariotti mirko@mirkomariotti.it
*/
package rulemancer

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (e *Engine) gameRoutes(r chi.Router) {
	r.Route("/", func(r chi.Router) {
		r.Get("/list", e.apiListGames)
		r.Get("/{id}", e.apiGetGame)
	})
}

func (e *Engine) apiGetGame(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if game, err := e.searchGame(id); err != nil {
		Error(w, http.StatusNotFound, "game not found")
		return
	} else {
		JSON(w, http.StatusOK, map[string]any{
			"id":          game.id,
			"name":        game.name,
			"description": game.description,
			"rules":       game.rulesLocation,
			"assertable":  game.assertable,
			"responses":   game.responses,
			"queryable":   game.queryable,
		})
	}
}

func (e *Engine) apiListGames(w http.ResponseWriter, r *http.Request) {
	gamesList := e.listGames()

	JSON(w, http.StatusOK, map[string]any{
		"games": gamesList,
	})
}
