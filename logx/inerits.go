package logx

import "log/slog"

const (
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
)

type (
	Logger = slog.Logger
	Level  = slog.Level
	Value  = slog.Value
	Record = slog.Record
)

var (
	sWith      = slog.With
	sNew       = slog.New
	SetDefault = slog.SetDefault
	Default    = slog.Default
)
