package business

import (
	"github.com/glass.plugin-anchor/server/models"
	"github.com/mattermost/mattermost-server/v6/model"
	"strings"
)

func GetChannelByDisplayName(c *models.Context, displayName string) (*model.Channel, *model.AppError) {
	channelName := strings.ToLower(strings.ReplaceAll(displayName, " ", "-"))

	channel, appErr := c.API.GetChannelByName(c.Team.Id, channelName, false)

	if appErr != nil {
		return nil, appErr
	}

	return channel, nil
}

func createChannelName(displayName string) string {
	return strings.ReplaceAll(strings.ToLower(displayName), " ", "-")
}

//func getChannelNameByID(c *models.Context, channelID string) (string, error) {
//	channel, err := c.API.GetChannel(channelID)
//	if err != nil {
//		return "", fmt.Errorf("could not retrieve channel with ID %s: %w", channelID, err)
//	}
//	return channel.Name, nil
//}
