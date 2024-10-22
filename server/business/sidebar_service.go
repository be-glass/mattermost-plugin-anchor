package business

import (
	"encoding/json"
	"fmt"
	"github.com/glass.plugin-anchor/server/config"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"io/ioutil"
	"net/http"
	"strings"
)

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

func getOrCreateSidebarCategory(api plugin.API, userID string, teamID string, categoryName string) (*model.SidebarCategoryWithChannels, error) {
	// Fetch the user's sidebar categories for the specified team
	categories, appErr := api.GetChannelSidebarCategories(userID, teamID)
	if appErr != nil {
		return nil, appErr
	}

	// Look for the category by name
	for _, category := range categories.Categories {
		if category.DisplayName == categoryName {
			return category, nil // Return the existing category
		}
	}

	// If the category does not exist, create it
	newCategory := &model.SidebarCategoryWithChannels{
		SidebarCategory: model.SidebarCategory{
			UserId:      userID,
			TeamId:      teamID,
			DisplayName: categoryName,
			Type:        model.SidebarCategoryCustom, // Custom category
		},
	}

	createdCategory, appErr := api.CreateChannelSidebarCategory(userID, teamID, newCategory)
	if appErr != nil {
		return nil, appErr
	}

	return createdCategory, nil
}

func addChannelToSidebarCategory(api plugin.API, userID string, teamID string, category *model.SidebarCategoryWithChannels, channelID string) *model.AppError {
	// Get the current list of channel IDs from the category using ChannelIds()
	channelIDs := category.ChannelIds()

	// Check if the channel is already in the category
	for _, existingChannelID := range channelIDs {
		if existingChannelID == channelID {
			// Channel is already in the category, no need to update
			return nil
		}
	}

	// Add the channel ID to the list of channel IDs
	newChannelIDs := append(channelIDs, channelID)

	// Create a new category with the updated channel list
	updatedCategory := &model.SidebarCategoryWithChannels{
		SidebarCategory: model.SidebarCategory{
			Id:          category.Id,
			UserId:      category.UserId,
			TeamId:      category.TeamId,
			DisplayName: category.DisplayName,
			Type:        category.Type,
		},
		Channels: newChannelIDs, // Assign the updated channel list
	}

	// Update the sidebar category with the new channel list
	_, appErr := api.UpdateChannelSidebarCategories(userID, teamID, []*model.SidebarCategoryWithChannels{updatedCategory})
	return appErr
}

func CheckAndJoinDefaultChannelStructure(api plugin.API, targetUser string, teamID string) string {
	var resultBuilder strings.Builder

	// Get the user ID based on the target username
	userID, err := GetUserIDByUsername(api, targetUser)
	if err != nil {
		return "Could not find that user"
	}

	// Add user to missing channels
	resultBuilder.WriteString(AddUserToMissingChannels(api, userID, teamID, config.PublicChannels))

	// Create missing sidebar categories
	sidebarCategories := createMissingSidebarCategories(api, userID, teamID, config.PublicChannels)

	// Assign channels to the created categories
	resultBuilder.WriteString(assignChannelsToCategories(api, userID, teamID, sidebarCategories, config.PublicChannels))

	// Reorder sidebar categories based on the order in config.PublicChannels
	categoryOrder := make([]string, 0, len(config.PublicChannels))
	for category := range config.PublicChannels {
		categoryOrder = append(categoryOrder, category)
	}
	resultBuilder.WriteString(reorderSidebarCategories(api, userID, teamID, categoryOrder))

	return resultBuilder.String()
}

func createMissingSidebarCategories(api plugin.API, userID string, teamID string, categoryChannels map[string][]string) map[string]*model.SidebarCategoryWithChannels {
	sidebarCategories := make(map[string]*model.SidebarCategoryWithChannels)

	for category := range categoryChannels {
		// Get or create the sidebar category for the user in the specified team
		sidebarCategory, err := getOrCreateSidebarCategory(api, userID, teamID, category)
		if err != nil {
			continue
		}

		// Store the category for later use in assigning channels
		sidebarCategories[category] = sidebarCategory
	}

	return sidebarCategories
}
func assignChannelsToCategories(api plugin.API, userID string, teamID string, sidebarCategories map[string]*model.SidebarCategoryWithChannels, categoryChannels map[string][]string) string {
	var resultBuilder strings.Builder

	// Loop through the categories and assign channels
	for category, channels := range categoryChannels {
		sidebarCategory := sidebarCategories[category]
		channelIDs := sidebarCategory.ChannelIds()

		// Collect all new channel IDs that need to be added to the category
		for _, displayName := range channels {
			// Get the channel by display name and team ID
			channel, appErr := GetChannelByDisplayName(api, teamID, displayName)
			if appErr != nil || channel == nil {
				resultBuilder.WriteString(fmt.Sprintf("Channel not found: %s\n", displayName))
				continue
			}

			// Check if the channel is already in the category
			channelAlreadyInCategory := false
			for _, existingChannelID := range channelIDs {
				if existingChannelID == channel.Id {
					channelAlreadyInCategory = true
					break
				}
			}

			// If the channel is not in the category, add it to the list
			if !channelAlreadyInCategory {
				channelIDs = append(channelIDs, channel.Id)
				resultBuilder.WriteString(fmt.Sprintf("Queued channel %s to be added to category %s\n", displayName, category))
			} else {
				resultBuilder.WriteString(fmt.Sprintf("Channel %s already in category %s\n", displayName, category))
			}
		}

		// Update the sidebar category with the complete list of channels at once
		sidebarCategoryWithUpdatedChannels := &model.SidebarCategoryWithChannels{
			SidebarCategory: sidebarCategory.SidebarCategory,
			Channels:        channelIDs, // Updated list of channels
		}

		// Apply the batch update
		_, appErr := api.UpdateChannelSidebarCategories(userID, teamID, []*model.SidebarCategoryWithChannels{sidebarCategoryWithUpdatedChannels})
		if appErr != nil {
			resultBuilder.WriteString(fmt.Sprintf("Failed to update sidebar category %s: %s\n", category, appErr.Error()))
		} else {
			resultBuilder.WriteString(fmt.Sprintf("Successfully updated sidebar category %s with all channels\n", category))
		}
	}

	return resultBuilder.String()
}

