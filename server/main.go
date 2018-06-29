package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/ubclaunchpad/bumper/server/arena"
	"github.com/ubclaunchpad/bumper/server/database"
	"github.com/ubclaunchpad/bumper/server/game"
	"github.com/ubclaunchpad/bumper/server/models"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	arena.MessageChannel = make(chan models.Message)
	game := game.CreateGame()

	database.ConnectDB("service-account.json")
	if database.DBC == nil {
		log.Println("DBClient not initialized correctly")
	}

	http.Handle("/", http.FileServer(http.Dir("./build")))
	http.Handle("/connect", game)
	game.StartGame()
	log.Println("Starting server on localhost:" + os.Getenv("PORT"))
	log.Println(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
