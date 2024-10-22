package models

import (
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
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
