package executor

import (
	"snail_job_go/dto"
	"snail_job_go/util"
)

type MapExecute interface {
	DoJobMapExecute(args *dto.MapArgs) dto.ExecuteResult
	BindMapExecute(child MapExecute)
}

type AbstractMapJobExecutor struct {
	AbstractJobExecutor
	Execute MapExecute
}

func (executor *AbstractMapJobExecutor) BindMapExecute(child MapExecute) {
	executor.Execute = child
}

// DoJobExecute 模板类
func (executor *AbstractMapJobExecutor) DoJobExecute(jobArgs dto.IJobArgs) dto.ExecuteResult {
	var mapArgs dto.MapArgs
	util.ToObj(util.ToByteArr(jobArgs), mapArgs)
	executor.Execute.DoJobMapExecute(&mapArgs)
	return dto.ExecuteResult{}
}
