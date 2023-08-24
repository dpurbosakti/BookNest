package config

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type HttpConf struct {
	Host string
	Port int
}

type EmailConf struct {
	Email    string
	Password string
	Host     string
	Port     int
}

type Config struct {
	DbConf     DbConf
	HttpConf   HttpConf
	EmailConf  EmailConf
	LoggerConf *logrus.Logger
}

func GetConfig() Config {
	var cfg *Config
	// main viper config
	// viper.SetConfigName("config")
	// viper.SetConfigType("yml")
	// viper.AddConfigPath(".")
	// viper.AutomaticEnv()
	viper.SetConfigFile("D:/Belajar/BE/Moko/BookNest/config.yml") // windows
	// viper.SetConfigFile("/mnt/d/Belajar/BE/learnEcho2/config.yml") // linux
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
	if err := viper.Unmarshal(&cfg); err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
		panic(err)

	}

	// done
	logrus.WithFields(logrus.Fields{
		"source":  "config",
		"status":  "done",
		"message": "config is loaded successfully",
	}).Info("loading config")

	return *cfg
}
