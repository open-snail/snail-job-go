package demo

import (
	"fmt"
	"snail_job_go/dto"
	"snail_job_go/executor"
	"snail_job_go/job"
	"time"
)

/**
  todo 待讨论Go的是支持结构体还是直接执行方法方式
*/

// TestMapJobExecutor 这是一个测试类
type TestMapJobExecutor struct {
	executor.AbstractMapJobExecutor
}

// NewTestMapJobExecutor 创建对象
func NewTestMapJobExecutor() *TestMapJobExecutor {
	executor := &TestMapJobExecutor{}
	executor.BindMapExecute(executor)
	executor.BindJobStrategy(executor)
	return executor
}

func (executor *TestMapJobExecutor) DoJobMapExecute(mpArgs *dto.MapArgs) dto.ExecuteResult {

	job.LocalLog.Info(fmt.Sprintf("TestMapJobExecutor 开始执行 DoJobMapExecute. jobId: [%d] TaskName:[%s] now:[%s]", mpArgs.GetJobId(), mpArgs.TaskName, time.Now().String()))
	time.Sleep(1 * time.Second)
	job.LocalLog.Info(fmt.Sprintf("TestMapJobExecutor 执行结束 DoJobMapExecute. jobId: [%d]  TaskName:[%s] now:[%s]", mpArgs.GetJobId(), mpArgs.TaskName, time.Now().String()))
	//panic("这是故意抛出的异常")
	num1 := 1
	num2 := 1
	num3 := num1 / num2
	fmt.Println(num3)
	return dto.ExecuteResult{}
}
