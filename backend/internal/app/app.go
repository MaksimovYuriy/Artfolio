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
	v1 "github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/v1"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/v1/response"
	kafkaproducer "github.com/maksimovyuriy/artfolio/backend/internal/kafka"
	contactkafka "github.com/maksimovyuriy/artfolio/backend/internal/kafka/contact"
	"github.com/maksimovyuriy/artfolio/backend/internal/lib/filestorage"
	artworkstorage "github.com/maksimovyuriy/artfolio/backend/internal/lib/filestorage/artwork"
	"github.com/maksimovyuriy/artfolio/backend/internal/lib/logger"
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
	contactusecase "github.com/maksimovyuriy/artfolio/backend/internal/usecase/contact"
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
	contact       usecase.ContactUseCase
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

	artworkStorage, err := artworkstorage.New(appCfg.FileStorage)
	if err != nil {
		return err
	}
	log.Info("Artwork storage started")

	producer, err := kafkaproducer.NewProducer(appCfg.Kafka)
	if err != nil {
		return err
	}
	defer producer.Close()
	log.Info("Kafka producer started")
	contactPublisher, err := contactkafka.NewPublisher(producer, appCfg.Kafka.Topics.ContactMessageSubmitted)
	if err != nil {
		return err
	}

	repositories := initRepositories(database)
	useCases := initUseCases(repositories, artworkStorage, contactPublisher, log)
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

func initUseCases(
	repositories repositories, artworkStorage filestorage.Artwork,
	contactPublisher *contactkafka.Publisher, log *slog.Logger,
) useCases {
	return useCases{
		session:       sessionusecase.NewUseCase(repositories.key, repositories.session),
		artistProfile: artistusecase.NewUseCase(repositories.artistProfile, repositories.socialLink),
		artwork:       artworkusecase.NewUseCase(repositories.artwork, artworkStorage, log),
		socialLink:    sociallinkusecase.NewUseCase(repositories.artistProfile, repositories.socialLink),
		contact:       contactusecase.NewUseCase(repositories.artistProfile, contactPublisher),
	}
}

func initServers(appCfg *config.Config, useCases useCases, log *slog.Logger) servers {
	artworkMapper := response.NewArtworkMapper(appCfg.FileStorage.PublicURL)
	controller := v1.NewController(
		useCases.session,
		useCases.artistProfile,
		useCases.artwork,
		useCases.socialLink,
		useCases.contact,
		artworkMapper,
	)
	authMiddleware := middleware.NewAuth(useCases.session)
	maxUploadBodySize := appCfg.FileStorage.MaxFileSize + multipartFormOverheadAllowance
	router := restapi.NewRouter(controller, authMiddleware, maxUploadBodySize, log)

	return servers{api: restapi.NewServer(appCfg.HTTP, router, log)}
}
