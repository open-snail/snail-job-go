package demo

import (
	"opensnail.com/snail-job/snail-job-go/dto"
	"opensnail.com/snail-job/snail-job-go/job"
)

// TestWorkflowJobExecutor 这是一个测试类
type TestWorkflowJobExecutor struct {
	job.BaseJobExecutor
}

func (executor *TestWorkflowJobExecutor) DoJobExecute(jobArgs dto.IJobArgs) dto.ExecuteResult {
	executor.LocalLogger.Infof("TestWorkflowJobExecutor. jobId: [%d] wfContext:[%+v]",
		jobArgs.GetJobId(), jobArgs.GetWfContext("name"))
	jobArgs.AppendContext("name", "xiaowoniu")
	return *dto.Success("")
}
