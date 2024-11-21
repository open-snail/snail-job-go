package snailjob

import (
	"fmt"
	"opensnail.com/snail-job/snail-job-go/constant"
	"opensnail.com/snail-job/snail-job-go/dto"
	"opensnail.com/snail-job/snail-job-go/job"
	"sync"
)

// SnailJobManager snail job 客户端启动者
type SnailJobManager struct {
	factory   LoggerFactory
	logger    Logger
	executors map[string]job.IJobExecutor
	client    job.SnailJobClient
	opts      *dto.Options
	lock      sync.Mutex
}

func NewSnailJobManager(opts *dto.Options, factory LoggerFactory) *SnailJobManager {
	return &SnailJobManager{
		factory:   factory,
		logger:    factory.GetLogger("executor"),
		executors: make(map[string]job.IJobExecutor),
		opts:      opts,
		client:    job.NewSnailJobClient(opts),
	}
}

func (e *SnailJobManager) GetLoggerFactory() LoggerFactory {
	return e.factory
}

func (e *SnailJobManager) GetClient() job.SnailJobClient {
	return e.client
}

func (e *SnailJobManager) Run() {
	e.logger.Info("Run SnailJob Client v%s", constant.VERSION)
	go e.client.SendHeartbeat()
	go job.NewHookLogService(e.client)
	job.RunServer(e.opts, e.client, e.executors)
}

func (e *SnailJobManager) Init() error {
	e.logger.Info("%s", "Init manager")
	return nil
}

// Register 注册job执行器入门
func (e *SnailJobManager) Register(name string, executor job.IJobExecutor) {
	e.lock.Lock()
	defer e.lock.Unlock()

	if _, exists := e.executors[name]; exists {
		panic(fmt.Sprintf("SnailJobManager [%s] already registered", name))
	}

	e.executors[name] = executor
	job.LocalLog.Info(fmt.Sprintf("Registered executor: %s", name))
}
