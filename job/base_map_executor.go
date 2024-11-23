package job

import (
	"opensnail.com/snail-job/snail-job-go/dto"
	"opensnail.com/snail-job/snail-job-go/util"
)

type JobContextKeType int

const (
	JobContextKey JobContextKeType = iota
)

type MapExecute interface {
	DoJobMapExecute(args *dto.MapArgs) dto.ExecuteResult
	bindMapExecute(child MapExecute)
}

type BaseMapJobExecutor struct {
	BaseJobExecutor
	execute MapExecute
}

func (executor *BaseMapJobExecutor) bindMapExecute(child MapExecute) {
	executor.execute = child
}

// DoJobExecute 模板类
func (executor *BaseMapJobExecutor) DoJobExecute(jobArgs dto.IJobArgs) dto.ExecuteResult {
	var mapArgs dto.MapArgs
	util.ToObj(util.ToByteArr(jobArgs), mapArgs)
	return executor.execute.DoJobMapExecute(&mapArgs)
}
