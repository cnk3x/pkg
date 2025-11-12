package configx

import (
	"os"
	"path/filepath"
	"slices"
)

const (
	BinaryPath       = "binary"
	CachePath        = "cache"
	TempPath         = "temp"
	LogPath          = "log"
	UIPath           = "ui"
	ConfigPath       = "config"
	PidPath          = "pid"
	RuntimePath      = "runtime"
	SubscriptionPath = "subscriptions"
)

var workPath string

func init() {
	workPath, _ = os.Getwd()
}

func SetWorkPath(path string) { workPath, _ = filepath.Abs(path) }

func WorkPath() string { return workPath }

func GetWorkPath(pType string, paths ...string) string {
	return filepath.Join(slices.Insert(paths, 0, workPath, pType)...)
}
