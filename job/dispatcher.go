package job

import (
	"context"
	"fmt"

	"opensnail.com/snail-job/snail-job-go/constant"
	"opensnail.com/snail-job/snail-job-go/dto"
)

type Dispatcher struct {
	client    SnailJobClient
	executors map[string]NewJobExecutor
	factory   LoggerFactory
	execCache executorCache
}

func Init(client SnailJobClient, executors map[string]NewJobExecutor, factory LoggerFactory) *Dispatcher {
	return &Dispatcher{client, executors, factory, *NewExecutorCache()}
}

func (e *Dispatcher) DispatchJob(dispatchJob dto.DispatchJobRequest) dto.Result {
	jobContext := buildJobContext(dispatchJob)
	cxt := context.WithValue(context.Background(), constant.JOB_CONTEXT_KEY, jobContext)
	remoteLogger := e.factory.GetRemoteLogger(dispatchJob.ExecutorInfo, cxt)
	localLogger := e.factory.GetLocalLogger(dispatchJob.ExecutorInfo)

	if dispatchJob.RetryCount > 0 {
		remoteLogger.Info("Task execution/scheduling failed, retrying. Retry count: [%d]", dispatchJob.RetryCount)
	}

	// 必须是GO客户端才能使用
	//if dispatchJob.ExecutorType != 3 {
	//	log.Printf("Non-Go executor type not supported. ExecutorType: [%s]", dispatchJob.ExecutorInfo)
	//	return dto.Result{Status: 1, Message: "不支持非Java类型的执行器", Data: false}
	//}

	jobExecute, _ := e.GetExecutor(dispatchJob.ExecutorInfo)
	if jobExecute == nil {
		remoteLogger.Info("Invalid executor configuration. ExecutorInfo: [%s]", dispatchJob.ExecutorInfo)
		return dto.Result{Status: 1, Message: "执行器配置有误", Data: false}
	}

	jobStrategy := jobExecute.(JobStrategy)
	jobStrategy.setExecutorCache(e.execCache)
	jobStrategy.bindJobStrategy(jobStrategy)
	jobStrategy.setClient(e.client)

	//logrus.WithContext()
	jobStrategy.setContext(cxt)
	jobStrategy.setLogger(localLogger, remoteLogger)
	jobStrategy.setLogrusLogger(e.client.LocalLog, e.client.RemoteLog)
	// 注册实例
	e.execCache.register(dispatchJob.TaskBatchId, jobStrategy)

	// bing executor
	if dispatchJob.TaskType == constant.MAP {
		mapExecute := jobExecute.(MapExecute)
		mapExecute.bindMapExecute(mapExecute)
	} else if dispatchJob.TaskType == constant.MAP_REDUCE {
		mapExecute := jobExecute.(MapExecute)
		mapExecute.bindMapExecute(mapExecute)
		mapReduceExecute := jobExecute.(MapReduceExecute)
		mapReduceExecute.BindMapReduceExecute(mapReduceExecute)
	}

	// 集群、 广播、静态分片 直接执行方法
	jobExecute.JobExecute(jobContext)

	return dto.Result{Status: constant.YES, Data: true}
}

func buildJobContext(dispatchJob dto.DispatchJobRequest) dto.JobContext {
	jobContext := dto.JobContext{
		JobId:               dispatchJob.JobId,
		ShardingTotal:       dispatchJob.ShardingTotal,
		ShardingIndex:       dispatchJob.ShardingIndex,
		NamespaceId:         dispatchJob.NamespaceId,
		TaskId:              dispatchJob.TaskId,
		TaskBatchId:         dispatchJob.TaskBatchId,
		GroupName:           dispatchJob.GroupName,
		ExecutorInfo:        dispatchJob.ExecutorInfo,
		ParallelNum:         dispatchJob.ParallelNum,
		TaskType:            dispatchJob.TaskType,
		ExecutorTimeout:     dispatchJob.ExecutorTimeout,
		WorkflowNodeId:      dispatchJob.WorkflowNodeId,
		WorkflowTaskBatchId: dispatchJob.WorkflowTaskBatchId,
		RetryStatus:         dispatchJob.RetryStatus,
		RetryScene:          dispatchJob.RetryScene,
		TaskName:            dispatchJob.TaskName,
		MrStage:             dispatchJob.MrStage,
	}

	// Parse ArgsStr and WfContext (simplified example)
	if dispatchJob.ArgsStr != "" {
		jobContext.JobArgsHolder = dto.JobArgsHolder{JobParams: dispatchJob.ArgsStr}
	}

	if dispatchJob.WfContext != "" {
		jobContext.WfContext = make(map[string]interface{})
		// JSON parse here
	}

	return jobContext
}

func (e *Dispatcher) GetExecutor(name string) (IJobExecutor, error) {
	executor, exists := e.executors[name]
	if !exists {
		return nil, fmt.Errorf("executor [%s] not found", name)
	}
	return executor(), nil
}

func (e *Dispatcher) Stop(stopJob dto.StopJob) dto.Result {

	if executors, found := e.execCache.executors[stopJob.TaskBatchId]; found {
		for _, handler := range executors {
			if handler != nil {
				handler.setContext(context.WithValue(handler.getContext(), constant.INTERRUPT_KEY, true))
			}
		}
	}

	// 删除缓存
	delete(e.execCache.executors, stopJob.TaskBatchId)
	return dto.Result{Status: 1, Data: true}
}
