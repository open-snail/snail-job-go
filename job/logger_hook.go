package job

import (
	"github.com/sirupsen/logrus"
	"opensnail.com/snail-job/snail-job-go/dto"
)

type LoggerHook struct {
	jobContext dto.JobContext
}

func (h *LoggerHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *LoggerHook) Fire(entry *logrus.Entry) error {
	//fmt.Printf(entry.Level.String())
	//fmt.Printf(entry.Message)
	return nil
}
