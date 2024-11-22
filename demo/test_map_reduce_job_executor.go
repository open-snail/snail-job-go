package demo

import (
	"time"

	"opensnail.com/snail-job/snail-job-go/dto"
	"opensnail.com/snail-job/snail-job-go/job"
)

// TestMapReduceJobExecutor 这是一个测试类
type TestMapReduceJobExecutor struct {
	job.BaseMapReduceJobExecutor
}

func (executor *TestMapReduceJobExecutor) DoJobMapExecute(mpArgs *dto.MapArgs) dto.ExecuteResult {
	logger := executor.GetLogger()
	logger.Info("TestMapReduceJobExecutor 开始执行 DoJobMapExecute. jobId: [%d] TaskName:[%s] now:[%s]", mpArgs.GetJobId(), mpArgs.TaskName, time.Now().String())
	time.Sleep(1 * time.Second)
	logger.Info("TestMapReduceJobExecutor 执行结束 DoJobMapExecute. jobId: [%d]  TaskName:[%s] now:[%s]", mpArgs.GetJobId(), mpArgs.TaskName, time.Now().String())
	//panic("这是故意抛出的异常")
	//num1 := 1
	//num2 := 1
	//num3 := num1 / num2
	//fmt.Println(num3)
	return dto.ExecuteResult{}
}

// DoReduceExecute 模板类
func (executor *TestMapReduceJobExecutor) DoReduceExecute(jobArgs *dto.ReduceArgs) dto.ExecuteResult {
	logger := executor.GetLogger()
	// todo 怎么把jobArgs 转成 mapArgs
	logger.Info("TestMapReduceJobExecutor 开始执行 DoReduceExecute.")

	return dto.ExecuteResult{}
}

func (executor *TestMapReduceJobExecutor) DoMergeReduceExecute(jobArgs *dto.MergeReduceArgs) dto.ExecuteResult {
	// todo 怎么把jobArgs 转成 mapArgs
	//executor.Execute.DoJobMapExecute(&dto.MapArgs{})
	logger := executor.GetLogger()
	logger.Info("TestMapReduceJobExecutor 开始执行 DoMergeReduceExecute.")

	return dto.ExecuteResult{}
}
