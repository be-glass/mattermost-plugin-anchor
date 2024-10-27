package business

import (
	"fmt"
	"github.com/glass.plugin-anchor/server/config"
	"github.com/glass.plugin-anchor/server/models"
	"github.com/mattermost/mattermost-server/v6/model"
	"strings"
)

type User struct {
	C *models.Context
	*model.User
}

// Constructors

func WrapUser(c *models.Context, user *model.User) *User {
	return &User{c, user}
}

func NewUser(c *models.Context, userName string) (*User, error) {
	user, err := c.API.GetUserByUsername(userName)
	if err != nil {
		return nil, err
	}
	return &User{C: c, User: user}, nil
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

func (u *User) CheckChannelStructure() string {

	return fmt.Sprintf("User: **%s** (%s):\n%s\n%s\n%s\n",
		u.Username,
		u.GetFullName(),
		u.checkSidebarCategories(),
		u.checkChannelSubscription(),
		u.checkChannelCategorization())
}

func (u *User) GetSubscribedPublicChannels() ([]*model.Channel, error) {
	var publicChannels []*model.Channel

	// Get all public channels for the team
	channels, appErr := u.C.API.GetPublicChannelsForTeam(u.C.Team.Id, 0, 10000) // Adjust the limit if necessary
	if appErr != nil {
		return nil, appErr
	}

	// Iterate over the channels and check if the user is a member of the public channel
	for _, channel := range channels {
		if channel.Type == model.ChannelTypeOpen { // Public channel type
			_, memberErr := u.C.API.GetChannelMember(channel.Id, u.Id)
			if memberErr == nil { // If the user is a member, add to the list
				publicChannels = append(publicChannels, channel)
			}
		}
	}

	return publicChannels, nil
}

func (u *User) checkChannelSubscription() string {
	// Get the list of public channels the user is subscribed to
	publicChannels, err := u.GetSubscribedPublicChannels()
	if err != nil {
		return "Unable to retrieve user subscribed public channels."
	}

	// Convert the list of public channel names the user is subscribed to into a map for easier lookup
	subscribedChannelNames := make(map[string]bool)
	for _, channel := range publicChannels {
		subscribedChannelNames[channel.DisplayName] = true
	}

	// Get the default channel names
	defaultChannelNames := config.ChannelNames()

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

func (u *User) JoinMissingChannels(categoryChannels map[string][]string) string {
	var resultBuilder strings.Builder

	// Loop through the categories and their corresponding channels
	for _, channels := range categoryChannels {
		for _, displayName := range channels {
			// Get the channel by display name and team ID
			channel, appErr := GetChannelByDisplayName(u.C, displayName)
			if appErr != nil || channel == nil {
				resultBuilder.WriteString(fmt.Sprintf("Channel not found: %s\n", displayName))
				continue
			}

			// Check if the user is already a member of the channel
			_, appErr = u.C.API.GetChannelMember(channel.Id, u.Id)
			if appErr != nil {
				// If the user is not a member, add them to the channel
				resultBuilder.WriteString(fmt.Sprintf("User is not a member of %s. Adding to channel...\n", displayName))
				_, addErr := u.C.API.AddChannelMember(channel.Id, u.Id)
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
