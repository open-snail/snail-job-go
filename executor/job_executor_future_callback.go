package executor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"snail_job_go/constant"
	"snail_job_go/dto"
	"snail_job_go/job"
)

type JobExecutorFutureCallback struct {
	jobContext dto.JobContext
}

func (executor JobExecutorFutureCallback) onCallback(result *dto.ExecuteResult) {

	// todo 这里要改成Remote日志
	job.LocalLog.Info(fmt.Sprintf("Success result: %v", result))

	if result == nil {
		result = dto.Success(nil)
	}

	var taskStatus constant.JobTaskStatusEnum
	if result.Status == 0 {
		taskStatus = constant.FAIL
	} else {
		taskStatus = constant.SUCCESS
	}

	request := buildDispatchJobResultRequest(result, taskStatus, executor.jobContext)
	if err := dispatchResult(request); err != nil {
		log.Printf("Error reporting execution result: %s, TaskID: %s\n", err.Error(), executor.jobContext.TaskId)
		//sendMessage(err.Error())
	}
}

func (executor JobExecutorFutureCallback) onFailure(t error) {
	if errors.Is(t, context.Canceled) {
		job.LocalLog.Info(fmt.Sprintf("Task has been canceled, not reporting status"))
		return
	}

	failure := dto.Failure(nil, t.Error())

	log.Printf("Task execution failed TaskBatchId: %s, Error: %s\n", executor.jobContext.TaskBatchId, t.Error())
	request := buildDispatchJobResultRequest(failure, constant.FAIL, executor.jobContext) // JobTaskStatusEnum.FAIL

	if err := dispatchResult(request); err != nil {
		job.LocalLog.Info(fmt.Sprintf("Error reporting execution result: %s, TaskID: %s\n", err.Error(), executor.jobContext.TaskId))
		//sendMessage(err.Error())
	}
}

func dispatchResult(req dto.DispatchJobResultRequest) error {
	// todo 请求服务端
	job.LocalLog.Info(fmt.Sprintf("request server: %v", req))
	return nil
}

func buildDispatchJobResultRequest(result *dto.ExecuteResult, status constant.JobTaskStatusEnum, jobContext dto.JobContext) dto.DispatchJobResultRequest {
	wfContext, _ := json.Marshal(jobContext.ChangeWfContext)
	return dto.DispatchJobResultRequest{
		TaskBatchId:         jobContext.TaskBatchId,
		GroupName:           jobContext.GroupName,
		JobId:               jobContext.JobId,
		TaskId:              jobContext.TaskId,
		WorkflowTaskBatchId: jobContext.WorkflowTaskBatchId,
		WorkflowNodeId:      jobContext.WorkflowNodeId,
		TaskType:            jobContext.TaskType,
		ExecuteResult:       *result,
		TaskStatus:          int(status),
		Retry:               jobContext.IsRetry,
		RetryScene:          jobContext.RetryScene,
		WfContext:           string(wfContext),
	}
}
