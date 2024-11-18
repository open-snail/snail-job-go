package executor

import "snail_job_go/dto"

type MapReduceExecute interface {
	DoReduceExecute(args *dto.ReduceArgs) dto.ExecuteResult
	DoMergeReduceExecute(args *dto.MergeReduceArgs) dto.ExecuteResult
	BindMapReduceExecute(child MapReduceExecute)
}

type AbstractMapReduceJobExecutor struct {
	AbstractMapJobExecutor
	MrExecute MapReduceExecute
}

func (executor *AbstractMapReduceJobExecutor) BindMapReduceExecute(child MapReduceExecute) {
	//executor.BindMapExecute(executor)
	executor.MrExecute = child
}

// DoJobExecute 模板类
func (executor *AbstractMapReduceJobExecutor) DoJobExecute(jobArgs dto.IJobArgs) dto.ExecuteResult {
	// todo 怎么把jobArgs 转成 mapArgs
	executor.Execute.DoJobMapExecute(&dto.MapArgs{})
	executor.MrExecute.DoReduceExecute(&dto.ReduceArgs{})
	executor.MrExecute.DoMergeReduceExecute(&dto.MergeReduceArgs{})
	return dto.ExecuteResult{}
}

func (executor *AbstractMapReduceJobExecutor) DoJobMapExecute(args *dto.MapArgs) dto.ExecuteResult {
	return executor.Execute.DoJobMapExecute(args)
}
