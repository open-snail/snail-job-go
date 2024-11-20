package main

import snailjob "opensnail.org/snail-job/snail-job-go"

func main() {

	exec := snailjob.NewExecutor(&snailjob.Options{}, snailjob.NewLoggerFactory())

	if nil == exec.Init() {
		exec.Run()
	}
}
