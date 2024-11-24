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
	LogEntryCh chan *dto.JobLogTask
	client     SnailJobClient
}

func NewHookLogService(client SnailJobClient) *HookLogService {
	hs := &HookLogService{
		Cache:      sync.Map{},
		LogEntryCh: make(chan *dto.JobLogTask),
		client:     client,
	}
	return hs
}

func (hs *HookLogService) Init() {
	for entry := range hs.LogEntryCh {
		hs.Wg.Add(1)
		go func() {
			defer hs.Wg.Done()
			var items []*dto.JobLogTask
			//record := hs.Transform(hs.jobContext, entry)
			items = append(items, entry)
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

func (hs *HookLogService) Transform(arg dto.JobContext, record *logrus.Entry) *dto.JobLogTask {
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
		NamespaceID: hs.client.opts.Namespace,
		GroupName:   hs.client.opts.GroupName,
		RealTime:    time.Now().UnixMilli(),
		FieldList:   fieldList,
		JobID:       arg.JobId,
		TaskBatchID: arg.TaskBatchId,
		TaskID:      arg.TaskId,
	}
}
