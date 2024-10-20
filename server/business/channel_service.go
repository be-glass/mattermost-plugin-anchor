package business

import (
	"fmt"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"strings"
)

func GetChannelIDByName(api plugin.API, teamID string, channelName string) (string, error) {
	channel, appErr := api.GetChannelByName(channelName, teamID, false)
	if appErr != nil {
		return "", appErr
	}
	return channel.Id, nil
}

func GetChannelsListString(api plugin.API, teamID string) string {
	// Get all channels in the given team
	channels, err := getAllChannelsInTeam(api, teamID)
	if err != nil {
		// Return the error message as part of the string
		return fmt.Sprintf("Error fetching channels: %v", err)
	}

	// Use a string builder for efficient string concatenation
	var builder strings.Builder

	// Loop through all channels and append their names to the builder
	for _, channel := range channels {
		builder.WriteString(fmt.Sprintf("%s\n", channel.Name))
	}

	// Convert the builder to a string and return it
	return builder.String()
}

// private

func getAllChannelsInTeam(api plugin.API, teamID string) ([]*model.Channel, error) {
	var allChannels []*model.Channel
	page := 0
	perPage := 100 // You can adjust this to change how many channels are fetched per page

	for {
		// Get channels for the current page in the team
		channels, appErr := api.GetPublicChannelsForTeam(teamID, page, perPage)
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
