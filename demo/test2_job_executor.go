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

// Test2JobExecutor 这是一个测试类
type Test2JobExecutor struct {
	executor.AbstractJobExecutor
}

func (executor *Test2JobExecutor) DoJobExecute(jobArgs dto.IJobArgs) dto.ExecuteResult {
	job.LocalLog.Info(fmt.Sprintf("TestJobExecutor 开始执行 DoJobExecute. jobId: [%d] now:[%s]", jobArgs.GetJobId(), time.Now().String()))
	time.Sleep(3 * time.Second)
	job.LocalLog.Info(fmt.Sprintf("TestJobExecutor 执行结束 DoJobExecute. jobId: [%d] now:[%s]", jobArgs.GetJobId(), time.Now().String()))
	//panic("这是故意抛出的异常")
	num1 := 1
	num2 := 0
	num3 := num1 / num2
	fmt.Println(num3)
	return dto.ExecuteResult{}
}
