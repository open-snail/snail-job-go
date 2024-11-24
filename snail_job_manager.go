package snailjob

import (
	"fmt"
	"sync"

	"opensnail.com/snail-job/snail-job-go/constant"
	"opensnail.com/snail-job/snail-job-go/dto"
	"opensnail.com/snail-job/snail-job-go/job"
)

// SnailJobManager snail job 客户端启动者
type SnailJobManager struct {
	factory   job.LoggerFactory
	logger    job.SnailJobLogger
	executors map[string]job.NewJobExecutor
	client    job.SnailJobClient
	opts      *dto.Options
	hls       *job.HookLogService
	lock      sync.Mutex
}

func NewSnailJobManager(opts *dto.Options) *SnailJobManager {
	factory := job.NewLoggerFactory(opts)
	logger := factory.GetLocalLogger("snail-job-manager")
	client := job.NewSnailJobClient(opts, factory)
	hls := job.NewHookLogService(client)
	client.RemoteLog.AddHook(&job.LoggerHook{Hls: hls})
	go hls.Init()

	return &SnailJobManager{
		factory:   factory,
		logger:    logger,
		executors: make(map[string]job.NewJobExecutor),
		opts:      opts,
		client:    client,
		hls:       hls,
	}
}

func (e *SnailJobManager) GetLoggerFactory() job.LoggerFactory {
	return e.factory
}

func (e *SnailJobManager) GetClient() job.SnailJobClient {
	return e.client
}

func (e *SnailJobManager) Run() {
	e.logger.Info("Run SnailJob Client v%s", constant.VERSION)
	go e.client.SendHeartbeat()
	// go e.hls.Init()
	job.RunServer(e.opts, e.client, e.executors, e.factory)
}

func (e *SnailJobManager) Init() error {
	e.logger.Info("%s", "Init manager")
	// 添加日志hook
	// e.factory.GetLogRus().AddHook(&job.LoggerHook{Hls: e.hls})
	// 日志添加调用者信息
	// e.factory.GetLogRus().SetReportCaller(true)
	// 设置日志级别
	// TODO: 配置level
	// e.factory.GetLogRus().SetLevel(logrus.DebugLevel)
	return nil
}

// Register 注册job执行器入门
func (e *SnailJobManager) Register(name string, executor job.NewJobExecutor) *SnailJobManager {
	e.lock.Lock()
	defer e.lock.Unlock()

	if _, exists := e.executors[name]; exists {
		panic(fmt.Sprintf("SnailJobManager [%s] already registered", name))
	}

	e.executors[name] = executor
	e.logger.Info(fmt.Sprintf("Registered executor: %s", name))
	return e
}
