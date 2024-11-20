package dto

import (
	"snail_job_go/constant"
	"sync"
)

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

type ExecuteResult struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

func Success(result interface{}) *ExecuteResult {
	return &ExecuteResult{1, "任务执行成功", result}
}

func Failure(result interface{}, msg string) *ExecuteResult {
	return &ExecuteResult{0, msg, result}
}

type JobArgsHolder struct {
	JobParams interface{} `json:"job_params"`
	Maps      interface{} `json:"maps"`
	Reduces   interface{} `json:"reduces"`
}

type IJobArgs interface {
	GetJobParams() interface{}
	GetExecutorInfo() string
	GetJobId() int64
	GetTaskBatchId() int64
}

// JobArgs job参数
type JobArgs struct {
	JobParams       interface{}            `json:"jobParams"`
	ExecutorInfo    string                 `json:"executorInfo"`
	TaskBatchId     int64                  `json:"taskBatchId"`
	JobId           int64                  `json:"jobId"`
	WfContext       map[string]interface{} `json:"wfContext"`
	ChangeWfContext map[string]interface{} `json:"changeWfContext"`
	mu              sync.Mutex
}

// AppendContext 工作流场景下 添加上下文参数
func (j *JobArgs) AppendContext(key string, value interface{}) {
	if j.WfContext == nil || key == "" || value == nil {
		return
	}

	j.mu.Lock()
	defer j.mu.Unlock()

	if j.ChangeWfContext == nil {
		j.ChangeWfContext = make(map[string]interface{})
	}
	j.ChangeWfContext[key] = value
}

// GetWfContext 获取工作流上下文
func (j *JobArgs) GetWfContext(key string) interface{} {
	if j.WfContext == nil || key == "" {
		return nil
	}

	j.mu.Lock()
	defer j.mu.Unlock()

	return j.WfContext[key]
}

func (j *JobArgs) GetJobParams() interface{} {
	if j.JobParams == nil {
		return nil
	}
	return j.JobParams
}

func (j *JobArgs) GetExecutorInfo() string {
	if j.ExecutorInfo == "" {
		return ""
	}
	return j.ExecutorInfo
}

func (j *JobArgs) GetTaskBatchId() int64 {
	if j.TaskBatchId == 0 {
		return 0
	}

	return j.TaskBatchId
}

func (j *JobArgs) GetJobId() int64 {
	if j.JobId == 0 {
		return 0
	}

	return j.JobId
}

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

type MapArgs struct {
	JobArgs
	TaskName  string      `json:"task_name" description:"任务名称"`
	MapResult interface{} `json:"map_result" description:"分片结果"`
}

type MergeReduceArgs struct {
	JobArgs
	Reduces []interface{} `json:"reduces" description:"reduce参数"`
}

type ReduceArgs struct {
	JobArgs
	MapResult []interface{} `json:"mapResult" description:"mapResult结果"`
}

type Result struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"result"`
}

type ShardingJobArgs struct {
	JobArgs
	ShardingTotal int `json:"sharding_total"`
	ShardingIndex int `json:"sharding_index"`
}
