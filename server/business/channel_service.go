package business

import (
	"fmt"
	"github.com/glass.plugin-anchor/server/config"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"strings"
)

//func GetChannelIDByName(api plugin.API, teamID string, channelName string) (string, error) {
//	channel, appErr := api.GetChannelByName(channelName, teamID, false)
//	if appErr != nil {
//		return "", appErr
//	}
//	return channel.Id, nil
//}

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
		builder.WriteString(fmt.Sprintf("%s\n", channel.DisplayName))
	}

	// Convert the builder to a string and return it
	return builder.String()
}

// private

func CheckUserChannelStructure(api plugin.API, userId string, teamId string) string {
	var resultBuilder strings.Builder

	// Append the user's name to the result
	user, appErr := api.GetUser(userId)
	if appErr != nil {
		return "Unable to retrieve user information."
	}
	resultBuilder.WriteString(fmt.Sprintf("User: **%s**\n", user.Username))

	// Call the function to check the sidebar categories
	sidebarResult := checkUserSidebarCategories(api, userId, teamId)
	resultBuilder.WriteString(sidebarResult + "\n")

	// Call the function to check the channel subscriptions
	channelResult := checkChannelSubscription(api, userId, teamId)
	resultBuilder.WriteString(channelResult + "\n")

	// Call the function to check the channel categorization
	categorizationResult := checkChannelCategorization(api, userId, teamId)
	resultBuilder.WriteString(categorizationResult + "\n")

	return resultBuilder.String()
}

func CheckUserChannelStructureForTeam(api plugin.API, teamId string) string {
	var resultBuilder strings.Builder

	// Get all users in the given team
	page := 0
	perPage := 100
	for {
		// Retrieve a page of users
		users, appErr := api.GetUsersInTeam(teamId, page, perPage)
		if appErr != nil {
			return "Unable to retrieve users in the team."
		}

		// If no users are found, stop the loop
		if len(users) == 0 {
			break
		}

		// Iterate over the users and check their channel structure
		for _, user := range users {
			userStructureResult := CheckUserChannelStructure(api, user.Id, teamId)
			resultBuilder.WriteString(userStructureResult + "\n")
		}

		// Increment the page to get the next set of users
		page++
	}

	return resultBuilder.String()
}

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

func GetUserSubscribedPublicChannels(api plugin.API, userId string, teamId string) ([]*model.Channel, error) {
	var publicChannels []*model.Channel

	// Get all public channels for the team
	channels, appErr := api.GetPublicChannelsForTeam(teamId, 0, 10000) // Adjust the limit if necessary
	if appErr != nil {
		return nil, appErr
	}

	// Iterate over the channels and check if the user is a member of the public channel
	for _, channel := range channels {
		if channel.Type == model.ChannelTypeOpen { // Public channel type
			_, memberErr := api.GetChannelMember(channel.Id, userId)
			if memberErr == nil { // If the user is a member, add to the list
				publicChannels = append(publicChannels, channel)
			}
		}
	}

	return publicChannels, nil
}

func GetUserSidebarCategoryNames(api plugin.API, userId string, teamId string) ([]string, error) {
	var categories []string

	// Use the Plugin API to get the sidebar categories for the user and team
	sidebarCategories, appErr := api.GetChannelSidebarCategories(userId, teamId)
	if appErr != nil {
		return nil, appErr
	}

	// Extract the display names of the categories
	for _, category := range sidebarCategories.Categories {
		categories = append(categories, category.DisplayName)
	}

	return categories, nil
}

