package demo

import (
	"fmt"
	"time"

	"github.com/open-snail/snail-job-go/dto"
	"github.com/open-snail/snail-job-go/job"
)

// Test2JobExecutor 这是一个测试类
type Test2JobExecutor struct {
	job.BaseJobExecutor
}

func (executor *Test2JobExecutor) DoJobExecute(jobArgs dto.IJobArgs) dto.ExecuteResult {
	executor.LocalLogger.Infof("TestJobExecutor 开始执行 DoJobExecute. jobId: [%d] now:[%s]", jobArgs.GetJobId(), time.Now().String())
	//time.Sleep(3 * time.Second)
	executor.RemoteLogger.Infof("TestJobExecutor 执行结束 DoJobExecute. jobId: [%d] now:[%s]", jobArgs.GetJobId(), time.Now().String())
	//panic("这是故意抛出的异常")
	num1 := 1
	num2 := 1
	num3 := num1 / num2
	fmt.Println(num3)
	return *dto.Success("hello 这是go客户端")
}
