package server

import (
	"JotunBack/model"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Hub struct {
	connections map[string]*websocket.Conn
}

func NewHub() *Hub {
	return &Hub{
		connections: make(map[string]*websocket.Conn),
	}
}

func (hub *Hub) AddConnection(conn *websocket.Conn, userName string) string {
	id := userName
	hub.connections[id] = conn
	return id
}

func (hub *Hub) GetConnectionByID(id string) *websocket.Conn {
	return hub.connections[id]
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		return
	}
	defer conn.Close()

	var settings model.AirConditionerSettings
	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Println("Error reading message:", err)
		return
	}
	err = json.Unmarshal(msg, &settings)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		return
	}
	log.Printf("Received settings: %+v\n", settings)

	id := hub.AddConnection(conn, settings.Username)
	defer delete(hub.connections, id)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
		log.Printf("Received message: %s\n", msg)

		err = conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("Error writing message:", err)
			break
		}
	}
}

func (hub *Hub) Unicast(id string, msg []byte) error {
	conn := hub.GetConnectionByID(id)
	if conn == nil {
		return nil
	}
	return conn.WriteMessage(websocket.TextMessage, msg)
}
