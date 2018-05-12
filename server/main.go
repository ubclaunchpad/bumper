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
	Arena       *game.Arena
	RefreshRate time.Duration
}

// MessageChannel is used by the server to emit messages to a client
var MessageChannel chan models.Message

// An instance of Upgrader that upgrades a connection to a WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		log.Printf("Accepting client from remote address %v\n", r.RemoteAddr)
		return true
	},
}

// ServeHTTP handles a connection from a client
// Upgrades client's connection to WebSocket and listens for messages
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
			log.Printf("%v", err)
			delete(g.Arena.Players, ws)
			break
		}

		switch msg.Type {
		case "spawn":
			var spawn models.SpawnHandlerMessage
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
			MessageChannel <- msg

		case "keyHandler":
			var kh models.KeyHandlerMessage
			err = json.Unmarshal([]byte(msg.Data.(string)), &kh)
			if err != nil {
				log.Printf("error: %v", err)
				continue
			}
			if _, ok := g.Arena.Players[ws]; ok {
				if kh.IsPressed == true {
					g.Arena.Players[ws].KeyDownHandler(kh.Key)
				} else {
					g.Arena.Players[ws].KeyUpHandler(kh.Key)
				}
			}

		default:
			log.Println("Unknown message type received")
		}
	}
}

func run(g *Game) {
	for {
		time.Sleep(g.RefreshRate)

		g.Arena.UpdatePositions()
		g.Arena.CollisionDetection()
	}
}

func tick(g *Game) {
	for {
		time.Sleep(g.RefreshRate)

		players := make([]models.Player, 0)
		for _, player := range g.Arena.Players {
			players = append(players, *player)
		}

		msg := models.Message{
			Type: "update",
			Data: struct {
				Holes   []*models.Hole  `json:"holes"`
				Junk    []*models.Junk  `json:"junk"`
				Players []models.Player `json:"players"`
			}{
				g.Arena.Holes,
				g.Arena.Junk,
				players,
			},
		}

		// update every client
		for client := range g.Arena.Players {
			p := g.Arena.Players[client]
			p.Mutex.Lock()
			err := client.WriteJSON(&msg)
			p.Mutex.Unlock()
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(g.Arena.Players, client)
			}
		}
	}
}

func messageEmitter(g *Game) {
	for {
		msg := <-MessageChannel

		switch msg.Type {
		case "initial":
			ws := msg.Data.(*websocket.Conn)
			initalMsg := models.Message{
				Type: "initial",
				Data: models.ConnectionMessage{
					g.Arena.Width,
					g.Arena.Height,
					g.Arena.Players[ws].Color,
				},
			}

			g.Arena.Players[ws].Mutex.Lock()
			error := ws.WriteJSON(&initalMsg)
			g.Arena.Players[ws].Mutex.Unlock()
			if error != nil {
				log.Printf("error: %v", error)
				ws.Close()
				delete(g.Arena.Players, ws)
			}

		case "death":
			ws := msg.Data.(*websocket.Conn)
			deathMsg := models.Message{
				Type: "death",
				Data: nil,
			}
			g.Arena.Players[ws].Mutex.Lock()
			error := ws.WriteJSON(&deathMsg)
			g.Arena.Players[ws].Mutex.Unlock()
			if error != nil {
				log.Printf("error: %v", error)
				ws.Close()
				delete(g.Arena.Players, ws)
			}
			delete(g.Arena.Players, ws)

		default:
			log.Println("Unknown message type to emit")
		}
	}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	MessageChannel = make(chan models.Message)

	game.MessageChannel = MessageChannel
	game := Game{
		Arena:       game.CreateArena(2400, 2800),
		RefreshRate: time.Millisecond * 17, // 60 Hz
	}

	http.Handle("/", http.FileServer(http.Dir("./build")))
	http.Handle("/connect", &game)
	go messageEmitter(&game)
	go run(&game)
	go tick(&game)

	log.Println("Starting server on localhost:9090")
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("Error starting server")
	}
}
