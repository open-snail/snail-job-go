<p align="center">
  <a href="https://snailjob.opensnail.com">
   <img alt="snail-job-Logo" src="doc/images/favicon.svg" width="200px">
  </a>
</p>

<p align="center">
    🔥🔥🔥 灵活，可靠和快速的分布式任务重试和分布式任务调度平台<br> <br/>
</p>

<p align="center">

> ✅️ 可重放，可管控、为提高分布式业务系统一致性的分布式任务重试平台 <br/>
> ✅️ 支持秒级、可中断、可编排的高性能分布式任务调度平台
</p>

# 简介

> SnailJob 是一个灵活、可靠且高效的分布式任务重试和任务调度平台。其核心采用分区模式实现，具备高度可伸缩性和容错性的分布式系统。拥有完善的权限管理、强大的告警监控功能和友好的界面交互。欢迎大家接入并使用。

## snail-job-go

snail-job 项目的 GO 客户端。[snail-job项目 java 后端](https://gitee.com/aizuda/snail-job)

采用GO原生语言开发的SnailJob客户端具备与SnailJob的Java客户端Job模块一样的能力包括(集群、广播、静态分片、Map、MapReuce、DAG工作流、实时日志等功能)

## 开始使用

1. 在go.mod文件中添加依赖
```shell
require  github.com/open-snail/snail-job-go v0.0.2
```
2. 配置客户端参数
```go
// 配置Options参数
exec := snailjob.NewSnailJobManager(&dto.Options{
		ServerHost:   "127.0.0.1",
		ServerPort:   "17888",
		HostIP:       "127.0.0.1",
		HostPort:     "17889",
		Namespace:    "764d604ec6fc45f68cd92514c40e9e1a",
		GroupName:    "snail_job_demo_group",
		Token:        "SJ_Wyz3dmsdbDOkDujOTSSoBjGQP1BMsVnj",
		Level:        logrus.InfoLevel,
		ReportCaller: true,
	})

    // 注册执行器
	exec.Register("testJobExecutor", func() job.IJobExecutor {
		return &Test3JobExecutor{}
	})

    // 初始化环境
	if nil == exec.Init() {
		// 启动客户端
		exec.Run()
	}
```
登录后台，能看到对应host-id 为 `go-xxxxxx` 的客户端

### 示例

#### 定时任务

```go
// TestJobExecutor 这是一个示例执行器
type TestJobExecutor struct {
   job.BaseJobExecutor
}

func (executor *Test2JobExecutor) DoJobExecute(jobArgs dto.IJobArgs) dto.ExecuteResult {
executor.RemoteLogger.Infof("TestJobExecutor 执行结束 DoJobExecute. jobId: [%d] now:[%s]", jobArgs.GetJobId(), time.Now().String())
return *dto.Success("hello 这是go客户端")
}

```

新建定时任务, 执行器类型选择【Go】，执行器名称填入【testJobExecutor】

#### 动态分片

```go
// TestMapJobExecutor 这是一个测试类
type TestMapJobExecutor struct {
	job.BaseMapJobExecutor
}

func (executor *TestMapJobExecutor) DoJobMapExecute(mpArgs *dto.MapArgs) dto.ExecuteResult {
	logger := executor.LocalLogger
	if mpArgs.TaskName == constant.ROOT_MAP {
		_, _ = executor.DoMap([]interface{}{1, 2, 3}, "secondTaskName")
		return *dto.Success(nil)
	}

	logger.Infof("TestMapJobExecutor执行 DoJobMapExecute. jobId: [%d] TaskName:[%s] ", mpArgs.GetJobId(), mpArgs.TaskName)
	return *dto.Success("这是动态分片")
}

```

#### MapReduce

```go
// TestMapReduceJobExecutor 这是一个测试类
type TestMapReduceJobExecutor struct {
	job.BaseMapReduceJobExecutor
}

func (executor *TestMapReduceJobExecutor) DoJobMapExecute(mpArgs *dto.MapArgs) dto.ExecuteResult {
	logger := executor.LocalLogger
	return *dto.Success("这是动态分片阶段")
}

// DoReduceExecute 模板类
func (executor *TestMapReduceJobExecutor) DoReduceExecute(jobArgs *dto.ReduceArgs) dto.ExecuteResult {
	logger := executor.LocalLogger
	logger.Infof("TestMapReduceJobExecutor 开始执行 DoReduceExecute.")

    return *dto.Success("这是Reduce阶段")
}

func (executor *TestMapReduceJobExecutor) DoMergeReduceExecute(jobArgs *dto.MergeReduceArgs) dto.ExecuteResult {
	logger := executor.LocalLogger
	logger.Info("TestMapReduceJobExecutor 开始执行 DoMergeReduceExecute.")

    return *dto.Success("这是merge阶段")
}

```

#### 响应停止事件

```go
type TestJobExecutor struct {
	job.BaseJobExecutor
}

// 测试超时时间
func (executor *TestJobExecutor) DoJobExecute(jobArgs dto.IJobArgs) dto.ExecuteResult {

	time.Sleep(1 * time.Second)
	interrupt := executor.Context().Value(constant.INTERRUPT_KEY)
	if interrupt != nil {
		executor.LocalLogger.Errorf("任务被中断. jobId: [%d] now:[%s]", jobArgs.GetJobId(), time.Now().String())
		return *dto.Failure(nil, "任务被中断")
	}
	
	return *dto.Success("hello 这是go客户端")
}

```

### 工作流

```go

// TestWorkflowJobExecutor 这是一个测试类
type TestWorkflowJobExecutor struct {
	job.BaseJobExecutor
}

func (executor *TestWorkflowJobExecutor) DoJobExecute(jobArgs dto.IJobArgs) dto.ExecuteResult {
	executor.LocalLogger.Infof("TestWorkflowJobExecutor. jobId: [%d] wfContext:[%+v]",
		jobArgs.GetJobId(), jobArgs.GetWfContext("name"))
	jobArgs.AppendContext("name", "xiaowoniu")
	return *dto.Success("测试工作流")
}

```