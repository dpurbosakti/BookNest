package app

import (
	"book-nest/config"
	"book-nest/migration"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

func init() {
	cobra.OnInitialize(initProject)

	rootCmd.AddCommand(migrateCmd)

}

func initProject() {
	Config = config.GetConfig()
	DB = config.InitDb(&Config)
}

var (
	rootCmd = &cobra.Command{
		Use: "book-nest",
	}
	migrateCmd = &cobra.Command{
		Use: "migrate",
		Run: migrate,
	}

	// global variable
	DB     *gorm.DB
	Config config.Config
)

func migrate(cmd *cobra.Command, args []string) {
	logger := logrus.WithField("func", "migrate")
	logger.Info("start migration")
	migration.Migrate(DB)
	// err := migrator()
	// if err != nil {
	// 	logger.WithError(err).Error("error where running migration")
	// 	panic(err)
	// }
	logger.Info("done")
}

func Execute() error {
	return rootCmd.Execute()
}
