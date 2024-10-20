package main

import "github.com/mattermost/mattermost/server/public/plugin"

type AnchorPlugin struct {
	plugin.MattermostPlugin
	configuration *Configuration // Add this field to store the configuration
}
