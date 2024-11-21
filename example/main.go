package main

import (
	snailjob "opensnail.com/snail-job/snail-job-go"
	"opensnail.com/snail-job/snail-job-go/demo"
	"opensnail.com/snail-job/snail-job-go/dto"
)

func main() {

	exec := snailjob.NewSnailJobManager(&dto.Options{
		ServerHost: "127.0.0.1",
		ServerPort: "17888",
		HostIP:     "127.0.0.1",
		HostPort:   "17889",
		Namespace:  "764d604ec6fc45f68cd92514c40e9e1a",
		GroupName:  "snail_job_demo_group",
		Token:      "SJ_Wyz3dmsdbDOkDujOTSSoBjGQP1BMsVnj",
	}, snailjob.NewLoggerFactory())

	exec.Register("testJobExecutor", &demo.TestJobExecutor{})
	exec.Register("test2JobExecutor", &demo.Test2JobExecutor{})
	exec.Register("testMapJobExecutor", &demo.TestMapJobExecutor{})
	exec.Register("testMapReduceJobExecutor", &demo.TestMapReduceJobExecutor{})

	// todo
	if nil == exec.Init() {
		exec.Run()
	}

}
