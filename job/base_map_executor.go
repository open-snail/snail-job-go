package job

import (
	"opensnail.com/snail-job/snail-job-go/dto"
	"opensnail.com/snail-job/snail-job-go/util"
)

type MapExecute interface {
	DoJobMapExecute(args *dto.MapArgs) dto.ExecuteResult
	BindMapExecute(child MapExecute)
}

type BaseMapJobExecutor struct {
	BaseJobExecutor
	Execute MapExecute
}

func (executor *BaseMapJobExecutor) BindMapExecute(child MapExecute) {
	executor.Execute = child
}

// DoJobExecute 模板类
func (executor *BaseMapJobExecutor) DoJobExecute(jobArgs dto.IJobArgs) dto.ExecuteResult {
	var mapArgs dto.MapArgs
	util.ToObj(util.ToByteArr(jobArgs), mapArgs)
	executor.Execute.DoJobMapExecute(&mapArgs)
	return dto.ExecuteResult{}
}
