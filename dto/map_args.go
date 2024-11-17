package dto

type MapArgs struct {
	JobArgs
	TaskName  string      `json:"task_name" description:"任务名称"`
	MapResult interface{} `json:"map_result" description:"分片结果"`
}
