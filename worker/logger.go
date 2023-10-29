package worker

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logger struct{}

func NewLogger() *Logger {
	return &Logger{}
}

func (logger *Logger) Printf(ctx context.Context, format string, v ...interface{}) {
	log.WithLevel(zerolog.DebugLevel).Msgf(format, v...)
}

func (logger *Logger) Print(level zerolog.Level, args ...any) {
	log.WithLevel(level).Msg(fmt.Sprint(args...))
}

func (logger *Logger) Debug(arg ...any) {
	logger.Print(zerolog.DebugLevel, arg...)
}

func (logger *Logger) Info(arg ...any) {
	logger.Print(zerolog.InfoLevel, arg...)
}

func (logger *Logger) Warn(arg ...any) {
	logger.Print(zerolog.WarnLevel, arg...)
}

func (logger *Logger) Error(arg ...any) {
	logger.Print(zerolog.ErrorLevel, arg...)
}

func (logger *Logger) Fatal(arg ...any) {
	logger.Print(zerolog.FatalLevel, arg...)
}
