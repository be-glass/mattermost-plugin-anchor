package main

import (
	"fmt"
	"github.com/glass.plugin-anchor/server/models"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/plugin"
)

type AnchorPlugin struct {
	plugin.MattermostPlugin
	Context *models.Context
}

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
