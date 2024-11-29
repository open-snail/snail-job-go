package job

import (
	"context"
	"encoding/json"
	"time"

	"github.com/open-snail/snail-job-go/constant"
	"github.com/open-snail/snail-job-go/dto"
	"github.com/sirupsen/logrus"
)

type NewJobExecutor func() IJobExecutor

// IJobExecutor 执行器接口
type IJobExecutor interface {
	JobExecute(context dto.JobContext)
}

type JobStrategy interface {
	DoJobExecute(dto.IJobArgs) dto.ExecuteResult
	bindJobStrategy(child JobStrategy)
	setClient(client SnailJobClient)
	setContext(ctx context.Context)
	getContext() context.Context
	setLogger(localLogger *logrus.Entry, remoteLogger *logrus.Entry)
	setExecutorCache(execCache executorCache)
}

type BaseJobExecutor struct {
	strategy     JobStrategy
	client       SnailJobClient
	ctx          context.Context
	LocalLogger  *logrus.Entry
	RemoteLogger *logrus.Entry
	execCache    executorCache
}

func (executor *BaseJobExecutor) bindJobStrategy(child JobStrategy) {
	executor.strategy = child
}

func (executor *BaseJobExecutor) setClient(client SnailJobClient) {
	executor.client = client
}

func (executor *BaseJobExecutor) setContext(ctx context.Context) {
	executor.ctx = ctx
}

func (executor *BaseJobExecutor) getContext() context.Context {
	return executor.ctx
}

func (executor *BaseJobExecutor) setExecutorCache(cache executorCache) {
	executor.execCache = cache
}

func (executor *BaseJobExecutor) setLogger(localLogger *logrus.Entry, remoteLogger *logrus.Entry) {
	executor.LocalLogger = localLogger
	executor.RemoteLogger = remoteLogger
}

func (executor *BaseJobExecutor) Context() context.Context {
	return executor.ctx
}

// JobExecute 模板类
func (executor *BaseJobExecutor) JobExecute(jobContext dto.JobContext) {
	resultChan := make(chan dto.ExecuteResult)
	timer := time.NewTimer(time.Duration(jobContext.ExecutorTimeout) * time.Second)
	defer timer.Stop()
	defer func() {

		if executors, found := executor.execCache.executors[jobContext.TaskBatchId]; found {
			i := 0
			for _, handler := range executors {
				if executor.strategy == handler {
					// 删除执行器
					executor.LocalLogger.Infof("delete executor cache jobTask:[%d]", jobContext.TaskId)
					//executors[i] = nil
					break
				}

				i++
			}

			// 若value没有值了  删除缓存
			// 遍历并删除满足条件的 key
			for _, value := range executor.execCache.executors {
				if allNil(value) {
					delete(executor.execCache.executors, jobContext.TaskBatchId)
				}
			}
			executor.LocalLogger.Infof("delete executor cache executors:[%+v]", executor.execCache.executors)

		}

	}()

	go func() {
		defer func() {
			err := recover()
			if err != nil {
				executor.LocalLogger.Error("job execute error", err)
				// 失败捕获异常
				resultChan <- *dto.Failure().WithMessage("执行失败").WithResult(err)
			}
		}()

		jobArgs := executor.buildJobArgsBasedOnType(jobContext)
		// 执行任务
		resultChan <- executor.strategy.DoJobExecute(jobArgs)
	}()

	// Wait for the result or timeout
	select {
	case <-timer.C:
		executor.LocalLogger.Warnf("任务执行超时. jobId: [%d] taskBatchId:[%d]", jobContext.JobId, jobContext.TaskBatchId)
		// 中断标志
		executor.ctx = context.WithValue(executor.ctx, constant.INTERRUPT_KEY, true)
	case result := <-resultChan:
		executor.LocalLogger.Debugf("BaseJobExecutor 执行了 JobExecute. jobId: [%d] result:[%s]", jobContext.JobId, result.Message)
		// 回调处理
		callback := &JobExecutorFutureCallback{jobContext, executor.LocalLogger, executor.RemoteLogger}
		callback.onCallback(executor.client, &result)
	}
}

