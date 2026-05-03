package logger

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LoggerOptions holds logger configuration
type LoggerOptions struct {
	Enabled    bool   `yaml:"enabled" env:"LOGGER_ENABLED"`
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

	logLevel := parseLogLevel(opt.Level)

	var writers []io.Writer
	if !opt.Enabled {
		writers = append(writers, io.Discard)
	} else {
		consoleWriter := formatWriter(opt.Format, baseWriter(opt.Output))
		writers = append(writers, consoleWriter)

		if opt.Path != "" {
			dir := filepath.Dir(opt.Path)
			if dir != "." && dir != "" {
				_ = os.MkdirAll(dir, 0o755)
			}

			fileLogger := &lumberjack.Logger{
				Filename:   opt.Path,
				MaxSize:    opt.MaxSize,
				MaxBackups: opt.MaxBackups,
				MaxAge:     opt.MaxAge,
				Compress:   opt.Compress,
			}
			writers = append(writers, fileLogger)
		}
	}

	writer := io.Discard
	if len(writers) == 1 {
		writer = writers[0]
	} else if len(writers) > 1 {
		writer = zerolog.MultiLevelWriter(writers...)
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
