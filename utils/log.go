package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	logrus.FieldLogger
}

type LoggingConfig struct {
	Level        LogLevel
	Formatter    LogFormatter
	ReportCaller bool
}

type LogLevel string
type LogFormatter string

const (
	PanicLevel LogLevel = "panic"
	FatalLevel          = "fatal"
	ErrorLevel          = "error"
	WarnLevel           = "warn"
	InfoLevel           = "info"
	DebugLevel          = "debug"
	TraceLevel          = "trace"

	TerminalFormat LogFormatter = "terminal"
	JSONFormat                  = "json"
	TextFormat                  = "text"
)

func CreateLogger(config *LoggingConfig) Logger {
	return &logrus.Logger{
		Out:          os.Stderr,
		Formatter:    getLogrusFormatter(config.Formatter),
		Hooks:        make(logrus.LevelHooks),
		Level:        getLogrusLevel(config.Level),
		ExitFunc:     os.Exit,
		ReportCaller: config.ReportCaller,
	}
}

func getLogrusLevel(level LogLevel) logrus.Level {
	switch level {
	case PanicLevel:
		return logrus.PanicLevel
	case FatalLevel:
		return logrus.FatalLevel
	case ErrorLevel:
		return logrus.ErrorLevel
	case WarnLevel:
		return logrus.WarnLevel
	case InfoLevel:
		return logrus.InfoLevel
	case DebugLevel:
		return logrus.DebugLevel
	default:
		return logrus.TraceLevel
	}
}

func getLogrusFormatter(format LogFormatter) logrus.Formatter {
	switch format {
	case TextFormat:
		return &logrus.TextFormatter{DisableColors: true, FullTimestamp: true}
	case JSONFormat:
		return &logrus.JSONFormatter{}
	default:
		return &logrus.TextFormatter{}
	}
}
