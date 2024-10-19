package main

import (
	"fmt"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"strings"
)

func (p *HelloPlugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {

	var response string

	switch strings.TrimSpace(args.Command) {
	case "/hello":
		response = fmt.Sprintf("Hello, %s! :) ", args.UserId)
	case "/users":
		response = "Weiss nicht"
	default:
		response = "Unknown command. Please try something else."
	}

	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeInChannel,
		Text:         response,
	}, nil

}