func (executor *BaseJobExecutor) buildJobArgsBasedOnType(jobContext dto.JobContext) dto.IJobArgs {
	var jobArgs dto.IJobArgs
	switch jobContext.TaskType {
	case constant.SHARDING:
		jobArgs = executor.buildShardingJobArgs(jobContext)
	case constant.MAP_REDUCE, constant.MAP:
		if jobContext.MrStage == constant.MAP_STAGE {
			jobArgs = executor.buildMapJobArgs(jobContext)
		} else if jobContext.MrStage == constant.REDUCE_STAGE {
			jobArgs = executor.buildReduceJobArgs(jobContext)
		} else {
			jobArgs = executor.buildMergeReduceJobArgs(jobContext)
		}
	default:
		jobArgs = executor.buildBasicJobArgs(jobContext)
	}

	return jobArgs
}

func (executor *BaseJobExecutor) buildBasicJobArgs(jobContext dto.JobContext) dto.IJobArgs {
	return &dto.JobArgs{
		JobParams:       jobContext.JobArgsHolder.JobParams,
		ExecutorInfo:    jobContext.ExecutorInfo,
		TaskBatchId:     jobContext.TaskBatchId,
		JobId:           jobContext.JobId,
		WfContext:       jobContext.WfContext,
		ChangeWfContext: jobContext.ChangeWfContext,
	}
}

// Build sharding job args
func (executor *BaseJobExecutor) buildShardingJobArgs(jobContext dto.JobContext) dto.IJobArgs {
	args := dto.ShardingJobArgs{
		ShardingIndex: jobContext.ShardingIndex,
		ShardingTotal: jobContext.ShardingTotal,
	}
	args.JobParams = jobContext.JobArgsHolder.JobParams
	args.ExecutorInfo = jobContext.ExecutorInfo
	args.WfContext = jobContext.WfContext
	args.JobId = jobContext.JobId
	args.ChangeWfContext = jobContext.ChangeWfContext
	return &args
}

// Build map job args
func (executor *BaseJobExecutor) buildMapJobArgs(jobContext dto.JobContext) dto.IJobArgs {
	args := dto.MapArgs{
		MapResult: jobContext.JobArgsHolder.Maps,
		TaskName:  jobContext.TaskName,
	}

	args.JobParams = jobContext.JobArgsHolder.JobParams
	args.ExecutorInfo = jobContext.ExecutorInfo
	args.TaskBatchId = jobContext.TaskBatchId
	args.WfContext = jobContext.WfContext
	args.JobId = jobContext.JobId
	args.ChangeWfContext = jobContext.ChangeWfContext
	return &args
}

// Build reduce job args
func (executor *BaseJobExecutor) buildReduceJobArgs(jobContext dto.JobContext) dto.IJobArgs {
	args := dto.ReduceArgs{}

	args.JobParams = jobContext.JobArgsHolder.JobParams
	args.ExecutorInfo = jobContext.ExecutorInfo
	args.TaskBatchId = jobContext.TaskBatchId
	args.WfContext = jobContext.WfContext
	args.JobId = jobContext.JobId
	args.ChangeWfContext = jobContext.ChangeWfContext
	if maps := jobContext.JobArgsHolder.Maps; maps != nil {
		args.MapResult = parseMapResult(maps, executor.RemoteLogger)
	}

	return &args
}

func parseMapResult(maps interface{}, l *logrus.Entry) []interface{} {
	var result []interface{}

	switch v := maps.(type) {
	case string:
		// If the input is a JSON string, attempt to parse it
		err := json.Unmarshal([]byte(v), &result)
		if err != nil {
			l.Error("Error parsing JSON string:", err)
			return nil
		}
	case []interface{}:
		// If the input is already a slice of interface{}, use it directly
		result = v
	default:
		// If the input is of an unexpected type, handle it appropriately
		l.Warn("Unexpected type for maps:", v)
		return nil
	}

	return result
}

// Build merge reduce job args
func (executor *BaseJobExecutor) buildMergeReduceJobArgs(jobContext dto.JobContext) dto.IJobArgs {
	args := dto.MergeReduceArgs{}
	args.JobParams = jobContext.JobArgsHolder.JobParams
	args.ExecutorInfo = jobContext.ExecutorInfo
	args.TaskBatchId = jobContext.TaskBatchId
	args.WfContext = jobContext.WfContext

	if reduces := jobContext.JobArgsHolder.Reduces; reduces != nil {
		args.Reduces = parseMapResult(reduces, nil)
	}
	return &args
}

// 判断切片是否全为 nil
func allNil(slice []JobStrategy) bool {
	for _, v := range slice {
		if v != nil {
			return false
		}
	}
	return true
}
