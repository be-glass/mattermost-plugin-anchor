package main

import (
	"fmt"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"strings"
)

type HelloPlugin struct {
	plugin.MattermostPlugin
	configuration *Configuration // Add this field to store the configuration
}

func (p *HelloPlugin) OnActivate() error {
	// Register the /hello command
	err := p.API.RegisterCommand(&model.Command{
		Trigger:          "hello",
		AutoComplete:     true,
		AutoCompleteDesc: "Respond with a greeting",
		AutoCompleteHint: "",
	})
	if err != nil {
		return err
	}
	return nil
}

func (p *HelloPlugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	command := strings.TrimSpace(args.Command)
	if command == "/hello" {
		// Respond with a greeting
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeInChannel,
			Text:         fmt.Sprintf("Hello, %s!", args.UserId),
		}, nil
	}

	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		//Text:         "Unknown command. Please use `/hello`.",
		Text: args.Command,
	}, nil
}
