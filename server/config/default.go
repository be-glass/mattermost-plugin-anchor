package config

func ChannelNames() []string {
	var channels []string

	for _, channelList := range PublicChannels {
		channels = append(channels, channelList...)
	}
	return channels
}

func CategoryNames() []string {
	var categories []string

	for category := range PublicChannels {
		categories = append(categories, category)
	}
	return categories
}
