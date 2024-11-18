package executor

import "snail_job_go/dto"

type MapExecute interface {
	DoJobMapExecute(args *dto.MapArgs) dto.ExecuteResult
	BindMapExecute(child MapExecute)
}

type AbstractMapJobExecutor struct {
	AbstractJobExecutor
	Execute MapExecute
}

func (executor *AbstractMapJobExecutor) BindMapExecute(child MapExecute) {
	//executor.BindJobStrategy(executor)
	executor.Execute = child
}

// DoJobExecute 模板类
func (executor *AbstractMapJobExecutor) DoJobExecute(jobArgs dto.IJobArgs) dto.ExecuteResult {
	// todo 怎么把jobArgs 转成 mapArgs
	executor.Execute.DoJobMapExecute(&dto.MapArgs{})
	return dto.ExecuteResult{}
}
