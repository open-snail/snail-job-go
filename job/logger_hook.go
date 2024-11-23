package job

import (
	"github.com/sirupsen/logrus"
)

type LoggerHook struct {
	hs *HookLogService
}

func (h *LoggerHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *LoggerHook) Fire(entry *logrus.Entry) error {
	h.hs.LogEntryCh <- entry
	return nil
}
