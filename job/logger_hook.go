package job

import (
	"fmt"
	"opensnail.com/snail-job/snail-job-go/util"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"opensnail.com/snail-job/snail-job-go/constant"
	"opensnail.com/snail-job/snail-job-go/dto"
)

type LoggerHook struct {
	Hls *HookLogService
}

func (h *LoggerHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *LoggerHook) Fire(entry *logrus.Entry) error {
	if entry.Context != nil {
		jobContext := entry.Context.Value(constant.JOB_CONTEXT_KEY).(dto.JobContext)
		h.Hls.LogEntryCh <- h.transform(&jobContext, entry)
		fmt.Printf("开始上报: %+v msg:%+v\n", jobContext.TaskBatchId, entry.Message)
	}

	return nil
}

func (h *LoggerHook) transform(ctx *dto.JobContext, entry *logrus.Entry) *dto.JobLogTask {
	if ctx == nil {
		return nil
	}
	if !entry.HasCaller() {
		panic("请设置 logrus 的 ReportCaller 为 true")
	}

	file := util.TrimProjectPath(entry.Caller.File, path)
	method := util.TrimProjectPath(entry.Caller.Function, moduleName)
	fieldList := []dto.TaskLogFieldDTO{
		{Name: "time_stamp", Value: fmt.Sprintf("%d", entry.Time.UnixMilli())},
		{Name: "level", Value: strings.ToUpper(entry.Level.String())},
		{Name: "thread", Value: file},
		{Name: "message", Value: entry.Message},
		{Name: "location", Value: fmt.Sprintf("%s:%d", method, entry.Caller.Line)},
		{Name: "throwable", Value: formatExcInfo(entry.Context.Err())},
		{Name: "host", Value: h.Hls.client.opts.HostIP},
		{Name: "port", Value: h.Hls.client.opts.HostPort},
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

func formatExcInfo(err error) string {
	if err == nil {
		return ""
	}
	return fmt.Sprintf("%s", err)
}
