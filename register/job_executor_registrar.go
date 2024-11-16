package register

import (
	"fmt"
	jobExecutor "snail_job_go/executor"
	"snail_job_go/job"
	"sync"
)

type ExecutorManager struct {
	hls       *job.HookLogService
	executors map[string]jobExecutor.IJobExecutor
	//executorPools map[int]*sync.Pool
	lock sync.Mutex
}

func NewExecutorManager(hls *job.HookLogService) *ExecutorManager {
	return &ExecutorManager{
		hls:       hls,
		executors: make(map[string]jobExecutor.IJobExecutor),
		//executorPools: make(map[int]*sync.Pool),
	}
}

// Register 注册job执行器入门
func (m *ExecutorManager) Register(name string, executor jobExecutor.IJobExecutor) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if _, exists := m.executors[name]; exists {
		panic(fmt.Sprintf("Executor [%s] already registered", name))
	}

	m.executors[name] = executor
	job.LocalLog.Info(fmt.Sprintf("Registered executor: %s", name))
}

func (m *ExecutorManager) GetExecutor(name string) (jobExecutor.IJobExecutor, error) {
	executor, exists := m.executors[name]
	if !exists {
		return nil, fmt.Errorf("Executor [%s] not found", name)
	}
	return executor, nil
}
