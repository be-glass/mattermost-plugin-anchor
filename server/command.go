package main

import (
	"fmt"
	"github.com/glass.plugin-anchor/server/business"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"strings"
)

func CheckSysAdmin(api plugin.API, userId string) bool {
	user, appErr := api.GetUser(userId)
	if appErr != nil {
		return false
	}
	if !strings.Contains(user.Roles, "system_admin") {
		return false
	}
	return true
}

func (p *AnchorPlugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         p.GetCommandResponse(c, args),
	}, nil
}

func (p *AnchorPlugin) GetCommandResponse(c *plugin.Context, args *model.CommandArgs) string {
	_ = c
	commandArgs := strings.Fields(args.Command)

	var command, target string

	p.API.LogWarn("XXXX 1")

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
	if !CheckSysAdmin(p.API, args.UserId) {
		return fmt.Sprintf("You do not have permission to execute this command.")
	}
	p.API.LogWarn("XXXX 2")

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
		return business.CheckAndJoinDefaultChannels(p.API, target, args.TeamId)

	case "create_channels":
		return business.CreateDefaultChannels(p.API, args.TeamId)

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
