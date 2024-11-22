package demo

import (
	"fmt"
	"time"

	"opensnail.com/snail-job/snail-job-go/dto"
	"opensnail.com/snail-job/snail-job-go/job"
)

// TestMapJobExecutor 这是一个测试类
type TestMapJobExecutor struct {
	job.BaseMapJobExecutor
}

func (executor *TestMapJobExecutor) DoJobMapExecute(mpArgs *dto.MapArgs) dto.ExecuteResult {
	logger := executor.GetLogger()
	logger.Info("TestMapJobExecutor 开始执行 DoJobMapExecute. jobId: [%d] TaskName:[%s] ", mpArgs.GetJobId(), mpArgs.TaskName)
	time.Sleep(1 * time.Second)
	logger.Info("TestMapJobExecutor 执行结束 DoJobMapExecute. jobId: [%d]  TaskName:[%s]", mpArgs.GetJobId(), mpArgs.TaskName)
	num1 := 1
	num2 := 1
	num3 := num1 / num2
	fmt.Println(num3)
	return dto.ExecuteResult{}
}
