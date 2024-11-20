package snailjob

import (
	"fmt"
	jobExecutor "opensnail.com/snail-job/snail-job-go/executor"
	"opensnail.com/snail-job/snail-job-go/job"
	"sync"
)

type Options struct {
	ServerHost string
	ServerPort int
	HostIP     string
	HostPort   int
	Namespace  string
	GroupName  string
}

type Executor struct {
	factory   LoggerFactory
	logger    Logger
	executors map[string]jobExecutor.IJobExecutor
	lock      sync.Mutex
}

func NewExecutor(opts *Options, factory LoggerFactory) *Executor {

	// factory *LoggerFactory = logger == nil ? NewLoggerFactory():logger
	return &Executor{
		factory:   factory,
		logger:    factory.GetLogger("executor"),
		executors: make(map[string]jobExecutor.IJobExecutor),
	}
}

func (e *Executor) GetLoggerFactory() LoggerFactory {
	return e.factory
}

func (e *Executor) Run() {
	e.logger.Info("%s", "")
}

func (e *Executor) Init() error {

	e.logger.Info("%s", "Init executor")
	return nil
}

// Register 注册job执行器入门
func (e *Executor) Register(name string, executor jobExecutor.IJobExecutor) {
	e.lock.Lock()
	defer e.lock.Unlock()

	if _, exists := e.executors[name]; exists {
		panic(fmt.Sprintf("Executor [%s] already registered", name))
	}

	e.executors[name] = executor
	job.LocalLog.Info(fmt.Sprintf("Registered executor: %s", name))
}

func (e *Executor) GetExecutor(name string) (jobExecutor.IJobExecutor, error) {
	executor, exists := e.executors[name]
	if !exists {
		return nil, fmt.Errorf("Executor [%s] not found", name)
	}
	return executor, nil
}
