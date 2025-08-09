package core

import "log"

type StdLoggerAdapter struct{ l *log.Logger }
func NewStdLoggerAdapter(l *log.Logger) StdLoggerAdapter { return StdLoggerAdapter{l:l} }
func (s StdLoggerAdapter) Debug(msg string, kv ...any) { s.l.Print(append([]any{"DEBUG:", msg}, kv...)...) }
func (s StdLoggerAdapter) Info(msg string, kv ...any)  { s.l.Print(append([]any{"INFO:", msg}, kv...)...) }
func (s StdLoggerAdapter) Warn(msg string, kv ...any)  { s.l.Print(append([]any{"WARN:", msg}, kv...)...) }
func (s StdLoggerAdapter) Error(msg string, kv ...any) { s.l.Print(append([]any{"ERROR:", msg}, kv...)...) }
