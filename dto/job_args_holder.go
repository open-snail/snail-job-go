package dto

type JobArgsHolder struct {
	JobParams interface{} `json:"job_params"`
	Maps      interface{} `json:"maps"`
	Reduces   interface{} `json:"reduces"`
}
