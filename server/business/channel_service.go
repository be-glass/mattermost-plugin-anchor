package business

import (
	"fmt"
	"github.com/glass.plugin-anchor/server/config"
	"github.com/glass.plugin-anchor/server/models"
	"github.com/mattermost/mattermost-server/v6/model"
	"strings"
)

func GetChannelsListString(c *models.Context) string {
	// Get all channels in the given team
	channels, err := getAllChannelsInTeam(c)
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

// private

func CheckUserChannelStructureForTeam(c *models.Context) string {
	var resultBuilder strings.Builder

	// Get all users in the given team
	page := 0
	perPage := 100
	for {
		// Retrieve a page of users
		users, appErr := c.API.GetUsersInTeam(c.Team.Id, page, perPage)
		if appErr != nil {
			return "Unable to retrieve users in the team."
		}

		// If no users are found, stop the loop
		if len(users) == 0 {
			break
		}

		// Iterate over the users and check their channel structure
		for _, user := range users {
			u := WrapUser(user)
			userStructureResult := u.CheckUserChannelStructure(c)
			resultBuilder.WriteString(userStructureResult + "\n")
		}

		// Increment the page to get the next set of users
		page++
	}

	return resultBuilder.String()
}

func getAllChannelsInTeam(c *models.Context) ([]*model.Channel, error) {
	var allChannels []*model.Channel
	page := 0
	perPage := 100 // You can adjust this to change how many channels are fetched per page

	for {
		// Get channels for the current page in the team
		channels, appErr := c.API.GetPublicChannelsForTeam(c.Team.Id, page, perPage)
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

func GetDefaultChannelNames() []string {
	var channels []string

	for _, channelList := range config.PublicChannels {
		channels = append(channels, channelList...)
	}
	return channels
}

func GetDefaultCategoryNames() []string {
	var categories []string

	for category := range config.PublicChannels {
		categories = append(categories, category)
	}
	return categories
}

func GetChannelByDisplayName(c *models.Context, displayName string) (*model.Channel, *model.AppError) {
	// Convert the display name to channel name format
	channelName := strings.ToLower(strings.ReplaceAll(displayName, " ", "-"))

	// Use the converted name to get the channel
	channel, appErr := c.API.GetChannelByName(c.Team.Id, channelName, false)

	if appErr != nil {
		c.API.LogWarn("DBG NOT FOUND", displayName, channelName, c.Team.Id, appErr)
		return nil, appErr // Return error if the channel is not found
	}

	c.API.LogWarn("DBG Found it :)", displayName, channelName)
	return channel, nil
}

// createChannelName converts a string to lowercase and replaces spaces with hyphens
func createChannelName(name string) string {
	// Convert to lowercase and replace spaces with hyphens
	return strings.ReplaceAll(strings.ToLower(name), " ", "-")
}

// CreateDefaultChannels creates default public channels for the given team ID
func CreateDefaultChannels(c *models.Context) string {
	var result string

	// Loop through the config's PublicChannels map
	for _, channels := range config.PublicChannels {
		for _, channelName := range channels {
			// Define a new channel to create
			channel := &model.Channel{
				TeamId:      c.Team.Id,
				Name:        createChannelName(channelName), // Convert name to a valid channel name
				DisplayName: channelName,
				Type:        model.ChannelTypeOpen, // Public channel
			}

			// Create the channel using the Mattermost API
			_, appErr := c.API.CreateChannel(channel)
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
