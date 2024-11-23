package constant

const (
	VERSION         = "1.3.0"
	JOB_CONTEXT_KEY = "jobContext"
	INTERRUPT_KEY   = "taskInterrupt"
	ROOT_MAP        = "ROOT_MAP"
)

// StatusEnum 是响应状态的枚举
type StatusEnum int

const (
	NO = iota + 0
	YES
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
