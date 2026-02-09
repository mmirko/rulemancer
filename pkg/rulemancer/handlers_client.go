/*
Copyright © 2026 Mirko Mariotti mirko@mirkomariotti.it
*/
package rulemancer

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	jwtauth "github.com/go-chi/jwtauth/v5"
)

func (e *Engine) clientRoutes(r chi.Router) {
	r.Use(jwtauth.Verifier(e.JWTAuth))
	r.Use(jwtauth.Authenticator(e.JWTAuth))
	r.Route("/", func(r chi.Router) {
		r.Post("/create", e.apiCreateClient)
		r.Get("/list", e.apiListClients)
		r.Get("/{id}", e.apiGetClient)
		r.Delete("/{id}", e.apiDeleteClient)
	})
}

type CreateClientRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (e *Engine) apiCreateClient(w http.ResponseWriter, r *http.Request) {
	var req CreateClientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiCreateClient]")+" ", 0)
			l.Printf("Invalid JSON: %v", err)
		}
		Error(w, http.StatusBadRequest, "invalid json")
		return
	}

	client := e.newClient(req.Name, req.Description)

	_, tokenString, _ := e.Encode(map[string]interface{}{"id": client.id})

	if e.Debug {
		l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, green("[rulemancer/apiCreateClient]")+" ", 0)
		l.Printf("Creating client: %s with ID: %s", req.Name, client.id)
	}

	JSON(w, http.StatusCreated, map[string]string{
		"id":        client.id,
		"api_token": tokenString,
	})
}

func (e *Engine) apiGetClient(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, claims, err := jwtauth.FromContext(r.Context())
	requester := ""
	if err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiGetClient]")+" ", 0)
			l.Printf("Unauthorized get clientattempt: %v", err)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	} else if clientID, ok := claims["id"].(string); !ok || (clientID != "admin" && clientID != id) {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiGetClient]")+" ", 0)
			l.Printf("Unauthorized get client attempt by %s with invalid token: %v", requester, claims)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	} else {
		requester = clientID
	}

	if client, err := e.searchClient(id); err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/apiGetClient]")+" ", 0)
			l.Printf("Client not found: %v", err)
		}
		Error(w, http.StatusNotFound, "client not found")
		return
	} else {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, green("[rulemancer/apiGetClient]")+" ", 0)
			l.Printf("Client requested by %s found: %s", requester, client.id)
		}
		JSON(w, http.StatusOK, map[string]any{
			"id":          client.id,
			"name":        client.name,
			"description": client.description,
		})
	}
}

func (e *Engine) apiDeleteClient(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if _, err := e.removeClient(id); err != nil {
		Error(w, http.StatusNotFound, "client not found")
		return
	} else {
		JSON(w, http.StatusOK, map[string]string{
			"status": "deleted",
		})
	}
}

func (e *Engine) apiListClients(w http.ResponseWriter, r *http.Request) {
	clientsList := e.listClients()

	JSON(w, http.StatusOK, map[string]any{
		"clients": clientsList,
	})
}
