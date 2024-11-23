package job

import (
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"opensnail.com/snail-job/snail-job-go/dto"
)

type HookLogService struct {
	Cache      sync.Map
	Wg         sync.WaitGroup //  优雅关闭
	LogEntryCh chan *logrus.Entry
	client     SnailJobClient
	jobContext dto.JobContext
}

func NewHookLogService(client SnailJobClient, jobContext dto.JobContext) *HookLogService {
	hs := &HookLogService{
		Cache:      sync.Map{},
		LogEntryCh: make(chan *logrus.Entry),
		client:     client,
		jobContext: jobContext,
	}
	return hs
}

func (hs *HookLogService) Init() {
	for entry := range hs.LogEntryCh {
		hs.Wg.Add(1)
		go func() {
			defer hs.Wg.Done()
			var items []*dto.JobLogTask
			record := Transform(hs.jobContext, entry)
			items = append(items, record)
			hs.client.SendBatchLogReport(items)
		}()
	}
}

func FormatExcInfo(err error) string {
	if err == nil {
		return ""
	}
	return fmt.Sprintf("%s", err)
}

func Transform(arg dto.JobContext, record *logrus.Entry) *dto.JobLogTask {
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
		NamespaceID: arg.NamespaceId,
		GroupName:   arg.GroupName,
		RealTime:    time.Now().UnixMilli(),
		FieldList:   fieldList,
		JobID:       int(arg.JobId),
		TaskBatchID: int(arg.TaskBatchId),
		TaskID:      int(arg.TaskId),
	}
}
