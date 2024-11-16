package business

import (
	"errors"
	"fmt"
	"github.com/glass.plugin-anchor/server/config"
	"github.com/glass.plugin-anchor/server/models"
	"github.com/glass.plugin-anchor/server/utils"
	"github.com/mattermost/mattermost-server/v6/model"
	"strings"
)

type SideBar struct {
	c          *models.Context
	u          *User
	User       *model.User
	categories *model.OrderedSidebarCategories
}

func NewSideBar(user *User) (*SideBar, error) {

	//user.c.API.LogDebug("NEWSIDEBAR XXXXX", user.Id, user.Username)

	sidebar := &SideBar{
		c:          user.c,
		u:          user,
		User:       user.User,
		categories: nil,
	}

	err := sidebar.fetch()
	if err != nil {
		return nil, err
	}

	return sidebar, nil
}

func (s *SideBar) fetch() *model.AppError {
	var err *model.AppError
	s.categories, err = s.c.API.GetChannelSidebarCategories(s.u.Id, s.c.Team.Id)
	return err
}

func (s *SideBar) SidebarCategoryNames() ([]string, error) {
	var categories []string

	for _, category := range s.categories.Categories {
		categories = append(categories, category.DisplayName)
	}

	return categories, nil
}

func (s *SideBar) checkSidebarCategories() string {
	// Get the list of category names in the user's sidebar for the given team
	userCategories, err := s.SidebarCategoryNames()
	if err != nil {
		return "Unable to retrieve user sidebar categories."
	}

	// Convert the list of user's sidebar category names into a map for easier lookup
	userCategoryMap := make(map[string]bool)
	for _, category := range userCategories {
		userCategoryMap[category] = true
	}

	// Get the default category names
	defaultCategoryNames := config.CategoryNames()

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

func (s *SideBar) checkChannelCategorization() string {
	// Get the list of public channels the user is subscribed to
	publicChannels, err := s.u.GetSubscribedPublicChannels()
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
	userSidebarCategories, appErr := s.c.API.GetChannelSidebarCategories(s.User.Id, s.c.Team.Id)
	if appErr != nil {
		return "Unable to retrieve user sidebar categories."
	}

	// Map actual categories from the sidebar for easier lookup
	userCategoryMap := make(map[string]string)
	for _, sidebarCategory := range userSidebarCategories.Categories {
		for _, channelId := range sidebarCategory.Channels {
			channel, err := s.c.API.GetChannel(channelId)
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

func (s *SideBar) getOrCreateSidebarCategory(categoryName string) (*model.SidebarCategoryWithChannels, error) {
	// Fetch the user's sidebar categories for the specified team
	categories, appErr := s.c.API.GetChannelSidebarCategories(s.User.Id, s.c.Team.Id)
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
			UserId:      s.User.Id,
			TeamId:      s.c.Team.Id,
			DisplayName: categoryName,
			Type:        model.SidebarCategoryCustom, // Custom category
		},
	}

	createdCategory, appErr := s.c.API.CreateChannelSidebarCategory(s.User.Id, s.c.Team.Id, newCategory)
	if appErr != nil {
		return nil, appErr
	}

	return createdCategory, nil
}

func (s *SideBar) addChannelToSidebarCategory(category *model.SidebarCategoryWithChannels, channelID string) *model.AppError {
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
	_, appErr := s.c.API.UpdateChannelSidebarCategories(s.User.Id, s.c.Team.Id, []*model.SidebarCategoryWithChannels{updatedCategory})
	return appErr
}

func (s *SideBar) CheckAndJoinDefaultChannelStructure() string {

	resultChannels := s.u.JoinMissingChannels(config.PublicChannels)
	resultJoin := s.u.JoinMissingChannels(config.PublicChannels)

	sidebarCategories := s.createMissingSidebarCategories(config.PublicChannels)
	resultCategories := s.assignChannelsToCategories(sidebarCategories, config.PublicChannels)

	resultReorder := s.ReorderSidebarCategories()

	return strings.Join([]string{
		resultChannels,
		resultJoin,
		resultCategories,
		resultReorder,
	}, "\n")

}

func (s *SideBar) createMissingSidebarCategories(categoryChannels map[string][]string) map[string]*model.SidebarCategoryWithChannels {
	sidebarCategories := make(map[string]*model.SidebarCategoryWithChannels)
	var orderedCategories []*model.SidebarCategoryWithChannels

	// Iterate over the categories in the order defined in PublicChannels
	for category, channels := range categoryChannels {
		// Get or create the sidebar category for the user in the specified team
		sidebarCategory, err := s.getOrCreateSidebarCategory(category)
		if err != nil {
			continue
		}

		// Create a slice to hold the channel IDs
		var channelIDs []string

		// Populate the slice with channel IDs by fetching them using their names
		for _, channelName := range channels {
			channel, appErr := s.c.API.GetChannelByNameForTeamName(s.c.Team.Name, channelName, false)
			if appErr != nil {
				// Log or handle the error, but skip this channel if it fails
				continue
			}
			channelIDs = append(channelIDs, channel.Id)
		}

		// Update the category's channel list with the appropriate channel IDs
		sidebarCategory.Channels = channelIDs

		// Store the category in the result map for returning later
		sidebarCategories[category] = sidebarCategory

		// Append the category to the ordered list to apply all at once
		orderedCategories = append(orderedCategories, sidebarCategory)
	}

	// Reverse the ordered categories slice
	for i, j := 0, len(orderedCategories)-1; i < j; i, j = i+1, j-1 {
		orderedCategories[i], orderedCategories[j] = orderedCategories[j], orderedCategories[i]
	}

	// Now, apply all the created/retrieved categories in one batch API call
	_, appErr := s.c.API.UpdateChannelSidebarCategories(s.User.Id, s.c.Team.Id, orderedCategories)
	if appErr != nil {
		// Handle error (logging, etc.)
	}

	return sidebarCategories
}

func (s *SideBar) assignChannelsToCategories(sidebarCategories map[string]*model.SidebarCategoryWithChannels, categoryChannels map[string][]string) string {
	var resultBuilder strings.Builder

	// Loop through the categories and assign channels
	for category, channels := range categoryChannels {
		sidebarCategory := sidebarCategories[category]
		channelIDs := sidebarCategory.ChannelIds()

		// Collect all new channel IDs that need to be added to the category
		for _, displayName := range channels {
			// Get the channel by display name and team ID
			channel, appErr := GetChannelByDisplayName(s.c, displayName)
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
		_, appErr := s.c.API.UpdateChannelSidebarCategories(s.User.Id, s.c.Team.Id, []*model.SidebarCategoryWithChannels{sidebarCategoryWithUpdatedChannels})
		if appErr != nil {
			resultBuilder.WriteString(fmt.Sprintf("Failed to update sidebar category %s: %s\n", category, appErr.Error()))
		} else {
			resultBuilder.WriteString(fmt.Sprintf("Successfully updated sidebar category %s with all channels\n", category))
		}
	}

	return resultBuilder.String()
}

func (s *SideBar) findCategory(categoryName string) (*model.SidebarCategoryWithChannels, error) {
	var matchedCategory *model.SidebarCategoryWithChannels

	actualCategories, err := s.c.API.GetChannelSidebarCategories(s.User.Id, s.c.Team.Id)

	if err != nil {
		return nil, err
	}

	for _, category := range actualCategories.Categories {
		if category.DisplayName == categoryName {
			matchedCategory = category
			break
		}
	}
	return matchedCategory, nil
}

func categoryChannelIDs(c *models.Context, categoryName string) ([]string, error) {
	var orderedChannelIDs []string

	channelNames, exists := config.AllChannels()[categoryName]
	if !exists {
		return nil, errors.New("category not found " + categoryName)
	}

	for _, channelName := range channelNames {

		channel, err := GetChannelByDisplayName(c, channelName)
		if err != nil {
			continue
		}
		orderedChannelIDs = append(orderedChannelIDs, channel.Id)
	}
	return orderedChannelIDs, nil
}

func (s *SideBar) ReorderSidebarCategories() string {
	var dbg []string

	s.c.API.LogDebug("XXXXX ReorderSidebarCategories", s.User.Username, s.c.Team.Id)

	var updatedCategories []*model.SidebarCategoryWithChannels

	for index, categoryName := range config.CategoryOrder {

		category, err := s.findCategory(categoryName)
		if err != nil {
			return "err.Error()"
		}

		orderedChannelIDs, err := categoryChannelIDs(s.c, categoryName)
		if err != nil {
			return err.Error()
		}

		s.c.API.LogDebug("XXXXX Loop", category.DisplayName, len(orderedChannelIDs), index)

		updatedCategory := newSidebarCategory(category, orderedChannelIDs, 10+10*index)
		updatedCategories = append(updatedCategories, updatedCategory)
	}

	for _, category := range updatedCategories {
		dbg = append(dbg, fmt.Sprintf("%s - %d", category.DisplayName, category.SortOrder))
	}
	dbg = append(dbg, "->>>")

	if _, err := s.c.API.UpdateChannelSidebarCategories(s.User.Id, s.c.Team.Id, updatedCategories); err != nil {
		return err.Error()
	}

	err := s.fetch()
	if err != nil {
		return err.Error()
	}

	for _, category := range s.categories.Categories {
		dbg = append(dbg, fmt.Sprintf("%s - %d", category.DisplayName, category.SortOrder))
	}

	return strings.Join(dbg, "\n")
}

//func (s *SideBar) ReorderSidebarCategories_OLD() string {
//
//	var updatedCategories []*model.SidebarCategoryWithChannels
//
//	channelMap := config.AllChannels()
//	for sortOrder, categoryName := range config.CategoryOrder {
//		channelNames, exists := channelMap[categoryName]
//		if !exists {
//			continue
//		}
//
//		category = findCategory()
//
//		// Skip if no matching category is found in actualCategories
//		if matchedCategory == nil {
//			continue
//		}
//
//		// Retrieve the current channel IDs for this category
//		currentChannelIDs := matchedCategory.ChannelIds()
//
//		// Order channels based on channelMap
//		var orderedChannelIDs []string
//		for _, channelName := range channelNames {
//			for _, channelID := range currentChannelIDs {
//				name, err := getChannelNameByID(s.c, channelID)
//				if err != nil {
//					continue // Skip if channel name cannot be retrieved
//				}
//				if name == channelName {
//					orderedChannelIDs = append(orderedChannelIDs, channelID)
//				}
//			}
//		}
//
//		newCategory := newSidebarCategory(matchedCategory, orderedChannelIDs, sortOrder)
//		updatedCategories = append(updatedCategories, newCategory)
//	}
//
//	if _, err := s.c.API.UpdateChannelSidebarCategories(u.Id, s.c.Team.Id, updatedCategories); err != nil {
//		return err.Error()
//	}
//
//	// Debug output to confirm order
//	var dbg []string
//	for _, category := range updatedCategories {
//		dbg = append(dbg, category.DisplayName)
//	}
//	return strings.Join(dbg, "\n")
//}

func newSidebarCategory(m *model.SidebarCategoryWithChannels, orderedChannelIDs []string, sortOrder int) *model.SidebarCategoryWithChannels {
	return &model.SidebarCategoryWithChannels{
		SidebarCategory: model.SidebarCategory{
			Id:          m.Id,
			UserId:      m.UserId,
			TeamId:      m.TeamId,
			SortOrder:   int64(sortOrder), // Set SortOrder based on the sortOrder
			Sorting:     "manual",
			Type:        m.Type,
			DisplayName: m.DisplayName,
			Muted:       m.Muted,
			Collapsed:   m.Collapsed,
		},
		Channels: orderedChannelIDs,
	}
}

func (s *SideBar) DeleteAllSidebarCategories() string {

	var names = []string{"Deleting... "}

	sidebarCategories, appErr := s.c.API.GetChannelSidebarCategories(s.User.Id, s.c.Team.Id)
	if appErr != nil {
		return appErr.DetailedError
	}

	for _, category := range sidebarCategories.Categories {

		if utils.Contains(config.DefaultCategories, category.DisplayName) {
			continue
		}

		names = append(names, category.DisplayName)
		_, err := s.DeleteCategory(category.Id)
		if err != nil {
			names = append(names, fmt.Sprintf("Could not delete **%s** because **%s**\n", category.DisplayName, err.Error()))
		}
	}

	return strings.Join(names, ", ")
}

func (s *SideBar) DeleteCategory(categoryID string) ([]byte, error) {
	path := fmt.Sprintf("users/%s/teams/%s/channels/categories/%s", s.User.Id, s.c.Team.Id, categoryID)
	return s.c.Rest.Delete(path)
}

func (s *SideBar) SetCategoryOrder(categoryIDsOrdered []string) ([]byte, error) {
	path := fmt.Sprintf("users/%s/teams/%s/channels/categories/order", s.User.Id, s.c.Team.Id)
	return s.c.Rest.Put(path, categoryIDsOrdered)
}

//func (s *SideBar) CategoryDetails() string {
//
//	categories, err := s.c.API.GetChannelSidebarCategories(s.User.Id, s.c.Team.Id)
//	if err != nil {
//		return "Could not retrieve some categories"
//	}
//
//	answer := []string{}
//
//	for i, category := range categories.Categories {
//		category.SortOrder = int64(100 - i)
//	}
//
//	s.c.API.UpdateChannelSidebarCategories(s.User.Id, s.c.Team.Id, categories)
//
//	categories2, err2 := s.c.API.GetChannelSidebarCategories(u.Id, s.c.Team.Id)
//	if err2 != nil {
//		return "Could not retrieve some categories"
//	}
//
//	for _, category := range categories2.Categories {
//		answer = append(answer, fmt.Sprintf("%s: %d", category.DisplayName, category.SortOrder))
//	}
//
//	return strings.Join(answer, "\n")
//}

func (s *SideBar) CheckChannelStructure() string {

	return fmt.Sprintf("User: **%s** (%s):\n%s\n%s\n%s\n",
		s.User.Username,
		s.User.GetFullName(),
		s.checkSidebarCategories(),
		s.u.checkChannelSubscription(),
		s.checkChannelCategorization())
}
