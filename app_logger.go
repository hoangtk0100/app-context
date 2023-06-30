package appctx

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/spf13/pflag"
)

var (
	defaultLevel = "info"
	defaultType  = "stdout"

	defaultLogger = newAppLogger(&loggerConfig{
		basePrefix:   "core",
		defaultLevel: "trace",
		defaultType:  defaultType,
	})
)

type AppLogger interface {
	GetLogger(prefix string) Logger
}

func GlobalLogger() AppLogger {
	return defaultLogger
}

type appLogger struct {
	config   loggerConfig
	logger   *zerolog.Logger
	logLevel string
	logType  string
	logPath  string
}

type loggerConfig struct {
	basePrefix   string
	defaultLevel string
	defaultType  string
	defaultPath  string
}

func newAppLogger(config *loggerConfig) *appLogger {
	if config == nil {
		config = &loggerConfig{}
	}

	if config.defaultLevel == "" {
		config.defaultLevel = defaultLevel
	}

	if config.defaultType == "" {
		config.defaultType = defaultType
	}

	level := parseLogLevel(config.defaultLevel)
	zerolog.SetGlobalLevel(level)

	logger := getNewLogger(prdEnv, config.defaultType, config.defaultPath)
	return &appLogger{
		config:   *config,
		logger:   logger,
		logLevel: config.defaultLevel,
		logType:  config.defaultType,
		logPath:  config.defaultPath,
	}
}

func parseLogLevel(input string) zerolog.Level {
	level, err := zerolog.ParseLevel(input)
	if err != nil {
		log.Fatal().Err(err).Msg("Error parsing log level")
	}

	return level
}

func getOutputFormat(env string, writer *os.File) io.Writer {
	if prdEnv == env {
		return writer
	}

	output := zerolog.ConsoleWriter{Out: writer, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("*** %s ***", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%-6s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}

	return output
}

func getNewLogger(env, logType, logPath string) *zerolog.Logger {
	var writer *os.File

	switch logType {
	case "stderr":
		writer = os.Stderr

	case "file":
		if logPath == "" {
			log.Fatal().Msg("Empty log path")
		}

		var err error
		writer, err = os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal().Err(err).Msg("Cannot write to log file")
		}

	default:
		writer = os.Stdout
	}

	output := getOutputFormat(env, writer)
	logger := zerolog.New(output).With().Timestamp().Logger()
	return &logger
}

func (l *appLogger) GetLogger(prefix string) Logger {
	var log *zerolog.Logger

	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		log = l.logger
	} else {
		prefix = fmt.Sprintf("%s.%s", l.config.basePrefix, prefix)
		lg := l.logger.With().Str("prefix", prefix).Logger()
		log = &lg
	}

	return &logger{
		Logger: log,
	}
}

func (l *appLogger) ID() string {
	return "logger"
}

func (l *appLogger) InitFlags() {
	pflag.StringVar(
		&l.logLevel,
		"log-level",
		l.config.defaultLevel,
		"Log level (panic | fatal | error | warn | info | debug | trace) - Default: trace",
	)

	pflag.StringVar(
		&l.logType,
		"log-type",
		l.config.defaultType,
		"Log type (stdout | stderr | file) - Default: stdout",
	)

	pflag.StringVar(
		&l.logPath,
		"log-path",
		l.config.defaultPath,
		"Log path (require if log type is file) - Ex: \"./app.log\"",
	)
}

func (l *appLogger) Run(ac AppContext) error {
	level := parseLogLevel(l.logLevel)
	zerolog.SetGlobalLevel(level)

	if level <= zerolog.DebugLevel {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	}

	l.logger = getNewLogger(ac.GetEnvName(), l.logType, l.logPath)
	return nil
}

func (l *appLogger) Stop() error {
	return nil
}
