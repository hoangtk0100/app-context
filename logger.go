package appctx

import (
	"fmt"

	"github.com/rs/zerolog"
)

type Logger interface {
	GetLevel() string

	Print(args ...interface{})
	Printf(format string, args ...interface{})

	Debug(args ...interface{})
	Debugf(format string, args ...interface{})

	Info(args ...interface{})
	Infof(format string, args ...interface{})

	Warn(args ...interface{})
	Warnf(format string, args ...interface{})

	Error(err error, args ...interface{})
	Errorf(err error, format string, args ...interface{})

	Fatal(err error, args ...interface{})
	Fatalf(err error, format string, args ...interface{})
}

type logger struct {
	*zerolog.Logger
}

func (lg *logger) GetLevel() string {
	return zerolog.GlobalLevel().String()
}

func (lg *logger) Print(args ...interface{}) {
	lg.Debug(fmt.Sprint(args...))
}

func (lg *logger) Printf(format string, args ...interface{}) {
	lg.Debugf(format, args...)
}

func (lg *logger) Debug(args ...interface{}) {
	lg.Logger.Debug().Msg(fmt.Sprint(args...))
}

func (lg *logger) Debugf(format string, args ...interface{}) {
	lg.Logger.Debug().Msgf(format, args...)
}

func (lg *logger) Info(args ...interface{}) {
	lg.Logger.Info().Msg(fmt.Sprint(args...))
}

func (lg *logger) Infof(format string, args ...interface{}) {
	lg.Logger.Info().Msgf(format, args...)
}

func (lg *logger) Warn(args ...interface{}) {
	lg.Logger.Warn().Msg(fmt.Sprint(args...))
}

func (lg *logger) Warnf(format string, args ...interface{}) {
	lg.Logger.Warn().Msgf(format, fmt.Sprint(args...))
}

func (lg *logger) Error(err error, args ...interface{}) {
	lg.Logger.Error().Err(err).Msg(fmt.Sprint(args...))
}

func (lg *logger) Errorf(err error, format string, args ...interface{}) {
	lg.Logger.Error().Err(err).Msgf(format, fmt.Sprint(args...))
}

func (lg *logger) Fatal(err error, args ...interface{}) {
	lg.Logger.Fatal().Err(err).Msg(fmt.Sprint(args...))
}

func (lg *logger) Fatalf(err error, format string, args ...interface{}) {
	lg.Logger.Fatal().Err(err).Msgf(format, fmt.Sprint(args...))
}
