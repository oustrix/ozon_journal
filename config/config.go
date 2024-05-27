package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config stores all application settings.
	Config struct {
		Environment Environment `yaml:"environment"`
		Storage     Storage     `yaml:"storage"`
		Postgres    Postgres    `yaml:"postgres"`
		Log         Log         `yaml:"log"`
		HTTP        HTTP        `yaml:"http"`
		Comment     Comment     `yaml:"comment"`
		Post        Post        `yaml:"post"`
	}

	// Environment contains settings for application environment.
	Environment string // valid values: "development", "production"

	// Storage containts settings for application data storage.
	Storage struct {
		Type string `yaml:"type" env:"STORAGE_TYPE" env-required:"true"` // valid values: "in-memory", "postgres"
	}

	Postgres struct {
		DSN          string `yaml:"dsn" env:"POSTGRES_DSN"`
		MaxPoolSize  uint   `yaml:"max_pool_size" env:"POSTGRES_MAX_POOL_SIZE"`
		ConnAttempts uint   `yaml:"conn_attempts" env:"POSTGRES_CONN_ATTEMPTS"`
		ConnTimeout  uint   `yaml:"conn_timeout" env:"POSTGRES_CONN_TIMEOUT"`
	}

	// Log contains settings for application logger.
	Log struct {
		Level string `yaml:"level" env:"LOG_LEVEL" env-required:"true"`
	}

	// HTTP contains settings for HTTP server.
	HTTP struct {
		Port string `yaml:"port" env:"HTTP_PORT" env-required:"true"`
	}

	// Comment contains settings for comment service.
	Comment struct {
		MaxCharacters uint `yaml:"max_characters" env:"COMMENT_AX_CHARACTERS" env-required:"true"`
		DefaultPage   uint `yaml:"default_page" env:"COMMENT_DEFAULT_PAGE" env-required:"true"`
		DefaultAmount uint `yaml:"default_amount" env:"COMMENT_DEFAULT_AMOUNT" env-required:"true"`
	}

	// Post contains settings for post service.
	Post struct {
		TitleMaxCharacters   uint `yaml:"title_max_characters" env:"POST_TITLE_MAX_CHARACTERS" env-required:"true"`
		ContentMaxCharacters uint `yaml:"content_max_characters" env:"POST_CONTENT_MAX_CHARACTERS" env-required:"true"`
		DefaultPage          uint `yaml:"default_page" env:"POST_DEFAULT_PAGE" env-required:"true"`
		DefaultAmount        uint `yaml:"default_amount" env:"POST_DEFAULT_AMOUNT" env-required:"true"`
	}
)

// NewConfig creates a new Config instance and reads the configuration from config/config.yml file.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	// Read the configuration from the file and environment variables.
	err := cleanenv.ReadConfig("config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("NewConfig - ReadConfig: %w", err)
	}

	// If the storage type is postgres, the DSN must be set.
	if cfg.Postgres.DSN == "" && cfg.Storage.Type == "postgres" {
		return nil, fmt.Errorf("NewConfig - DSN is empty")
	}

	return cfg, nil
}
