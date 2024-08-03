# SnailJob 客户端 golang

## 任务

- [x] 注册到服务器 (go-xxxx)
- [x] 实现`/job/dispatch/v1`任务调度
- [ ] 日志对象`SnailRemoteLog`实现上报log到服务器，字段如下：

```txt
time_stamp
level
thread
message
location
throwable
```

另外，每次执行任务需要向goroutine通过context传递如下三个变量：

```txt
jobId
taskBatchId
taskId
```

- [ ] 本地日志对象`SnailLocalLog` 
- [ ] 建package
