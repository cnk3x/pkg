package pg

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/cnk3x/pkg/gormx"
)

func Open(cfg *gormx.Config) (gorm.Dialector, error) {
	return postgres.Open(cfg.DSN), nil
}

func init() {
	gormx.Register("postgres", Open)
	gormx.Register("postgresql", Open)
	gormx.Register("pg", Open)
}
