package client

import (
	"dnd-backend-go/internal/common"
	"errors"
	"log/slog"
)

type Command struct {
	Name    string
	Execute func(ctx *WsClient, infrastructure common.Infrastructure, params map[string]any) (map[string]any, error)
}

type CommandRouter struct {
	commands map[string]*Command
}

func NewCommandRouter() *CommandRouter {
	return &CommandRouter{
		commands: make(map[string]*Command),
	}
}

func (obj *CommandRouter) ValidatePayload(payload map[string]any) error {
	if payload["command"] == nil {
		return errors.New("command is required")
	}
	if payload["data"] == nil {
		return errors.New("data is required")
	}
	// if payload["type"] == nil {
	// 	return errors.New("type is required")
	// }
	return nil
}

func (obj *CommandRouter) RegisterCommand(command *Command) error {

	existingCommand := obj.commands[command.Name]
	if existingCommand != nil {
		slog.Error("Command already registered:", "command", command.Name)
		return errors.New("command already registered")
	}

	obj.commands[command.Name] = command
	return nil
}

func (obj *CommandRouter) ExecuteCommand(ctx *WsClient, infrastructure common.Infrastructure, name string, params map[string]any) (map[string]any, error) {
	command := obj.commands[name]
	if command == nil {
		return nil, errors.New("command not found")
	}
	return command.Execute(ctx, infrastructure, params)
}
