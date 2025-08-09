package core

import "log/slog"

type SlogAdapter struct{ l *slog.Logger }
func NewSlogAdapter(l *slog.Logger) SlogAdapter { return SlogAdapter{l:l} }
func (s SlogAdapter) Debug(msg string, kv ...any) { s.l.Debug(msg, kv...) }
func (s SlogAdapter) Info(msg string, kv ...any)  { s.l.Info(msg, kv...) }
func (s SlogAdapter) Warn(msg string, kv ...any)  { s.l.Warn(msg, kv...) }
func (s SlogAdapter) Error(msg string, kv ...any) { s.l.Error(msg, kv...) }
