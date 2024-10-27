package config

var CategoryOrder = []string{"Club Life", "Racing", "Cruising", "Fleet", "Training"}

var PublicChannels = map[string][]string{
	"Club Life": {"Town Square", "Club News", "Club House", "Crew Finder", "Market Place", "Car Pool", "Off-Topic"},
	"Racing":    {"Monday Races", "Seven Bars", "Kaag Cup", "ESA Cup", "Arianes Cup", "Other Races"},
	"Cruising":  {"Cruising"},
	"Fleet":     {"Wayfarer", "Randmeer", "Venture", "Laser", "Buzz", "Fox", "Safety Boat", "Booking"},
	"Training":  {"Sign Up"},
}

var PrivateChannels = map[string][]string{
	"Club Life": {"Committee"},
	"Racing":    {},
	"Cruising":  {},
	"Fleet":     {"Fox maintenance and management"},
	"Training":  {"Instructors", "Training 2024 B", "Training 2024 A", "Training 2023 B"},
}

var DefaultCategories = []string{"Favorites", "Channels", "Direct Messages"} // cannot delete them

var (
	ServerURL = "http://localhost:8065"
	Headers   = map[string]string{
		"Authorization": "Bearer " + AuthConfig.AuthToken,
		"Content-Type":  "application/json",
	}
)
