package business

import (
	"fmt"
	"github.com/glass.plugin-anchor/server/models"
	"github.com/mattermost/mattermost-server/v6/model"
	"strings"
)

// public

func GetUserListString(c *models.Context) string {

	users, err := listAllUsers(c)

	if err != nil {
		return "Don't know!"
	} else if len(users) == 0 {
		return "No users"
	} else {
		var userNames []string
		for _, user := range users {
			userNames = append(userNames, user.Username)
		}
		return strings.Join(userNames, "\n")
	}
}

func listAllUsers(c *models.Context) ([]*model.User, error) {
	var allUsers []*model.User
	page := 0
	perPage := 50 // number of users per page

	for {
		users, appErr := c.API.GetUsers(&model.UserGetOptions{
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

func GetUserIDByUsername(c *models.Context, username string) (string, *model.AppError) {
	// Retrieve the user by username
	user, appErr := c.API.GetUserByUsername(username)
	if appErr != nil {
		return "", appErr // Return error if the user is not found or there is an issue
	}

	// Return the user's ID
	return user.Id, nil
}
