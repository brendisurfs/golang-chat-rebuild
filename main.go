package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) //connected clients map.
var broadcast = make(chan Message)           // channel to host messages.

// upgrader for upgrading http to websocket.
var upgrader = websocket.Upgrader{}

//STRUCTS
//	|
//	v
type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

func main() {
	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)

	// request handling
	http.HandleFunc("/ws", handleConnections)
	//Listen for incoming chat messages.
	//	|
	//	v
	go handleMessages()

	//Server setup
	//	|
	//	v
	log.Println("server started on localhost:8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// must upgrade to websocket first.
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	//Make sure to close the connection after return
	//	|
	//	v
	defer ws.Close()

	//Register client
	//	|
	//	v
	clients[ws] = true

	for {
		var msg Message

		// read input as json and convert it to a Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error while reading json: %v", err)
			delete(clients, ws)
			break
		}
		// send the new message to the broadcast channel
		broadcast <- msg
	}
}
