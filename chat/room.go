package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/matryer/goblueprints/chapter1/trace"
)

// room é uma estrutura que representa uma sala de chat.
type room struct {

	// forward é um canal que contém as mensagens recebidas
	//// que deve ser encaminhado para o outro clients.
	forward chan []byte

	// join é um canal para clients desejando ingressar na sala.
	join chan *client

	// leave é um canal para clients desejam sair da sala.
	leave chan *client

	// clients  é um mapa que armazena todos os clientes atuais na sala de chat.
	clients map[*client]bool

	// tracer receberá informações de rastreamento de atividade
	// no room.
	tracer trace.Tracer
}

// newRoom é uma função que cria uma nova instância de uma sala de chat e a retorna.
// Ela inicializa os canais forward, join e leave.
// Inicializa o mapa clients.
// Define o rastreador como trace.Off(), que é um rastreador vazio (sem ação).
func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Off(),
	}
}

// run Ele implementa um loop infinito que aguarda mensagens de três canais: join, leave e forward.
// Quando um cliente é recebido no canal join, ele é adicionado ao mapa clients
// e uma mensagem de rastreamento é registrada.
// Quando um cliente é recebido no canal leave, ele é removido do mapa clients,
// o canal send do cliente é fechado e uma mensagem de rastreamento é registrada.
// Quando uma mensagem é recebida no canal forward, ela é encaminhada para todos os clientes
// /**/no mapa clients através do canal send de cada cliente, e uma mensagem de rastreamento é registrada.
func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			// joining
			r.clients[client] = true
			r.tracer.Trace("New client joined")
		case client := <-r.leave:
			// leaving
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("Client left")
		case msg := <-r.forward:
			r.tracer.Trace("Message received: ", string(msg))
			// forward message to all clients
			for client := range r.clients {
				client.send <- msg
				r.tracer.Trace(" -- sent to client")
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
