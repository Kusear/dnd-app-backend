package client

type Hub struct {
	clients    map[*WsClient]bool
	broadcast  chan []byte
	register   chan *WsClient
	unregister chan *WsClient
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*WsClient]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *WsClient),
		unregister: make(chan *WsClient),
	}
}

func (obj *Hub) Run() {
	for {
		select {
		case client := <-obj.register:
			obj.clients[client] = true
		case client := <-obj.unregister:
			delete(obj.clients, client)
		case message := <-obj.broadcast:
			for client := range obj.clients {
				client.Send(message)
			}
		}
	}
}

func (obj *Hub) Register(client *WsClient) {
	obj.register <- client
}

func (obj *Hub) Unregister(client *WsClient) {
	obj.unregister <- client
}

func (obj *Hub) Broadcast(message []byte) {
	obj.broadcast <- message
}
