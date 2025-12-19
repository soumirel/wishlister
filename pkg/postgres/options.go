package postgres

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DbConfig struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
}

type options struct {
	db                    DbConfig
	maxConnLifetime       time.Duration
	maxConnLifetimeJitter time.Duration
	maxConnIdleTime       time.Duration
	maxConns              int32
	minConns              int32
	minIdleConns          int32
	healthCheckPeriod     time.Duration
}

var (
	defaultOptions = options{
		maxConnLifetime:       time.Hour,
		maxConnLifetimeJitter: 0,
		maxConnIdleTime:       30 * time.Minute,
		maxConns:              4,
		minConns:              0,
		minIdleConns:          0,
		healthCheckPeriod:     time.Minute,
	}
)

func (o *options) toPgxPoolConfig() (*pgxpool.Config, error) {
	connStr := fmt.Sprintf(`host=%s port=%s dbname=%s user=%s password=%s sslmode=disable`,
		o.db.Host, o.db.Port, o.db.Database, o.db.User, o.db.Password,
	)
	c, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}
	c.MaxConnLifetime = o.maxConnLifetime
	c.MaxConnLifetimeJitter = o.maxConnLifetimeJitter
	c.MaxConnIdleTime = o.maxConnIdleTime
	c.MaxConns = o.maxConns
	c.MinConns = o.minConns
	c.MinIdleConns = o.minIdleConns
	c.HealthCheckPeriod = o.healthCheckPeriod
	return c, nil
}

type Option interface {
	apply(*options) error
}

type optionFunc func(*options) error

func (of optionFunc) apply(opts *options) error {
	return of(opts)
}

func WithMaxConnLifetime(maxLifetime time.Duration) Option {
	return optionFunc(func(opts *options) error {
		opts.maxConnLifetime = maxLifetime
		return nil
	})
}

func WithConnLifetimeJitter(jitter time.Duration) Option {
	return optionFunc(func(opts *options) error {
		opts.maxConnLifetimeJitter = jitter
		return nil
	})
}

func WithMaxConnIdleTime(idleTime time.Duration) Option {
	return optionFunc(func(opts *options) error {
		opts.maxConnIdleTime = idleTime
		return nil
	})
}

func WithMaxConns(maxConns int32) Option {
	return optionFunc(func(opts *options) error {
		opts.maxConns = maxConns
		return nil
	})
}

func WithMinConns(minConns int32) Option {
	return optionFunc(func(opts *options) error {
		opts.minConns = minConns
		return nil
	})
}

func WithMinIdleConns(minIdleConns int32) Option {
	return optionFunc(func(opts *options) error {
		opts.minIdleConns = minIdleConns
		return nil
	})
}

func WithHealthCheckPeriod(healthCheckPeriod time.Duration) Option {
	return optionFunc(func(opts *options) error {
		opts.healthCheckPeriod = healthCheckPeriod
		return nil
	})
}
