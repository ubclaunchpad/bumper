package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/ubclaunchpad/bumper/server/game"
	"github.com/ubclaunchpad/bumper/server/models"
)

// Game represents a session
type Game struct {
	Arena   *game.Arena
	Clients map[*websocket.Conn]*models.Player
}

// Message is the schema for client/server communication
type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

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

	g.Clients[ws] = g.Arena.AddPlayer()
	initialMsg := Message{
		Type: "initial",
		Data: g.Clients[ws],
	}

	// send initial message back to client with id
	err = ws.WriteJSON(&initialMsg)
	if err != nil {
		log.Printf("error: %v", err)
		ws.Close()
		delete(g.Clients, ws)
		return
	}

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(g.Clients, ws)
			break
		}

		// TODO: receive controls here
		// player can be accessed through clients[ws]
	}
}

func run(g *Game) {
	for {
		g.Arena.UpdatePositions()
		g.Arena.CollisionDetection()
		time.Sleep(time.Second * 30)
	}
}

func tick(g *Game) {
	for {
		time.Sleep(time.Millisecond * 17) // 60 Hz
		msg := Message{
			Type: "update",
			Data: g.Arena,
		}
		// update every client
		for client := range g.Clients {
			err := client.WriteJSON(&msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(g.Clients, client)
			}
		}
	}
}

func main() {
	game := Game{
		Arena:   game.CreateArena(400, 400),
		Clients: make(map[*websocket.Conn]*models.Player),
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "HELLO, is Inertia working yet?\n")
	})
	http.Handle("/connect", &game)
	go run(&game)
	go tick(&game)

	log.Println("Starting server on localhost:9090")
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("Error starting server")
	}
}
