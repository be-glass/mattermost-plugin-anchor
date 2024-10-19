package main

import (
	"fmt"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
)

type HelloPlugin struct {
	plugin.MattermostPlugin
	configuration *Configuration // Add this field to store the configuration
}

//func (p *HelloPlugin) OnActivate() error {
//	// Register the /hello command
//	err := p.API.RegisterCommand(&model.Command{
//		Trigger:          "hello",
//		AutoComplete:     true,
//		AutoCompleteDesc: "Respond with a greeting",
//		AutoCompleteHint: "",
//	})
//	if err != nil {
//		return err
//	}
//	return nil
//}

func (p *HelloPlugin) OnActivate() error {
	commands := []*model.Command{
		{
			Trigger:          "hello",
			AutoComplete:     true,
			AutoCompleteDesc: "Respond with a greeting",
		},
		{
			Trigger:     "users",
			Description: "List all users",
			DisplayName: "Users",
		},
	}

	for _, command := range commands {
		if err := p.API.RegisterCommand(command); err != nil {
			return fmt.Errorf("failed to register command %s: %w", command.Trigger, err)
		}
	}

	return nil
}
