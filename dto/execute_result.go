package dto

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
