package client

import (
	"dnd-backend-go/internal/common"
	"errors"
	"log/slog"
)

type Command struct {
	Name            string
	Execute         func(socket *WsClient, infrastructure common.Infrastructure, params map[string]any) (interface{}, error)
	ValidatePayload func(payload map[string]any) error
}

type CommandRouter struct {
	commands map[string]*Command
}

func NewCommandRouter() *CommandRouter {
	return &CommandRouter{
		commands: make(map[string]*Command),
	}
}

func (obj *CommandRouter) ValidateBasicPayload(payload map[string]any) error {
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
	slog.Info("Command registered:", "command", command.Name)
	return nil
}

func (obj *CommandRouter) GetCommand(name string) (*Command, error) {
	command := obj.commands[name]
	if command == nil {
		return nil, errors.New("command not found")
	}
	return command, nil
}

// func (obj *CommandRouter) ValidateCommandPayload(command *Command, params map[string]any) error {
// 	return command.ValidatePayload(params)
// }

func (obj *CommandRouter) ExecuteCommand(socket *WsClient, infrastructure common.Infrastructure, name string, params map[string]any) (interface{}, error) {
	command, err := obj.GetCommand(name)
	if err != nil {
		return nil, err
	}
	response, err := command.Execute(socket, infrastructure, params)
	if err != nil {
		return nil, err
	}
	return response, nil
}
