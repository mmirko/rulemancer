/*
Copyright © 2026 Mirko Mariotti mirko@mirkomariotti.it
*/
package rulemancer

import (
	"log"
	"net/http"
	"os"

	chi "github.com/go-chi/chi/v5"
	jwtauth "github.com/go-chi/jwtauth/v5"
)

func (e *Engine) joinRoutes(r chi.Router) {
	r.Use(jwtauth.Verifier(e.JWTAuth))
	r.Use(jwtauth.Authenticator(e.JWTAuth))
	r.Route("/", func(r chi.Router) {
		r.Post("/available/{gameRef}", e.availableRoom) // Join the first available room for the specified game
		r.Post("/room/{roomID}", e.joinRoom)            // Join a specific room by ID
		r.Post("/new/{gameRef}", e.newGameRoom)         // Create a new room for the specified game and join it
	})
}

func (e *Engine) availableRoom(w http.ResponseWriter, r *http.Request) {
	// TODO
	JSON(w, http.StatusOK, map[string]string{"status": "OK"})
}

func (e *Engine) joinRoom(w http.ResponseWriter, r *http.Request) {
	roomId := chi.URLParam(r, "roomID")
	_, claims, err := jwtauth.FromContext(r.Context())

	if err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/joinRoom]")+" ", 0)
			l.Printf("Unauthorized join room attempt: %v", err)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	} else if clientID, ok := claims["id"].(string); !ok {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/joinRoom]")+" ", 0)
			l.Printf("Unauthorized join room attempt with invalid token: %v", claims)
		}
		Error(w, http.StatusUnauthorized, "unauthorized")
		return
	} else {
		if room, err := e.searchRoom(roomId); err != nil {
			// Room existence
			if e.Debug {
				l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/joinRoom]")+" ", 0)
				l.Printf("Room not found: %s", roomId)
			}
			Error(w, http.StatusNotFound, "room not found")
			return
		} else if client, err := e.searchClient(clientID); err != nil {
			// Client existence
			if e.Debug {
				l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/joinRoom]")+" ", 0)
				l.Printf("Client not found: %s", clientID)
			}
			Error(w, http.StatusNotFound, "client not found")
			return
		} else {
			// Start locking the room
			room.clientsMutex.Lock()
			defer room.clientsMutex.Unlock()

			if len(room.clients) >= room.maxClients {
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/joinRoom]")+" ", 0)
					l.Printf("Room is full: %s", roomId)
				}
				Error(w, http.StatusForbidden, "room is full")
				return
			}
			if _, exists := room.clients[clientID]; exists {
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/joinRoom]")+" ", 0)
					l.Printf("Client already in room: %s", roomId)
				}
				Error(w, http.StatusConflict, "client already in room")
				return
			}

			// Ok the room! now lock the client
			client.roomsMutex.Lock()
			defer client.roomsMutex.Unlock()

			if _, exists := client.playingRooms[roomId]; exists {
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/joinRoom]")+" ", 0)
					l.Printf("Client already playing in room: %s", roomId)
				}
				Error(w, http.StatusConflict, "client already playing in room")
				return
			}

			// Apply the join to both the room and the client
			room.clients[clientID] = client
			client.playingRooms[roomId] = room
			if e.Debug {
				l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, green("[rulemancer/joinRoom]")+" ", 0)
				l.Printf("Client joined room: %s", roomId)
			}
			JSON(w, http.StatusOK, map[string]string{"status": "joined"})
			return
		}
	}

}

func (e *Engine) newGameRoom(w http.ResponseWriter, r *http.Request) {
	// TODO
	JSON(w, http.StatusOK, map[string]string{"status": "OK"})

}
