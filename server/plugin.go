package main

import (
	"fmt"
	"github.com/mattermost/mattermost/server/public/model"
)

func (p *AnchorPlugin) OnActivate() error {
	commands := []*model.Command{
		{
			Trigger:          "hello",
			AutoComplete:     false,
			AutoCompleteDesc: "Respond with a greeting",
		},
		{
			Trigger:          "users",
			AutoComplete:     false,
			AutoCompleteDesc: "List users",
		},
		{
			Trigger:          "channels",
			AutoComplete:     false,
			AutoCompleteDesc: "List channels",
		},
		{
			Trigger:          "teams",
			AutoComplete:     false,
			AutoCompleteDesc: "List teams",
		},
		{
			Trigger:          "cleanup",
			AutoComplete:     false,
			AutoCompleteDesc: "Find and remove unwanted posts",
		},
		{
			Trigger:          "check",
			AutoComplete:     false,
			AutoCompleteDesc: "check channel structure of users",
		},
	}

	for _, command := range commands {
		if err := p.API.RegisterCommand(command); err != nil {
			return fmt.Errorf("failed to register command %s: %w", command.Trigger, err)
		}
	}

	return nil
}
