package rulemancer

import "errors"

type Client struct {
	name        string
	description string
	id          string
	rooms       []*Room
}

func (e *Engine) newClient(name, description string) *Client {
	e.clientsMutex.Lock()
	defer e.clientsMutex.Unlock()
	client := &Client{
		name:        name,
		description: description,
		id:          e.generateClientUniqueID(),
		rooms:       make([]*Room, 0),
	}

	e.clients[client.id] = client
	return client
}

func (e *Engine) generateClientUniqueID() string {
	for {
		newId := RandStringBytes(16)
		if _, exists := e.clients[newId]; !exists {
			return newId
		}
	}
}

func (e *Engine) searchClient(id string) (*Client, error) {
	e.clientsMutex.RLock()
	defer e.clientsMutex.RUnlock()
	client, exists := e.clients[id]
	if !exists {
		return nil, errors.New("client not found")
	}
	return client, nil
}

func (e *Engine) removeClient(id string) (*Client, error) {
	e.clientsMutex.Lock()
	defer e.clientsMutex.Unlock()
	if client, exists := e.clients[id]; exists {
		delete(e.clients, id)
		return client, nil
	}
	return nil, errors.New("client not found")
}

func (e *Engine) listClients() []string {
	e.clientsMutex.RLock()
	defer e.clientsMutex.RUnlock()
	clients := make([]string, 0, len(e.clients))
	for _, client := range e.clients {
		clients = append(clients, client.id)
	}
	return clients
}
