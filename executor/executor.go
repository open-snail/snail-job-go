package executor

import (
	"snail_job_go/dto"
)

// IJobExecutor 执行器接口
type IJobExecutor interface {
	JobExecute(context dto.JobContext)
}
