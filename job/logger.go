package job

import (
	"context"
	"sync"

	"github.com/open-snail/snail-job-go/dto"
	"github.com/sirupsen/logrus"
)

type LoggerFactory interface {
	GetRemoteLogger(name string, ctx context.Context) *logrus.Entry
	GetLocalLogger(name string) *logrus.Entry
	GetLogRus() *logrus.Logger
	Init(hls *HookLogService)
}

type loggerFactory struct {
	lock   sync.Mutex
	opts   *dto.Options
	logger *logrus.Logger
}

func (e *loggerFactory) GetRemoteLogger(name string, ctx context.Context) *logrus.Entry {
	e.lock.Lock()
	defer e.lock.Unlock()
	return e.logger.WithFields(logrus.Fields{"logger": name}).WithContext(ctx)
}

func (e *loggerFactory) GetLocalLogger(name string) *logrus.Entry {
	e.lock.Lock()
	defer e.lock.Unlock()
	return e.logger.WithFields(logrus.Fields{"logger": name})
}

func NewLoggerFactory(opts *dto.Options) LoggerFactory {
	return &loggerFactory{
		opts:   opts,
		logger: logrus.New(),
	}
}

func (e *loggerFactory) Init(hls *HookLogService) {
	log := e.logger
	// 添加日志hook
	log.AddHook(&LoggerHook{Hls: hls})
	// 日志添加调用者信息
	log.SetReportCaller(e.opts.ReportCaller)
	// 设置日志级别
	log.SetLevel(e.opts.Level)
	// 设置日志格式
	if e.opts.Formatter != nil {
		log.SetFormatter(e.opts.Formatter)
	} else {
		log.SetFormatter(&DefaultFormatter{})
	}
}

func (e *loggerFactory) GetLogRus() *logrus.Logger {
	return e.logger
}
