package app

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/maksimovyuriy/artfolio/backend/internal/config"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/middleware"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/v1"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/v1/response"
	"github.com/maksimovyuriy/artfolio/backend/internal/lib/logger"
	"github.com/maksimovyuriy/artfolio/backend/internal/lib/storage"
	artworkstorage "github.com/maksimovyuriy/artfolio/backend/internal/lib/storage/artwork"
	"github.com/maksimovyuriy/artfolio/backend/internal/repo"
	artistprofile "github.com/maksimovyuriy/artfolio/backend/internal/repo/artist_profile"
	artworkrepo "github.com/maksimovyuriy/artfolio/backend/internal/repo/artwork"
	keyrepo "github.com/maksimovyuriy/artfolio/backend/internal/repo/key"
	sessionrepo "github.com/maksimovyuriy/artfolio/backend/internal/repo/session"
	sociallinkrepo "github.com/maksimovyuriy/artfolio/backend/internal/repo/social_link"
	"github.com/maksimovyuriy/artfolio/backend/internal/storage/postgres"

	"github.com/maksimovyuriy/artfolio/backend/internal/usecase"
	artistusecase "github.com/maksimovyuriy/artfolio/backend/internal/usecase/artist_profile"
	artworkusecase "github.com/maksimovyuriy/artfolio/backend/internal/usecase/artwork"
	sessionusecase "github.com/maksimovyuriy/artfolio/backend/internal/usecase/session"
	sociallinkusecase "github.com/maksimovyuriy/artfolio/backend/internal/usecase/social_link"
)

const multipartFormOverheadAllowance = 1 << 20

type repositories struct {
	key           repo.KeyRepository
	session       repo.SessionRepository
	artistProfile repo.ArtistProfileRepository
	artwork       repo.ArtworkRepository
	socialLink    repo.SocialLinkRepository
}

type useCases struct {
	session       usecase.SessionUseCase
	artistProfile usecase.ArtistProfileUseCase
	artwork       usecase.ArtworkUseCase
	socialLink    usecase.SocialLinkUseCase
}

type servers struct {
	api *http.Server
}

func Run() error {
	appCtx, appCancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer appCancel()

	appCfg, err := config.Load()
	if err != nil {
		return err
	}

	log := logger.New(appCfg.App.Env)
	slog.SetDefault(log)
	log.Info("Logger started")

	database, err := postgres.New(appCfg.DB)
	if err != nil {
		return err
	}
	defer database.Close()
	log.Info("Database started")

	artworkStorage, err := artworkstorage.New(appCfg.Storage)
	if err != nil {
		return err
	}
	log.Info("Artwork storage started")

	repositories := initRepositories(database)
	useCases := initUseCases(repositories, artworkStorage, log)
	servers := initServers(appCfg, useCases, log)

	apiErrors := make(chan error, 1)
	go func() {
		log.Info("API started", slog.String("addr", servers.api.Addr))
		apiErrors <- servers.api.ListenAndServe()
	}()

	select {
	case err := <-apiErrors:
		return err
	case <-appCtx.Done():
		log.Info("API shutting down")
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := servers.api.Shutdown(shutdownCtx); err != nil {
		return err
	}
	log.Info("API stopped")

	return nil
}

func initRepositories(database *sql.DB) repositories {
	return repositories{
		key:           keyrepo.NewRepo(database),
		session:       sessionrepo.NewRepo(database),
		artistProfile: artistprofile.NewRepo(database),
		artwork:       artworkrepo.NewRepo(database),
		socialLink:    sociallinkrepo.NewRepo(database),
	}
}

func initUseCases(repositories repositories, artworkStorage storage.Artwork, log *slog.Logger) useCases {
	return useCases{
		session:       sessionusecase.NewUseCase(repositories.key, repositories.session),
		artistProfile: artistusecase.NewUseCase(repositories.artistProfile, repositories.socialLink),
		artwork:       artworkusecase.NewUseCase(repositories.artwork, artworkStorage, log),
		socialLink:    sociallinkusecase.NewUseCase(repositories.artistProfile, repositories.socialLink),
	}
}

func initServers(appCfg *config.Config, useCases useCases, log *slog.Logger) servers {
	artworkMapper := response.NewArtworkMapper(appCfg.Storage.PublicURL)
	controller := v1.NewController(
		useCases.session,
		useCases.artistProfile,
		useCases.artwork,
		useCases.socialLink,
		artworkMapper,
	)
	authMiddleware := middleware.NewAuth(useCases.session)
	maxUploadBodySize := appCfg.Storage.MaxFileSize + multipartFormOverheadAllowance
	router := restapi.NewRouter(controller, authMiddleware, maxUploadBodySize, log)

	return servers{api: restapi.NewServer(appCfg.HTTP, router, log)}
}
