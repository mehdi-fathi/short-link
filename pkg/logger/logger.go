package logger

import (
	"os"

	nested "github.com/antonfisher/nested-logrus-formatter"
	graylog "github.com/gemnasium/logrus-graylog-hook/v3"
	"github.com/sirupsen/logrus"
)

type StandardLogger struct {
	*logrus.Logger
}

type Config struct {
	LogLevel      string `yaml:"LOG_LEVEL" envconfig:"LOGGER_LEVEL"`
	GrayLogActive bool   `yaml:"GRAYLOG_ACTIVE" envconfig:"LOGGER_GRAYLOG_ACTIVE"`
	GrayLogServer string `yaml:"GRAYLOG_SERVER" envconfig:"LOGGER_GRAYLOG_SERVER"`
	GrayLogStream string `yaml:"GRAYLOG_STREAM" envconfig:"LOGGER_GRAYLOG_STREAM"`
}

var StandardLoggerVar *logrus.Logger

// CreateLogger create an instance of logger
func CreateLogger(cfg Config) *StandardLogger {
	level, _ := logrus.ParseLevel(cfg.LogLevel)

	logger := logrus.New()
	logger.Out = os.Stdout
	logger.SetFormatter(&nested.Formatter{
		FieldsOrder: []string{"kind", "operation"},
		HideKeys:    true,
	})
	logger.SetLevel(level)

	if cfg.GrayLogActive {
		hook := graylog.NewGraylogHook(cfg.GrayLogServer, map[string]interface{}{"stream": cfg.GrayLogStream})
		logger.AddHook(hook)
	}

	standardLogger := &StandardLogger{logger}
	StandardLoggerVar = logger

	return standardLogger
}

// DefaultLogger configs logrus default instance
func DefaultLogger(cfg Config) {
	level, _ := logrus.ParseLevel(cfg.LogLevel)
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&nested.Formatter{
		FieldsOrder: []string{"hostname", "operation"},
		HideKeys:    true,
	})
	logrus.SetLevel(level)

	if cfg.GrayLogActive {
		hook := graylog.NewGraylogHook(cfg.GrayLogServer, map[string]interface{}{"stream": cfg.GrayLogStream})
		logrus.AddHook(hook)
	}
}

func CreateLogInfo(text string) {

	StandardLoggerVar.Info(text)
}

func CreateLogError(text string) {

	StandardLoggerVar.Error(text)
}
