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

	command := strings.TrimSpace(args.Command)

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

	default:
		response = "Unknown command. Please try something else."
	}

	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         response,
	}, nil

}
