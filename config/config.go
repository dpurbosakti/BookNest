package config

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type HttpConf struct {
	Host string
	Port string
}

type EmailConf struct {
	Email    string
	Password string
	Host     string
	Port     int
}

type JwtConf struct {
	SecretKey         string
	ExpiredTimeInHour uint
}

type Config struct {
	DbConf       DbConf
	HttpConf     HttpConf
	EmailConf    EmailConf
	LoggerConf   *logrus.Logger
	GoogleConf   GoogleConf
	TwitterConf  TwitterConf
	GithubConf   GithubConf
	JwtConf      JwtConf
	MidtransConf MidtransConf
	BiteshipConf BiteshipConf
	XenditConf   XenditConf
}

var Cfg *Config

func GetConfig() {
	if Cfg == nil {
		viper.SetConfigFile("config.yml")
		// default values
		viper.SetDefault("FullName", "mokotest")
		viper.SetDefault("Version", "0.0.1")
		viper.SetDefault("HttpConf.Host", "127.0.0.1")
		viper.SetDefault("HttpConf.Port", "8000")

		// read the file
		if err := viper.ReadInConfig(); err != nil {
			fmt.Printf("Error reading config file, %s", err)
			panic(err)
		}

		// map to app
		if err := viper.Unmarshal(&Cfg); err != nil {
			fmt.Printf("Unable to decode into struct, %v", err)
			panic(err)

		}

		// done
		logrus.WithFields(logrus.Fields{
			"source":  "config",
			"status":  "done",
			"message": "config is loaded successfully",
		}).Info("loading config")
	}

}
