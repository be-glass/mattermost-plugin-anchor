package business

import (
	"fmt"
	"github.com/glass.plugin-anchor/server/config"
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

		// Iterate over the users and check their sidebar categories and channel subscriptions
		for _, user := range users {
			// Append the user's name to the result
			resultBuilder.WriteString(fmt.Sprintf("User: **%s**\n", user.Username))

			// Call the function to check the sidebar categories
			sidebarResult := checkUserSidebarCategories(api, user.Id, teamId)
			resultBuilder.WriteString(sidebarResult + "\n")

			// Call the function to check the channel subscriptions
			channelResult := checkChannelSubscription(api, user.Id, teamId)
			resultBuilder.WriteString(channelResult + "\n")

			categorizationResult := checkChannelCategorization(api, user.Id, teamId)
			resultBuilder.WriteString(categorizationResult + "\n")

			resultBuilder.WriteString("\n")
		}

		// Increment the page to get the next set of users
		page++
	}

	// Return the accumulated results as a single string
	return resultBuilder.String()
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

func getDefaultChannelNames() []string {
	var channels []string

	for _, channelList := range config.ChannelTree {
		channels = append(channels, channelList...)
	}
	return channels
}

func getDefaultCategoryNames() []string {
	var categories []string

	for category := range config.ChannelTree {
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
	defaultChannelNames := getDefaultChannelNames()

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
	defaultCategoryNames := getDefaultCategoryNames()

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

	// Create a map to hold the expected category for each channel from channelTree
	expectedCategoryMap := make(map[string]string)
	for category, channels := range config.ChannelTree {
		for _, channel := range channels {
			expectedCategoryMap[channel] = category
		}
	}

	// Create a slice to store any wrongly categorized channels
	var wronglyCategorized []string

	// Get the user's sidebar categories (to check actual categorization)
	userCategories, err := GetUserSidebarCategoryNames(api, userId, teamId)
	if err != nil {
		return "Unable to retrieve user sidebar categories."
	}

	// Convert user's sidebar categories to a map for easier lookup
	userCategoryMap := make(map[string]string)
	for _, category := range userCategories {
		for _, channel := range publicChannels {
			userCategoryMap[channel.DisplayName] = category // Assuming you can map categories to channels
		}
	}

	// Check if each subscribed channel is in the expected category
	for _, channel := range publicChannels {
		expectedCategory, exists := expectedCategoryMap[channel.DisplayName]
		if !exists {
			continue // If the channel is not in the channelTree, skip the check
		}

		actualCategory, isCategorized := userCategoryMap[channel.DisplayName]
		if isCategorized && actualCategory != expectedCategory {
			wronglyCategorized = append(wronglyCategorized, channel.DisplayName+" (expected: "+expectedCategory+", got: "+actualCategory+")")
		}
	}

	// If no channels are wrongly categorized, return a success message
	if len(wronglyCategorized) == 0 {
		//return "All channels are categorized correctly."
		return "."
	}

	// Return the wrongly categorized channels as a comma-separated string
	return "Wrongly categorized channels: " + strings.Join(wronglyCategorized, ", ")
}
