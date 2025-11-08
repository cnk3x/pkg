//go:build cgo

package sqlite

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/cnk3x/gopkg/gormx"
)

func Open(cfg *gormx.Config) (gorm.Dialector, error) {
	return sqlite.Open(cfg.DSN), nil
}

func init() {
	gormx.Register("sqlite", Open)
}
