package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

var pool Pool

type Pool struct {
	hubs map[string]*Hub
}

func (self *Pool) Init() {
	self.hubs = make(map[string]*Hub)
}

func (self *Pool) Get(name string) *Hub {
	if _, ok := self.hubs[name]; !ok {
		self.hubs[name] = &Hub{}
		self.hubs[name].Init()
	}
	return self.hubs[name]
}

func init() {
	pool = Pool{}
	pool.Init()
}

// Define our message object
type Message struct {
	Data string `json:"data"`
}

type Hub struct {
	guard     sync.Mutex
	clients   map[*websocket.Conn]bool
	broadcast chan Message
}

func (self *Hub) Init() {
	self.guard = sync.Mutex{}
	self.clients = make(map[*websocket.Conn]bool) // connected clients
	self.broadcast = make(chan Message)           // broadcast channel
	go self.broadcaster()
}

func (self *Hub) add(ws *websocket.Conn) {
	guard.Lock()
	logger.Debug("Adding client")
	self.clients[ws] = true
	guard.Unlock()
}

func (self *Hub) remove(ws *websocket.Conn) {
	guard.Lock()
	logger.Debug("Removing client")
	ws.Close()
	delete(self.clients, ws)
	guard.Unlock()
}

func (self *Hub) Add(ws *websocket.Conn) {
	go func() {
		// defer ws.Close()
		defer self.remove(ws)
		self.add(ws)
		for {
			var msg Message
			// Read in a new message as JSON and map it to a Message object
			err := ws.ReadJSON(&msg)
			if err != nil {
				logger.Error(err)
				// self.remove(ws)
				break
			}
			// Send the newly received message to the broadcast channel
			self.broadcast <- msg
		}
	}()
}

func (self *Hub) broadcaster() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-self.broadcast
		// Send it out to every client that is currently connected
		for client := range self.clients {
			err := client.WriteJSON(msg)
			if err != nil {
				logger.Error(err)
				// client.Close()
				self.remove(client)
			}
		}
	}
}