func reorderSidebarCategories(api plugin.API, userID string, teamID string, categoryOrder []string) string {
	var resultBuilder strings.Builder

	// Get the current sidebar categories for the user
	sidebarCategories, appErr := api.GetChannelSidebarCategories(userID, teamID)
	if appErr != nil {
		resultBuilder.WriteString(fmt.Sprintf("Failed to retrieve sidebar categories: %s\n", appErr.Error()))
		return resultBuilder.String()
	}

	// Create a map of category by display name for quick lookup
	categoryMap := make(map[string]*model.SidebarCategoryWithChannels)
	for _, category := range sidebarCategories.Categories {
		categoryMap[category.DisplayName] = category
	}

	// Build the ordered list of categories based on config.PublicChannels order
	var orderedCategories []*model.SidebarCategoryWithChannels
	for _, category := range categoryOrder {
		if cat, exists := categoryMap[category]; exists {
			orderedCategories = append(orderedCategories, cat)
		} else {
			resultBuilder.WriteString(fmt.Sprintf("Category not found for reordering: %s\n", category))
		}
	}

	// Apply the batch update with the reordered categories
	_, appErr = api.UpdateChannelSidebarCategories(userID, teamID, orderedCategories)
	if appErr != nil {
		resultBuilder.WriteString(fmt.Sprintf("Failed to reorder sidebar categories: %s\n", appErr.Error()))
	} else {
		resultBuilder.WriteString("Successfully reordered sidebar categories\n")
	}

	return resultBuilder.String()
}

// Function to delete all sidebar categories for a given user and team using the HTTP API
func DeleteAllSidebarCategories(api plugin.API, userID string, teamID string, token string) string {
	var resultBuilder strings.Builder

	client := &http.Client{}

	// Fetch all sidebar categories using the HTTP API
	categoriesEndpoint := fmt.Sprintf("http://localhost:8065/api/v4/users/%s/teams/%s/channels/categories", userID, teamID)
	req, err := http.NewRequest("GET", categoriesEndpoint, nil)
	if err != nil {
		return fmt.Sprintf("Failed to create request: %s\n", err.Error())
	}
	req.Header.Add("Authorization", "Bearer "+token)

	// Make the request to Mattermost API
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Sprintf("Failed to fetch sidebar categories: %s\n", err.Error())
	}
	defer resp.Body.Close()

	// Read and parse the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("Failed to read response body: %s\n", err.Error())
	}

	// Parse the response JSON to get category information
	categories, err := parseCategories(body)
	if err != nil {
		return fmt.Sprintf("Failed to parse categories: %s\n", err.Error())
	}

	// Loop through and delete each category
	for _, category := range categories {
		deleteEndpoint := fmt.Sprintf("http://localhost:8065/api/v4/users/%s/teams/%s/channels/categories/%s", userID, teamID, category.Id)
		req, err := http.NewRequest("DELETE", deleteEndpoint, nil)
		if err != nil {
			resultBuilder.WriteString(fmt.Sprintf("Failed to create delete request for category %s: %s\n", category.DisplayName, err.Error()))
			continue
		}
		req.Header.Add("Authorization", "Bearer "+token)

		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			resultBuilder.WriteString(fmt.Sprintf("Failed to delete category %s: %s\n", category.DisplayName, err.Error()))
		} else {
			resultBuilder.WriteString(fmt.Sprintf("Successfully deleted category: %s\n", category.DisplayName))
		}
	}

	return resultBuilder.String()
}
func parseCategories(body []byte) ([]*model.SidebarCategoryWithChannels, error) {
	var categories []*model.SidebarCategoryWithChannels
	err := json.Unmarshal(body, &categories)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (api *MattermostAPI) DeleteUserCategory(userID, teamID, categoryID string) ([]byte, error) {
	path := fmt.Sprintf("users/%s/teams/%s/channels/categories/%s", userID, teamID, categoryID)
	return api.Delete(path)
}
