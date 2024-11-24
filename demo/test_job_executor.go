package demo

import (
	"fmt"
	"time"

	"opensnail.com/snail-job/snail-job-go/constant"

	"opensnail.com/snail-job/snail-job-go/dto"
	"opensnail.com/snail-job/snail-job-go/job"
)

// TestJobExecutor 这是一个测试类
type TestJobExecutor struct {
	job.BaseJobExecutor
}

func (executor *TestJobExecutor) DoJobExecute(jobArgs dto.IJobArgs) dto.ExecuteResult {

	time.Sleep(1 * time.Second)
	interrupt := executor.Context().Value(constant.INTERRUPT_KEY)
	if interrupt != nil {
		executor.RemoteLog.Errorf("任务被中断. jobId: [%d] now:[%s]", jobArgs.GetJobId(), time.Now().String())
		return *dto.Failure(nil, "任务被中断")
	}

	num1 := 1
	num2 := 1
	num3 := num1 / num2
	fmt.Println(num3)
	executor.LocalLog.Infof("任务执行结束. jobId: [%d] now:[%s]", jobArgs.GetJobId(), time.Now().String())
	executor.RemoteLog.Infof("任务执行结束. jobId: [%d] now:[%s]", jobArgs.GetJobId(), time.Now().String())
	return *dto.Success("hello 这是go客户端")
}
