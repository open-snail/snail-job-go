package main

import (
	"snail_job_go/constant"
	"snail_job_go/demo"
	"snail_job_go/endpoint"
	"snail_job_go/job"
	"snail_job_go/register"
)

type key string

const (
	EXECUTOR_MANAGER_KEY key = "executor_manager"
)

//func BuildDispatchJobResult(dispatchJobRequest job.DispatchJobRequest, executeResult job.ExecuteResult) job.DispatchJobResult {
//	args := dispatchJobRequest.Args[0]
//	taskStatus := job.BATCH_STATUS_FAIL
//	if executeResult.Success == job.STATUS_SUCCESS {
//		taskStatus = job.BATCH_STATUS_SUCCESS
//	}
//	return job.DispatchJobResult{
//		JobId:               args.JobId,
//		TaskBatchId:         args.TaskBatchId,
//		WorkflowTaskBatchId: args.WorkflowTaskBatchId,
//		WorkflowNodeId:      args.WorkflowNodeId,
//		TaskId:              args.TaskId,
//		TaskType:            args.TaskType,
//		GroupName:           args.GroupName,
//		TaskStatus:          taskStatus,
//		ExecuteResult:       executeResult,
//		RetryScene:          args.RetryScene,
//		IsRetry:             args.IsRetry,
//	}
//}
//
//type ExecutorManager struct {
//	hls           *job.HookLogService
//	executors     map[string]func(*job.HookLogService, job.DispatchJobRequest) job.ExecuteResult
//	executorPools map[int]*sync.Pool
//	lock          sync.Mutex
//}
//
//func NewExecutorManager(hookLogService *job.HookLogService) *ExecutorManager {
//	return &ExecutorManager{
//		hls:           hookLogService,
//		executors:     make(map[string]func(*job.HookLogService, job.DispatchJobRequest) job.ExecuteResult),
//		executorPools: make(map[int]*sync.Pool),
//	}
//}
//
//func (m *ExecutorManager) Register(name string, executor func(*job.HookLogService, job.DispatchJobRequest) job.ExecuteResult) {
//	m.lock.Lock()
//	defer m.lock.Unlock()
//
//	if _, exists := m.executors[name]; exists {
//		panic(fmt.Sprintf("Executor [%s] already registered", name))
//	}
//
//	m.executors[name] = executor
//	job.LocalLog.Info(fmt.Sprintf("Registered executor: %s", name))
//}
//
//func (m *ExecutorManager) Execute(hls *job.HookLogService, req job.DispatchJobRequest) job.ExecuteResult {
//	defer func() {
//		if err := recover(); err != nil {
//			fmt.Println(err)
//		}
//	}()
//
//	if len(req.Args) == 0 {
//		return job.ExecuteResult{Success: job.STATUS_SUCCESS, Message: "args cannot be empty"}
//	}
//
//	args := req.Args[0]
//	executor, exists := m.executors[args.ExecutorInfo]
//	if !exists {
//		return job.ExecuteResult{Success: job.STATUS_SUCCESS, Message: fmt.Sprintf("Executor not found: %s", args.ExecutorInfo)}
//	}
//
//	job.LocalLog.Info(fmt.Sprintf("Executing with executor: %s", args.ExecutorInfo))
//
//	result := executor(hls, req)
//
//	job.SendDispatchResult(BuildDispatchJobResult(req, result))
//
//	return result
//}
//
//func (m *ExecutorManager) Stop(req job.StopJobRequest) {
//	if len(req.Args) == 0 {
//		job.LocalLog.Info("args cannot be empty")
//		return
//	}
//
//	args := req.Args[0]
//	m.lock.Lock()
//	defer m.lock.Unlock()
//
//	if pool, exists := m.executorPools[args.TaskBatchID]; exists {
//		pool.Put(nil)
//		delete(m.executorPools, args.TaskBatchID)
//		job.LocalLog.Info(fmt.Sprintf("Stopped task batch: %d", args.TaskBatchID))
//	}
//}
//
//func HandleDispatch(w http.ResponseWriter, r *http.Request) {
//	var req job.DispatchJobRequest
//	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
//		job.LocalLog.Fatal(err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	job.LocalLog.Info(fmt.Sprintf("Received job dispatch request: reqId=%d", req.ReqID))
//	manager := r.Context().Value(EXECUTOR_MANAGER_KEY).(*ExecutorManager)
//
//	go manager.Execute(manager.hls, req)
//
//	json.NewEncoder(w).Encode(job.NettyResult{
//		Status: job.STATUS_SUCCESS,
//		ReqID:  req.ReqID,
//		Data:   true,
//	})
//}
//
//func HandleStop(w http.ResponseWriter, r *http.Request) {
//	var req job.StopJobRequest
//	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
//		job.LocalLog.Fatal(err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	job.LocalLog.Info(fmt.Sprintf("Received job stop request: reqId=%d", req.ReqID))
//	manager := r.Context().Value(EXECUTOR_MANAGER_KEY).(*ExecutorManager)
//	manager.Stop(req)
//
//	json.NewEncoder(w).Encode(job.NettyResult{
//		Status: job.STATUS_SUCCESS,
//		ReqID:  req.ReqID,
//		Data:   true,
//	})
//}
//
//func LoggingMiddleware(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		start := time.Now()
//		job.LocalLog.Printf("Started %s %s", r.Method, r.RequestURI)
//		job.LocalLog.Printf("Request Headers: %v", r.Header)
//
//		// 检查是否有正确Token的头部
//		if r.Header.Get("Sj-Token") != job.HEADERS["token"] {
//			http.Error(w, "Method Not Allowed", http.StatusNonAuthoritativeInfo)
//			return
//		}
//
//		// 检查是否为 Post
//		if r.Method != http.MethodPost {
//			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
//			return
//		}
//
//		next.ServeHTTP(w, r)
//
//		job.LocalLog.Printf("Completed %s in %v", r.RequestURI, time.Since(start))
//	})
//}
//
//func RunServer(manager *ExecutorManager) {
//	mux := http.NewServeMux()
//
//	mux.HandleFunc("/job/dispatch/v1", func(w http.ResponseWriter, r *http.Request) {
//		ctx := context.WithValue(r.Context(), EXECUTOR_MANAGER_KEY, manager)
//		r = r.WithContext(ctx)
//		HandleDispatch(w, r)
//	})
//
//	mux.HandleFunc("/job/stop/v1", func(w http.ResponseWriter, r *http.Request) {
//		ctx := context.WithValue(r.Context(), EXECUTOR_MANAGER_KEY, manager)
//		r = r.WithContext(ctx)
//		HandleStop(w, r)
//	})
//
//	loggedMux := LoggingMiddleware(mux)
//
//	job.LocalLog.Info("Starting server")
//	if err := http.ListenAndServe(":1789", loggedMux); err != nil {
//		job.LocalLog.Fatal(err)
//	}
//}
//
//var LogContextKey = &ContextKey{"SnailLogContext"}
//
//type ContextKey struct {
//	name string
//}
//
//func (k *ContextKey) String() string {
//	return "snailjob context value " + k.name
//}
//
//func FormatExcInfo(err error) string {
//	if err == nil {
//		return ""
//	}
//	return fmt.Sprintf("%s", err)
//}
//
//func Transform(arg job.DispatchJobArgs, record *job.LogRecord) *job.JobLogTask {
//	fieldList := []job.TaskLogFieldDTO{
//		{"time_stamp", fmt.Sprintf("%d", record.TimeStamp.UnixMilli())},
//		{"level", record.Level},
//		{"thread", record.Thread},
//		{"message", record.Message},
//		{"location", fmt.Sprintf("%s:%s:%d", record.Module, record.FuncName, record.Lineno)},
//		{"throwable", FormatExcInfo(record.ExcInfo)},
//		{"host", job.SNAIL_HOST_IP},
//		{"port", job.SNAIL_HOST_PORT},
//	}
//
//	return &job.JobLogTask{
//		LogType:     "JOB",
//		NamespaceID: job.SNAIL_NAMESPACE,
//		GroupName:   job.SNAIL_GROUP_NAME,
//		RealTime:    time.Now().UnixMilli(),
//		FieldList:   fieldList,
//		JobID:       arg.JobId,
//		TaskBatchID: arg.TaskBatchId,
//		TaskID:      arg.TaskId,
//	}
//}

