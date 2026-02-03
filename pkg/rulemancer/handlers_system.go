/*
Copyright © 2026 Mirko Mariotti mirko@mirkomariotti.it
*/
package rulemancer

import (
	"net/http"
	"syscall"

	chi "github.com/go-chi/chi/v5"
	jwtauth "github.com/go-chi/jwtauth/v5"
)

func (e *Engine) systemRoutes(r chi.Router) {
	r.Use(jwtauth.Verifier(e.JWTAuth))
	r.Use(jwtauth.Authenticator(e.JWTAuth))
	r.Route("/", func(r chi.Router) {
		r.Get("/health", e.health)
		r.Post("/quit", e.quit)
	})
}

func (e *Engine) health(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	if clientID, ok := claims["id"].(string); !ok || clientID != "admin" {
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	JSON(w, http.StatusOK, map[string]string{"status": "OK"})
}

func (e *Engine) quit(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	if clientID, ok := claims["id"].(string); !ok || clientID != "admin" {
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	e.stopChan <- syscall.SIGTERM
	JSON(w, http.StatusOK, map[string]string{"status": "shutting down"})
}
