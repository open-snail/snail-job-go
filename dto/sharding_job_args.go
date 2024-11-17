package dto

type ShardingJobArgs struct {
	JobArgs
	ShardingTotal int `json:"sharding_total"`
	ShardingIndex int `json:"sharding_index"`
}
