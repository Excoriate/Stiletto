package logger

import (
	log "github.com/sirupsen/logrus"
	"os"
)

var (
	logLevel  = os.Getenv("STILETTO_LOG_LEVEL")
	isEnabled = os.Getenv("STILETTO_LOG_ENABLED")
	format    = os.Getenv("STILETTO_LOG_FORMAT")
)

type StilettoLog struct {
}

type Logger interface {
	LogInfo(message, details string, args ...interface{})
	LogWarn(message, details string, args ...interface{})
	LogError(message, details string, args ...interface{})
	LogFatal(message, details string, args ...interface{})
	LogDebug(message, details string, args ...interface{})
	InitLogger()
}

func IsEnabled() bool {
	if isEnabled == "" {
		isEnabled = "false"
	}

	return isEnabled == "true"
}

func getLogger(message, details string) *log.Entry {
	logger := log.WithFields(log.Fields{
		"message": message,
		"details": details,
	})

	return logger
}

func setLogLevel(level string) {
	switch level {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	}
}

func (l *StilettoLog) InitLogger() {
	if !IsEnabled() {
		return
	}

	if format == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{})
	}

	if logLevel == "" {
		logLevel = "error"
	}

	setLogLevel(logLevel)
}

func (l *StilettoLog) LogInfo(message, details string, args ...interface{}) {
	if !IsEnabled() {
		return
	}
	logger := getLogger(message, details)
	logger.Info(args...)
}

func (l *StilettoLog) LogWarn(message, details string, args ...interface{}) {
	if !IsEnabled() {
		return
	}

	logger := getLogger(message, details)
	logger.Warn(args...)
}

func (l *StilettoLog) LogError(message, details string, args ...interface{}) {
	if !IsEnabled() {
		return
	}

	logger := getLogger(message, details)
	logger.Error(args...)
}

func (l *StilettoLog) LogFatal(message, details string, args ...interface{}) {
	if !IsEnabled() {
		return
	}

	logger := getLogger(message, details)
	logger.Fatal(args...)
}

func (l *StilettoLog) LogDebug(message, details string, args ...interface{}) {
	if !IsEnabled() {
		return
	}

	logger := getLogger(message, details)
	logger.Debug(args...)
}
