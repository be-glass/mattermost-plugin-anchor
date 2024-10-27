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

func AllChannels() map[string][]string {
	merged := make(map[string][]string)

	// First, add channels from PublicChannels in order
	for category, channels := range PublicChannels {
		merged[category] = append(merged[category], channels...)
		// Add channels from PrivateChannels if the category exists
		if privateChannels, exists := PrivateChannels[category]; exists {
			merged[category] = append(merged[category], privateChannels...)
		}
	}

	// Next, add any categories from PrivateChannels that are not in PublicChannels
	for category, channels := range PrivateChannels {
		if _, exists := merged[category]; !exists {
			merged[category] = append(merged[category], channels...)
		}
	}

	return merged
}
