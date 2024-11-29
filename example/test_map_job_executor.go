package main

import (
	"time"

	"github.com/open-snail/snail-job-go/constant"

	"github.com/open-snail/snail-job-go/dto"
	"github.com/open-snail/snail-job-go/job"
)

// TestMapJobExecutor 这是一个测试类
type TestMapJobExecutor struct {
	job.BaseMapJobExecutor
}

func (executor *TestMapJobExecutor) DoJobMapExecute(mpArgs *dto.MapArgs) dto.ExecuteResult {
	taskList := []interface{}{1, 2, 3} // 示例任务列表
	logger := executor.LocalLogger
	if mpArgs.TaskName == constant.ROOT_MAP {
		_, _ = executor.DoMap(taskList, "secondTaskName")
		return *dto.Success()
	}

	logger.Infof("TestMapJobExecutor 开始执行 DoJobMapExecute. jobId: [%d] TaskName:[%s] ", mpArgs.GetJobId(), mpArgs.TaskName)
	time.Sleep(1 * time.Second)
	num1 := 1
	num2 := 1
	num3 := num1 / num2
	logger.Infof("TestMapJobExecutor 执行结束 DoJobMapExecute. jobId: [%d] TaskName:[%s] num3:[%d]", mpArgs.GetJobId(), mpArgs.TaskName, num3)
	return *dto.Success()
}
