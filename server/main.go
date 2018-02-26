package main

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

// Position x y position
type Position struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

// Message is the schema for client/server communication
type Message struct {
	Type     string   `json:"type"`
	ID       int      `json:"id"`
	Position Position `json:"pos"` //Position variable takes Position struct as datatype
	Message  string   `json:"message"`
}

// ServerState map of states for all players
type ServerState struct {
	Players []ObjectState
}

// ObjectState of an object, position, velocity, id
type ObjectState struct {
	ID       int
	Position Position
}

var clients = make(map[*websocket.Conn]*ObjectState)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		log.Printf("Accepting client from remote address %v\n", r.RemoteAddr)
		return true
	},
}

/* 	handleConnection handles received messages from a client
Upgrades the connection to be persistent
Initializes the client connection to a map of clients
Listens for messages and acts on different message formats
*/
func handleConnection(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()
	log.Println("accepted client")

	// record this connection in our map
	// initialize state struct
	clients[ws] = &ObjectState{
		0,
		Position{0, 0},
	}

	// infinite loop that receives msgs from clients
	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		log.Printf("Message Received: %+v\n", msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}

		// initial message received
		if msg.Type == "initial" {
			reply := &Message{
				Type: "initial",
				ID:   rand.Intn(1000),
			}

			err := ws.WriteJSON(reply)
			if err != nil {
				log.Printf("error sending message: %v", err)
				ws.Close()
				delete(clients, ws)
			}
			//add player to map
			clients[ws] = &ObjectState{reply.ID, Position{}}
		} else { //update player in map
			clients[ws].Position.X = msg.Position.X
			clients[ws].Position.Y = msg.Position.Y
		}
		log.Printf("Client %d State: %+v\n", msg.ID, *clients[ws])
		log.Printf("Message Type: %+v\n", msg.Type)
	}
}

func tick() {

	// for client := range clients {
	// 	if clients[client].ID == msg.ID {
	// 		clients[client] = ObjectState{msg.ID, Position{5, 5}}}
	// 	}
	// }

	tickCount := 0
	for {
		time.Sleep(time.Second * 5)
		msg := Message{
			Message: "tick" + strconv.Itoa(tickCount),
		}
		// update every client
		for client := range clients {
			err := client.WriteJSON(&msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
		tickCount++
	}
}

func main() {
	http.HandleFunc("/connect", handleConnection)
	go tick()

	log.Println("Starting server on localhost:9090")
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("Error starting server")
	}
}
