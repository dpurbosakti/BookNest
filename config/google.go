package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func GetGoogleConfig() *oauth2.Config {
	conf := &oauth2.Config{
		ClientID:     Cfg.GoogleConf.ClientID,
		ClientSecret: Cfg.GoogleConf.ClientSecret,
		RedirectURL:  Cfg.GoogleConf.RedirectUrl,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	return conf
}
