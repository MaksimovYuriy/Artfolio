package migrator

import (
	"github.com/maksimovyuriy/artfolio/backend/internal/config"
	"github.com/maksimovyuriy/artfolio/backend/internal/lib/logger"
	"github.com/maksimovyuriy/artfolio/backend/internal/storage/postgres"
	"github.com/pressly/goose"
)

func Run() error {
	appCfg, err := config.Load()
	if err != nil {
		return err
	}

	log := logger.New(appCfg.App.Env)
	log.Info("Logger started")

	database, err := postgres.New(appCfg.DB)
	if err != nil {
		log.Error("Database connection failed")
		return err
	}
	defer database.Close()
	log.Info("Database connected")

	if err := goose.SetDialect("postgres"); err != nil {
		log.Error("Failed to set postgresql dialect")
		return err
	}

	if err := goose.Up(database, "./migrations"); err != nil {
		log.Error("Failed to apply migrations")
		return err
	}

	log.Info("Migrations applied")
	return nil
}
