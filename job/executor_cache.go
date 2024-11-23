package job

type executorCache struct {
	executors map[int64][]IJobExecutor
}

func NewExecutorCache() *executorCache {
	return &executorCache{
		executors: make(map[int64][]IJobExecutor),
	}
}

func (receiver executorCache) register(jobBatchId int64, prototype IJobExecutor) {
	if _, ok := receiver.executors[jobBatchId]; !ok {
		receiver.executors[jobBatchId] = []IJobExecutor{prototype}
	} else {
		receiver.executors[jobBatchId] = append(receiver.executors[jobBatchId], prototype)
	}
}
