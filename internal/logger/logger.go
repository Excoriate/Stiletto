package logger

import (
	"github.com/Excoriate/stiletto/internal/common"
	"github.com/hashicorp/go-hclog"
	"os"
)

var (
	logLevel  = os.Getenv("PIPELINE_LOG_LEVEL")
	isEnabled = os.Getenv("PIPELINE_LOG_ENABLED")
)

type PipelineLogger struct {
}

type Logger interface {
	LogInfo(action, message string, args ...interface{})
	LogWarn(action, message string, args ...interface{})
	LogError(action, message string, args ...interface{})
	LogDebug(action, message string, args ...interface{})
	InitLogger() hclog.Logger
}

func IsEnabled() bool {
	if isEnabled == "" {
		isEnabled = "false"
	}

	return isEnabled == "true"
}

func getLogger() hclog.Logger {
	appLogger := hclog.New(&hclog.LoggerOptions{
		Name:  "stiletto",
		Level: hclog.LevelFromString(common.NormaliseStringUpper(logLevel)),
	})

	return appLogger
}

func (l *PipelineLogger) InitLogger() hclog.Logger {
	if !IsEnabled() {
		return nil
	}

	return getLogger()
}

func (l *PipelineLogger) LogInfo(action, message string, args ...interface{}) {
	if !IsEnabled() {
		return
	}

	logger := getLogger()
	if action != "" {
		logger = logger.With("action")
	}

	logger.Info(message, args...)
}

func (l *PipelineLogger) LogWarn(action, message string, args ...interface{}) {
	if !IsEnabled() {
		return
	}

	logger := getLogger()
	if action != "" {
		logger = logger.With("action")
	}

	logger.Warn(message, args...)
}

func (l *PipelineLogger) LogError(action, message string, args ...interface{}) {
	if !IsEnabled() {
		return
	}

	logger := getLogger()
	if action != "" {
		logger = logger.With("action")
	}

	logger.Error(message, args...)
}

func (l *PipelineLogger) LogDebug(action, message string, args ...interface{}) {
	if !IsEnabled() {
		return
	}

	logger := getLogger()
	if action != "" {
		logger = logger.With("action")
	}

	logger.Debug(message, args...)
}

func NewLogger() Logger {
	return &PipelineLogger{}
}
