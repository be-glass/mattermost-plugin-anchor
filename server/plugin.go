package main

import (
	"encoding/json"
	"fmt"
	"github.com/glass.plugin-anchor/server/models"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/plugin"
	"os"
	"path/filepath"
)

type AnchorPlugin struct {
	plugin.MattermostPlugin
	Context *models.Context
}

type PluginManifest struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

func (p *AnchorPlugin) OnActivate() error {
	commands := []*model.Command{
		{
			Trigger:          "anchor",
			AutoComplete:     false,
			AutoCompleteDesc: "plugin commands",
		},
		{
			Trigger:          "q",
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

func (p *AnchorPlugin) GetVersion() (string, error) {

	bundlePath, err := p.API.GetBundlePath()
	if err != nil {
		return "", fmt.Errorf("failed to get bundle path: %v", err)
	}

	filePath := filepath.Join(bundlePath, "plugin.json")

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var manifest PluginManifest
	err = json.NewDecoder(file).Decode(&manifest)
	if err != nil {
		return "", err
	}

	return manifest.Version, nil
}
