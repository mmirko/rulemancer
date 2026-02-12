package rulemancer

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"

	chi "github.com/go-chi/chi/v5"
	jwtauth "github.com/go-chi/jwtauth/v5"
	"github.com/gorilla/websocket"
)

func (e *Engine) Monitor(url string) error {

	dialer := websocket.Dialer{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	// Read API_TOKEN from environment
	apiToken := os.Getenv("API_TOKEN")
	if apiToken == "" {
		return fmt.Errorf("API_TOKEN environment variable is not set")
	}

	header := http.Header{}
	header.Set("Authorization", "Bearer "+apiToken)

	conn, _, err := dialer.Dial(url, header)
	if err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[cmd/monitor]")+" ", 0)
			l.Println("dial error:", err)
		}
		return fmt.Errorf("dial error: %w", err)
	}
	defer conn.Close()

	if e.Debug {
		l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[cmd/monitor]")+" ", 0)
		l.Println("connected to", url)
	}

	// reader async
	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[cmd/monitor]")+" ", 0)
					l.Println("read error:", err)
				}
				return
			}
			fmt.Println("server:", string(msg))
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		text := scanner.Text()

		err := conn.WriteMessage(websocket.TextMessage, []byte(text))
		if err != nil {
			if e.Debug {
				l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[cmd/monitor]")+" ", 0)
				l.Println("write error:", err)
			}
			return fmt.Errorf("write error: %w", err)
		}
	}
	return nil
}

func (e *Engine) systemMonitor(w http.ResponseWriter, r *http.Request) {

	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			_, claims, err := jwtauth.FromContext(r.Context())
			if err != nil {
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/systemMonitor]")+" ", 0)
					l.Println("JWT error:", err)
				}
				return false
			} else if clientID, ok := claims["id"].(string); !ok || clientID != "admin" {
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/systemMonitor]")+" ", 0)
					l.Println("Unauthorized client ID:", clientID)
				}
				return false
			}
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/systemMonitor]")+" ", 0)
			l.Println("upgrade error:", err)
		}
		return
	}
	defer conn.Close()

	if e.Debug {
		l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/systemMonitor]")+" ", 0)
		l.Println("client connected")
	}

	ctx, cancel := context.WithCancel(context.Background())

	wsIn := make(chan []byte)
	wsOut := make(chan []byte)
	wsErr := make(chan error)

	// reader async
	go func() {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/systemMonitor]")+" ", 0)
			l.Println("reader started")
		}
	loop:
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/systemMonitor]")+" ", 0)
					l.Println("read error:", err)
				}
				select {
				case <-ctx.Done():
					break loop
				case wsErr <- err:
				}
				break loop
			}
			select {
			case <-ctx.Done():
				break loop
			case wsIn <- msg:
			}
		}
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/systemMonitor]")+" ", 0)
			l.Println("reader stopped")
		}
	}()

	// writer async
	go func() {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/systemMonitor]")+" ", 0)
			l.Println("writer started")
		}
	loop:
		for {
			select {
			case <-ctx.Done():
				break loop
			case msg := <-wsOut:
				err := conn.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					if e.Debug {
						l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/systemMonitor]")+" ", 0)
						l.Println("write error:", err)
					}
					select {
					case <-ctx.Done():
						break loop
					case wsErr <- err:
					}
					break loop
				}
			}
		}
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/systemMonitor]")+" ", 0)
			l.Println("writer stopped")
		}
	}()

	// error handler
	go func() {
		<-wsErr
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/systemMonitor]")+" ", 0)
			l.Println("connection error, closing monitor")
		}
		cancel()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-wsIn:
			if e.Debug {
				l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/systemMonitor]")+" ", 0)
				l.Printf("websocket message received: %s\n", msg)
			}
			wsOut <- []byte("Message received, but repl not implemented")
		}
	}
}

func (e *Engine) roomMonitor(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var room *Room

	if r, err := e.searchRoom(id); err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/roomMonitor]")+" ", 0)
			l.Println("room not found:", id)
		}
		Error(w, http.StatusNotFound, "room not found")
		return
	} else {
		room = r
	}

	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			_, claims, err := jwtauth.FromContext(r.Context())
			if err != nil {
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/roomMonitor]")+" ", 0)
					l.Println("JWT error:", err)
				}
				return false
			} else if clientID, ok := claims["id"].(string); !ok {
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/roomMonitor]")+" ", 0)
					l.Println("Unauthorized client ID:", clientID)
				}
				return false

			} else {

				canQuery := false
				room.clientsMutex.RLock()
				room.watchersMutex.RLock()
				if _, ok := room.clients[clientID]; ok {
					canQuery = true
				}
				if _, ok := room.watchers[clientID]; ok {
					canQuery = true
				}
				room.watchersMutex.RUnlock()
				room.clientsMutex.RUnlock()

				return canQuery
			}
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/roomMonitor]")+" ", 0)
			l.Println("upgrade error:", err)
		}
		return
	}

	recvChan := make(socketChan)

	room.socketsMutex.Lock()
	room.sockets[conn] = recvChan
	room.socketsMutex.Unlock()
	defer func() {
		room.socketsMutex.Lock()
		delete(room.sockets, conn)
		room.socketsMutex.Unlock()
		conn.Close()
	}()

	if e.Debug {
		l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/roomMonitor]")+" ", 0)
		l.Println("client connected")
	}

	ctx, cancel := context.WithCancel(context.Background())

	wsIn := make(chan []byte)
	wsOut := make(chan []byte)
	wsErr := make(chan error)

	// reader async
	go func() {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/roomMonitor]")+" ", 0)
			l.Println("reader started for room", id)
		}
	loop:
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				if e.Debug {
					l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/roomMonitor]")+" ", 0)
					l.Println("read error:", err)
				}
				select {
				case <-ctx.Done():
					break loop
				case wsErr <- err:
				}
				break loop
			}
			select {
			case <-ctx.Done():
				break loop
			case wsIn <- msg:
			}
		}
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/roomMonitor]")+" ", 0)
			l.Println("reader stopped for room", id)
		}
	}()

	// writer async
	go func() {
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/roomMonitor]")+" ", 0)
			l.Println("writer started for room", id)
		}
	loop:
		for {
			select {
			case <-ctx.Done():
				break loop
			case msg := <-wsOut:
				err := conn.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					if e.Debug {
						l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/roomMonitor]")+" ", 0)
						l.Println("write error:", err)
					}
					select {
					case <-ctx.Done():
						break loop
					case wsErr <- err:
					}
					break loop
				}
			}
		}
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/roomMonitor]")+" ", 0)
			l.Println("writer stopped for room", id)
		}
	}()

	// error handler
	go func() {
		<-wsErr
		if e.Debug {
			l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, red("[rulemancer/roomMonitor]")+" ", 0)
			l.Println("connection error, closing monitor for room", id)
		}
		cancel()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-wsIn:
			if e.Debug {
				l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/roomMonitor]")+" ", 0)
				l.Printf("websocket message received for room %s: %s\n", id, msg)
			}
			wsOut <- []byte("Message received, but repl not implemented")
		case msg := <-recvChan:
			if e.Debug {
				l := log.New(&writer{os.Stdout, "2006-01-02 15:04:05 "}, yellow("[rulemancer/roomMonitor]")+" ", 0)
				l.Printf("message received from room %s: %s\n", id, msg)
			}
			wsOut <- msg.message
		}
	}
}
