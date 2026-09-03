package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/maksimovyuriy/artfolio/mail-service/internal/config"
	"github.com/maksimovyuriy/artfolio/mail-service/internal/handlers"
	"github.com/maksimovyuriy/artfolio/mail-service/internal/kafka"
)

func Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(log)

	consumer, err := kafka.NewConsumer(cfg.Kafka)
	if err != nil {
		return err
	}
	defer consumer.Close()

	emailHandler := handlers.NewEmail(log)
	log.Info("Kafka consumer started",
		slog.String("topic", cfg.Kafka.Topic),
		slog.String("group_id", cfg.Kafka.GroupID),
	)

	err = consumer.Consume(ctx, emailHandler.Handle)
	if ctx.Err() != nil {
		log.Info("Kafka consumer stopped")
		return nil
	}
	return err
}
