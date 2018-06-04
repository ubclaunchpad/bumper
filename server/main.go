package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/ubclaunchpad/bumper/server/arena"
	"github.com/ubclaunchpad/bumper/server/models"
)

// Game represents a session
type Game struct {
	Arena       *arena.Arena
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
	var initialMsg models.Message
	id := models.GenUniqueID()

	err = g.Arena.AddPlayer(id, ws)
	if err != nil {
		log.Printf("Error adding player:\n%v", err)
	} else {
		initialMsg = models.Message{
			Type: "connect",
			Data: id,
		}
		MessageChannel <- initialMsg
	}
	for {
		var msg models.Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("%v\n", err)
			delete(g.Arena.Players, id)
			break
		}
		switch msg.Type {
		case "spawn":

			var spawn models.SpawnHandlerMessage
			err = json.Unmarshal([]byte(msg.Data.(string)), &spawn)
			if err != nil {
				log.Printf("%v\n", err)
				continue
			}

			err := g.Arena.SpawnPlayer(id, spawn.Name)
			if err != nil {
				log.Printf("Error spawning player:\n%v", err)
				continue
			}
		case "reconnect":
			err = g.Arena.AddPlayer(id, ws)
			if err != nil {
				log.Printf("Error adding player:\n%v", err)
			} else {
				connectMsg := models.Message{
					Type: "connect",
					Data: id,
				}
				MessageChannel <- connectMsg
			}
		case "keyHandler":
			var kh models.KeyHandlerMessage
			err = json.Unmarshal([]byte(msg.Data.(string)), &kh)
			if err != nil {
				log.Printf("%v\n", err)
				continue
			}
			if _, ok := g.Arena.Players[id]; ok {
				if kh.IsPressed {
					g.Arena.Players[id].KeyDownHandler(kh.Key)
				} else {
					g.Arena.Players[id].KeyUpHandler(kh.Key)
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

		msg := models.Message{
			Type: "update",
			Data: g.Arena.GetState(),
		}

		// update every client
		for id := range g.Arena.Players {
			p := g.Arena.Players[id]
			err := p.SendJSON(&msg)
			if err != nil {
				log.Printf("error: %v", err)
				p.Close()
				delete(g.Arena.Players, id)
			}
		}
	}
}

func messageEmitter(g *Game) {
	for {
		msg := <-MessageChannel

		switch msg.Type {
		case "connect":
			id := msg.Data.(string)
			p := g.Arena.Players[id]

			initalMsg := models.Message{
				Type: "initial",
				Data: models.ConnectionMessage{
					ArenaWidth:  g.Arena.Width,
					ArenaHeight: g.Arena.Height,
					PlayerID:    id,
				},
			}

			err := p.SendJSON(&initalMsg)
			if err != nil {
				log.Printf("error: %v", err)
				p.Close()
				delete(g.Arena.Players, id)
			}

		case "death":
			id := msg.Data.(string)
			deathMsg := models.Message{
				Type: "death",
				Data: nil,
			}

			p := g.Arena.Players[id]
			err := p.SendJSON(&deathMsg)
			if err != nil {
				log.Printf("error: %v", err)
				p.Close()
			}
			delete(g.Arena.Players, id)

		default:
			log.Println("Unknown message type to emit")
		}
	}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	MessageChannel = make(chan models.Message)

	arena.MessageChannel = MessageChannel
	game := Game{
		Arena:       arena.CreateArena(2400, 2800, 20, 30),
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
