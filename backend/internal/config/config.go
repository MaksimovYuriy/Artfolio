package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTP        HTTPConfig        `env-prefix:"HTTP_"`
	DB          DBConfig          `env-prefix:"DB_"`
	App         AppConfig         `env-prefix:"APP_"`
	FileStorage FileStorageConfig `env-prefix:"STORAGE_"`
	Kafka       KafkaConfig       `env-prefix:"KAFKA_"`
}

type HTTPConfig struct {
	Port              string        `env:"PORT" env-default:"8081"`
	Address           string        `env:"ADDRESS" env-default:"0.0.0.0"`
	ReadTimeout       time.Duration `env:"READ_TIMEOUT" env-default:"60s"`
	WriteTimeout      time.Duration `env:"WRITE_TIMEOUT" env-default:"60s"`
	ReadHeaderTimeout time.Duration `env:"READ_HEADER_TIMEOUT" env-default:"5s"`
	IdleTimeout       time.Duration `env:"IDLE_TIMEOUT" env-default:"60s"`
}

type DBConfig struct {
	Host     string `env:"HOST" env-default:"localhost"`
	Port     string `env:"PORT" env-default:"5432"`
	User     string `env:"USER" env-default:"artfolio"`
	Password string `env:"PASSWORD"`
	Name     string `env:"NAME" env-default:"artfolio"`
	SSLMode  string `env:"SSL_MODE" env-default:"disable"`
}

type AppConfig struct {
	Env string `env:"ENV" env-default:"prod"`
}

type FileStorageConfig struct {
	Path        string `env:"PATH" env-default:"./media"`
	PublicURL   string `env:"PUBLIC_URL" env-default:"/media"`
	MaxFileSize int64  `env:"MAX_FILE_SIZE" env-default:"12582912"`
	MaxPixels   int64  `env:"MAX_PIXELS" env-default:"25000000"`
}

type KafkaConfig struct {
	Brokers []string `env:"BROKERS" env-default:"localhost:9092" env-separator:","`
	Topics  KafkaTopicsConfig
}

type KafkaTopicsConfig struct {
	ContactMessageSubmitted string `env:"CONTACT_MESSAGE_SUBMITTED_TOPIC" env-default:"contact-message.submitted.v1"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
