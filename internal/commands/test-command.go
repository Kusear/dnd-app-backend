package commands

import (
	"dnd-backend-go/internal/client"
	"dnd-backend-go/internal/common"
	"dnd-backend-go/internal/utils"
	"log/slog"
)

func NewTestCommand() *client.Command {
	return &client.Command{
		Name:    "test",
		Execute: executeTestCommand,
	}
}

func executeTestCommand(socket *client.WsClient, infrastructure common.Infrastructure, params map[string]any) (interface{}, error) {

	slog.Info("Test command executed")

	mockMessage, err := utils.ConvertMessageToJson(map[string]any{
		"message": "Test command executed",
	})
	if err != nil {
		return nil, err
	}

	socket.BroadcastToHub(mockMessage)

	// Return nil, nil to complete command without sending any response to the CURRENT client
	// If used like this, need to send response to current client manually
	return nil, nil
}
