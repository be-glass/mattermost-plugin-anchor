package main

import (
	"github.com/glass.plugin-anchor/server/api"
	"github.com/glass.plugin-anchor/server/config"
	"github.com/glass.plugin-anchor/server/models"
	"github.com/mattermost/mattermost-server/v6/model"
)

func (p *AnchorPlugin) SetContextFromCommandArgs(args *model.CommandArgs) *model.AppError {

	team, appErr := p.API.GetTeam(args.TeamId)
	if appErr != nil {
		p.API.LogError("Failed to get team", "teamId", "error", appErr.Error())
		return appErr
	}

	// Retrieve the User
	user, appErr := p.API.GetUser(args.UserId)
	if appErr != nil {
		p.API.LogError("Failed to get user", "userId", args.UserId, "error", appErr.Error())
		return appErr
	}

	// Retrieve the Channel
	channel, appErr := p.API.GetChannel(args.ChannelId)
	if appErr != nil {
		p.API.LogError("Failed to get channel", "channelId", args.ChannelId, "error", appErr.Error())
		return appErr
	}

	// Initialize p.Context if it's nil
	if p.Context == nil {
		p.Context = &models.Context{}
	}

	// Set the retrieved objects in the Context
	p.Context.Team = team
	p.Context.Channel = channel
	p.Context.User = user

	// Optionally set other fields
	p.Context.API = p.API
	p.Context.Auth = config.AuthConfig
	p.Context.Rest = api.NewRestClient(config.ServerURL, p.Context.Auth.AuthToken, config.Headers)

	p.API.LogWarn("AUTH CONFIG", "ServerURL", config.ServerURL, "AuthToken", p.Context.Auth.AuthToken, "Headers", config.Headers)

	return nil
}
