package constant

const (
	SNAIL_SERVER_HOST = "127.0.0.1"
	SNAIL_SERVER_PORT = "1788"
	SNAIL_HOST_IP     = "127.0.0.1"
	SNAIL_HOST_PORT   = "1789"
	SNAIL_NAMESPACE   = "764d604ec6fc45f68cd92514c40e9e1a"
	SNAIL_GROUP_NAME  = "snail_job_demo_group"

	SNAIL_LOG_LOCAL_FILENAME     = "snail_job.log"
	SNAIL_LOG_REMOTE_BUFFER_SIZE = 10
	SNAIL_LOG_REMOTE_INTERVAL    = 10
)

// JobTaskTypeEnum 定义 JobTaskTypeEnum 枚举类型
type JobTaskTypeEnum int

const (
	CLUSTER JobTaskTypeEnum = iota + 1
	BROADCAST
	SHARDING
	MAP
	MAP_REDUCE
)

type MapReduceStageEnum int

const (
	MAP_STAGE MapReduceStageEnum = iota + 1
	REDUCE_STAGE
	MERGE_REDUCE_STAGE
	// Other stages...
)

type JobTaskStatusEnum int

const (
	RUNNING JobTaskStatusEnum = iota + 2
	SUCCESS
	FAIL
	STOP
	CANCEL
)
