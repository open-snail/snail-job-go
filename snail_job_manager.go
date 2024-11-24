package snailjob

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
	"opensnail.com/snail-job/snail-job-go/constant"
	"opensnail.com/snail-job/snail-job-go/dto"
	"opensnail.com/snail-job/snail-job-go/job"
)

// SnailJobManager snail job 客户端启动者
type SnailJobManager struct {
	factory   job.LoggerFactory
	logger    *logrus.Entry
	executors map[string]job.NewJobExecutor
	client    job.SnailJobClient
	opts      *dto.Options
	hls       *job.HookLogService
	lock      sync.Mutex
}

func NewSnailJobManager(opts *dto.Options) *SnailJobManager {
	factory := job.NewLoggerFactory(opts)
	client := job.NewSnailJobClient(opts, factory)
	hls := job.NewHookLogService(client)
	factory.Init(hls)
	return &SnailJobManager{
		factory:   factory,
		logger:    factory.GetLocalLogger("snail-job-manager"),
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
	e.logger.Infof("Run SnailJob Client v%s", constant.VERSION)
	go e.client.SendHeartbeat()
	go e.hls.Init()
	job.RunServer(e.opts, e.client, e.executors, e.factory)
}

func (e *SnailJobManager) Init() error {
	e.logger.Infof("%s", "Init manager")
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
	e.logger.Infof("Registered executor: %s", name)
	return e
}
