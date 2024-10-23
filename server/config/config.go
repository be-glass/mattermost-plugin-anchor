package config

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

var (
	BotToken  = "bg68g5ddytgump6xehbdt4c6nw"
	ServerURL = "http://localhost:8065"
	AuthToken = "your-auth-token"
	Headers   = map[string]string{
		"Authorization": "Bearer your-auth-token",
		"Content-Type":  "application/json",
	}
)
