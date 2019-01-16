package message

import (
	"github.com/sirupsen/logrus"
	"sync"
)

type Hub interface {
	Broadcast(*Message)
	Run()
	Register(Client)
	Unregister(Client)
	Close()
}

type hub struct{
	clients map[Client]bool
	// Inbound messages from the clients.
	broadcast chan *Message

	// Register requests from the clients.
	register chan Client

	// Unregister requests from clients.
	unregister chan Client

	exit chan bool

	wg sync.WaitGroup
}

func NewHub() Hub {
	bchan := make(chan *Message)
	rchan := make(chan Client)
	urchan := make(chan Client)
	echan := make(chan bool)

	cmap := make(map[Client]bool)

	return &hub{
		clients: cmap,
		broadcast: bchan,
		register: rchan,
		unregister: urchan,
		exit: echan,
	}
}

func (h *hub) Register(client Client) {
	logrus.Info("received registration request")
	h.register <- client
}

func (h *hub) Unregister(client Client) {
	h.unregister <- client
}

func (h *hub) Broadcast(message *Message) {
	h.broadcast <- message
}

func (h *hub) Close() {
	h.exit <- true

	close(h.unregister)
	close(h.register)
	close(h.broadcast)
	h.wg.Wait()
}

func (h *hub) Run() {
	h.wg.Add(1)
	go func() {
		for {
			select {
			case client := <-h.register:
				h.clients[client] = true
				logrus.WithField("clients", len(h.clients)).Info("client registered")
			case client := <-h.unregister:
				if _, ok := h.clients[client]; ok {
					delete(h.clients, client)
					logrus.WithField("clients", len(h.clients)).Info("client removed")
				}
			case message := <-h.broadcast:
				logrus.Info("message received")
				for client := range h.clients {
					client.Send(message)
				}
			case <- h.exit:
				return
			}
		}
		h.wg.Done()
	}()
}
