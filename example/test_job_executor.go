package main

import (
	"fmt"
	"time"

	"github.com/open-snail/snail-job-go/constant"

	"github.com/open-snail/snail-job-go/dto"
	"github.com/open-snail/snail-job-go/job"
)

// TestJobExecutor 这是一个测试类
type TestJobExecutor struct {
	job.BaseJobExecutor
}

// 测试超时时间
func (executor *TestJobExecutor) DoJobExecute(jobArgs dto.IJobArgs) dto.ExecuteResult {

	time.Sleep(1 * time.Second)
	interrupt := executor.Context().Value(constant.INTERRUPT_KEY)
	if interrupt != nil {
		executor.LocalLogger.Errorf("任务被中断. jobId: [%d] now:[%s]", jobArgs.GetJobId(), time.Now().String())
		return *dto.Failure().WithMessage("任务被中断")
	}

	for i := 0; i < 10; i++ {
		executor.RemoteLogger.Infof("Sleeping for a while, batchId: [%d] now:[%s]", jobArgs.GetTaskBatchId(), time.Now().String())
		time.Sleep(3 * time.Second)
	}

	num1 := 1
	num2 := 1
	num3 := num1 / num2
	fmt.Println(num3)
	executor.LocalLogger.Infof("任务执行结束. jobId: [%d] now:[%s]", jobArgs.GetJobId(), time.Now().String())
	executor.RemoteLogger.Infof("任务执行结束. jobId: [%d] now:[%s]", jobArgs.GetJobId(), time.Now().String())
	return *dto.Success().WithMessage("hello 这是go客户端")
}
