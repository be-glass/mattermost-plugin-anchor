package main

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
)

func (p *AnchorPlugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {

	var response string

	command := strings.TrimSpace(args.Command)

	switch command {

	case "/hello":
		p.API.LogInfo("YOU SAID HELLO, THAT IS VERY KIND!", "NICE", "GUY")
		response = fmt.Sprintf("Hello, %s! :) ", args.UserId)

	case "/users":
		users, err := ListAllUsers(p)

		//response = strconv.Itoa(len(users)) + " users"
		//
		//response = users[0].Username

		if err != nil {
			response = "Don't know!"
		} else if len(users) == 0 {
			response = "No users"
		} else {
			var userNames []string
			for _, user := range users {
				userNames = append(userNames, user.Username)
			}
			response = strings.Join(userNames, "\n")
		}

	default:
		response = "Unknown command. Please try something else."
	}

	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeInChannel,
		Text:         response,
	}, nil

}

func ListAllUsers(p *AnchorPlugin) ([]*model.User, error) {
	var allUsers []*model.User
	page := 0
	perPage := 50 // number of users per page

	for {
		users, appErr := p.API.GetUsers(&model.UserGetOptions{
			Page:    page,
			PerPage: perPage,
		})
		if appErr != nil {
			return nil, fmt.Errorf("failed to get users: %w", appErr)
		}

		if len(users) == 0 {
			break // no more users to fetch
		}

		allUsers = append(allUsers, users...)
		page++
	}

	return allUsers, nil
}
