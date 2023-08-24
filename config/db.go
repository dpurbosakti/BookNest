package config

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DbConf struct {
	Host         string
	Port         string
	User         string
	Password     string
	DataBaseName string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
	Dialect      string
}

func InitDb(c *Config) (DB *gorm.DB) {
	logger := logrus.WithField("configuring db", "initiate db")
	if c.DbConf.Dialect != "mysql" && c.DbConf.Dialect != "postgres" {
		logger.WithFields(logrus.Fields{
			"type":    "db",
			"source":  "gorm",
			"status":  "unset",
			"message": "no proper dialect provided",
		}).Info("instantiation")
		return
	}
	var gormD gorm.Dialector
	if c.DbConf.Dialect == "mysql" {
		gormD = mysql.Open(mysqlDsnBuilder(c.DbConf))
		logger.WithField("dsn", mysqlDsnBuilder(c.DbConf)).Info()
	} else if c.DbConf.Dialect == "postgres" {
		gormD = postgres.Open(postgresDsnBuilder(c.DbConf))
	}

	db, err := gorm.Open(gormD, &gorm.Config{})
	if err != nil {
		logger.WithFields(logrus.Fields{
			"type":    "db",
			"source":  "gorm",
			"status":  "panic",
			"message": "Failed to connect to database!",
		}).Error("instantiation")
		logger.Panic(err)
	} else {
		DB = db
		logger.WithFields(logrus.Fields{
			"type":   "db",
			"source": "gorm",
			"status": "done",
		}).Info("instantiation")
	}

	return DB
}

func mysqlDsnBuilder(c DbConf) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", c.User, c.Password, c.Host, c.Port, c.DataBaseName)
}

func postgresDsnBuilder(c DbConf) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Jakarta", c.Host, c.Port, c.User, c.Password, c.DataBaseName)
}
