package keygen

import (
	"context"
	"fmt"

	"github.com/maksimovyuriy/artfolio/backend/internal/config"
	"github.com/maksimovyuriy/artfolio/backend/internal/lib/logger"
	keyrepo "github.com/maksimovyuriy/artfolio/backend/internal/repo/key"
	"github.com/maksimovyuriy/artfolio/backend/internal/storage/postgres"
	keyusecase "github.com/maksimovyuriy/artfolio/backend/internal/usecase/key"
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

	repository := keyrepo.NewRepo(database)
	useCase := keyusecase.NewUseCase(repository)

	accessKey, err := useCase.Create(context.Background())
	if err != nil {
		return err
	}
	log.Info("Admin key created")
	fmt.Println("Admin key generated. Save it now:")
	fmt.Println(accessKey)

	return nil
}
