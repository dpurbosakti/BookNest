package app

import (
	"book-nest/config"
	"book-nest/internal/http"
	"book-nest/migration"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

func init() {
	cobra.OnInitialize(initProject)

	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(serveCmd)

}

func initProject() {
	config.GetConfig()
	DB = config.InitDb()
	server = http.Serve
}

var (
	rootCmd = &cobra.Command{
		Use: "book-nest",
	}
	migrateCmd = &cobra.Command{
		Use: "migrate",
		Run: migrate,
	}

	serveCmd = &cobra.Command{
		Use: "serve",
		Run: serve,
	}

	// global variable
	DB     *gorm.DB
	Config *config.Config
	server http.Server
)

func migrate(cmd *cobra.Command, args []string) {
	logger := logrus.WithField("func", "migrate")
	logger.Info("start migration")
	err := migration.Migrate(DB)
	if err != nil {
		logger.WithError(err).Error("error where running migration")
		panic(err)
	}
	logger.Info("done")
}

func serve(cmd *cobra.Command, args []string) {
	logger := logrus.WithField("func", "serve")
	logger.Info("serve")
	server(DB, config.Cfg.HttpConf.Port)
	logger.Info("done")
}

func Execute() error {
	return rootCmd.Execute()
}
