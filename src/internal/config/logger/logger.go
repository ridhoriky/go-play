package logger

import (
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LoggerOptions holds logger configuration
type LoggerOptions struct {
	Enabled    bool   `yaml:"enabled" env:"LOGGER_ENABLED" env-default:"true"`
	Level      string `yaml:"level" env:"LOGGER_LEVEL" env-default:"info"`
	Format     string `yaml:"format" env:"LOGGER_FORMAT" env-default:"json"`
	Output     string `yaml:"output" env:"LOGGER_OUTPUT" env-default:"stdout"`
	Path       string `yaml:"path" env:"LOGGER_PATH" env-default:"./logs/app.log"`
	MaxSize    int    `yaml:"max_size" env:"LOGGER_MAX_SIZE" env-default:"100"`
	MaxBackups int    `yaml:"max_backups" env:"LOGGER_MAX_BACKUPS" env-default:"7"`
	MaxAge     int    `yaml:"max_age" env:"LOGGER_MAX_AGE" env-default:"30"`
	Compress   bool   `yaml:"compress" env:"LOGGER_COMPRESS" env-default:"false"`
}

// InitLogger initializes the logger
func InitLogger(opt LoggerOptions) *zerolog.Logger {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = time.RFC3339

	// Add caller info to all log levels
	logLevel := parseLogLevel(opt.Level)

	// Base writer setup
	writer := baseWriter(opt.Output)

	// Apply format to the base writer
	writer = formatWriter(opt.Format, writer)

	// Add file writer if enabled and path is provided
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

	logInst := zerolog.New(writer).
		Level(logLevel).
		With().
		Timestamp().
		Caller().
		Logger()

	return &logInst
}

func parseLogLevel(defaultLevel string) zerolog.Level {
	if lvl, err := zerolog.ParseLevel(strings.ToLower(defaultLevel)); err == nil {
		return lvl
	}

	if intVal, err := strconv.Atoi(defaultLevel); err == nil {
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
