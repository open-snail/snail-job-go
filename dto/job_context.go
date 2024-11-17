package dto

import "snail_job_go/constant"

type JobContext struct {
	JobId               int64                       `json:"jobId"`
	TaskBatchId         int64                       `json:"taskBatchId"`
	WorkflowTaskBatchId int64                       `json:"workflowTaskBatchId"`
	WorkflowNodeId      int64                       `json:"workflowNodeId"`
	TaskId              int64                       `json:"taskId"`
	NamespaceId         string                      `json:"namespaceId"`
	GroupName           string                      `json:"groupName"`
	ExecutorInfo        string                      `json:"executorInfo"`
	TaskType            constant.JobTaskTypeEnum    `json:"taskType"`
	ParallelNum         int                         `json:"parallelNum"`
	ShardingTotal       int                         `json:"shardingTotal"`
	ShardingIndex       int                         `json:"shardingIndex"`
	ExecutorTimeout     int                         `json:"executorTimeout"`
	RetryScene          int                         `json:"retryScene"` // 0=auto, 1=manual
	IsRetry             bool                        `json:"isRetry"`
	TaskList            []interface{}               `json:"taskList"`
	TaskName            string                      `json:"taskName"`
	MrStage             constant.MapReduceStageEnum `json:"mrStage"`
	WfContext           map[string]interface{}      `json:"wfContext"`
	ChangeWfContext     map[string]interface{}      `json:"changeWfContext"`
	JobArgsHolder       JobArgsHolder               `json:"jobArgsHolder"`
}