// 这是一个执行器样例
//func TestJobExecutor(hls *job.HookLogService, req job.DispatchJobRequest) job.ExecuteResult {
//	job.LocalLog.Info(fmt.Sprintf("Executing exampleExecutor with args: %d", req.ReqID))
//	// FIXME: 获取不到实际参数
//	fmt.Printf("%v", req.Args)
//
//	for i := 0; i < 10; i++ {
//		// 会在服务器上的任务批次日志中提现
//		hls.MessageCh <- Transform(req.Args[0], &job.LogRecord{
//			TimeStamp: time.Now(),
//			Level:     "1",
//			Message:   fmt.Sprintf("这是一个循环体 %d", i),
//			FuncName:  "test",
//			Lineno:    1,
//			ExcInfo:   nil,
//		})
//	}
//	return job.ExecuteResult{Success: job.STATUS_SUCCESS, Message: "Execution successful"}
//}

// 这是一个执行器样例, 他会抛出异常
//func TestJobExecutorFailed(hls *job.HookLogService, argsStr job.DispatchJobRequest) job.ExecuteResult {
//	panic("这是故意抛出的异常")
//}

func main() {

	//go job.SendHeartbeat()

	hls := job.NewHookLogService()

	manager := register.NewExecutorManager(hls)

	// 初始化EndPoint

	jbEndPoint := endpoint.Init(manager)

	go hls.Init()

	//// 注册示例执行器
	manager.Register("testJobExecutor", demo.NewTestJobExecutor())
	manager.Register("test2JobExecutor", &demo.Test2JobExecutor{})
	manager.Register("testMapJobExecutor", demo.NewTestMapJobExecutor())
	manager.Register("testMapReduceJobExecutor", &demo.TestMapReduceJobExecutor{})

	// ToDo 模拟远程调用endpoint
	jbEndPoint.DispatchJob(job.DispatchJobRequest{NamespaceId: "12", JobId: 1,
		TaskType: constant.MAP, ExecutorType: 3, ExecutorInfo: "testMapJobExecutor",
		ExecutorTimeout: 3, TaskName: "ROOT", MrStage: constant.MAP_STAGE})

	//testJobExecutor, _ := manager.GetExecutor("testJobExecutor")

	////println(testJobExecutor_)
	////println(error.Error())
	//
	//testJobExecutor.JobExecute(dto.JobContext{JobId: 1})
	//RunServer(manager)

	//// 创建子类实例
	//context := dto.JobContext{}
	//
	//// 子类实例化
	//jobExecutor := executor.NewAbstractMapJobExecutor()
	//
	////jobExecutor.BindJobStrategy(jobExecutor)
	//
	//// 调用父类的 JobExecute 方法
	//jobExecutor.JobExecute(context)

}
