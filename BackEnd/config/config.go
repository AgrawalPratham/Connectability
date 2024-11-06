package config

import (
	"github.com/agrawalpratham/Connectability/BackEnd/database"
	"github.com/gorilla/sessions"
)

type Application struct {
	PortNumber   string
	DBConnection *database.DB
	Session *sessions.CookieStore
	UserEmail string
}

var App = &Application{}