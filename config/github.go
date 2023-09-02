package config

import (
	"golang.org/x/oauth2"
)

type GithubConf struct {
	ClientID       string
	ClientSecret   string
	RedirectUrl    string
	TokenAccessUrl string
}

func GetGithubConfig() *oauth2.Config {
	conf := &oauth2.Config{
		ClientID:     Cfg.GithubConf.ClientID,
		ClientSecret: Cfg.GithubConf.ClientSecret,
		RedirectURL:  Cfg.GithubConf.RedirectUrl,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}
	return conf
}
