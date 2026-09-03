package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Kafka  KafkaConfig  `env-prefix:"KAFKA_"`
	Resend ResendConfig `env-prefix:"RESEND_"`
}

type KafkaConfig struct {
	Brokers []string `env:"BROKERS" env-default:"localhost:9092" env-separator:","`
	Topic   string   `env:"CONTACT_MESSAGE_SUBMITTED_TOPIC" env-default:"contact-message.submitted.v1"`
	GroupID string   `env:"CONSUMER_GROUP" env-default:"mail-service"`
}

type ResendConfig struct {
	APIKey string `env:"API_KEY"`
	From   string `env:"FROM" env-default:"Artfolio <onboarding@resend.dev>"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
