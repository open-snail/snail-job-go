package dto

import (
	"github.com/sirupsen/logrus"
	"sync"
	"time"

	"opensnail.com/snail-job/snail-job-go/constant"
)

type Options struct {
	ServerHost string
	ServerPort string
	HostIP     string
	HostPort   string
	Namespace  string
	GroupName  string
	Token      string
	Level      logrus.Level
}

type DispatchJobResultRequest struct {
	TaskBatchId         int64                    `json:"taskBatchId"`
	GroupName           string                   `json:"groupName"`
	JobId               int64                    `json:"jobId"`
	TaskId              int64                    `json:"taskId"`
	WorkflowTaskBatchId int64                    `json:"workflowTaskBatchId"`
	WorkflowNodeId      int64                    `json:"workflowNodeId"`
	TaskType            constant.JobTaskTypeEnum `json:"taskType"`
	ExecuteResult       ExecuteResult            `json:"executeResult"`
	TaskStatus          int                      `json:"taskStatus"`
	Retry               bool                     `json:"retry"`
	RetryScene          int                      `json:"retryScene"`
	WfContext           string                   `json:"wfContext"`
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
	Status  int32       `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"result"`
}

type ShardingJobArgs struct {
	JobArgs
	ShardingTotal int `json:"sharding_total"`
	ShardingIndex int `json:"sharding_index"`
}

// SnailJobRequest 定义 SnailJobRequest 和 Metadata 的数据结构
type SnailJobRequest struct {
	ReqId    int64
	Body     string
	Metadata Metadata
}

type Metadata struct {
	Uri     string
	Headers map[string]string
}

type DispatchJobRequest struct {
	NamespaceId         string                      `json:"namespaceId" description:"namespaceId 不能为空"`
	JobId               int64                       `json:"jobId" description:"jobId 不能为空"`
	TaskBatchId         int64                       `json:"taskBatchId" description:"taskBatchId 不能为空"`
	TaskId              int64                       `json:"taskId" description:"taskId 不能为空"`
	TaskType            constant.JobTaskTypeEnum    `json:"taskType" description:"taskType 不能为空"`
	GroupName           string                      `json:"groupName" description:"group 不能为空"`
	ParallelNum         int                         `json:"parallelNum" description:"parallelNum 不能为空"`
	ExecutorType        int                         `json:"executorType" description:"executorType 不能为空"`
	ExecutorInfo        string                      `json:"executorInfo" description:"executorInfo 不能为空"`
	ExecutorTimeout     int                         `json:"executorTimeout" description:"executorTimeout 不能为空"`
	ArgsStr             string                      `json:"argsStr,omitempty"`
	ShardingTotal       int                         `json:"shardingTotal,omitempty"`
	ShardingIndex       int                         `json:"shardingIndex,omitempty"`
	WorkflowTaskBatchId int64                       `json:"workflowTaskBatchId,omitempty"`
	WorkflowNodeId      int64                       `json:"workflowNodeId,omitempty"`
	RetryCount          int                         `json:"retryCount,omitempty"`
	RetryScene          int                         `json:"retryScene,omitempty" description:"重试场景 auto、manual"`
	IsRetry             bool                        `json:"isRetry" description:"是否是重试流量"`
	WfContext           string                      `json:"wfContext" description:"工作流上下文"`
	TaskName            string                      `json:"taskName"`
	MrStage             constant.MapReduceStageEnum `json:"mrStage"`
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

type NettyResult struct {
	ReqID  int64               `json:"reqId"`
	Status constant.StatusEnum `json:"status"`
	Data   interface{}         `json:"data"`
}

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

type LogRecord struct {
	TimeStamp time.Time
	Level     string
	Thread    string
	Message   string
	Module    string
	FuncName  string
	Lineno    int
	ExcInfo   error
}
