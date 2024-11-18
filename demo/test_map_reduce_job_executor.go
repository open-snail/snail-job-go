package demo

import (
	"fmt"
	"snail_job_go/dto"
	"snail_job_go/executor"
	"snail_job_go/job"
	"time"
)

/**
  todo 待讨论Go的是支持结构体还是直接执行方法方式
*/

// TestMapReduceJobExecutor 这是一个测试类
type TestMapReduceJobExecutor struct {
	executor.AbstractMapReduceJobExecutor
}

// NewTestMapJobExecutor 创建对象
func NewTestMapReduceJobExecutor() *TestMapReduceJobExecutor {
	executor := &TestMapReduceJobExecutor{}
	executor.BindMapReduceExecute(executor)
	executor.BindMapExecute(executor)
	executor.BindJobStrategy(executor)
	return executor
}

func (executor *TestMapReduceJobExecutor) DoJobMapExecute(mpArgs *dto.MapArgs) dto.ExecuteResult {

	job.LocalLog.Info(fmt.Sprintf("TestMapReduceJobExecutor 开始执行 DoJobMapExecute. jobId: [%d] TaskName:[%s] now:[%s]", mpArgs.GetJobId(), mpArgs.TaskName, time.Now().String()))
	time.Sleep(1 * time.Second)
	job.LocalLog.Info(fmt.Sprintf("TestMapReduceJobExecutor 执行结束 DoJobMapExecute. jobId: [%d]  TaskName:[%s] now:[%s]", mpArgs.GetJobId(), mpArgs.TaskName, time.Now().String()))
	//panic("这是故意抛出的异常")
	num1 := 1
	num2 := 1
	num3 := num1 / num2
	fmt.Println(num3)
	return dto.ExecuteResult{}
}

// DoReduceExecute 模板类
func (executor *TestMapReduceJobExecutor) DoReduceExecute(jobArgs *dto.ReduceArgs) dto.ExecuteResult {
	// todo 怎么把jobArgs 转成 mapArgs
	fmt.Printf("TestMapReduceJobExecutor 开始执行 DoReduceExecute.")

	return dto.ExecuteResult{}
}

func (executor *TestMapReduceJobExecutor) DoMergeReduceExecute(jobArgs *dto.MergeReduceArgs) dto.ExecuteResult {
	// todo 怎么把jobArgs 转成 mapArgs
	//executor.Execute.DoJobMapExecute(&dto.MapArgs{})
	fmt.Printf("TestMapReduceJobExecutor 开始执行 DoMergeReduceExecute.")

	return dto.ExecuteResult{}
}
