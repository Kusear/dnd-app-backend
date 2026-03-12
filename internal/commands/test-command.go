package commands

import (
	"dnd-backend-go/internal/client"
	"dnd-backend-go/internal/common"
	"log/slog"
)

func NewTestCommand() *client.Command {
	return &client.Command{
		Name:    "test",
		Execute: executeTestCommand,
	}
}

func executeTestCommand(ctx *client.WsClient, infrastructure common.Infrastructure, params map[string]any) (map[string]any, error) {

	slog.Info("Test command executed")
	return nil, nil
}
