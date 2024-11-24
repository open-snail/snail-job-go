package job

import (
	"fmt"
	"sync"

	"opensnail.com/snail-job/snail-job-go/dto"
)

type HookLogService struct {
	Cache      sync.Map
	Wg         sync.WaitGroup //  优雅关闭
	LogEntryCh chan *dto.JobLogTask
	client     SnailJobClient
}

func NewHookLogService(client SnailJobClient) *HookLogService {
	hls := &HookLogService{
		Cache:      sync.Map{},
		LogEntryCh: make(chan *dto.JobLogTask),
		client:     client,
	}
	return hls
}

func (hls *HookLogService) Init() {
	for entry := range hls.LogEntryCh {
		hls.Wg.Add(1)
		go func() {
			defer hls.Wg.Done()
			var items []*dto.JobLogTask
			items = append(items, entry)
			hls.client.SendBatchLogReport(items)
		}()
	}
}

func FormatExcInfo(err error) string {
	if err == nil {
		return ""
	}
	return fmt.Sprintf("%s", err)
}
