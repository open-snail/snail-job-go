package job

import (
	"opensnail.com/snail-job/snail-job-go/dto"
	"sync"

	"github.com/sirupsen/logrus"
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
	GetRemoteLogger(name string, h logrus.Hook) Logger
	GetLocalLogger(name string) Logger
}

type loggerFactory struct {
	lock    sync.Mutex
	loggers map[string]Logger
	opts    *dto.Options
}

func (e *loggerFactory) GetRemoteLogger(name string, h logrus.Hook) Logger {
	e.lock.Lock()
	defer e.lock.Unlock()
	if e.loggers[name] == nil {
		e.loggers[name] = &logger{
			Name: name,
			//Domain: "",
			Level: e.opts.Level,
		}
	}

	if h != nil {
		logrus.AddHook(h)
	}
	return e.loggers[name]
}

func (e *loggerFactory) GetLocalLogger(name string) Logger {
	e.lock.Lock()
	defer e.lock.Unlock()
	if e.loggers[name] == nil {
		e.loggers[name] = &logger{
			Name: name,
			//Domain: "",
			Level: e.opts.Level,
		}
	}
	return e.loggers[name]
}

func NewLoggerFactory(opts *dto.Options) LoggerFactory {
	return &loggerFactory{
		loggers: make(map[string]Logger),
		opts:    opts,
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
