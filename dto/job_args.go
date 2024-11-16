package dto

import "sync"

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
