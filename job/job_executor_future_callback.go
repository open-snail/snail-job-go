package job

import (
	"encoding/json"
	"fmt"
	"log"
	"opensnail.com/snail-job/snail-job-go/constant"
	"opensnail.com/snail-job/snail-job-go/dto"
)

type JobExecutorFutureCallback struct {
	jobContext dto.JobContext
}

func (executor JobExecutorFutureCallback) onCallback(client SnailJobClient, result *dto.ExecuteResult) {

	// todo 这里要改成Remote日志
	LocalLog.Info(fmt.Sprintf("Success result: %v", result))

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
	if err := dispatchResult(client, request); err != nil {
		log.Printf("Error reporting execution result: %s, TaskID: %s\n", err.Error(), executor.jobContext.TaskId)
		//sendMessage(err.Error())
	}
}

func dispatchResult(client SnailJobClient, req dto.DispatchJobResultRequest) error {
	LocalLog.Info(fmt.Sprintf("request server: %v", req))
	client.SendDispatchResult(req)
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
