package dto

type ReduceArgs struct {
	JobArgs
	MapResult []interface{} `json:"mapResult" description:"mapResult结果"`
}
