package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/ubclaunchpad/bumper/server/game"
	"github.com/ubclaunchpad/bumper/server/models"
)

// Game represents a session
type Game struct {
	Arena *game.Arena
}

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

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		log.Printf("Accepting client from remote address %v\n", r.RemoteAddr)
		return true
	},
}

/*	ServeHTTP handles received messages from a client
Upgrades the connection to be persistent
Initializes the client connection to a map of clients
Listens for messages and acts on different message formats
*/
func (g *Game) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	name := r.URL.Query().Get("name")
	g.Arena.AddPlayer(ws)
	g.Arena.Players[ws].Name = name

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(g.Arena.Players, ws)
			break
		}

		if msg.Type == "keyHandler" {
			var kh KeyHandler
			err = json.Unmarshal([]byte(msg.Data.(string)), &kh)
			if err != nil {
				log.Printf("error: %v", err)
				continue
			}

			if kh.Pressed == true {
				g.Arena.Players[ws].KeyDownHandler(kh.Key)
			} else {
				g.Arena.Players[ws].KeyUpHandler(kh.Key)
			}
		}
	}
}

func run(g *Game) {
	for {
		g.Arena.UpdatePositions()
		g.Arena.CollisionDetection()
		time.Sleep(time.Millisecond * 17)
	}
}

func tick(g *Game) {
	for {
		time.Sleep(time.Millisecond * 17) // 60 Hz

		slice := make([]models.Player, 0)
		for _, val := range g.Arena.Players {
			slice = append(slice, *val)
		}

		msg := Message{
			Type: "update",
			Data: struct {
				Holes   []models.Hole   `json:"holes"`
				Junk    []models.Junk   `json:"junk"`
				Players []models.Player `json:"players"`
			}{
				g.Arena.Holes,
				g.Arena.Junk,
				slice,
			},
		}
		// update every client
		for client := range g.Arena.Players {
			err := client.WriteJSON(&msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(g.Arena.Players, client)
			}
		}
	}
}

func main() {
	game := Game{
		Arena: game.CreateArena(800, 1000),
	}

	http.Handle("/", http.FileServer(http.Dir("./build")))
	http.Handle("/connect", &game)
	go run(&game)
	go tick(&game)

	log.Println("Starting server on localhost:9090")
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("Error starting server")
	}
}
