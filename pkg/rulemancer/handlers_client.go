/*
Copyright © 2026 Mirko Mariotti mirko@mirkomariotti.it
*/
package rulemancer

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (e *Engine) clientRoutes(r chi.Router) {
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
		Error(w, http.StatusBadRequest, "invalid json")
		return
	}

	client := e.newClient(req.Name, req.Description)

	_, tokenString, _ := e.Encode(map[string]interface{}{"id": client.id})

	JSON(w, http.StatusCreated, map[string]string{
		"id":        client.id,
		"api_token": tokenString,
	})
}

func (e *Engine) apiGetClient(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if client, err := e.searchClient(id); err != nil {
		Error(w, http.StatusNotFound, "client not found")
		return
	} else {
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
