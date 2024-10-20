package business

import (
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"strings"
)

func GetTeamIDByName(api plugin.API, teamName string) (string, error) {
	team, appErr := api.GetTeamByName(teamName)
	if appErr != nil {
		return "", appErr
	}
	return team.Id, nil
}

func GetTeamsListString(api plugin.API) string {

	api.LogWarn("1")

	teams, err := listAllTeams(api)

	api.LogWarn("2")

	if err != nil {
		return "Don't know!"
	} else if len(teams) == 0 {
		return "No teams"
	} else {
		var teamNames []string
		for _, team := range teams {
			teamNames = append(teamNames, team.Name)
		}
		return strings.Join(teamNames, "\n")
	}
}

// private

func listAllTeams(api plugin.API) ([]*model.Team, error) {

	// Fetch teams for the current page
	teams, appErr := api.GetTeams()
	if appErr != nil {
		return nil, appErr
	}

	return teams, nil
}
