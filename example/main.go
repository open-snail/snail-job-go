package main

import (
	snailjob "opensnail.com/snail-job/snail-job-go"
	"opensnail.com/snail-job/snail-job-go/demo"
)

func main() {

	exec := snailjob.NewExecutor(&snailjob.Options{}, snailjob.NewLoggerFactory())
	config := snailjob.NewConfig()

	exec.Register("testJobExecutor", &demo.TestJobExecutor{})
	exec.Register("test2JobExecutor", &demo.Test2JobExecutor{})
	exec.Register("testMapJobExecutor", &demo.TestMapJobExecutor{})
	exec.Register("testMapReduceJobExecutor", &demo.TestMapReduceJobExecutor{})

	// todo
	if nil == exec.Init() {
		exec.Run()
	}

}
