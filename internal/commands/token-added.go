package commands

import (
	"dnd-backend-go/internal/client"
	"dnd-backend-go/internal/common"
	"dnd-backend-go/internal/utils"
	"errors"
)

const TOKEN_ADDED_COMMAND = "token-added"

func NewTokenAddedCommand() *client.Command {
	return &client.Command{
		Name:            TOKEN_ADDED_COMMAND,
		Execute:         executeTokenAddedCommand,
		ValidatePayload: validateTokenAddedPayload,
	}
}

func executeTokenAddedCommand(socket *client.WsClient, infrastructure common.Infrastructure, params map[string]any) (interface{}, error) {
	response := map[string]any{
		"command": TOKEN_ADDED_COMMAND,
		"data":    params,
	}

	message, err := utils.ConvertMessageToJson(response)
	if err != nil {
		return nil, err
	}
	socket.BroadcastToHubExceptCurrentClient(message)

	return nil, nil
}

func validateTokenAddedPayload(payload map[string]any) error {
	if payload["tokenId"] == nil {
		return errors.New("token is required")
	}
	if payload["token"] == nil {
		return errors.New("token is required")
	}
	if payload["position"] == nil {
		return errors.New("position is required")
	}
	return nil
}

// {
// 	"command": "token-added",
// 	"data": {
// 	  "tokenId": "123",
// 	  "token": {
// 		"type": "hero",
// 		"label": "Knight",
// 		"width": 48,
// 		"height": 48,
// 		"color": "#1d4ed8"
// 	  },
// 	  "position": { "x": 250, "y": 180 }
// 	}
//   }
