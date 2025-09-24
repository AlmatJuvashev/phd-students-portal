package logging

import (
	"log"
)

// For starter simplicity, wrap stdlib log. In production, swap to zerolog/zap.
func Info(msg string, kv ...any)  { log.Println(append([]any{"INFO", msg}, kv...)...) }
func Warn(msg string, kv ...any)  { log.Println(append([]any{"WARN", msg}, kv...)...) }
func Error(msg string, kv ...any) { log.Println(append([]any{"ERROR", msg}, kv...)...) }
