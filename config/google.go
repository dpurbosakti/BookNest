package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleConf struct {
	ClientID       string
	ClientSecret   string
	RedirectUrl    string
	State          string
	TokenAccessUrl string
}

func GetGoogleConfig() *oauth2.Config {
	conf := &oauth2.Config{
		ClientID:     Cfg.GoogleConf.ClientID,
		ClientSecret: Cfg.GoogleConf.ClientSecret,
		RedirectURL:  Cfg.GoogleConf.RedirectUrl,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/calendar",
		},
		Endpoint: google.Endpoint,
	}

	return conf
}

func GetGoogleConfigCal() *oauth2.Config {
	conf := &oauth2.Config{
		ClientID:     Cfg.GoogleConf.ClientID,
		ClientSecret: Cfg.GoogleConf.ClientSecret,
		RedirectURL:  Cfg.GoogleConf.RedirectUrl,
		Scopes: []string{
			"https://www.googleapis.com/auth/calendar",
		},
		Endpoint: google.Endpoint,
	}

	return conf
}
