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

	rootCmd.AddCommand(migrateUpCmd)
	rootCmd.AddCommand(migrateDownCmd)
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
	migrateUpCmd = &cobra.Command{
		Use: "migrate_up",
		Run: migrateUp,
	}
	migrateDownCmd = &cobra.Command{
		Use: "migrate_down",
		Run: migrateDown,
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

func migrateUp(cmd *cobra.Command, args []string) {
	logger := logrus.WithField("func", "migrate_up")
	logger.Info("start migration")
	err := migration.MigrateUp(DB)
	if err != nil {
		logger.WithError(err).Error("error when running migration")
		panic(err)
	}
	logger.Info("done")
}

func migrateDown(cmd *cobra.Command, args []string) {
	logger := logrus.WithField("func", "migrate_down")
	logger.Info("start migration")
	err := migration.MigrateDown(DB)
	if err != nil {
		logger.WithError(err).Error("error when running migration")
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
