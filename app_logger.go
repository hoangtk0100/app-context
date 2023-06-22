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

	defaultLogger = newAppLogger(&Config{
		BasePrefix:   "core",
		DefaultLevel: "trace",
		DefaultType:  defaultType,
	})
)

type AppLogger interface {
	GetLogger(prefix string) Logger
}

func GlobalLogger() AppLogger {
	return defaultLogger
}

type appLogger struct {
	config   Config
	logger   *zerolog.Logger
	logLevel string
	logType  string
	logPath  string
}

type Config struct {
	BasePrefix   string
	DefaultLevel string
	DefaultType  string
	DefaultPath  string
}

func newAppLogger(config *Config) *appLogger {
	if config == nil {
		config = &Config{}
	}

	if config.DefaultLevel == "" {
		config.DefaultLevel = defaultLevel
	}

	if config.DefaultType == "" {
		config.DefaultType = defaultType
	}

	level := parseLogLevel(config.DefaultLevel)
	zerolog.SetGlobalLevel(level)

	logger := getNewLogger(config.DefaultType, config.DefaultPath)
	return &appLogger{
		config:   *config,
		logger:   logger,
		logLevel: config.DefaultLevel,
		logType:  config.DefaultType,
		logPath:  config.DefaultPath,
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
		prefix = fmt.Sprintf("%s.%s", l.config.BasePrefix, prefix)
		logger := l.logger.With().Str("prefix", prefix).Logger()
		log = &logger
	}

	return &logger{
		Logger: log,
	}
}

func (l *appLogger) ID() string {
	return "logger"
}

func (l *appLogger) InitFlags() {
	pflag.StringVar(&l.logLevel, "log-level", l.config.DefaultLevel, "Log level: panic | fatal | error | warn | info | debug | trace")
	pflag.StringVar(&l.logType, "log-type", l.config.DefaultType, "Log type: stdout | stderr | file")
	pflag.StringVar(&l.logPath, "log-path", l.config.DefaultPath, "Log path (require if log type is file), Ex: ./")
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
