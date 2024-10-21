package main

import (
	"fmt"
	"github.com/mattermost/mattermost/server/public/model"
)

func (p *AnchorPlugin) OnActivate() error {
	commands := []*model.Command{
		{
			Trigger:          "anchor",
			AutoComplete:     false,
			AutoCompleteDesc: "plugin commands",
		},
	}

	for _, command := range commands {
		if err := p.API.RegisterCommand(command); err != nil {
			return fmt.Errorf("failed to register command %s: %w", command.Trigger, err)
		}
	}

	return nil
}
