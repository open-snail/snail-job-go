package job

import (
	"encoding/json"
	"fmt"
	"log"

	"opensnail.com/snail-job/snail-job-go/constant"
	"opensnail.com/snail-job/snail-job-go/dto"
	"opensnail.com/snail-job/snail-job-go/util"
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

func (executor *BaseMapJobExecutor) DoMap(taskList []interface{}, nextTaskName string) (*dto.ExecuteResult, error) {
	logger := executor.RemoteLogger

	// 检查 nextTaskName
	if nextTaskName == "" {
		logger.Error("The next task name can not blank or null {%s}", nextTaskName)
		return dto.Failure(nil, ""), fmt.Errorf("the next task name can not blank or null {%s}", nextTaskName)
	}

	// 检查 taskList
	if len(taskList) == 0 {
		logger.Error("The next task name can not blank or null {%s}", nextTaskName)
		return dto.Failure(nil, ""), fmt.Errorf("the next task name can not blank or null {%s}", nextTaskName)
	}

	// 检查任务数量
	if len(taskList) > 200 {
		logger.Warn("[%s] map task size is too large, network maybe overload... please try to split the tasks.\n", nextTaskName)
	}

	if len(taskList) > 500 {
		return dto.Failure(nil, ""), fmt.Errorf("[%s] map task size is too large, network maybe overload... please try to split the tasks", nextTaskName)
	}

	// 检查任务名是否为 ROOT_MAP
	if nextTaskName == constant.ROOT_MAP {
		logger.Error("The Next taskName cannot be %s", "ROOT_MAP")
		return dto.Failure(nil, ""), fmt.Errorf("the Next taskName cannot be %s", "ROOT_MAP")
	}

	jobContext := executor.ctx.Value(constant.JOB_CONTEXT_KEY).(dto.JobContext)

	// 构造 MapTaskRequest
	mapTaskRequest := dto.MapTaskRequest{
		JobId:               jobContext.JobId,
		TaskBatchId:         jobContext.TaskBatchId,
		TaskName:            nextTaskName,
		SubTask:             taskList,
		ParentId:            jobContext.TaskId,
		WorkflowTaskBatchId: jobContext.WorkflowTaskBatchId,
		WorkflowNodeId:      jobContext.WorkflowNodeId,
	}

	if changeWfContext := jobContext.ChangeWfContext; changeWfContext != nil {
		contextJson, err := json.Marshal(changeWfContext)
		if err != nil {
			log.Fatal(err)
		}
		mapTaskRequest.WfContext = string(contextJson)
	}

	status := executor.client.SendBatchReportMapTask(mapTaskRequest)
	if status == constant.NO {
		logger.Error("map failed for task: %s", nextTaskName)
		return dto.Failure(nil, ""), fmt.Errorf("map failed for task: %s", nextTaskName)
	} else {
		logger.Info("Map task create successfully!. taskName:[%s] TaskId:[%d]", nextTaskName, jobContext.TaskId)
	}

	return dto.Success("分片成功"), nil

}

// DoJobExecute 模板类
func (executor *BaseMapJobExecutor) DoJobExecute(jobArgs dto.IJobArgs) dto.ExecuteResult {
	if executor.execute == nil {
		return *dto.Failure(nil, "执行器类型错误")
	}

	var mapArgs dto.MapArgs
	util.ToObj(util.ToByteArr(jobArgs), &mapArgs)
	return executor.execute.DoJobMapExecute(&mapArgs)
}
