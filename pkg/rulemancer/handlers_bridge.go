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

func (e *Engine) bridgeRoutes(r chi.Router) {
	r.Use(jwtauth.Verifier(e.JWTAuth))
	r.Use(jwtauth.Authenticator(e.JWTAuth))
	r.Route("/", func(r chi.Router) {
		r.Get("/list", e.apiListBridges)
		r.Get("/{id}", e.apiGetBridge)
	})
}

func (e *Engine) apiGetBridge(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiGetBridge]")+" ", 0)
			l.Printf("Unauthorized get bridge attempt: %v", err)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	} else if clientID, ok := claims["id"].(string); !ok || clientID != "admin" {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiGetBridge]")+" ", 0)
			l.Printf("Unauthorized get bridge attempt with invalid token: %v", claims)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if bridge, err := e.searchBridge(id); err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiGetBridge]")+" ", 0)
			l.Printf("Bridge not found: %v", err)
		}
		Error(w, http.StatusNotFound, "bridge not found")
		return
	} else {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, green("[rulemancer/apiGetBridge]")+" ", 0)
			l.Printf("Bridge %s info provided to client admin", id)
		}
		JSON(w, http.StatusOK, map[string]any{
			"id":    bridge.id,
			"name":  bridge.name,
			"rules": bridge.rulesLocation,
		})
	}
}

func (e *Engine) apiListBridges(w http.ResponseWriter, r *http.Request) {

	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiListBridges]")+" ", 0)
			l.Printf("Unauthorized list bridges attempt: %v", err)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	} else if _, ok := claims["id"].(string); !ok {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiListBridges]")+" ", 0)
			l.Printf("Unauthorized list bridges attempt with invalid token: %v", claims)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if e.Debug {
		l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, green("[rulemancer/apiListBridges]")+" ", 0)
		l.Printf("Listing all bridges")
	}

	bridgesList := e.listBridges()

	JSON(w, http.StatusOK, map[string]any{
		"bridges": bridgesList,
	})
}
