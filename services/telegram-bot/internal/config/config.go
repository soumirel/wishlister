package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

const (
	envConfigPath     = "CONFIG_PATH"
	defaultConfigPath = "./config.yaml"
)

type Config struct {
	Server      Server      `mapstructure:"server"`
	Services    Services    `mapstructure:"services"`
	TelegramBot TelegramBot `mapstructure:"telegram_bot"`
}

type Server struct {
	HTTPAddr string `mapstructure:"http_addr"`
}

type Services struct {
	WishlistGrpcAddr string `mapstructure:"wishlist_grpc_addr"`
}

type TelegramBot struct {
	Token string
}

// Load reads configuration from defaults, optional config file, and env.
// Env vars override file; file overrides defaults. Use WISHLIST_* for env (e.g. WISHLIST_DB_PASSWORD).
func Load() (*Config, error) {
	v := viper.New()

	// Defaults
	v.SetDefault("server.http_addr", ":8080")
	v.SetDefault("services.wishlist_grpc_addr", "localhost:8081")
	v.SetDefault("telegram_bot.token", "")

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

	// viper env settings
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetEnvPrefix("BOT")
	v.AutomaticEnv()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}
	return &cfg, nil
}
