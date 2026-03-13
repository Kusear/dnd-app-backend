package commands

import (
	"dnd-backend-go/internal/client"
	"dnd-backend-go/internal/common"
	"dnd-backend-go/internal/utils"
	"errors"
)

const MAP_CHANGED_COMMAND = "map-changed"

func NewMapChangedCommand() *client.Command {
	return &client.Command{
		Name:            MAP_CHANGED_COMMAND,
		Execute:         executeMapChangedCommand,
		ValidatePayload: validateMapChangedPayload,
	}
}

func executeMapChangedCommand(socket *client.WsClient, infrastructure common.Infrastructure, params map[string]any) (interface{}, error) {
	response := map[string]any{
		"command": MAP_CHANGED_COMMAND,
		"data":    params,
	}

	message, err := utils.ConvertMessageToJson(response)
	if err != nil {
		return nil, err
	}

	socket.BroadcastToHubExceptCurrentClient(message)

	return nil, nil
}

func validateMapChangedPayload(payload map[string]any) error {
	// if payload["mapId"] == nil {
	// 	return errors.New("map is required")
	// }
	if payload["prevMap"] == nil {
		return errors.New("map is required")
	}
	if payload["newMap"] == nil {
		return errors.New("map is required")
	}
	return nil
}
