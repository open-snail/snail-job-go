package job

import (
	"github.com/open-snail/snail-job-go/util"

	"github.com/open-snail/snail-job-go/constant"
	"github.com/open-snail/snail-job-go/dto"
	"github.com/sirupsen/logrus"
)

type JobExecutorFutureCallback struct {
	jobContext   dto.JobContext
	localLogger  *logrus.Entry
	remoteLogger *logrus.Entry
}

func (executor JobExecutorFutureCallback) onCallback(client SnailJobClient, result *dto.ExecuteResult) {
	executor.localLogger.Infof("Success result: %v", result)

	if result == nil {
		result = dto.Success()
	}

	var taskStatus constant.JobTaskStatusEnum
	if result.Status == 0 {
		taskStatus = constant.FAIL
	} else {
		taskStatus = constant.SUCCESS
	}

	request := buildDispatchJobResultRequest(result, taskStatus, executor.jobContext)
	if err := dispatchResult(client, request); err != nil {
		executor.remoteLogger.Errorf("Error reporting execution result: %v, TaskID: %d\n", err.Error(), executor.jobContext.TaskId)
		//sendMessage(err.Error())
	}
}

func dispatchResult(client SnailJobClient, req dto.DispatchJobResultRequest) error {
	client.log.Debugf("request server: %+v", req)
	client.SendDispatchResult(req)
	return nil
}

func buildDispatchJobResultRequest(result *dto.ExecuteResult, status constant.JobTaskStatusEnum, jobContext dto.JobContext) dto.DispatchJobResultRequest {
	var wfContext = ""
	if jobContext.ChangeWfContext != nil {
		j, _ := util.ToByteArr(jobContext.ChangeWfContext)
		wfContext = string(j)
	}
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
		RetryStatus:         jobContext.RetryStatus,
		RetryScene:          jobContext.RetryScene,
		WfContext:           wfContext,
	}
}
