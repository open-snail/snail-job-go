package job

type executorCache struct {
	executors map[int64][]JobStrategy
}

func NewExecutorCache() *executorCache {
	return &executorCache{
		executors: make(map[int64][]JobStrategy),
	}
}

func (receiver executorCache) register(jobBatchId int64, prototype JobStrategy) {
	if _, ok := receiver.executors[jobBatchId]; !ok {
		receiver.executors[jobBatchId] = []JobStrategy{prototype}
	} else {
		receiver.executors[jobBatchId] = append(receiver.executors[jobBatchId], prototype)
	}
}
