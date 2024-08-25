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

type Graylog struct {
	Level  string `envconfig:"LEVEL"`
	Active bool   `envconfig:"ACTIVE"`
	Server string `envconfig:"SERVER"`
	Stream string `envconfig:"STREAM"`
}

var StandardLoggerVar *logrus.Logger

// CreateLogger create an instance of logger
func CreateLogger(cfg Graylog) *StandardLogger {
	level, _ := logrus.ParseLevel(cfg.Level)

	logger := logrus.New()
	logger.Out = os.Stdout
	logger.SetFormatter(&nested.Formatter{
		FieldsOrder: []string{"kind", "operation"},
		HideKeys:    true,
	})
	logger.SetLevel(level)

	if cfg.Active {
		hook := graylog.NewGraylogHook(cfg.Server, map[string]interface{}{"stream": cfg.Stream})
		logger.AddHook(hook)
	}

	standardLogger := &StandardLogger{logger}
	StandardLoggerVar = logger

	return standardLogger
}

// DefaultLogger configs logrus default instance
func DefaultLogger(cfg Graylog) {
	level, _ := logrus.ParseLevel(cfg.Level)
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&nested.Formatter{
		FieldsOrder: []string{"hostname", "operation"},
		HideKeys:    true,
	})
	logrus.SetLevel(level)

	if cfg.Active {
		hook := graylog.NewGraylogHook(cfg.Server, map[string]interface{}{"stream": cfg.Stream})
		logrus.AddHook(hook)
	}
}

func CreateLogInfo(text interface{}) {

	StandardLoggerVar.Info(text)
}

func CreateLogError(text string) {

	StandardLoggerVar.Error(text)
}
