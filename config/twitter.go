package config

import (
	"golang.org/x/oauth2"
)

type TwitterConf struct {
	ClientID     string
	ClientSecret string
	RedirectUrl  string
	ApiEndpoint  string
}

func GetTwitterConfig() *oauth2.Config {
	conf := &oauth2.Config{
		ClientID:     Cfg.TwitterConf.ClientID,
		ClientSecret: Cfg.TwitterConf.ClientSecret,
		RedirectURL:  Cfg.TwitterConf.RedirectUrl,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://api.twitter.com/oauth/authenticate",
			TokenURL: "https://api.twitter.com/oauth/access_token",
		},
	}
	return conf
}
