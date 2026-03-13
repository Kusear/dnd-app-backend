package commands

import (
	"dnd-backend-go/internal/client"
	"dnd-backend-go/internal/common"
	"dnd-backend-go/internal/utils"
	"errors"
)

const TOKEN_MOVED_COMMAND = "token-moved"

func NewTokenMovedCommand() *client.Command {
	return &client.Command{
		Name:            TOKEN_MOVED_COMMAND,
		Execute:         executeTokenMovedCommand,
		ValidatePayload: validateTokenMovedPayload,
	}
}

func executeTokenMovedCommand(socket *client.WsClient, infrastructure common.Infrastructure, params map[string]any) (interface{}, error) {
	response := map[string]any{
		"command": TOKEN_MOVED_COMMAND,
		"data":    params,
	}

	message, err := utils.ConvertMessageToJson(response)
	if err != nil {
		return nil, err
	}
	socket.BroadcastToHubExceptCurrentClient(message)

	return nil, nil
}

func validateTokenMovedPayload(payload map[string]any) error {
	if payload["tokenId"] == nil {
		return errors.New("token is required")
	}
	return nil
}
