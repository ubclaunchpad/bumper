package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/ubclaunchpad/bumper/server/arena"
	"github.com/ubclaunchpad/bumper/server/game"
	"github.com/ubclaunchpad/bumper/server/models"
)

func run(g *game.Game) {
	for {
		time.Sleep(g.RefreshRate)

		g.Arena.UpdatePositions()
		g.Arena.CollisionDetection()
	}
}

func tick(g *game.Game) {
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

func messageEmitter(g *game.Game) {
	for {
		msg := <-arena.MessageChannel

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
	arena.MessageChannel = make(chan models.Message)
	game := game.CreateGame()

	http.Handle("/", http.FileServer(http.Dir("./build")))
	http.Handle("/connect", game)
	go messageEmitter(game)
	go run(game)
	go tick(game)

	log.Println("Starting server on localhost:" + os.Getenv("PORT"))
	log.Println(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
