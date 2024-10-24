package main

import (
	"fmt"
	"github.com/glass.plugin-anchor/server/business"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/plugin"
	"strings"
)

func (p *AnchorPlugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {

	var response string

	err := p.SetContextFromCommandArgs(args)
	if err != nil {
		response = err.DetailedError
	} else {
		response = p.GetCommandResponse(c, args.Command)
	}

	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         response,
	}, nil
}

func (p *AnchorPlugin) GetCommandResponse(_ *plugin.Context, commandLine string) string {

	arguments := strings.Fields(commandLine)

	var command, target string

	var c = p.Context

	if len(arguments) < 2 {
		return "missing a command"
	}
	if len(arguments) > 1 {
		command = arguments[1]
	}
	if len(arguments) > 2 {
		target = arguments[2]
	}
	if arguments[0] != "/anchor" {
		return fmt.Sprintf("invalid command: %s", arguments[0])
	}
	if !c.User.IsSystemAdmin() {
		return fmt.Sprintf("You do not have permission to execute this command.")
	}

	switch command {

	case "hello":
		return fmt.Sprintf("Hello, %s! :) ", c.User.GetFullName())

	case "users":
		return business.GetUserListString(c)

	case "cleanup":
		return business.CleanPosts(c, c.Channel.Id, true)

	case "teams":
		return business.GetTeamsListString(c)

	case "channels":
		return business.GetChannelsListString(c)

	case "check":
		return business.CheckUserOrAll(c, target)

	case "onboard":
		return business.CheckAndJoinDefaultChannelStructure(c, target)

	case "create_channels":
		return business.CreateDefaultChannels(c)

	case "delete_sidebar":
		return business.DeleteAllSidebarCategories(c, target)

	case "debug":

		actualCategories, _ := business.GetUserSidebarCategoryNames(c, c.User.Id)

		return strings.Join([]string{
			"**Default Channels:**",
			strings.Join(business.GetDefaultChannelNames(), "\n"),
			"\n**Subscribed Channels**",
			business.GetChannelsListString(c),
			"\n**Default Categories:**",
			strings.Join(business.GetDefaultCategoryNames(), "\n"),
			"\n**Actual Categories:**",
			strings.Join(actualCategories, "\n"),
		}, "\n")

	default:
		return "Unknown command. Please try something else."
	}
}
