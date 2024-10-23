package models

import (
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/plugin"
)

type RestAPI struct {
	ServerURL string
	AuthToken string
	Headers   map[string]string
}

type Context struct {
	Team    *model.Team
	Channel *model.Channel
	User    *model.User

	API  plugin.API
	Rest RestAPI
}

type Auth struct {
	AuthToken   string
	SQLPassword string
}
