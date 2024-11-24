package job

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"opensnail.com/snail-job/snail-job-go/constant"
	"opensnail.com/snail-job/snail-job-go/dto"
	"time"
)

type LoggerHook struct {
	Hs *HookLogService
}

func (h *LoggerHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *LoggerHook) Fire(entry *logrus.Entry) error {
	if entry.Context != nil {
		jobContext := entry.Context.Value(constant.JOB_CONTEXT_KEY).(dto.JobContext)
		h.Hs.LogEntryCh <- h.transform(&jobContext, entry)
		fmt.Println(fmt.Sprintf("开始上报: %+v msg:%+v", jobContext.TaskBatchId, entry.Message))
	}

	return nil
}

func (h *LoggerHook) transform(ctx *dto.JobContext, record *logrus.Entry) *dto.JobLogTask {
	if ctx == nil {
		return nil
	}

	fieldList := []dto.TaskLogFieldDTO{
		{Name: "time_stamp", Value: fmt.Sprintf("%d", record.Time.UnixMilli())},
		{Name: "level", Value: record.Level.String()},
		{Name: "thread", Value: ""}, //record.Caller.File},
		{Name: "message", Value: record.Message},
		{Name: "location", Value: "unknown"}, //fmt.Sprintf("%s:%s:%d", record.Caller.File, record.Caller.Function, record.Caller.Line)},
		{Name: "throwable", Value: ""},       //FormatExcInfo(record.Context.Err())},
		{Name: "host", Value: "localhost"},
		{Name: "port", Value: "17889"},
	}

	return &dto.JobLogTask{
		LogType:     "JOB",
		NamespaceID: ctx.NamespaceId,
		GroupName:   ctx.GroupName,
		RealTime:    time.Now().UnixMilli(),
		FieldList:   fieldList,
		JobID:       ctx.JobId,
		TaskBatchID: ctx.TaskBatchId,
		TaskID:      ctx.TaskId,
	}
}
