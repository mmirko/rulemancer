/*
Copyright © 2026 Mirko Mariotti mirko@mirkomariotti.it
*/
package rulemancer

import (
	"net/http"
	"syscall"

	chi "github.com/go-chi/chi/v5"
)

func (e *Engine) systemRoutes(r chi.Router) {
	r.Route("/", func(r chi.Router) {
		r.Get("/health", e.health)
		r.Post("/quit", e.quit)
	})
}

func (e *Engine) health(w http.ResponseWriter, r *http.Request) {
	JSON(w, http.StatusOK, map[string]string{"status": "OK"})
}

func (e *Engine) quit(w http.ResponseWriter, r *http.Request) {
	e.stopChan <- syscall.SIGTERM
	JSON(w, http.StatusOK, map[string]string{"status": "shutting down"})
}
