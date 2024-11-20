package demo

import (
	"fmt"
	"opensnail.com/snail-job/snail-job-go/dto"
	"opensnail.com/snail-job/snail-job-go/executor"
	"opensnail.com/snail-job/snail-job-go/job"
	"time"
)

// TestMapJobExecutor 这是一个测试类
type TestMapJobExecutor struct {
	executor.AbstractMapJobExecutor
}

// NewTestMapJobExecutor 创建对象
func NewTestMapJobExecutor() *TestMapJobExecutor {
	executor := &TestMapJobExecutor{}
	//executor.BindMapExecute(executor)
	//executor.BindJobStrategy(executor)
	return executor
}

func (executor *TestMapJobExecutor) DoJobMapExecute(mpArgs *dto.MapArgs) dto.ExecuteResult {

	job.LocalLog.Info(fmt.Sprintf("TestMapJobExecutor 开始执行 DoJobMapExecute. jobId: [%d] TaskName:[%s] ", mpArgs.GetJobId(), mpArgs.TaskName))
	time.Sleep(1 * time.Second)
	job.LocalLog.Info(fmt.Sprintf("TestMapJobExecutor 执行结束 DoJobMapExecute. jobId: [%d]  TaskName:[%s]", mpArgs.GetJobId(), mpArgs.TaskName))
	num1 := 1
	num2 := 1
	num3 := num1 / num2
	fmt.Println(num3)
	return dto.ExecuteResult{}
}
