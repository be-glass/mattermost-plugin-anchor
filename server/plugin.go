package main

import (
	"fmt"
	"github.com/mattermost/mattermost/server/public/model"
)

func (p *AnchorPlugin) OnActivate() error {
	commands := []*model.Command{
		{
			Trigger:          "hello",
			AutoComplete:     true,
			AutoCompleteDesc: "Respond with a greeting",
		},
		{
			Trigger:          "users",
			AutoComplete:     true,
			AutoCompleteDesc: "List users",
		},
		{
			Trigger:          "channels",
			AutoComplete:     true,
			AutoCompleteDesc: "List channels",
		},
		{
			Trigger:          "teams",
			AutoComplete:     true,
			AutoCompleteDesc: "List teams",
		},
		{
			Trigger:          "cleanup",
			AutoComplete:     true,
			AutoCompleteDesc: "Find and remove unwanted posts",
		},
	}

	for _, command := range commands {
		if err := p.API.RegisterCommand(command); err != nil {
			return fmt.Errorf("failed to register command %s: %w", command.Trigger, err)
		}
	}

	return nil
}
