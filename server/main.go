package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/ubclaunchpad/bumper/server/game"
	"github.com/ubclaunchpad/bumper/server/models"
)

// Message is the schema for client/server communication
type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// KeyHandler is the schema for client/server key handling communication
type KeyHandler struct {
	PlayerID int  `json:"playerID"`
	Key      int  `json:"key"`
	Pressed  bool `json:"pressed"`
} //TODO move to player?

var clients = make(map[*websocket.Conn]bool)

var players = []models.Player{
	models.Player{ID: 1, Position: models.Position{X: 200, Y: 200}, Velocity: models.Velocity{Dx: 0, Dy: 0}, Color: "blue"},
	models.Player{ID: 2, Position: models.Position{X: 250, Y: 250}, Velocity: models.Velocity{Dx: 0, Dy: 0}, Color: "red"},
	models.Player{ID: 3, Position: models.Position{X: 300, Y: 300}, Velocity: models.Velocity{Dx: 0, Dy: 0}, Color: "green"},
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		log.Printf("Accepting client from remote address %v\n", r.RemoteAddr)
		return true
	},
}

/*	handleConnection handles received messages from a client
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

	clients[ws] = true

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}

		if msg.Type == "keyHandler" {

			var kh KeyHandler
			err = json.Unmarshal([]byte(msg.Data.(string)), &kh)
			if err != nil {
				log.Printf("error: %v", err)
				continue
			}

			fmt.Println(kh)
			for _, player := range players {
				if player.ID == kh.PlayerID {
					if kh.Pressed == true {
						player.KeyDownHandler(kh.Key)
					} else {
						player.KeyUpHandler(kh.Key)
					}
				}
			}
		}
	}
}

func runGame() {
	a := game.CreateArena(400, 400)
	for {
		a.UpdatePositions()
		a.CollisionDetection()
	}
}

func tick() {
	holes := []models.Hole{
		models.Hole{Position: models.Position{X: 150, Y: 150}, Radius: 15},
		models.Hole{Position: models.Position{X: 100, Y: 100}, Radius: 15},
	}
	junk := []models.Junk{
		models.Junk{Position: models.Position{X: 50, Y: 50}, Velocity: models.Velocity{Dx: 0, Dy: 0}},
		models.Junk{Position: models.Position{X: 25, Y: 25}, Velocity: models.Velocity{Dx: 0, Dy: 0}},
	}
	for {
		time.Sleep(time.Millisecond * 17) // 60 Hz
		serverState := struct {
			Players []models.Player `json:"players"`
			Holes   []models.Hole   `json:"holes"`
			Junk    []models.Junk   `json:"junk"`
		}{
			players,
			holes,
			junk,
		}
		msg := Message{
			Type: "update",
			Data: serverState,
		}

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
		fmt.Fprint(w, "HELLO, is Inertia working yet?\n")
	})
	http.HandleFunc("/connect", handleConnection)
	go runGame()
	go tick()

	log.Println("Starting server on localhost:9090")
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("Error starting server")
	}
}
