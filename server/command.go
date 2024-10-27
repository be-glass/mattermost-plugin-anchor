package main

import (
	"errors"
	"fmt"
	"github.com/glass.plugin-anchor/server/business"
	"github.com/glass.plugin-anchor/server/config"
	"github.com/glass.plugin-anchor/server/models"
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

func checkCommand(c *models.Context, line string) error {
	arguments := strings.Fields(line)

	if arguments[0] != "/anchor" && arguments[0] != "/q" {
		return errors.New("invalid command: " + arguments[0])
	}
	if !c.User.IsSystemAdmin() {
		return errors.New("You do not have permission to execute this command.")
	}
	if arguments[0] == "/q" {
		return nil
	}
	if len(arguments) < 2 {
		return errors.New("missing a command")
	}

	return nil
}

func parseCommand(c *models.Context, line string) (string, *business.Team, *business.User, *business.SideBar, error) {

	var err error
	var command string
	var user *business.User
	var sidebar *business.SideBar
	var team = business.WrapTeam(c, c.Team)

	err = checkCommand(c, line)
	if err != nil {
		return "", nil, nil, nil, err
	}

	arguments := strings.Fields(line)

	switch arguments[0] {

	case "/q":

		user, err = business.NewUser(c, "boris")
		if err != nil {
			return "", nil, nil, nil, err
		}
		command = "reorder"
		sidebar, err = business.NewSideBar(user)
		if err != nil {
			return "", nil, nil, nil, err
		}

	case "/anchor":

		command = arguments[1]

		if len(arguments) > 2 {
			user, err = business.NewUser(c, arguments[2])
			if err != nil {
				return "", nil, nil, nil, err
			}
			sidebar, err = business.NewSideBar(user)
			if err != nil {
				return "", nil, nil, nil, err
			}
		}
	default:
		return "", nil, nil, nil, errors.New("no valid command recognized")
	}

	return command, team, user, sidebar, err

}

func (p *AnchorPlugin) GetCommandResponse(_ *plugin.Context, commandLine string) string {

	var c = p.Context

	command, team, user, sideBar, err := parseCommand(p.Context, commandLine)
	if err != nil {
		return err.Error()
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
			return sideBar.CheckChannelStructure()
		} else {
			return team.CheckUserChannelStructure()
		}

	case "onboard":
		if user == nil {
			return "Missing user name"
		}
		return sideBar.CheckAndJoinDefaultChannelStructure()

	case "create_channels":
		return team.CreateDefaultChannels()

	case "delete_sidebar":
		if user == nil {
			return "Missing user name"
		}
		return sideBar.DeleteAllSidebarCategories()

	case "reorder":

		if sideBar == nil {
			return "Missing user name"
		}

		return sideBar.ReorderSidebarCategories()

	case "debug":

		if user == nil {
			return "Missing user name"
		}

		actualCategories, err := sideBar.SidebarCategoryNames()
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
