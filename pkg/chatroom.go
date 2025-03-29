package pkg

import (
	"errors"
	"log"
	"net"
	"sync"
)

type ChatRoom struct {
	clients   map[string]bool
	broadcast chan []byte
	leave     chan string
	mu        sync.Mutex
}

func NewChatRoom() *ChatRoom {
	c := &ChatRoom{
		clients:   make(map[string]bool),
		broadcast: make(chan []byte),
	}

	go func() {
		for {
			select {
			case clientId := <-c.leave:
				err := c.RemoveClient(clientId)
				if err != nil {
					log.Println(err)
					return
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case msg := <-c.broadcast:
				for client := range c.clients {
					if !c.clients[client] {
						err := c.Leave(client)
						if err != nil {
							return
						}
						continue
					}
					conn, err := net.Dial("tcp", client)
					if err != nil {
						log.Fatalln(err)
					}

					_, err = conn.Write(msg)

					if err != nil {
						log.Fatalln(err)
					}

					err = conn.Close()
					if err != nil {
						log.Fatalln(err)
					}
				}
			}
		}
	}()

	return c
}

func (cr *ChatRoom) Broadcast(message string) {
	cr.broadcast <- []byte(message)
}

func (cr *ChatRoom) GetClients() map[string]bool {
	return cr.clients
}

func (cr *ChatRoom) Join(clientId string) error {
	if _, ok := cr.clients[clientId]; ok {
		return errors.New("client already exists")
	}
	cr.clients[clientId] = true
	return nil
}

func (cr *ChatRoom) Leave(clientId string) error {
	cr.leave <- clientId
	return nil
}

func (cr *ChatRoom) RemoveClient(clientId string) error {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	if _, ok := cr.clients[clientId]; !ok {
		return errors.New("client does not exist")
	}

	delete(cr.clients, clientId)
	return nil
}
