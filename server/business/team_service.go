package business

import (
	"github.com/glass.plugin-anchor/server/models"
	"github.com/mattermost/mattermost-server/v6/model"
	"strings"
)

//func GetTeamIDByName(c *models.Context, teamName string) (string, error) {
//	team, appErr := c.API.GetTeamByName(teamName)
//	if appErr != nil {
//		return "", appErr
//	}
//	return team.Id, nil
//}

func GetTeamsListString(c *models.Context) string {

	c.API.LogWarn("1")

	teams, err := listAllTeams(c)

	c.API.LogWarn("2")

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

func listAllTeams(c *models.Context) ([]*model.Team, error) {

	// Fetch teams for the current page
	teams, appErr := c.API.GetTeams()
	if appErr != nil {
		return nil, appErr
	}

	return teams, nil
}
