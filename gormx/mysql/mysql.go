package mysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/cnk3x/pkg/gormx"
)

func Open(cfg *gormx.Config) (gorm.Dialector, error) {
	return mysql.Open(cfg.DSN), nil
}

func init() {
	gormx.Register("mysql", Open)
}
