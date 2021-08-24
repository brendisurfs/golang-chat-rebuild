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
