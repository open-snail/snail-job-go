package endpoint

import (
	"log"
	snailjob "opensnail.com/snail-job/snail-job-go"
	"opensnail.com/snail-job/snail-job-go/constant"
	"opensnail.com/snail-job/snail-job-go/dto"
	"opensnail.com/snail-job/snail-job-go/executor"
	"opensnail.com/snail-job/snail-job-go/job"
)

type JobEndPoint struct {
	manager *snailjob.Executor
}

func Init(manager *snailjob.Executor) *JobEndPoint {
	return &JobEndPoint{manager}
}

func (e *JobEndPoint) DispatchJob(dispatchJob job.DispatchJobRequest) dto.Result {

	if dispatchJob.IsRetry {
		job.LocalLog.Info("Task execution/scheduling failed, retrying. Retry count: [%d]", dispatchJob.IsRetry)
	}

	// 必须是GO客户端才能使用
	if dispatchJob.ExecutorType != 3 {
		log.Printf("Non-Go executor type not supported. ExecutorType: [%s]", dispatchJob.ExecutorInfo)
		return dto.Result{Status: 1, Message: "不支持非Java类型的执行器", Data: false}
	}

	jobExecute, _ := e.manager.GetExecutor(dispatchJob.ExecutorInfo)
	if jobExecute == nil {
		job.LocalLog.Info("Invalid executor configuration. ExecutorInfo: [%s]", dispatchJob.ExecutorInfo)
		return dto.Result{Status: 1, Message: "执行器配置有误", Data: false}
	}

	jobContext := buildJobContext(dispatchJob)

	jobStrategy := jobExecute.(executor.JobStrategy)
	jobStrategy.BindJobStrategy(jobStrategy)

	// bing executor
	if dispatchJob.TaskType == constant.MAP {
		mapExecute := jobExecute.(executor.MapExecute)
		mapExecute.BindMapExecute(mapExecute)
	} else if dispatchJob.TaskType == constant.MAP_REDUCE {
		mapExecute := jobExecute.(executor.MapExecute)
		mapExecute.BindMapExecute(mapExecute)
		mapReduceExecute := jobExecute.(executor.MapReduceExecute)
		mapReduceExecute.BindMapReduceExecute(mapReduceExecute)
	}

	// 集群、 广播、静态分片 直接执行方法
	jobExecute.JobExecute(jobContext)

	return dto.Result{Status: 1, Data: true}
}

func buildJobContext(dispatchJob job.DispatchJobRequest) dto.JobContext {
	jobContext := dto.JobContext{
		JobId:               dispatchJob.JobId,
		ShardingTotal:       dispatchJob.ShardingTotal,
		ShardingIndex:       dispatchJob.ShardingIndex,
		NamespaceId:         dispatchJob.NamespaceId,
		TaskId:              dispatchJob.TaskId,
		TaskBatchId:         dispatchJob.TaskBatchId,
		GroupName:           dispatchJob.GroupName,
		ExecutorInfo:        dispatchJob.ExecutorInfo,
		ParallelNum:         dispatchJob.ParallelNum,
		TaskType:            dispatchJob.TaskType,
		ExecutorTimeout:     dispatchJob.ExecutorTimeout,
		WorkflowNodeId:      dispatchJob.WorkflowNodeId,
		WorkflowTaskBatchId: dispatchJob.WorkflowTaskBatchId,
		IsRetry:             dispatchJob.IsRetry,
		RetryScene:          dispatchJob.RetryScene,
		TaskName:            dispatchJob.TaskName,
		MrStage:             constant.MapReduceStageEnum(dispatchJob.MrStage),
	}

	// Parse ArgsStr and WfContext (simplified example)
	if dispatchJob.ArgsStr != "" {
		jobContext.JobArgsHolder = dto.JobArgsHolder{JobParams: dispatchJob.ArgsStr}
	}

	if dispatchJob.WfContext != "" {
		jobContext.WfContext = make(map[string]interface{})
		// JSON parse here
	}

	return jobContext
}
