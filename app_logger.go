package appctx

import (
	"errors"
	"fmt"
	glog "log"
	"os"
	"strings"

	"github.com/rs/zerolog"
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

	logger := getNewLogger(config.defaultType, config.defaultPath)
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
		glog.Fatal("Error parsing log level")
	}

	return level
}

func getNewLogger(logType, logPath string) *zerolog.Logger {
	var writer *os.File

	switch logType {
	case "stderr":
		writer = os.Stderr

	case "file":
		if logPath == "" {
			glog.Fatal(errors.New("empty log path"))
		}

		var err error
		writer, err = os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			glog.Fatal(err)
		}

	default:
		writer = os.Stdout
	}

	logger := zerolog.New(writer).With().Timestamp().Logger()
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

func (l *appLogger) Run(_ AppContext) error {
	level := parseLogLevel(l.logLevel)
	zerolog.SetGlobalLevel(level)

	if level <= zerolog.DebugLevel {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	}

	l.logger = getNewLogger(l.logType, l.logPath)
	return nil
}

func (l *appLogger) Stop() error {
	return nil
}
