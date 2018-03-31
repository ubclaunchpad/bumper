package main

import (
	"encoding/json"
	"log"
	"math/rand"
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
	Key     int  `json:"key"`
	Pressed bool `json:"pressed"`
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
		log.Printf("%v\n", err)
		return
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

		junks := make([]models.Junk, 0)
		for _, junk := range g.Arena.Junk {
			junks = append(junks, *junk)
		}

		holes := make([]models.Hole, 0)
		for _, hole := range g.Arena.Holes {
			holes = append(holes, *hole)
		}

		players := make([]models.Player, 0)
		for _, player := range g.Arena.Players {
			players = append(players, *player)
		}

		msg := Message{
			Type: "update",
			Data: struct {
				Holes   []models.Hole   `json:"holes"`
				Junk    []models.Junk   `json:"junk"`
				Players []models.Player `json:"players"`
			}{
				holes,
				junks,
				players,
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
	rand.Seed(time.Now().UTC().UnixNano())
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
