package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/soumirel/wishlister/services/wishlist/pkg/postgres"

	"github.com/spf13/viper"
)

const envConfigPath = "CONFIG_PATH"
const defaultConfigPath = "./config.yaml"

// Config holds all application configuration.
// Priority: env > config file > defaults.
type Config struct {
	Server Server `mapstructure:"server"`
	DB     DB     `mapstructure:"db"`
	Paths  Paths  `mapstructure:"paths"`
}

type Server struct {
	HTTPAddr string `mapstructure:"http_addr"`
	GRPCAddr string `mapstructure:"grpc_addr"`
}

type DB struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Database string `mapstructure:"database"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type Paths struct {
	Migrations  string `mapstructure:"migrations"`
	RefreshData string `mapstructure:"refresh_data"`
}

// Load reads configuration from defaults, optional config file, and env.
// Env vars override file; file overrides defaults. Use WISHLIST_* for env (e.g. WISHLIST_DB_PASSWORD).
func Load() (*Config, error) {
	v := viper.New()

	// Defaults
	v.SetDefault("server.http_addr", ":8080")
	v.SetDefault("server.grpc_addr", ":8081")
	v.SetDefault("db.host", "localhost")
	v.SetDefault("db.port", "5432")
	v.SetDefault("db.database", "wishlister")
	v.SetDefault("db.user", "wishlister")
	v.SetDefault("db.password", "")
	v.SetDefault("paths.migrations", "./init.sql")
	v.SetDefault("paths.refresh_data", "./refresh_test_data.sql")

	// Optional config file
	configPath := os.Getenv(envConfigPath)
	if configPath == "" {
		configPath = defaultConfigPath
	}
	v.SetConfigFile(configPath)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("read config: %w", err)
		}
		// Config file not found; continue with defaults and env
	}

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetEnvPrefix("WISHLIST")
	v.AutomaticEnv()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}
	return &cfg, nil
}

// DbConfig returns postgres.DbConfig from Config.DB.
func (c *Config) DbConfig() postgres.DbConfig {
	return postgres.DbConfig{
		Host:     c.DB.Host,
		Port:     c.DB.Port,
		Database: c.DB.Database,
		User:     c.DB.User,
		Password: c.DB.Password,
	}
}
