package main

// Configuration holds all the settings that can be configured for the plugin.
type Configuration struct {
	SettingOne string
	SettingTwo int
	// Add more fields here based on your plugin's settings.
}

// SetConfiguration applies the given configuration to the plugin.
func (p *AnchorPlugin) SetConfiguration(configuration *Configuration) {
	p.configuration = configuration
}

// GetConfiguration retrieves the current configuration for the plugin.
func (p *AnchorPlugin) GetConfiguration() *Configuration {
	if p.configuration == nil {
		return &Configuration{}
	}

	return p.configuration
}
