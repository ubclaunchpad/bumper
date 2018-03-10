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
	players := []models.Player{
		models.Player{ID: 1, Position: models.Position{X: 200, Y: 200}, Velocity: models.Velocity{Dx: 0, Dy: 0}, Color: "blue"},
		models.Player{ID: 2, Position: models.Position{X: 250, Y: 250}, Velocity: models.Velocity{Dx: 0, Dy: 0}, Color: "red"},
		models.Player{ID: 3, Position: models.Position{X: 300, Y: 300}, Velocity: models.Velocity{Dx: 0, Dy: 0}, Color: "green"},
	}
	holes := []models.Hole{
		models.Hole{Position: models.Position{X: 150, Y: 150}, Radius: 10},
		models.Hole{Position: models.Position{X: 100, Y: 100}, Radius: 10},
	}
	junk := []models.Junk{
		models.Junk{Position: models.Position{X: 50, Y: 50}, Velocity: models.Velocity{Dx: 0, Dy: 0}},
		models.Junk{Position: models.Position{X: 25, Y: 25}, Velocity: models.Velocity{Dx: 0, Dy: 0}},
	}
	for {
		time.Sleep(time.Millisecond * 17) // 60 Hz
		serverState := struct {
			Type    string          `json:"type"`
			Players []models.Player `json:"players"`
			Holes   []models.Hole   `json:"holes"`
			Junk    []models.Junk   `json:"junk"`
		}{
			"update",
			players,
			holes,
			junk,
		}
		// update every client
		for client := range clients {
			err := client.WriteJSON(&serverState)
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
		fmt.Fprint(w, "is Inertia working yet?\n")
	})
	http.HandleFunc("/connect", handleConnection)
	go tick()

	log.Println("Starting server on localhost:9090")
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("Error starting server")
	}
}
