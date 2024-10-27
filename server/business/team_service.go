package business

import (
	"fmt"
	"github.com/glass.plugin-anchor/server/config"
	"github.com/glass.plugin-anchor/server/models"
	"github.com/mattermost/mattermost-server/v6/model"
	"strings"
)

type Team struct {
	c *models.Context
	*model.Team
}

// Constructors

func WrapTeam(c *models.Context, team *model.Team) *Team {
	return &Team{c, team}
}

func (t *Team) GetChannelsListString() string {

	channels, err := t.getChannels()
	if err != nil {
		// Return the error message as part of the string
		return fmt.Sprintf("Error fetching channels: %v", err)
	}

	// Use a string builder for efficient string concatenation
	var builder strings.Builder

	// Loop through all channels and append their names to the builder
	for _, channel := range channels {
		builder.WriteString(fmt.Sprintf("%s\n", channel.DisplayName))
	}

	// Convert the builder to a string and return it
	return builder.String()
}

func (t *Team) getChannels() ([]*model.Channel, error) {
	var allChannels []*model.Channel
	page := 0
	perPage := 100 // You can adjust this to change how many channels are fetched per page

	for {
		// Get channels for the current page in the team
		channels, appErr := t.c.API.GetPublicChannelsForTeam(t.Team.Id, page, perPage)
		if appErr != nil {
			return nil, appErr
		}

		// If no channels are returned, we've retrieved all of them
		if len(channels) == 0 {
			break
		}

		// Append the retrieved channels to the final list
		allChannels = append(allChannels, channels...)

		// Move to the next page
		page++
	}

	return allChannels, nil
}

func (t *Team) CheckUserChannelStructure() string {
	var resultBuilder strings.Builder

	page := 0
	perPage := 100
	for {
		users, appErr := t.c.API.GetUsersInTeam(t.Team.Id, page, perPage)
		if appErr != nil {
			return "Unable to retrieve users in the team."
		}

		if len(users) == 0 {
			break
		}

		for _, user := range users {
			u := WrapUser(t.c, user)
			s, err := NewSideBar(u)
			if err != nil {
				return fmt.Sprintf("Error creating side-bar: %v", err)
			}
			userStructureResult := s.CheckChannelStructure()
			resultBuilder.WriteString(userStructureResult + "\n")
		}

		page++
	}

	return resultBuilder.String()
}

func GetTeamsListString(c *models.Context) string {

	teams, err := c.API.GetTeams()

	if err != nil {
		return "Don't know!"
	} else if len(teams) == 0 {
		return "No teams"
	} else {
		var teamNames []string
		for _, team := range teams {
			teamNames = append(teamNames, team.Name)
		}
		return strings.Join(teamNames, "\n")
	}
}

// private

func (t *Team) CreateDefaultChannels() string {
	var result string

	// Loop through the config's PublicChannels map
	for _, channels := range config.PublicChannels {
		for _, channelName := range channels {
			// Define a new channel to create
			channel := &model.Channel{
				TeamId:      t.Team.Id,
				Name:        createChannelName(channelName), // Convert name to a valid channel name
				DisplayName: channelName,
				Type:        model.ChannelTypeOpen, // Public channel
			}

			// Create the channel using the Mattermost API
			_, appErr := t.c.API.CreateChannel(channel)
			if appErr != nil {
				// If an error occurs, append it to the result and continue
				result += fmt.Sprintf("Failed to create channel %s: %v\n", channelName, appErr.Error())
				continue
			}

			// Append success message to the result
			result += fmt.Sprintf("Created channel: %s\n", channelName)
		}
	}

	return result
}
