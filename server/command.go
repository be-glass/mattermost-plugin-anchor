package main

import (
	"fmt"
	"github.com/glass.plugin-anchor/server/business"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"strings"
)

func (p *AnchorPlugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {

	p.API.LogWarn("Entering the command section")

	var response string

	// check user - only accept system admin
	user, appErr := p.API.GetUser(args.UserId)
	if appErr != nil {
		return nil, appErr
	}
	if !strings.Contains(user.Roles, "system_admin") {
		return &model.CommandResponse{
			Text: "You do not have permission to execute this command.",
		}, nil
	}

	//command := strings.TrimSpace(args.Command)
	commandArgs := strings.Fields(args.Command)
	command := commandArgs[0]

	switch command {

	case "/hello":
		response = fmt.Sprintf("Hello, %s! :) ", args.UserId)

	case "/users":
		response = business.GetUserListString(p.API)

	case "/cleanup":
		posts := business.CleanPosts(p.API, args.ChannelId)

		response = posts

	case "/teams":
		p.API.LogWarn("found teams command")

		response = business.GetTeamsListString(p.API)

	case "/channels":
		response = business.GetChannelsListString(p.API, args.TeamId)

	case "/check":

		var targetUser string

		if len(commandArgs) > 0 {
			targetUser = commandArgs[1]
		} else {
			targetUser = "all"
		}

		response = business.CheckUserOrAll(p.API, targetUser, args.TeamId)

	case "/onboard":

		var targetUser string

		if len(commandArgs) > 0 {
			targetUser = commandArgs[1]
		} else {
			targetUser = "all"
		}

		response = business.CheckAndJoinDefaultChannels(p.API, targetUser, args.TeamId)

	case "/debug":

		actualCategories, _ := business.GetUserSidebarCategoryNames(p.API, args.UserId, args.TeamId)

		response = strings.Join([]string{
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
		response = "Unknown command. Please try something else."
	}

	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         response,
	}, nil

}
