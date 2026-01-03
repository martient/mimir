package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig    `yaml:"server"`
	Database DatabaseConfig  `yaml:"database"`
	opencode opencodeConfig  `yaml:"opencode"`
	Projects []ProjectConfig `yaml:"projects"`
	Webhooks WebhooksConfig  `yaml:"webhooks"`
	Cron     []CronJob       `yaml:"cron"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	HTTPPort int `yaml:"http_port"`
	WSPort   int `yaml:"ws_port"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Path          string `yaml:"path"`
	EncryptionKey string `yaml:"encryption_key"`
}

// opencodeConfig holds opencode instance configuration
type opencodeConfig struct {
	DefaultModel string `yaml:"default_model"`
}

// ProjectConfig holds project configuration
type ProjectConfig struct {
	Name         string `yaml:"name"`
	Path         string `yaml:"path"`
	opencodePort int    `yaml:"opencode_port"`
}

// WebhooksConfig holds webhook configuration
type WebhooksConfig struct {
	Sentry SentryWebhookConfig `yaml:"sentry"`
}

// SentryWebhookConfig holds Sentry webhook configuration
type SentryWebhookConfig struct {
	Secret   string             `yaml:"secret"`
	Projects []SentryProjectMap `yaml:"projects"`
}

// SentryProjectMap maps Sentry projects to Mimir projects
type SentryProjectMap struct {
	SentryProject string `yaml:"sentry_project"`
	MimirProject  string `yaml:"mimir_project"`
}

// CronJob holds cron job configuration
type CronJob struct {
	Name     string `yaml:"name"`
	Schedule string `yaml:"schedule"`
	Project  string `yaml:"project"`
	Action   string `yaml:"action"`
}

// Load loads configuration from file or returns defaults
func Load() (*Config, error) {
	// Default config path
	configPath := os.Getenv("MIMIR_CONFIG")
	if configPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		configPath = fmt.Sprintf("%s/.mimir/config.yaml", homeDir)
	}

	// Ensure config directory exists
	configDir := os.Getenv("MIMIR_CONFIG_DIR")
	if configDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		configDir = fmt.Sprintf("%s/.mimir", homeDir)
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return default config if file doesn't exist
			return DefaultConfig(), nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse YAML
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Process environment variables
	if config.Database.EncryptionKey == "" {
		config.Database.EncryptionKey = os.Getenv("MIMIR_DB_KEY")
	}

	if config.Server.HTTPPort == 0 {
		config.Server.HTTPPort = 8080
	}

	if config.Server.WSPort == 0 {
		config.Server.WSPort = 8081
	}

	if config.opencode.DefaultModel == "" {
		config.opencode.DefaultModel = "anthropic/claude-sonnet-4"
	}

	return &config, nil
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			HTTPPort: 8080,
			WSPort:   8081,
		},
		Database: DatabaseConfig{
			Path: "~/.mimir/mimir.db",
		},
		opencode: opencodeConfig{
			DefaultModel: "anthropic/claude-sonnet-4",
		},
	}
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Server.HTTPPort <= 0 || c.Server.HTTPPort > 65535 {
		return fmt.Errorf("invalid HTTP port: %d", c.Server.HTTPPort)
	}

	if c.Server.WSPort <= 0 || c.Server.WSPort > 65535 {
		return fmt.Errorf("invalid WebSocket port: %d", c.Server.WSPort)
	}

	if c.Database.Path == "" {
		return fmt.Errorf("database path is required")
	}

	for _, project := range c.Projects {
		if project.Name == "" {
			return fmt.Errorf("project name is required")
		}
		if project.Path == "" {
			return fmt.Errorf("project path is required")
		}
	}

	return nil
}
