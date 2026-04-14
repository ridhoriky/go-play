package logger

import (
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LoggerOptions holds logger configuration
type LoggerOptions struct {
	Enabled    bool   `yaml:"enabled"`
	Level      string `yaml:"level"`
	Format     string `yaml:"format"`
	Output     string `yaml:"output"`
	Path       string `yaml:"path"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
	Compress   bool   `yaml:"compress"`
}

var (
	onceLogger = sync.Once{}
	logInst    zerolog.Logger
)

// InitLogger initializes the logger
func InitLogger(opt LoggerOptions) zerolog.Logger {
	onceLogger.Do(func() {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		zerolog.TimeFieldFormat = time.RFC3339

		logLevel := parseLogLevel(opt.Level)
		writer := baseWriter(opt.Output)
		writer = formatWriter(opt.Format, writer)

		if opt.Enabled && opt.Path != "" {
			fileLogger := &lumberjack.Logger{
				Filename:   opt.Path,
				MaxSize:    opt.MaxSize,
				MaxBackups: opt.MaxBackups,
				MaxAge:     opt.MaxAge,
				Compress:   opt.Compress,
			}
			writer = zerolog.MultiLevelWriter(writer, fileLogger)
		}

		logInst = zerolog.New(writer).
			Level(logLevel).
			With().
			Timestamp().
			Caller().
			Logger()
	})

	return logInst
}

func parseLogLevel(defaultLevel string) zerolog.Level {
	levelValue := strings.TrimSpace(os.Getenv("LOG_LEVEL"))
	if levelValue == "" {
		levelValue = strings.TrimSpace(defaultLevel)
	}

	if lvl, err := zerolog.ParseLevel(strings.ToLower(levelValue)); err == nil {
		return lvl
	}

	if intVal, err := strconv.Atoi(levelValue); err == nil {
		if intVal >= int(zerolog.NoLevel) && intVal <= int(zerolog.FatalLevel) {
			return zerolog.Level(intVal)
		}
	}

	return zerolog.InfoLevel
}

func baseWriter(output string) io.Writer {
	switch strings.ToLower(strings.TrimSpace(output)) {
	case "stderr":
		return os.Stderr
	default:
		return os.Stdout
	}
}

func formatWriter(format string, writer io.Writer) io.Writer {
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "console":
		return zerolog.ConsoleWriter{
			Out:        writer,
			TimeFormat: time.RFC3339,
		}
	default:
		return writer
	}
}
