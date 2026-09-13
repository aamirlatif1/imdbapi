package config

import (
	"fmt"
	"log/slog"
	"net/url"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	HTTPPort string `env:"HTTP_PORT" envDefault:"4000"`
	Env      string `env:"ENV" envDefault:"development"`
	Database Database
}

type Database struct {
	Host     string `env:"DB_HOST" envDefault:"localhost"`
	Port     string `env:"DB_PORT" envDefault:"5432"`
	User     string `env:"DB_USER" envDefault:"imdbuser"`
	Password string `env:"DB_PASSWORD" envDefault:"imdbpassword"`
	Name     string `env:"DB_NAME" envDefault:"imdb"`
	SSLMode  string `env:"DB_SSLMODE" envDefault:"disable"`
}

func Load() (Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (cfg *Database) ConnectionString() string {
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Path:     cfg.Name,
		RawQuery: fmt.Sprintf("sslmode=%s", cfg.SSLMode),
	}
	return u.String()
}

func (cfg *Database) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("host", cfg.Host),
		slog.String("port", cfg.Port),
		slog.String("user", cfg.User),
		slog.String("password", "[REDACTED]"),
	)
}
