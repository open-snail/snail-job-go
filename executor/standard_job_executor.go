package executor

//
//import (
//	"snail_job_go/dto"
//	"snail_job_go/register"
//)
//
//type StandardJobExecutor struct {
//	AbstractJobExecutor
//	manager *register.ExecutorManager
//}
//
//// NewStandardJobExecutor 创建对象
//func NewStandardJobExecutor(manager *register.ExecutorManager) *StandardJobExecutor {
//	executor := &StandardJobExecutor{}
//	executor.BindJobStrategy(executor)
//	executor.manager = manager
//	return executor
//}
//
//// DoJobExecute 模板类
//func (executor *StandardJobExecutor) DoJobExecute(context dto.JobContext) dto.ExecuteResult {
//	// TODO 执行任务
//	jobExecute, _ := executor.manager.GetExecutor(context.ExecutorInfo)
//	jobExecute.JobExecute(dto.JobContext{JobId: 1})
//	return dto.ExecuteResult{}
//}
