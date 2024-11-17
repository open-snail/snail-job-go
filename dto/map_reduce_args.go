package dto

type MergeReduceArgs struct {
	JobArgs
	Reduces []interface{} `json:"reduces" description:"reduce参数"`
}
