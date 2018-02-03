package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Message is the schema for client/server communication
type Message struct {
	Message string `json:"message"`
}

// global variable is fine for now, all we need is reference to connection
var clients = make(map[*websocket.Conn]bool)

// global channel for message receiving
var broadcast = make(chan Message)

// this 'upgrades' a normal HTTP connection to a persistent TCP connection (socket)
var upgrader = websocket.Upgrader{}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello from Go!")
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	// record this connection in our map
	clients[ws] = true

	// infinite loop that receives msgs from clients
	for {
		// ReadJSON blocks until a message is received
		var msg Message
		err := ws.ReadJSON(&msg)
		// terminate connection if error occurs
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}

		// pass received message to the global channel
		broadcast <- msg
	}
}

func tick() {

}

func main() {
	http.HandleFunc("/", handler)

	log.Println("Starting server on localhost:9090")
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("Error starting server")
	}
}
