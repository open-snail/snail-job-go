package demo

import (
	"snail_job_go/dto"
	"snail_job_go/executor"
)

// TestJobExecutor 这是一个测试类
type TestJobExecutor struct {
	executor.AbstractJobExecutor
}

func (e *TestJobExecutor) JobExecute(context dto.JobContext) {
	println("执行了", context.JobId)
}
