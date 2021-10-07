package logger

import (
	"os"

	"github.com/quyenphamkhac/gmd-productsrv/config"
	"github.com/sirupsen/logrus"
)

type Logger interface {
	InitLogger()
	Info(message string, fields LogFields)
}

type LogFields map[string]interface{}

type serviceLogger struct {
	cfg *config.Config
	log *logrus.Logger
}

func NewServiceLogger(cfg *config.Config) (*serviceLogger, error) {
	log := logrus.New()
	return &serviceLogger{
		cfg: cfg,
		log: log,
	}, nil
}

func (l *serviceLogger) InitLogger() {
	l.log.SetFormatter(&logrus.JSONFormatter{})
	l.log.SetOutput(os.Stderr)
}

func (l *serviceLogger) Info(message string, fields LogFields) {
	if fields != nil {
		l.log.WithFields(logrus.Fields(fields)).Info(message)
		return
	}
	l.log.Info(message)
}