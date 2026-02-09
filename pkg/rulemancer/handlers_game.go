/*
Copyright © 2026 Mirko Mariotti mirko@mirkomariotti.it
*/
package rulemancer

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	jwtauth "github.com/go-chi/jwtauth/v5"
)

func (e *Engine) gameRoutes(r chi.Router) {
	r.Use(jwtauth.Verifier(e.JWTAuth))
	r.Use(jwtauth.Authenticator(e.JWTAuth))
	r.Route("/", func(r chi.Router) {
		r.Get("/list", e.apiListGames)
		r.Get("/{id}", e.apiGetGame)
	})
}

func (e *Engine) apiGetGame(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if game, err := e.searchGame(id); err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiGetGame]")+" ", 0)
			l.Printf("Game not found: %v", err)
		}
		Error(w, http.StatusNotFound, "game not found")
		return
	} else {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, green("[rulemancer/apiGetGame]")+" ", 0)
			l.Printf("Game %s info provided to client %s", game.id, id)
		}
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
