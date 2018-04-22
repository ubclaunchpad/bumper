package main

import (
	"encoding/json"
	"fmt"
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

// KeyHandler is the schema for client/server key handling communication
type KeyHandler struct {
	Key     int  `json:"key"`
	Pressed bool `json:"pressed"`
} //TODO move to player?

// SpawnHandler is the schema for client/server key handling communication
type SpawnHandler struct {
	Name string `json:"name"`
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

	for {
		var msg models.Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(g.Arena.Players, ws)
			break
		}
		if msg.Type == "spawn" {
			var spawn SpawnHandler
			err = json.Unmarshal([]byte(msg.Data.(string)), &spawn)
			if err != nil {
				log.Printf("error: %v", err)
				continue
			}
			name := spawn.Name
			g.Arena.AddPlayer(ws)
			g.Arena.Players[ws].Name = name

			msg := models.Message{
				"initial",
				ws,
			}
			game.MessageChannel <- msg
		}

		if msg.Type == "keyHandler" {
			var kh KeyHandler
			err = json.Unmarshal([]byte(msg.Data.(string)), &kh)
			if err != nil {
				log.Printf("error: %v", err)
				continue
			}
			if _, ok := g.Arena.Players[ws]; ok {
				if kh.Pressed == true {
					g.Arena.Players[ws].KeyDownHandler(kh.Key)
				} else {
					g.Arena.Players[ws].KeyUpHandler(kh.Key)
				}
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

		msg := models.Message{
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
			p := g.Arena.Players[client]
			p.SocketLock.Lock()
			err := client.WriteJSON(&msg)
			p.SocketLock.Unlock()
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(g.Arena.Players, client)
			}
		}
	}
}

func sendMessage(g *Game) {
	for {
		msg := <-game.MessageChannel

		switch msg.Type {
		case "initial":
			ws := msg.Data.(*websocket.Conn)
			initalMsg := models.Message{
				Type: "initial",
				Data: struct {
					ArenaWidth  float64 `json:"arenawidth"`
					ArenaHeight float64 `json:"arenaheight"`
					PlayerID    string  `json:"playerid"`
				}{
					g.Arena.Width,
					g.Arena.Height,
					g.Arena.Players[ws].Color,
				},
			}
			g.Arena.Players[ws].SocketLock.Lock()
			error := ws.WriteJSON(&initalMsg)
			g.Arena.Players[ws].SocketLock.Unlock()
			if error != nil {
				log.Printf("error: %v", error)
				ws.Close()
				delete(g.Arena.Players, ws)
			}
		case "death":
			fmt.Println("GOT DEATH message in channel")
			ws := msg.Data.(*websocket.Conn)
			deathMsg := models.Message{
				Type: "death",
				Data: nil,
			}
			g.Arena.Players[ws].SocketLock.Lock()
			error := ws.WriteJSON(&deathMsg)
			g.Arena.Players[ws].SocketLock.Unlock()
			if error != nil {
				log.Printf("error: %v", error)
				ws.Close()
				delete(g.Arena.Players, ws)
			}
			delete(g.Arena.Players, ws)
		}
		fmt.Println(msg.Type)
	}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	game.MessageChannel = make(chan models.Message)
	game := Game{
		Arena: game.CreateArena(800, 1000),
	}

	http.Handle("/", http.FileServer(http.Dir("./build")))
	http.Handle("/connect", &game)
	go sendMessage(&game)
	go run(&game)
	go tick(&game)

	log.Println("Starting server on localhost:9090")
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("Error starting server")
	}
}
