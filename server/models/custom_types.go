package models

import (
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/plugin"
)

type RestAPI interface {
	Delete(path string) ([]byte, error)
	Get(path string) ([]byte, error)
	Post(path string, data interface{}) ([]byte, error)
	Put(path string, data interface{}) ([]byte, error)
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
