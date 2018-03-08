package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/ubclaunchpad/bumper/server/models"
)

// Message is the schema for client/server communication
type Message struct {
	Type     string          `json:"type"`
	ID       int             `json:"id"`
	Position models.Position `json:"position"` //Position variable takes Position struct as datatype
	Message  string          `json:"message"`
	Color    string          `json:"color"`
}

var p models.Player //DELETE THIS LATER, USED AS WORKAROUND FOR IMPORTING

// ServerState map of states for all players
type ServerState struct {
	Type    string        `json:"type"`
	Players []ObjectState `json:"players"`
}

// ObjectState of an object, position, velocity, id
type ObjectState struct {
	ID       int             `json:"id"`
	Position models.Position `json:"position"`
	Color    string          `json:"color"`
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

	// record this connection in our map
	// initialize state struct
	clients[ws] = &ObjectState{
		0,
		models.Position{0, 0},
		"",
	}

	// infinite loop that receives msgs from clients
	for {
		var msg Message
		err := ws.ReadJSON(&msg)
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
			clients[ws].ID = reply.ID
			clients[ws].Color = msg.Color
		} else {
			//update player in map
			clients[ws].Position.X = msg.Position.X
			clients[ws].Position.Y = msg.Position.Y
		}
	}
}

func tick() {
	for {
		time.Sleep(time.Millisecond * 17) // 60 Hz
		var objectarray []ObjectState
		msg := ServerState{Type: "update"}

		for client := range clients {
			objectarray = append(objectarray, *clients[client])
		}

		msg.Players = objectarray
		// update every client
		for client := range clients {
			err := client.WriteJSON(&msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Hello is Inertia working yet?\n")
	})
	http.HandleFunc("/connect", handleConnection)
	go tick()

	log.Println("Starting server on localhost:9090")
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("Error starting server")
	}
}
