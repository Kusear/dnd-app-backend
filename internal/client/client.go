package client

import (
	"dnd-backend-go/internal/common"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type WsClient struct {
	hub *Hub

	conn *websocket.Conn
	send chan []byte

	commandRouter *CommandRouter
}

func InitClient(wsConn *websocket.Conn, hub *Hub, commandRouter *CommandRouter) *WsClient {
	slog.Info("Initializing new WebSocket client")
	client := &WsClient{
		hub:           hub,
		conn:          wsConn,
		send:          make(chan []byte),
		commandRouter: commandRouter,
	}

	hub.Register(client)
	slog.Debug("Client registered to hub", "client", client)

	go client.ListenMessages()
	go client.WriteMessages()

	return client
}

// Converts a map[string]any to a JSON byte array
func (obj *WsClient) convertMessageToJson(parsedMessage map[string]any) ([]byte, error) {
	jsonMessage, err := json.Marshal(parsedMessage)
	if err != nil {
		slog.Error("error: " + err.Error())
		return nil, err
	}
	return jsonMessage, nil
}

// Converts a JSON byte array to a map[string]any
func (obj *WsClient) convertJsonToMessage(jsonMessage []byte) (map[string]any, error) {
	var parsedMessage map[string]any
	err := json.Unmarshal(jsonMessage, &parsedMessage)
	if err != nil {
		slog.Error("error: " + err.Error())
		return nil, err
	}
	return parsedMessage, nil
}

// Listens for messages from the WebSocket connection.
//
// Converts messages to json and sends them to the send channel.
func (obj *WsClient) ListenMessages() {
	defer func() {
		obj.hub.Unregister(obj)
		obj.conn.Close()
	}()
	obj.conn.SetReadLimit(maxMessageSize)
	obj.conn.SetReadDeadline(time.Now().Add(pongWait))
	obj.conn.SetPongHandler(func(string) error { obj.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		mt, message, err := obj.conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error("Error reading message:" + err.Error())
			}
			break
		}

		obj.conn.SetReadDeadline(time.Now().Add(pongWait))
		obj.conn.SetPongHandler(func(string) error { obj.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

		slog.Info("Message received:", "message type", mt, "message", message)

		parsedMessage, err := obj.convertJsonToMessage(message)
		if err != nil {
			slog.Error("error: " + err.Error())
			// TODO add error response to the client (400)
			continue
		}

		parsedMessage["type"] = "response"

		err = obj.commandRouter.ValidatePayload(parsedMessage)
		if err != nil {
			slog.Error("error: " + err.Error())
			// TODO add error response to the client (400)
			continue
		}

		response, err := obj.commandRouter.ExecuteCommand(obj, common.Infrastructure{}, parsedMessage["command"].(string), parsedMessage["data"].(map[string]any))
		if err != nil {
			slog.Error("error: " + err.Error())
			// TODO add error response to the client (can be custom error message)
			continue
		}

		jsonMessage, err := obj.convertMessageToJson(response)
		if err != nil {
			slog.Error("error: " + err.Error())
			// TODO add error response to the client (500)
			continue
		}

		if len(jsonMessage) == 0 {
			continue
		}

		obj.Send(jsonMessage)
	}
}

// Writes messages to the WebSocket connection.
//
// Reads messages from the send channel and writes them to the WebSocket connection.
func (obj *WsClient) WriteMessages() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		obj.conn.Close()
	}()
	for {
		select {
		case message := <-obj.send:
			err := obj.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				slog.Error("Error writing message:", err)
				return
			}
		case <-ticker.C:
			obj.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := obj.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (obj *WsClient) Send(message []byte) {
	obj.send <- message
}
