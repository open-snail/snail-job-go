package job

import (
	"snail_job_go/constant"
	"sync"
	"time"
)

type DispatchJobResult struct {
	JobId               int                      `json:"jobId"`
	TaskBatchId         int                      `json:"taskBatchId"`
	WorkflowTaskBatchId *int                     `json:"workflowTaskBatchId,omitempty"`
	WorkflowNodeId      *int                     `json:"workflowNodeId,omitempty"`
	TaskId              int                      `json:"taskId"`
	TaskType            constant.JobTaskTypeEnum `json:"taskType"`
	GroupName           string                   `json:"groupName"`
	TaskStatus          JobTaskBatchStatusEnum   `json:"taskStatus"`
	ExecuteResult       ExecuteResult            `json:"executeResult"`
	RetryScene          *int                     `json:"retryScene,omitempty"`
	IsRetry             bool                     `json:"isRetry"`
}

type JobTaskBatchStatusEnum int

const (
	BATCH_STATUS_WAITING JobTaskBatchStatusEnum = iota + 1
	BATCH_STATUS_RUNNING
	BATCH_STATUS_SUCCESS
	BATCH_STATUS_FAIL
	BATCH_STATUS_STOP
	BATCH_STATUS_CANCEL
)

type DispatchJobRequest struct {
	NamespaceId         string                   `json:"namespaceId" description:"namespaceId 不能为空"`
	JobId               int64                    `json:"jobId" description:"jobId 不能为空"`
	TaskBatchId         int64                    `json:"taskBatchId" description:"taskBatchId 不能为空"`
	TaskId              int64                    `json:"taskId" description:"taskId 不能为空"`
	TaskType            constant.JobTaskTypeEnum `json:"taskType" description:"taskType 不能为空"`
	GroupName           string                   `json:"groupName" description:"group 不能为空"`
	ParallelNum         int                      `json:"parallelNum" description:"parallelNum 不能为空"`
	ExecutorType        int                      `json:"executorType" description:"executorType 不能为空"`
	ExecutorInfo        string                   `json:"executorInfo" description:"executorInfo 不能为空"`
	ExecutorTimeout     int                      `json:"executorTimeout" description:"executorTimeout 不能为空"`
	ArgsStr             string                   `json:"argsStr,omitempty"`
	ShardingTotal       int                      `json:"shardingTotal,omitempty"`
	ShardingIndex       int                      `json:"shardingIndex,omitempty"`
	WorkflowTaskBatchId int64                    `json:"workflowTaskBatchId,omitempty"`
	WorkflowNodeId      int64                    `json:"workflowNodeId,omitempty"`
	RetryCount          int                      `json:"retryCount,omitempty"`
	RetryScene          int                      `json:"retryScene,omitempty" description:"重试场景 auto、manual"`
	IsRetry             bool                     `json:"isRetry" description:"是否是重试流量"`
	WfContext           string                   `json:"wfContext" description:"工作流上下文"`
	TaskName            string                   `json:"taskName"`
	MrStage             int                      `json:"mrStage"`
}

type DispatchJobArgs struct {
	NamespaceId         string                   `json:"namespaceId" description:"namespaceId 不能为空"`
	JobId               int                      `json:"jobId" description:"jobId 不能为空"`
	TaskBatchId         int                      `json:"taskBatchId" description:"taskBatchId 不能为空"`
	TaskId              int                      `json:"taskId" description:"taskId 不能为空"`
	TaskType            constant.JobTaskTypeEnum `json:"taskType" description:"taskType 不能为空"`
	GroupName           string                   `json:"groupName" description:"group 不能为空"`
	ParallelNum         int                      `json:"parallelNum" description:"parallelNum 不能为空"`
	ExecutorType        int                      `json:"executorType" description:"executorType 不能为空"`
	ExecutorInfo        string                   `json:"executorInfo" description:"executorInfo 不能为空"`
	ExecutorTimeout     int                      `json:"executorTimeout" description:"executorTimeout 不能为空"`
	ArgsStr             *string                  `json:"argsStr,omitempty"`
	ShardingTotal       *int                     `json:"shardingTotal,omitempty"`
	ShardingIndex       *int                     `json:"shardingIndex,omitempty"`
	WorkflowTaskBatchId *int                     `json:"workflowTaskBatchId,omitempty"`
	WorkflowNodeId      *int                     `json:"workflowNodeId,omitempty"`
	RetryCount          *int                     `json:"retryCount,omitempty"`
	RetryScene          *int                     `json:"retryScene,omitempty" description:"重试场景 auto、manual"`
	IsRetry             bool                     `json:"isRetry" description:"是否是重试流量"`
}

type StopJobRequest struct {
	ReqID int64 `json:"reqId"`
	Args  []struct {
		TaskBatchID int `json:"taskBatchId"`
	} `json:"args"`
}

type ExecuteResult struct {
	Success StatusEnum `json:"success"`
	Message string     `json:"message"`
}

type SnailJobRequest struct {
	ReqID int64       `json:"reqId"`
	Args  interface{} `json:"args"`
}

type NettyResult struct {
	ReqID  int64       `json:"reqId"`
	Status StatusEnum  `json:"status"`
	Data   interface{} `json:"data"`
}

type StatusEnum int

const (
	STATUS_FAILED  StatusEnum = 0
	STATUS_SUCCESS StatusEnum = 1
)

type SnailHttpLogHandler struct {
	mu       sync.Mutex
	capacity int
	interval time.Duration
	buffer   chan *JobLogTask
	timer    *time.Timer
}

type TaskLogFieldDTO struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// 上报服务器日志结构
type JobLogTask struct {
	LogType     string            `json:"logType"`
	NamespaceID string            `json:"namespaceId"`
	GroupName   string            `json:"groupName"`
	RealTime    int64             `json:"realTime"`
	FieldList   []TaskLogFieldDTO `json:"fieldList"`
	JobID       int               `json:"jobId"`
	TaskBatchID int               `json:"taskBatchId"`
	TaskID      int               `json:"taskId"`
}
