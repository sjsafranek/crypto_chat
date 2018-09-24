package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	if !HasSession(r) {
		logger.Warn("Not authenticated")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Return results
	template_file := "./templates/index.html"
	tmpl, _ := template.ParseFiles(template_file)
	err := tmpl.Execute(w, GetUsernameFromSession(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error(err)
		return
	}

	vars := mux.Vars(r)
	chatroom := vars["chatroom"]

	hub := pool.Get(chatroom)
	hub.Add(ws)
}
