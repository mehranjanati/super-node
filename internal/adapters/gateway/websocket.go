package gateway

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for dev; restrict in prod
		},
	}
)

// WebSocketHandler manages real-time connections
type WebSocketHandler struct {
	clients   map[*websocket.Conn]bool
	broadcast chan []byte
	mu        sync.Mutex
}

// NewWebSocketHandler creates a new handler
func NewWebSocketHandler() *WebSocketHandler {
	return &WebSocketHandler{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan []byte),
	}
}

// HandleConnection handles incoming websocket requests
func (h *WebSocketHandler) HandleConnection(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	h.mu.Lock()
	h.clients[ws] = true
	h.mu.Unlock()

	log.Println("New WebSocket client connected")

	// Keep connection alive and listen for close
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			h.mu.Lock()
			delete(h.clients, ws)
			h.mu.Unlock()
			break
		}
	}
	return nil
}

// Broadcast sends a message to all connected clients
func (h *WebSocketHandler) Broadcast(message []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for client := range h.clients {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("WebSocket write error: %v", err)
			client.Close()
			delete(h.clients, client)
		}
	}
}
