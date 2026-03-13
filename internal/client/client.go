package client

import (
	"dnd-backend-go/internal/common"
	"dnd-backend-go/internal/utils"
	"fmt"
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

	conn     *websocket.Conn
	send     chan []byte
	sendJson chan map[string]any

	commandRouter *CommandRouter
}

func InitClient(wsConn *websocket.Conn, hub *Hub, commandRouter *CommandRouter) *WsClient {
	client := &WsClient{
		hub:           hub,
		conn:          wsConn,
		send:          make(chan []byte),
		sendJson:      make(chan map[string]any),
		commandRouter: commandRouter,
	}

	hub.Register(client)
	slog.Debug("Client registered to hub", "client", client)

	go client.ListenMessages()
	go client.WriteMessages()

	return client
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
		// change to ReadJSON
		mt, message, err := obj.conn.ReadMessage()

		// TODO sync req-res log messages
		slog.Info("Reading message", "message type", mt, "message", message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error("Error reading message:" + err.Error())
			}
			break
		}

		requestUid := utils.RandomString(10)
		obj.conn.SetReadDeadline(time.Now().Add(pongWait))
		obj.conn.SetPongHandler(func(string) error { obj.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

		slog.Info(fmt.Sprintf("[rId: %s] (message type = %d) Message: %s", requestUid, mt, message))

		parsedMessage, err := utils.ConvertJsonToMessage(message)
		if err != nil {
			slog.Error("error: " + err.Error())
			// TODO add error response to the client (400)
			continue
		}

		parsedMessage["type"] = "response"

		err = obj.commandRouter.ValidateBasicPayload(parsedMessage)
		if err != nil {
			slog.Error("error: " + err.Error())
			// TODO add error response to the client (400)
			continue
		}

		// Execute command in a separate goroutine to avoid blocking other requests
		go func() {
			command, err := obj.commandRouter.GetCommand(parsedMessage["command"].(string))
			if err != nil {
				slog.Error("error: " + err.Error())
				obj.SendErrorResponse(common.NotFoundError, err.Error())
				return
			}

			err = command.ValidatePayload(parsedMessage["data"].(map[string]any))
			if err != nil {
				slog.Error("error: " + err.Error())
				obj.SendErrorResponse(common.BadRequestError, err.Error())
				return
			}

			response, err := command.Execute(obj, common.Infrastructure{}, parsedMessage["data"].(map[string]any))
			// response, err := obj.commandRouter.ExecuteCommand(obj, common.Infrastructure{}, parsedMessage["command"].(string), parsedMessage["data"].(map[string]any))
			if err != nil {
				slog.Error("error: " + err.Error())
				// TODO add error response to the client (can be custom error message)
				obj.SendErrorResponse(common.BadRequestError, err.Error())
				return
			}

			// jsonMessage, err := utils.ConvertMessageToJson(response)
			// if err != nil {
			// 	slog.Error("error: " + err.Error())
			// 	// TODO add error response to the client (500)
			// 	return
			// }

			if response == nil {
				slog.Info("No response to send to client")
				return
			}

			slog.Info(fmt.Sprintf("[rId: %s] Sending response to client", requestUid), response)

			obj.SendJson(response.(map[string]any))
		}()

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
				slog.Error("Error writing message:", "error", err)
				return
			}
		case message := <-obj.sendJson:
			err := obj.conn.WriteJSON(message)
			if err != nil {
				slog.Error("Error writing message:", "error", err)
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

func (obj *WsClient) SendJson(message map[string]any) {
	obj.sendJson <- message
}

func (obj *WsClient) BroadcastToHub(message []byte) {
	obj.hub.Broadcast(message)
}

func (obj *WsClient) BroadcastToHubJson(message map[string]any) {
	obj.hub.BroadcastJson(message)
}

func (obj *WsClient) BroadcastToHubExceptCurrentClient(message []byte) {
	obj.hub.BroadcastExceptCurrentClient(message, obj)
}

func (obj *WsClient) BroadcastToHubExceptCurrentClientJson(message map[string]any) {
	obj.hub.BroadcastJsonExceptCurrentClient(message, obj)
}

func (obj *WsClient) SendErrorResponse(errorType common.ErrorType, message string) {
	// jsonMessage, err := utils.ConvertMessageToJson(errorType.GetErrorMessage(message))
	// if err != nil {
	// 	slog.Error("error: " + err.Error())
	// 	return
	// }
	// obj.Send(jsonMessage)

	obj.sendJson <- errorType.GetErrorMessage(message)
}
