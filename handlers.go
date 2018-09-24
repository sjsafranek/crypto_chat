package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

func chatroomHandler(w http.ResponseWriter, r *http.Request) {
	// Return results
	template_file := "./templates/index.html"
	tmpl, _ := template.ParseFiles(template_file)
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error(err)
		return
	}
	// Make sure we close the connection when the function returns
	// defer ws.Close()

	vars := mux.Vars(r)
	chatroom := vars["chatroom"]

	hub := pool.Get(chatroom)
	hub.Add(ws)

	// // Register our new client
	// guard.Lock()
	// clients[ws] = true
	// guard.Unlock()
	//
	// for {
	// 	var msg Message
	// 	// Read in a new message as JSON and map it to a Message object
	// 	err := ws.ReadJSON(&msg)
	// 	if err != nil {
	// 		logger.Error(err)
	// 		guard.Lock()
	// 		delete(clients, ws)
	// 		guard.Unlock()
	// 		break
	// 	}
	// 	// Send the newly received message to the broadcast channel
	// 	broadcast <- msg
	// }
}

// func handleMessages() {
// 	for {
// 		// Grab the next message from the broadcast channel
// 		msg := <-broadcast
// 		// Send it out to every client that is currently connected
// 		for client := range clients {
// 			err := client.WriteJSON(msg)
// 			if err != nil {
// 				logger.Error(err)
// 				client.Close()
// 				guard.Lock()
// 				delete(clients, client)
// 				guard.Unlock()
// 			}
// 		}
// 	}
// }
