package snailjob

import (
	"sync"

	"github.com/sirupsen/logrus"
)

type Level logrus.Level

const (
	Panic = Level(logrus.PanicLevel)
	Fatal = Level(logrus.FatalLevel)
	Error = Level(logrus.ErrorLevel)
	Warn  = Level(logrus.WarnLevel)
	Info  = Level(logrus.InfoLevel)
	Debug = Level(logrus.DebugLevel)
	Trace = Level(logrus.TraceLevel)
)

type logger struct {
	Name   string
	Domain string
	Level  Level
}

type Logger interface {
	Trace(format string, args ...interface{})
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	Fatal(format string, args ...interface{})
	Panic(format string, args ...interface{})
}

type LoggerFactory interface {
	GetLogger(name string) Logger
}

type loggerfactory struct {
	lock    sync.Mutex
	loggers map[string]Logger
}

func (e *loggerfactory) GetLogger(name string) Logger {
	e.lock.Lock()
	defer e.lock.Unlock()
	if e.loggers[name] == nil {
		e.loggers[name] = &logger{
			Name:   name,
			Domain: "",
			Level:  Info,
		}
	}

	return e.loggers[name]
}

func NewLoggerFactory() LoggerFactory {
	return &loggerfactory{
		loggers: make(map[string]Logger),
	}
}

func (l *logger) Info(format string, args ...interface{}) {

	if l.Level < Info {
		return
	}

	logrus.WithFields(logrus.Fields{"domain": l.Domain, "logger": l.Name}).Infof(format, args...)

}

func (l *logger) Fatal(format string, args ...interface{}) {
	if l.Level < Fatal {
		return
	}
	logrus.WithFields(logrus.Fields{"domain": l.Domain, "logger": l.Name}).Fatalf(format, args...)
}

func (l *logger) Warn(format string, args ...interface{}) {
	if l.Level < Warn {
		return
	}
	logrus.WithFields(logrus.Fields{"domain": l.Domain, "logger": l.Name}).Warnf(format, args...)
}

func (l *logger) Error(format string, args ...interface{}) {
	if l.Level < Error {
		return
	}
	logrus.WithFields(logrus.Fields{"domain": l.Domain, "logger": l.Name}).Errorf(format, args...)
}

func (l *logger) Debug(format string, args ...interface{}) {
	if l.Level < Debug {
		return
	}
	logrus.WithFields(logrus.Fields{"domain": l.Domain, "logger": l.Name}).Debugf(format, args...)
}

func (l *logger) Trace(format string, args ...interface{}) {
	if l.Level < Trace {
		return
	}
	logrus.WithFields(logrus.Fields{"domain": l.Domain, "logger": l.Name}).Tracef(format, args...)
}

func (l *logger) Panic(format string, args ...interface{}) {
	if l.Level < Panic {
		return
	}
	logrus.WithFields(logrus.Fields{"domain": l.Domain, "logger": l.Name}).Panicf(format, args...)
}
