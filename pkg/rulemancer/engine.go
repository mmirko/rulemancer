/*
Copyright © 2026 Mirko Mariotti mirko@mirkomariotti.it
*/
package rulemancer

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Engine struct {
	*Config
	games        map[string]*Game
	gamesMutex   sync.RWMutex
	rooms        map[string]*Room
	roomsMutex   sync.RWMutex
	clients      map[string]*Client
	clientsMutex sync.RWMutex
	router       chi.Router
	stopChan     chan os.Signal
}

func NewEngine() *Engine {
	return &Engine{
		Config:       NewConfig(),
		games:        make(map[string]*Game),
		gamesMutex:   sync.RWMutex{},
		rooms:        make(map[string]*Room),
		roomsMutex:   sync.RWMutex{},
		clients:      make(map[string]*Client),
		clientsMutex: sync.RWMutex{},
		router:       chi.NewRouter(),
		stopChan:     make(chan os.Signal, 1),
	}
}

func (e *Engine) SpawnEngine() error {
	// Implement the logic to spawn and run the CLIPS engine
	// using the provided configuration and rule pool directory

	e.loadGames()

	r := e.router
	c := e.Config

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	if e.Debug {
		r.Use(middleware.Logger)
	}

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/system", e.systemRoutes)
		r.Route("/room", e.roomRoutes)
		r.Route("/game", e.gameRoutes)
	})

	srv := &http.Server{
		Addr:    ":3000",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServeTLS(c.TLSCertFile, c.TLSKeyFile); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %v\n", err)
		}
	}()

	signal.Notify(e.stopChan, os.Interrupt, syscall.SIGTERM)
	<-e.stopChan

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")

	return nil
}
