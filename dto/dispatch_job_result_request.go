package dto

import "snail_job_go/constant"

type DispatchJobResultRequest struct {
	TaskBatchId         int64
	GroupName           string
	JobId               int64
	TaskId              int64
	WorkflowTaskBatchId int64
	WorkflowNodeId      int64
	TaskType            constant.JobTaskTypeEnum
	ExecuteResult       ExecuteResult
	TaskStatus          int
	Retry               bool
	RetryScene          int
	WfContext           string
}
