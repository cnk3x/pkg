//go:build !cgo

package sqlite

import (
	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/ncruces/go-sqlite3/gormlite"
	"gorm.io/gorm"

	"github.com/cnk3x/gopkg/gormx"
)

func Open(cfg *gormx.Config) (gorm.Dialector, error) {
	return gormlite.Open(cfg.DSN), nil
}

func init() {
	gormx.Register("sqlite", Open)
}
