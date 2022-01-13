// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"encoding/json"
	"log"
	"sync"
	"time"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
	lastId     int
	idLock     sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			log.Println("client exit or unregister")
			if _, ok := h.clients[client]; ok {

				//{"type":"logout","client_id":xxx,"time":"xxx"}
				hData := make(map[string]interface{})
				hData["type"] = "logout"
				hData["time"] = time.Now().Format(DateFormat)
				hData["from_client_id"] = client.ClientId
				hData["from_client_name"] = client.ClientName
				log.Println("hData:\n", hData)
				if sendData, err2 := json.Marshal(hData); err2 != nil {
					log.Printf("\n json.Marshal err:%v \n", err2)
					return
				} else {
					for client2 := range h.clients {
						if client == client2 {
							continue
						}
						select {
						case client2.send <- sendData:
						default:
							close(client2.send)
							delete(h.clients, client2)
						}
					}

				}
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			// log.Printf("\n Hub message channel :%s \n", message)
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
