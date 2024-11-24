package demo

import (
	"opensnail.com/snail-job/snail-job-go/constant"
	"time"

	"opensnail.com/snail-job/snail-job-go/dto"
	"opensnail.com/snail-job/snail-job-go/job"
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
		return *dto.Success(nil)
	}

	logger.Info("TestMapJobExecutor 开始执行 DoJobMapExecute. jobId: [%d] TaskName:[%s] ", mpArgs.GetJobId(), mpArgs.TaskName)
	time.Sleep(1 * time.Second)
	num1 := 1
	num2 := 1
	num3 := num1 / num2
	logger.Info("TestMapJobExecutor 执行结束 DoJobMapExecute. jobId: [%d] TaskName:[%s] num3:[%s]", mpArgs.GetJobId(), mpArgs.TaskName, num3)
	return *dto.Success("")
}
