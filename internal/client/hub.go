package client

type BroadcastMessage struct {
	Message []byte
	Client  *WsClient
}

type BroadcastJsonMessage struct {
	Message map[string]any
	Client  *WsClient
}

type Hub struct {
	clients       map[*WsClient]bool
	register      chan *WsClient
	unregister    chan *WsClient
	broadcast     chan BroadcastMessage
	broadcastJson chan BroadcastJsonMessage
}

func NewHub() *Hub {
	return &Hub{
		clients:       make(map[*WsClient]bool),
		register:      make(chan *WsClient),
		unregister:    make(chan *WsClient),
		broadcast:     make(chan BroadcastMessage),
		broadcastJson: make(chan BroadcastJsonMessage),
	}
}

func (obj *Hub) Run() {
	for {
		select {
		case client := <-obj.register:
			obj.clients[client] = true
		case client := <-obj.unregister:
			delete(obj.clients, client)
		case broadcastMessage := <-obj.broadcast:
			for client := range obj.clients {
				if broadcastMessage.Client != nil && broadcastMessage.Client == client {
					continue
				}
				client.Send(broadcastMessage.Message)
			}
		case broadcastJsonMessage := <-obj.broadcastJson:
			for client := range obj.clients {
				if broadcastJsonMessage.Client != nil && broadcastJsonMessage.Client == client {
					continue
				}
				client.SendJson(broadcastJsonMessage.Message)
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
	obj.broadcast <- BroadcastMessage{
		Message: message,
		Client:  nil,
	}
}

func (obj *Hub) BroadcastJson(message map[string]any) {
	obj.broadcastJson <- BroadcastJsonMessage{
		Message: message,
		Client:  nil,
	}
}

func (obj *Hub) BroadcastExceptCurrentClient(message []byte, client *WsClient) {
	obj.broadcast <- BroadcastMessage{
		Message: message,
		Client:  client,
	}
}

func (obj *Hub) BroadcastJsonExceptCurrentClient(message map[string]any, client *WsClient) {
	obj.broadcastJson <- BroadcastJsonMessage{
		Message: message,
		Client:  client,
	}
}
