package job

import (
	"context"
	"encoding/json"
	"fmt"
	"opensnail.com/snail-job/snail-job-go/constant"
	"opensnail.com/snail-job/snail-job-go/dto"
	"time"
)

// IJobExecutor 执行器接口
type IJobExecutor interface {
	JobExecute(context dto.JobContext)
}

type JobStrategy interface {
	DoJobExecute(dto.IJobArgs) dto.ExecuteResult
	BindJobStrategy(child JobStrategy)
	SetClient(client SnailJobClient)
}

type BaseJobExecutor struct {
	Strategy JobStrategy
	client   SnailJobClient
}

func (executor *BaseJobExecutor) BindJobStrategy(child JobStrategy) {
	executor.Strategy = child
}

func (executor *BaseJobExecutor) SetClient(client SnailJobClient) {
	executor.client = client
}

func (executor *BaseJobExecutor) GetClient(client SnailJobClient) {
	executor.client = client
}

// JobExecute 模板类
func (executor *BaseJobExecutor) JobExecute(jobContext dto.JobContext) {

	resultChan := make(chan dto.ExecuteResult)
	// Add a stop task to the timer to stop execution upon timeout
	timer := time.NewTimer(time.Duration(jobContext.ExecutorTimeout) * time.Second)
	defer timer.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		defer func() {
			err := recover()
			if err != nil {
				LocalLog.Error(err)
				// 失败捕获异常
				resultChan <- *dto.Failure(err, "执行失败")
			}
		}()

		jobArgs := executor.buildJobArgsBasedOnType(jobContext)
		// 执行任务
		resultChan <- executor.Strategy.DoJobExecute(jobArgs)
	}()

	// Wait for the result or timeout
	select {
	case <-ctx.Done():
		LocalLog.Warnf(fmt.Sprintf("BaseJobExecutor 任务被取消. jobId: [%d]", jobContext.JobId))
	case <-timer.C:
		cancel() // Cancel the job execution
		LocalLog.Warnf(fmt.Sprintf("BaseJobExecutor 任务执行超时. jobId: [%d]", jobContext.JobId))
	case result := <-resultChan:
		LocalLog.Info(fmt.Sprintf("BaseJobExecutor 执行了 JobExecute. jobId: [%d] result:[%s]", jobContext.JobId, result.Message))
		// 回调处理
		callback := &JobExecutorFutureCallback{}
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
		JobParams:    jobContext.JobArgsHolder.JobParams,
		ExecutorInfo: jobContext.ExecutorInfo,
		TaskBatchId:  jobContext.TaskBatchId,
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

	return &args
}

// Build reduce job args
func (executor *BaseJobExecutor) buildReduceJobArgs(jobContext dto.JobContext) dto.IJobArgs {
	args := dto.ReduceArgs{}

	args.JobParams = jobContext.JobArgsHolder.JobParams
	args.ExecutorInfo = jobContext.ExecutorInfo
	args.TaskBatchId = jobContext.TaskBatchId
	args.WfContext = jobContext.WfContext
	if maps := jobContext.JobArgsHolder.Maps; maps != nil {
		args.MapResult = parseMapResult(maps)
	}

	return &args
}

func parseMapResult(maps interface{}) []interface{} {
	var result []interface{}

	switch v := maps.(type) {
	case string:
		// If the input is a JSON string, attempt to parse it
		err := json.Unmarshal([]byte(v), &result)
		if err != nil {
			fmt.Println("Error parsing JSON string:", err)
			return nil
		}
	case []interface{}:
		// If the input is already a slice of interface{}, use it directly
		result = v
	default:
		// If the input is of an unexpected type, handle it appropriately
		fmt.Println("Unexpected type for maps:", v)
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
		args.Reduces = parseMapResult(reduces)
	}
	return &args
}
