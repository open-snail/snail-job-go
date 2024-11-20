package snailjob

import (
	"fmt"
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

// TaskFunc 任务执行函数
type TaskFunc func() (string, error)

type executor struct {
	factory   LoggerFactory
	logger    Logger
	executors map[string]TaskFunc
	lock      sync.Mutex
}

func NewExecutor(opts *Options, factory LoggerFactory) *executor {

	// factory *LoggerFactory = logger == nil ? NewLoggerFactory():logger
	return &executor{
		factory:   factory,
		logger:    factory.GetLogger("executor"),
		executors: make(map[string]TaskFunc),
	}
}

func (e *executor) GetLoggerFactory() LoggerFactory {
	return e.factory
}

func (e *executor) Run() {

}

func (e *executor) Init() error {

	e.logger.Info("%s", "Init executor")
	return nil
}

func (e *executor) Register(name string, fun TaskFunc) {
	e.lock.Lock()
	defer e.lock.Unlock()

	if _, exists := e.executors[name]; exists {
		panic(fmt.Sprintf("Executor [%s] already registered", name))
	}

	e.executors[name] = fun
	e.logger.Info("Registered executor: %s", name)
}
