package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

func init() {
	pool = Pool{}
	pool.Init()
}

var pool Pool

type Pool struct {
	guard     sync.Mutex
	hubs map[string]*Hub
}

func (self *Pool) Init() {
	self.guard = sync.Mutex{}
	self.hubs = make(map[string]*Hub)
}

func (self *Pool) Get(name string) *Hub {
	self.guard.Lock()
	defer self.guard.Unlock()
	if _, ok := self.hubs[name]; !ok {
		logger.Warnf("Adding hub: %v", name)
		self.hubs[name] = &Hub{id:name}
		self.hubs[name].Init()
	}
	return self.hubs[name]
}

func (self *Pool) Remove(name string) {
	self.guard.Lock()
	defer self.guard.Unlock()
	if _, ok := self.hubs[name]; ok {
		logger.Warnf("Removing hub: %v", name)
		delete(self.hubs, name)
	}
}

type Hub struct {
	id string
	count int
	guard     sync.Mutex
	clients   map[*websocket.Conn]bool
	broadcast chan map[string]interface{}
}

func (self *Hub) Init() {
	self.guard = sync.Mutex{}
	self.clients = make(map[*websocket.Conn]bool)		 // connected clients
	self.broadcast = make(chan map[string]interface{})   // broadcast channel
	go self.broadcaster()
}

func (self *Hub) add(ws *websocket.Conn) {
	self.guard.Lock()
	logger.Debug("Adding client")
	self.count++
	self.clients[ws] = true
	self.guard.Unlock()
}

func (self *Hub) remove(ws *websocket.Conn) {
	self.guard.Lock()
	logger.Debug("Removing client")
	self.count--
	ws.Close()
	delete(self.clients, ws)

	if 0 == self.count {
		pool.Remove(self.id)
	}

	self.guard.Unlock()
}

func (self *Hub) Add(ws *websocket.Conn) {
	go func() {
		defer self.remove(ws)
		self.add(ws)
		for {
			var msg map[string]interface{}
			// Read in a new message as JSON and map it to a Message object
			err := ws.ReadJSON(&msg)
			if err != nil {
				logger.Warn(err)
				break
			}
			logger.Debug("Received message from client")
			// Send the newly received message to the broadcast channel
			self.broadcast <- msg
		}
	}()
}

func (self *Hub) broadcaster() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-self.broadcast
		logger.Debug("Broadcasting message to clients")
		// Send it out to every client that is currently connected
		for client := range self.clients {
			err := client.WriteJSON(msg)
			if err != nil {
				logger.Error(err)
				self.remove(client)
			}
		}
	}
}
