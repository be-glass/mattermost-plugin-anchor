package main

import (
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/plugin"
)

func (p *AnchorPlugin) UserHasJoinedChannel(c *plugin.Context, channelMember *model.ChannelMember, user *model.User) {
	// Get the channel by ID
	channel, appErr := p.API.GetChannel(channelMember.ChannelId)
	if appErr != nil {
		p.API.LogError("Failed to get channel", "channel_id", channelMember.ChannelId, "error", appErr.Error())
		return
	}

	p.API.LogInfo("HOOK CALLED: UserHasJoinedChannel")

	// Check if the channel is the "Master" channel
	if channel.Name == "master" {
		p.API.LogInfo("User has joined the Master channel", "user_id", user.Id, "channel_id", channel.Id)

		// Find the "Follower" channel by name
		followerChannel, appErr := p.API.GetChannelByNameForTeamName("lbw", "follower", false)
		if appErr != nil {
			p.API.LogError("Failed to find Follower channel", "error", appErr.Error())
			return
		}

		// Subscribe the user to the "Follower" channel
		_, appErr = p.API.AddChannelMember(followerChannel.Id, user.Id)
		if appErr != nil {
			p.API.LogError("Failed to add user to Follower channel", "user_id", user.Id, "channel_id", followerChannel.Id, "error", appErr.Error())
			return
		}

		p.API.LogInfo("User successfully subscribed to the Follower channel", "user_id", user.Id, "channel_id", followerChannel.Id)

	}
}
