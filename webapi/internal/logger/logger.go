package logger

import (
	"os"
	"webapi/config"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	Error(args ...interface{})
	Info(args ...interface{})
	Debug(args ...interface{})
}

type logger struct {
	log *logrus.Logger
}

func New(cfg config.LoggingConfig) Logger {
	log := logrus.New()

	// File output
	file, err := os.OpenFile(cfg.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	}

	// TODO: Add Elasticsearch hook

	return &logger{log: log}
}

func (l *logger) Error(args ...interface{}) {
	l.log.Error(args...)
}

func (l *logger) Info(args ...interface{}) {
	l.log.Info(args...)
}

func (l *logger) Debug(args ...interface{}) {
	l.log.Debug(args...)
}
