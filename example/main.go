package main

import (
	snailjob "github.com/open-snail/snail-job-go"
	"github.com/open-snail/snail-job-go/dto"
	"github.com/open-snail/snail-job-go/job"
	"github.com/sirupsen/logrus"
)

func main() {

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

	exec.Register("testJobExecutor", func() job.IJobExecutor {
		return &TestJobExecutor{}
	}).Register("test2JobExecutor", func() job.IJobExecutor {
		return &Test2JobExecutor{}
	}).Register("testMapJobExecutor", func() job.IJobExecutor {
		return &TestMapJobExecutor{}
	}).Register("testMapReduceJobExecutor", func() job.IJobExecutor {
		return &TestMapReduceJobExecutor{}
	}).Register("testWorkflowJobExecutor", func() job.IJobExecutor {
		return &TestWorkflowJobExecutor{}
	})

	if nil == exec.Init() {
		exec.Run()
	}

}
