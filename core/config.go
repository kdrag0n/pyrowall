package core

import (
	"fmt"
	"io/ioutil"

	"github.com/pelletier/go-toml"
)

// TelegramConfig holds the Telegram-related configuration data for a bot instance.
type TelegramConfig struct {
	Token string `toml:"token"`

	UseWebhooks            bool   `toml:"use_webhooks"`
	WebhookAddress         string `toml:"webhook_address"`
	WebhookPort            int    `toml:"webhook_port"`
	WebhookPath            string `toml:"webhook_path"`
	WebhookURL             string `toml:"webhook_url"`
	WebhookConnectionLimit int    `toml:"webhook_limit"`
}

// PprofConfig holds the pprof-related configuration data for a bot instance.
type PprofConfig struct {
	EnableServer  bool   `toml:"enable_server"`
	ServerAddress string `toml:"server_address"`
	ServerPort    int    `toml:"server_port"`
}

// LoggingConfig holds the logging-related configuration data for a bot instance.
type LoggingConfig struct {
	Enable bool   `toml:"enable"`
	Format string `toml:"format"`
	Level  string `toml:"level"`
}

// SentryConfig holds the Sentry-related configuration data for a bot instance.
type SentryConfig struct {
	Enable      bool   `toml:"enable"`
	DSN         string `toml:"dsn"`
	ServerName  string `toml:"server_name"`
	Release     string `toml:"release"`
	Environment string `toml:"environment"`
}

// DatabaseConfig holds the database-related configuration data for a bot instance.
type DatabaseConfig struct {
	Type     string `toml:"type"`
	Protocol string `toml:"protocol"`
	Address  string `toml:"address"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Database string `toml:"database"`
}

// Config holds all the configuration data for a bot instance.
type Config struct {
	Telegram TelegramConfig `toml:"telegram"`
	Pprof    PprofConfig    `toml:"pprof"`
	Logging  LoggingConfig  `toml:"logging"`
	Sentry   SentryConfig   `toml:"sentry"`
	Database DatabaseConfig `toml:"database"`
}

// ParseConfig parses the given data into a Config.
func ParseConfig(data []byte) (*Config, error) {
	var config Config
	err := toml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	if config.Telegram.WebhookPath == "" {
		config.Telegram.WebhookPath = config.Telegram.Token
	}

	return &config, nil
}

// ReadConfigFile reads a Config from the given file.
func ReadConfigFile(path string) (config *Config, err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("read %s: %w", path, err)
		return
	}

	config, err = ParseConfig(data)
	return
}
