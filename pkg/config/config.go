package config

import (
	"log"

	"github.com/alexedwards/scs/v2"
)

// This package stores system wide configurations.

type AppConfig struct {
	InfoLog   *log.Logger
	Session   *scs.SessionManager
	CSRFToken string
}

// CSRF: Cross site forgery requests occur whenever malicious code is used to trigger unwanted actions after the user is already authenticated. With each CSRFToken, each request is going to be verified before execution. And we also use NoSurf which is a CSRF protection middleware.
