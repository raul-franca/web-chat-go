package main

import (
	"github.com/gorilla/websocket"
)

// client representa um único usuário do chat.
type client struct {

	// socket é o socket da web para este cliente.
	//é o objeto de conexão WebSocket para esse cliente.
	socket *websocket.Conn

	// send é um canal no qual as mensagens são enviadas.
	//é um canal em que as mensagens são enviadas do servidor para o cliente.
	send chan []byte

	// room é a sala em que este client está conversando.
	room *room
}

// read é um método do cliente que lê continuamente as mensagens
// recebidas do cliente através do WebSocket.
func (c *client) read() {
	defer c.socket.Close()
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		c.room.forward <- msg
	}
}

// write  é um método do cliente que envia continuamente as mensagens
// para o cliente através do WebSocket.
func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
