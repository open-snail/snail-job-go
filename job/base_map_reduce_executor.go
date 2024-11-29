package job

import (
	"github.com/open-snail/snail-job-go/constant"
	"github.com/open-snail/snail-job-go/dto"
	"github.com/open-snail/snail-job-go/util"
)

type MapReduceExecute interface {
	DoReduceExecute(args *dto.ReduceArgs) dto.ExecuteResult
	DoMergeReduceExecute(args *dto.MergeReduceArgs) dto.ExecuteResult
	BindMapReduceExecute(child MapReduceExecute)
}

type BaseMapReduceJobExecutor struct {
	BaseMapJobExecutor
	mrExecute MapReduceExecute
}

func (executor *BaseMapReduceJobExecutor) BindMapReduceExecute(child MapReduceExecute) {
	executor.mrExecute = child
}

// DoJobExecute 模板类
func (executor *BaseMapReduceJobExecutor) DoJobExecute(jobArgs dto.IJobArgs) dto.ExecuteResult {
	// 将 User 转换为 JSON
	var mapArgs dto.MapArgs
	arr, toByteArrErr := util.ToByteArr(jobArgs)
	if toByteArrErr != nil {
		return *dto.Failure().WithMessage("参数解析错误")
	}

	toObjErr := util.ToObj(arr, &mapArgs)
	if toObjErr != nil {
		return *dto.Failure().WithMessage("参数解析错误")
	}

	jobContext := executor.ctx.Value(constant.JOB_CONTEXT_KEY).(dto.JobContext)
	mrStage := jobContext.MrStage
	if mrStage == constant.MAP_STAGE {
		return executor.execute.DoJobMapExecute(&mapArgs)
	} else if mrStage == constant.REDUCE_STAGE {
		return executor.mrExecute.DoReduceExecute(&dto.ReduceArgs{})
	} else {
		return executor.mrExecute.DoMergeReduceExecute(&dto.MergeReduceArgs{})
	}
}

func (executor *BaseMapReduceJobExecutor) DoJobMapExecute(args *dto.MapArgs) dto.ExecuteResult {
	return executor.execute.DoJobMapExecute(args)
}
