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
		response = p.GetCommandResponse(c, args.Command)
	}

	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         response,
	}, nil
}

func (p *AnchorPlugin) GetCommandResponse(_ *plugin.Context, commandLine string) string {

	arguments := strings.Fields(commandLine)

	var command string
	var user *business.User
	var team *business.Team
	var err error

	var c = p.Context

	if arguments[0] != "/anchor" {
		return fmt.Sprintf("invalid command: %s", arguments[0])
	}
	if !c.User.IsSystemAdmin() {
		return fmt.Sprintf("You do not have permission to execute this command.")
	}
	if len(arguments) < 2 {
		return "missing a command"
	}
	if len(arguments) > 1 {
		command = arguments[1]
		team = business.WrapTeam(*c, c.Team)
	}
	if len(arguments) > 2 {
		user, err = business.NewUser(c, arguments[2])
		if err != nil {
			return err.Error()
		}
	}

	switch command {

	case "hello":
		return fmt.Sprintf("Hello, %s! :) ",
			c.User.GetFullName())

	case "users":
		return business.GetUserListString(c)

	case "cleanup":
		return business.CleanPosts(c, c.Channel.Id, true)

	case "teams":
		return business.GetTeamsListString(c)

	case "channels":
		return team.GetChannelsListString()

	case "check":
		if user != nil {
			return user.CheckChannelStructure()
		} else {
			return team.CheckUserChannelStructure()
		}

	case "onboard":
		if user == nil {
			return "Missing user name"
		}
		return user.CheckAndJoinDefaultChannelStructure()

	case "create_channels":
		return team.CreateDefaultChannels()

	case "delete_sidebar":
		if user == nil {
			return "Missing user name"
		}
		return user.DeleteAllSidebarCategories()

	case "debug":

		if user == nil {
			return "Missing user name"
		}

		if user.C == nil {
			return "BS"
		}

		actualCategories, err := user.SidebarCategoryNames()
		if err != nil {
			return err.Error()
		}

		return strings.Join([]string{
			"**Default Channels:**",
			strings.Join(config.ChannelNames(), "\n"),
			"\n**Subscribed Channels**",
			team.GetChannelsListString(),
			"\n**Default Categories:**",
			strings.Join(config.CategoryNames(), "\n"),
			"\n**Actual Categories:**",
			strings.Join(actualCategories, "\n"),
		}, "\n")

	default:
		return "Unknown command. Please try something else."
	}
}