func checkChannelSubscription(api plugin.API, userId string, teamId string) string {
	// Get the list of public channels the user is subscribed to
	publicChannels, err := GetUserSubscribedPublicChannels(api, userId, teamId)
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

func checkUserSidebarCategories(api plugin.API, userId string, teamId string) string {
	// Get the list of category names in the user's sidebar for the given team
	userCategories, err := GetUserSidebarCategoryNames(api, userId, teamId)
	if err != nil {
		return "Unable to retrieve user sidebar categories."
	}

	// Convert the list of user's sidebar category names into a map for easier lookup
	userCategoryMap := make(map[string]bool)
	for _, category := range userCategories {
		userCategoryMap[category] = true
	}

	// Get the default category names
	defaultCategoryNames := GetDefaultCategoryNames()

	// Create a slice to accumulate missing categories
	var missingCategories []string

	// Check if all default categories are present in the user's sidebar categories
	for _, defaultCategory := range defaultCategoryNames {
		if !userCategoryMap[defaultCategory] {
			missingCategories = append(missingCategories, defaultCategory)
		}
	}

	// If no categories are missing, return a success message
	if len(missingCategories) == 0 {
		//return "All required categories are present."
		return "."

	}

	// Return the missing categories as a comma-separated string
	return "Missing required categories: " + strings.Join(missingCategories, ", ")
}

func checkChannelCategorization(api plugin.API, userId string, teamId string) string {
	// Get the list of public channels the user is subscribed to
	publicChannels, err := GetUserSubscribedPublicChannels(api, userId, teamId)
	if err != nil {
		return "Unable to retrieve user subscribed public channels."
	}

	// Create a map to hold the expected category for each channel from ChannelTree
	expectedCategoryMap := make(map[string]string)
	for category, channels := range config.PublicChannels {
		for _, channel := range channels {
			expectedCategoryMap[channel] = category
		}
	}

	// Create a slice to store any wrongly categorized channels
	var wronglyCategorized []string

	// Get the user's actual sidebar categories from the API
	userSidebarCategories, appErr := api.GetChannelSidebarCategories(userId, teamId)
	if appErr != nil {
		return "Unable to retrieve user sidebar categories."
	}

	// Map actual categories from the sidebar for easier lookup
	userCategoryMap := make(map[string]string)
	for _, sidebarCategory := range userSidebarCategories.Categories {
		for _, channelId := range sidebarCategory.Channels {
			channel, err := api.GetChannel(channelId)
			if err == nil {
				userCategoryMap[channel.DisplayName] = sidebarCategory.DisplayName
			}
		}
	}

	// Check if each subscribed channel is in the expected category
	for _, channel := range publicChannels {
		expectedCategory, exists := expectedCategoryMap[channel.DisplayName]
		if !exists {
			continue // If the channel is not in the ChannelTree, skip the check
		}

		actualCategory, isCategorized := userCategoryMap[channel.DisplayName]
		if isCategorized && actualCategory != expectedCategory {
			wronglyCategorized = append(wronglyCategorized, channel.DisplayName+" (expected: "+expectedCategory+", got: "+actualCategory+")")
		}
	}

	// If no channels are wrongly categorized, return a success message
	if len(wronglyCategorized) == 0 {
		return "."
	}

	// Return the wrongly categorized channels as a comma-separated string
	return "Wrongly categorized channels: " + strings.Join(wronglyCategorized, ", ")
}

func CheckUserOrAll(api plugin.API, targetUser string, teamId string) string {

	if targetUser == "all" {
		return CheckUserChannelStructureForTeam(api, teamId)
	} else {
		userID, err := GetUserIDByUsername(api, targetUser)
		if err != nil {
			return "Could not find that user"
		} else {
			return CheckUserChannelStructure(api, userID, teamId)
		}
	}
}

func GetChannelByDisplayName(api plugin.API, teamId string, displayName string) (*model.Channel, *model.AppError) {
	// Convert the display name to channel name format
	channelName := strings.ToLower(strings.ReplaceAll(displayName, " ", "-"))

	// Use the converted name to get the channel
	channel, appErr := api.GetChannelByName(teamId, channelName, false)

	if appErr != nil {
		api.LogWarn("DBG NOT FOUND", displayName, channelName, teamId, appErr)
		return nil, appErr // Return error if the channel is not found
	}

	api.LogWarn("DBG Found it :)", displayName, channelName)
	return channel, nil
}

func CheckAndJoinDefaultChannels(api plugin.API, targetUser string, teamId string) string {
	var resultBuilder strings.Builder

	userID, err := GetUserIDByUsername(api, targetUser)
	if err != nil {
		return "Could not find that user"
	}

	// Loop through the default channel categories and their corresponding channels
	for category, channels := range config.PublicChannels {
		resultBuilder.WriteString(fmt.Sprintf("Checking category: %s\n", category))

		for _, displayName := range channels {
			// Get the channel by name and team ID
			channel, appErr := GetChannelByDisplayName(api, teamId, displayName)
			if appErr != nil || channel == nil {
				resultBuilder.WriteString(fmt.Sprintf("Channel not found: %s\n", displayName))
				continue
			}

			// Check if the user is already a member of the channel
			_, appErr = api.GetChannelMember(channel.Id, userID)
			if appErr != nil {
				// If the user is not a member, add them to the channel
				resultBuilder.WriteString(fmt.Sprintf("User is not a member of %s. Adding to channel...\n", displayName))
				_, addErr := api.AddChannelMember(channel.Id, userID)
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

// CreateDefaultChannels creates default public channels for the given team ID
func CreateDefaultChannels(api plugin.API, teamID string) string {
	var result string

	// Loop through the config's PublicChannels map
	for _, channels := range config.PublicChannels {
		for _, channelName := range channels {
			// Define a new channel to create
			channel := &model.Channel{
				TeamId:      teamID,
				Name:        createChannelName(channelName), // Convert name to a valid channel name
				DisplayName: channelName,
				Type:        model.ChannelTypeOpen, // Public channel
			}

			// Create the channel using the Mattermost API
			_, appErr := api.CreateChannel(channel)
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

// createChannelName converts a string to lowercase and replaces spaces with hyphens
func createChannelName(name string) string {
	// Convert to lowercase and replace spaces with hyphens
	return strings.ReplaceAll(strings.ToLower(name), " ", "-")
}
