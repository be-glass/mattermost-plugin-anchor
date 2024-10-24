package business

import (
	"fmt"
	"github.com/glass.plugin-anchor/server/models"
	"github.com/mattermost/mattermost-server/v6/model"
	"strings"
)

type User struct {
	*model.User
}

// Constructors

func WrapUser(user *model.User) *User {
	return &User{user}
}

func NewUser(c *models.Context, userName string) (*User, error) {
	user, err := c.API.GetUserByUsername(userName)
	if err != nil {
		return nil, err
	}
	return &User{User: user}, nil
}

// static

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

// members

func (u *User) CheckUserChannelStructure(c *models.Context) string {
	var resultBuilder strings.Builder

	resultBuilder.WriteString(fmt.Sprintf("User: **%s**\n", u.Username))

	// Call the function to check the sidebar categories
	sidebarResult := u.checkUserSidebarCategories(c)
	resultBuilder.WriteString(sidebarResult + "\n")

	// Call the function to check the channel subscriptions
	channelResult := u.checkChannelSubscription(c)
	resultBuilder.WriteString(channelResult + "\n")

	// Call the function to check the channel categorization
	categorizationResult := u.checkChannelCategorization(c)
	resultBuilder.WriteString(categorizationResult + "\n")

	return resultBuilder.String()
}

func (u *User) GetUserSubscribedPublicChannels(c *models.Context) ([]*model.Channel, error) {
	var publicChannels []*model.Channel

	// Get all public channels for the team
	channels, appErr := c.API.GetPublicChannelsForTeam(c.Team.Id, 0, 10000) // Adjust the limit if necessary
	if appErr != nil {
		return nil, appErr
	}

	// Iterate over the channels and check if the user is a member of the public channel
	for _, channel := range channels {
		if channel.Type == model.ChannelTypeOpen { // Public channel type
			_, memberErr := c.API.GetChannelMember(channel.Id, u.Id)
			if memberErr == nil { // If the user is a member, add to the list
				publicChannels = append(publicChannels, channel)
			}
		}
	}

	return publicChannels, nil
}

func (u *User) checkChannelSubscription(c *models.Context) string {
	// Get the list of public channels the user is subscribed to
	publicChannels, err := u.GetUserSubscribedPublicChannels(c)
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

func (u *User) CheckUserOrAll(c *models.Context) string {

	if u == nil {
		return CheckUserChannelStructureForTeam(c)
	} else {
		return u.CheckUserChannelStructure(c)
	}
}

func (u *User) AddUserToMissingChannels(c *models.Context, categoryChannels map[string][]string) string {
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
			_, appErr = c.API.GetChannelMember(channel.Id, u.Id)
			if appErr != nil {
				// If the user is not a member, add them to the channel
				resultBuilder.WriteString(fmt.Sprintf("User is not a member of %s. Adding to channel...\n", displayName))
				_, addErr := c.API.AddChannelMember(channel.Id, u.Id)
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
