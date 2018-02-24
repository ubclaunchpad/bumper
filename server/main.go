package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

// Message is the schema for client/server communication
type Message struct {
	Message string `json:"message"`
}

// Coordinate x y position
type Position struct {
	x int
	y int
}

type Velocity struct {
	dx float32
	dy float32
}

// State of an object, position and velocity
type State struct {
	position Position
	velocity Velocity
}

// global variable is fine for now, all we need is reference to connection
var clients = make(map[*websocket.Conn]bool)

// global channel for message receiving
var broadcast = make(chan Message)

// this 'upgrades' a normal HTTP connection to a persistent TCP connection (socket)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		log.Printf("Accepting client from remote address %v\n", r.RemoteAddr)
		return true
	},
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	// record this connection in our map
	clients[ws] = true

	var msg Message
	// infinite loop that receives msgs from clients
	for {
		// ReadJSON blocks until a message is received
		log.Printf("%+v\n", msg)
		err := ws.ReadJSON(&msg)
		// terminate connection if error occurs
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}

		log.Printf("%+v\n", msg.Message)
		// pass received message to the global channel
		//broadcast <- msg
	}
}

func tick() {
	tickCount := 0
	for {
		time.Sleep(time.Second * 5)
		msg := Message{"tick" + strconv.Itoa(tickCount)}
		// update every client
		for client := range clients {
			err := client.WriteJSON(&msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
		tickCount++
	}
}

func main() {

	// main thread that will listen for connections
	http.HandleFunc("/connect", handleConnection)
	// separate thread that will handle updates
	go tick()

	log.Println("Starting server on localhost:9090")
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("Error starting server")
	}
}
