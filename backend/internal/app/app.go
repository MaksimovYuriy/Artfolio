package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/maksimovyuriy/artfolio/backend/internal/config"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi"
	"github.com/maksimovyuriy/artfolio/backend/internal/lib/logger"
	"github.com/maksimovyuriy/artfolio/backend/internal/repo/key"
	"github.com/maksimovyuriy/artfolio/backend/internal/repo/session"
	"github.com/maksimovyuriy/artfolio/backend/internal/storage/postgres"

	sessionusecase "github.com/maksimovyuriy/artfolio/backend/internal/usecase/session"
)

func Run() error {
	appCtx, appCancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer appCancel()

	appCfg, err := config.Load()
	if err != nil {
		return err
	}

	log := logger.New(appCfg.App.Env)
	log.Info("Logger started")

	database, err := postgres.New(appCfg.DB)
	if err != nil {
		return err
	}
	defer database.Close()
	log.Info("Database started")

	keyRepo := key.NewRepo(database)
	sessionRepo := session.NewRepo(database)

	sessionUseCase := sessionusecase.NewUseCase(keyRepo, sessionRepo)

	sessionController := restapi.NewSessionController(sessionUseCase)

	router := restapi.NewRouter(sessionController)
	server := restapi.NewServer(appCfg.HTTP, router)

	apiErrors := make(chan error, 1)
	go func() {
		log.Info("API started", slog.String("addr", server.Addr))
		apiErrors <- server.ListenAndServe()
	}()

	select {
	case err := <-apiErrors:
		return err
	case <-appCtx.Done():
		log.Info("API shutting down")
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return err
	}
	log.Info("API stopped")

	return nil
}
