package server

import (
	"JotunBack/model"
	"JotunBack/repository"
	"encoding/json"
	"log"
	"net/http"
	"time"

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

func handleMessage(conn *websocket.Conn, hub *Hub, id string, userRepo *repository.UserRepository) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
		log.Printf("Received message from connection %s: %s\n", id, msg)

		//try to decode as Temp
		var temp model.Temp
		err = json.Unmarshal(msg, &temp)
		if err == nil {
			log.Printf("Received temperature: %+v\n", temp)
			var tempDB model.TempDB
			tempDB.Temperature = temp.Temperature
			tempDB.TimeStamp = time.Now()
			err := userRepo.CreateTemp(tempDB, id)
			if err != nil {
				log.Println("Error creating temp:", err)
				continue
			}
			continue
		}

		//try to decode as AirConditionerConfig
		var acConfig model.AirConditionerConfig
		err = json.Unmarshal(msg, &acConfig)
		if err == nil {
			log.Printf("Received settings: %+v\n", acConfig)
			//update protocol in Firestore
			currentState, err := userRepo.GetACState(id)
			if err != nil {
				log.Println("Error getting AC state:", err)
				continue
			}
			currentState.Protocol = acConfig.Protocol
			err = userRepo.UpdateACState(currentState)
			if err != nil {
				log.Println("Error updating AC state:", err)
				continue
			}
			continue
		}
	}
}

func HandleWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request, userRepo *repository.UserRepository) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		return
	}
	defer conn.Close()

	var settings model.AirConditionerConfig
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

	//save settings to Firestore if it doesn't exist
	_, err = userRepo.GetACState(id)
	if err != nil {
		err = userRepo.CreateACState(settings)
		if err != nil {
			log.Println("Error creating AC state:", err)
			return
		}
	}

	go handleMessage(conn, hub, id, userRepo)
}

func (hub *Hub) SendACConfig(acConfig model.AirConditionerConfig) {
	conn := hub.GetConnectionByID(acConfig.Username)
	if conn == nil {
		log.Printf("User %s not connected\n", acConfig.Username)
		return
	}
	data, err := json.Marshal(acConfig)
	if err != nil {
		log.Println("Error encoding JSON:", err)
		return
	}
	err = conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		log.Println("Error writing message:", err)
	}
}
