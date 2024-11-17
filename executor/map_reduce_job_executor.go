package executor

type MapReduceJobExecutor struct {
	AbstractJobExecutor
	//manager *register.ExecutorManager
}

//// NewMapJobExecutor 创建对象
//func NewMapReduceJobExecutor(manager *register.ExecutorManager) *MapReduceJobExecutor {
//	executor := &MapReduceJobExecutor{}
//	executor.BindJobStrategy(executor)
//	executor.manager = manager
//	return executor
//}
//
//// DoJobExecute 模板类
//func (executor *MapReduceJobExecutor) DoJobExecute(context dto.JobContext) dto.ExecuteResult {
//	return dto.ExecuteResult{}
//}
