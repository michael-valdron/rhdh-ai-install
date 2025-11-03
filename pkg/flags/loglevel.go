package flags

import (
	"fmt"
	"log/slog"

	"github.com/spf13/pflag"
)

type LogLevel struct {
	pflag.Value
	obj *slog.Level
	str string
}

func NewLogLevel(level *slog.Level) *LogLevel {
	return &LogLevel{obj: level}
}

func (l *LogLevel) Set(value string) error {
	switch value {
	case "info":
		*l.obj = slog.LevelInfo
	case "debug":
		*l.obj = slog.LevelDebug
	case "warn":
		*l.obj = slog.LevelWarn
	case "error":
		*l.obj = slog.LevelError
	default:
		return fmt.Errorf("unrecognized log level '%s'", value)
	}
	l.str = value
	return nil
}

func (l *LogLevel) String() string {
	return l.str
}

func (l *LogLevel) Type() string {
	return "slog.Level"
}
