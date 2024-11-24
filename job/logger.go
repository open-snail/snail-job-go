package job

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"
	"opensnail.com/snail-job/snail-job-go/dto"
)

const (
	Panic = logrus.PanicLevel
	Fatal = logrus.FatalLevel
	Error = logrus.ErrorLevel
	Warn  = logrus.WarnLevel
	Info  = logrus.InfoLevel
	Debug = logrus.DebugLevel
	Trace = logrus.TraceLevel
)

type logger struct {
	Name   string
	Domain string
	Level  logrus.Level
	logger *logrus.Entry
}

type SnailJobLogger interface {
	Trace(format string, args ...interface{})
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	Fatal(format string, args ...interface{})
	Panic(format string, args ...interface{})
}

type LoggerFactory interface {
	GetRemoteLogger(name string, ctx context.Context) SnailJobLogger
	GetLocalLogger(name string) SnailJobLogger
	GetLogRus() *logrus.Logger
}

type loggerFactory struct {
	lock    sync.Mutex
	loggers map[string]SnailJobLogger
	opts    *dto.Options
	logger  *logrus.Logger
}

func (e *loggerFactory) GetRemoteLogger(name string, ctx context.Context) SnailJobLogger {
	e.lock.Lock()
	defer e.lock.Unlock()
	//if e.loggers[name] == nil {
	//	e.loggers[name] =
	//}

	return &logger{
		Name: name,
		//Domain: "",
		Level:  e.opts.Level,
		logger: e.logger.WithContext(ctx),
	}
}

func (e *loggerFactory) GetLocalLogger(name string) SnailJobLogger {
	e.lock.Lock()
	defer e.lock.Unlock()
	//if e.loggers[name] == nil {
	//	e.loggers[name] =
	//}

	return &logger{
		Name: name,
		//Domain: "",
		Level:  e.opts.Level,
		logger: e.logger.WithContext(nil),
	}
}

func NewLoggerFactory(opts *dto.Options) LoggerFactory {
	logrus := logrus.New()
	return &loggerFactory{
		loggers: make(map[string]SnailJobLogger),
		opts:    opts,
		logger:  logrus,
	}
}

func (e *loggerFactory) GetLogRus() *logrus.Logger {
	return e.logger
}

func (l *logger) Info(format string, args ...interface{}) {
	if l.Level < Info {
		return
	}
	l.logger.WithFields(logrus.Fields{"domain": l.Domain, "logger": l.Name}).Infof(format, args...)

}

func (l *logger) Fatal(format string, args ...interface{}) {
	if l.Level < Fatal {
		return
	}
	l.logger.WithFields(logrus.Fields{"domain": l.Domain, "logger": l.Name}).Fatalf(format, args...)
}

func (l *logger) Warn(format string, args ...interface{}) {
	if l.Level < Warn {
		return
	}
	l.logger.WithFields(logrus.Fields{"domain": l.Domain, "logger": l.Name}).Warnf(format, args...)
}

func (l *logger) Error(format string, args ...interface{}) {
	if l.Level < Error {
		return
	}
	l.logger.WithFields(logrus.Fields{"domain": l.Domain, "logger": l.Name}).Errorf(format, args...)
}

func (l *logger) Debug(format string, args ...interface{}) {
	if l.Level < Debug {
		return
	}
	l.logger.WithFields(logrus.Fields{"domain": l.Domain, "logger": l.Name}).Debugf(format, args...)
}

func (l *logger) Trace(format string, args ...interface{}) {
	if l.Level < Trace {
		return
	}
	l.logger.WithFields(logrus.Fields{"domain": l.Domain, "logger": l.Name}).Tracef(format, args...)
}

func (l *logger) Panic(format string, args ...interface{}) {
	if l.Level < Panic {
		return
	}
	l.logger.WithFields(logrus.Fields{"domain": l.Domain, "logger": l.Name}).Panicf(format, args...)
}
