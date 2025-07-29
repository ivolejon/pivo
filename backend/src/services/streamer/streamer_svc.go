package streamer_svc

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var streamHub *Hub

func GetHubInstance() *Hub {
	return streamHub
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for simplicity in this example.
		// In a production environment, you should restrict this to your frontend's domain.
		return true
	},
}

type Client struct {
	ID   string // Unique client ID (UUID)
	Conn *websocket.Conn
	Send chan []byte // Buffered channel for outbound messages
}

type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Inbound messages from the applications to be broadcast.
	broadcast chan []byte

	// Mutex to protect client map operations
	mu sync.RWMutex
}

func CreateStreamHub() *Hub {
	if streamHub != nil {
		log.Println("Stream hub already exists, skipping creation.")
		return streamHub
	}
	streamHub = &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
	return streamHub
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("Client registered: %s (Total clients: %d)", client.ID, len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if h.clients[client] {
				delete(h.clients, client)
				close(client.Send)
			}
			h.mu.Unlock()
			log.Printf("Client unregistered: %s (Total clients: %d)", client.ID, len(h.clients))

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.Send <- message:
					// Message sent successfully
				default:
					// If sending fails (e.g., channel full, client gone), unregister
					close(client.Send)
					delete(h.clients, client)
					log.Printf("Failed to send message to client %s, unregistering.", client.ID)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) PublishToClient(clientID string, message []byte) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	found := false
	for client := range h.clients {
		if client.ID == clientID {
			select {
			case client.Send <- message:
				log.Printf("Message sent to specific client %s: %s", clientID, string(message))
				found = true
			default:
				log.Printf("Failed to send message to specific client %s (channel full/closed)", clientID)
			}
			break // Client found, no need to continue iterating
		}
	}
	if !found {
		log.Printf("Client with ID %s not found.", clientID)
	}
	return found
}

func (c *Client) writePump() {
	ticker := time.NewTicker(time.Second * 10) // Ping interval
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second)) // Set a write deadline
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// **Modification Starts Here**
			// Send only the current message.
			// Do not attempt to drain and send multiple messages in one WebSocket frame.
			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("Error writing message for client %s: %v", c.ID, err)
				return
			}
			// **Modification Ends Here**

		case <-ticker.C:
			// Send a ping message to keep the connection alive
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Ping failed for client %s: %v", c.ID, err)
				return // Connection likely broken
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func ServeWs(hub *Hub, c *gin.Context) {
	// Get clientID from query parameter
	clientID := c.Query("client_id")
	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "clientID query parameter is required"})
		return
	}
	// Validate if clientID is a valid UUID
	if _, err := uuid.Parse(clientID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid clientID format (must be a UUID)"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{ID: clientID, Conn: conn, Send: make(chan []byte, 256)}
	hub.register <- client

	// Start a goroutine to handle sending messages to the client
	go client.writePump()

	// Since the requirement is only to publish messages to the frontend,
	// we don't need a readPump for incoming messages from the client.
	// However, we need to keep the connection alive and handle disconnects.
	// We can do this by continuously reading from the connection, even if we discard the messages.
	// This also allows the client to send close messages.
	defer func() {
		hub.unregister <- client
		conn.Close()
	}()
	for {
		// Read a message to detect if the connection is closed by the client
		// or if there's an error. We don't care about the content.
		_, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading from client %s: %v", client.ID, err)
			}
			break // Exit the loop on error or close
		}
		// If we receive a message, we just discard it as per requirement
	}
}
