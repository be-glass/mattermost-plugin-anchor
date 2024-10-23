package main

import (
	"fmt"
	"github.com/glass.plugin-anchor/server/business"
	"github.com/glass.plugin-anchor/server/config"
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
		response = p.GetCommandResponse(c, args)
	}

	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         response,
	}, nil
}

func (p *AnchorPlugin) GetCommandResponse(c *plugin.Context, args *model.CommandArgs) string {
	_ = c
	commandArgs := strings.Fields(args.Command)

	var command, target string

	if len(commandArgs) < 2 {
		return "missing a command"
	}
	if len(commandArgs) > 1 {
		command = commandArgs[1]
	}
	if len(commandArgs) > 2 {
		target = commandArgs[2]
	}
	if commandArgs[0] != "/anchor" {
		return fmt.Sprintf("invalid command: %s", commandArgs[0])
	}
	if !business.CheckSysAdmin(p.API, args.UserId) {
		return fmt.Sprintf("You do not have permission to execute this command.")
	}

	switch command {

	case "hello":
		return fmt.Sprintf("Hello, %s! :) ", args.UserId)

	case "users":
		return business.GetUserListString(p.API)

	case "cleanup":
		return business.CleanPosts(p.API, args.ChannelId, true)

	case "teams":
		return business.GetTeamsListString(p.API)

	case "channels":
		return business.GetChannelsListString(p.API, args.TeamId)

	case "check":
		return business.CheckUserOrAll(p.API, target, args.TeamId)

	case "onboard":
		return business.CheckAndJoinDefaultChannelStructure(p.API, target, args.TeamId)

	case "create_channels":
		return business.CreateDefaultChannels(p.API, args.TeamId)

	case "delete_sidebar":
		return business.DeleteAllSidebarCategories(p.API, target, args.TeamId, config.BotToken)

	case "debug":

		actualCategories, _ := business.GetUserSidebarCategoryNames(p.API, args.UserId, args.TeamId)

		return strings.Join([]string{
			"**Default Channels:**",
			strings.Join(business.GetDefaultChannelNames(), "\n"),
			"\n**Subscribed Channels**",
			business.GetChannelsListString(p.API, args.TeamId),
			"\n**Default Categories:**",
			strings.Join(business.GetDefaultCategoryNames(), "\n"),
			"\n**Actual Categories:**",
			strings.Join(actualCategories, "\n"),
		}, "\n")

	default:
		return "Unknown command. Please try something else."
	}
}
