package gormx

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var drivers = map[string]func(*Config) (gorm.Dialector, error){}

func Register(name string, driverFactory func(*Config) (gorm.Dialector, error)) {
	drivers[name] = driverFactory
}

type Config struct {
	DSN         string `json:"dsn,omitempty"`
	Debug       bool   `json:"debug,omitempty"`
	TablePrefix string `json:"table_prefix,omitempty"`
}

func Open(cfg *Config) (*gorm.DB, error) {
	driver, _, _ := strings.Cut(cfg.DSN, "://")

	factory, ok := drivers[driver]
	if !ok {
		return nil, fmt.Errorf("driver %s not found", driver)
	}

	dialector, err := factory(cfg)
	if err != nil {
		return nil, err
	}

	level := logger.Error
	if cfg.Debug {
		level = logger.Info
	}

	return gorm.Open(dialector, &gorm.Config{
		DisableAutomaticPing: true,
		NamingStrategy:       &schema.NamingStrategy{TablePrefix: cfg.TablePrefix},
		Logger:               logger.Default.LogMode(level),
	})
}
