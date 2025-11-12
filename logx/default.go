package logx

import "log/slog"

var (
	Debug = slog.Debug
	Info  = slog.Info
	Warn  = slog.Warn
	Error = slog.Error

	DebugContext = slog.DebugContext
	InfoContext  = slog.InfoContext
	WarnContext  = slog.WarnContext
	ErrorContext = slog.ErrorContext
)
