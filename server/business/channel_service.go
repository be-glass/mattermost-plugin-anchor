package business

import (
	"fmt"
	"github.com/glass.plugin-anchor/server/config"
	"github.com/glass.plugin-anchor/server/models"
	"github.com/mattermost/mattermost-server/v6/model"
	"strings"
)

//func GetChannelIDByName(c *models.Context, channelName string) (string, error) {
//	channel, appErr := c.API.GetChannelByName(channelName, teamID, false)
//	if appErr != nil {
//		return "", appErr
//	}
//	return channel.Id, nil
//}

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

func CheckUserChannelStructure(c *models.Context, userId string) string {
	var resultBuilder strings.Builder

	// Append the user's name to the result
	user, appErr := c.API.GetUser(userId)
	if appErr != nil {
		return "Unable to retrieve user information."
	}
	resultBuilder.WriteString(fmt.Sprintf("User: **%s**\n", user.Username))

	// Call the function to check the sidebar categories
	sidebarResult := checkUserSidebarCategories(c, userId)
	resultBuilder.WriteString(sidebarResult + "\n")

	// Call the function to check the channel subscriptions
	channelResult := checkChannelSubscription(c, userId)
	resultBuilder.WriteString(channelResult + "\n")

	// Call the function to check the channel categorization
	categorizationResult := checkChannelCategorization(c, userId)
	resultBuilder.WriteString(categorizationResult + "\n")

	return resultBuilder.String()
}

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
			userStructureResult := CheckUserChannelStructure(c, user.Id)
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

func GetUserSubscribedPublicChannels(c *models.Context, userId string) ([]*model.Channel, error) {
	var publicChannels []*model.Channel

	// Get all public channels for the team
	channels, appErr := c.API.GetPublicChannelsForTeam(c.Team.Id, 0, 10000) // Adjust the limit if necessary
	if appErr != nil {
		return nil, appErr
	}

	// Iterate over the channels and check if the user is a member of the public channel
	for _, channel := range channels {
		if channel.Type == model.ChannelTypeOpen { // Public channel type
			_, memberErr := c.API.GetChannelMember(channel.Id, userId)
			if memberErr == nil { // If the user is a member, add to the list
				publicChannels = append(publicChannels, channel)
			}
		}
	}

	return publicChannels, nil
}

func checkChannelSubscription(c *models.Context, userId string) string {
	// Get the list of public channels the user is subscribed to
	publicChannels, err := GetUserSubscribedPublicChannels(c, userId)
	if err != nil {
		return "Unable to retrieve user subscribed public channels."
	}

	// Convert the list of public channel names the user is subscribed to into a map for easier lookup
	subscribedChannelNames := make(map[string]bool)
	for _, channel := range publicChannels {
		subscribedChannelNames[channel.DisplayName] = true
	}

	// Get the default channel names
	defaultChannelNames := GetDefaultChannelNames()

	// Create a slice to accumulate missing channels
	var missingChannels []string

	// Check if all default channels are present in the user's subscribed public channels
	for _, defaultChannel := range defaultChannelNames {
		if !subscribedChannelNames[defaultChannel] {
			missingChannels = append(missingChannels, defaultChannel)
		}
	}

	// If no channels are missing, return a success message
	if len(missingChannels) == 0 {
		//return "User is subscribed to all required channels."
		return "."

	}

	// Return the missing channels as a comma-separated string
	return "Missing required channels: " + strings.Join(missingChannels, ", ")
}

func CheckUserOrAll(c *models.Context, targetUser string) string {

	if targetUser == "all" {
		return CheckUserChannelStructureForTeam(c)
	} else {
		userID, err := GetUserIDByUsername(c, targetUser)
		if err != nil {
			return "Could not find that user"
		} else {
			return CheckUserChannelStructure(c, userID)
		}
	}
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

func AddUserToMissingChannels(c *models.Context, userID string, categoryChannels map[string][]string) string {
	var resultBuilder strings.Builder

	// Loop through the categories and their corresponding channels
	for _, channels := range categoryChannels {
		for _, displayName := range channels {
			// Get the channel by display name and team ID
			channel, appErr := GetChannelByDisplayName(c, displayName)
			if appErr != nil || channel == nil {
				resultBuilder.WriteString(fmt.Sprintf("Channel not found: %s\n", displayName))
				continue
			}

			// Check if the user is already a member of the channel
			_, appErr = c.API.GetChannelMember(channel.Id, userID)
			if appErr != nil {
				// If the user is not a member, add them to the channel
				resultBuilder.WriteString(fmt.Sprintf("User is not a member of %s. Adding to channel...\n", displayName))
				_, addErr := c.API.AddChannelMember(channel.Id, userID)
				if addErr != nil {
					resultBuilder.WriteString(fmt.Sprintf("Failed to add user to channel: %s\n", displayName))
				} else {
					resultBuilder.WriteString(fmt.Sprintf("Successfully added user to channel: %s\n", displayName))
				}
			} else {
				resultBuilder.WriteString(fmt.Sprintf("User is already a member of channel: %s\n", displayName))
			}
		}
	}

	return resultBuilder.String()
}
