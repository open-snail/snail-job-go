package job

import (
	"opensnail.com/snail-job/snail-job-go/dto"
	"sync"
)

type HookLogService struct {
	Cache     sync.Map
	Wg        sync.WaitGroup //  优雅关闭
	MessageCh chan *dto.JobLogTask
	client    SnailJobClient
}

func NewHookLogService(client SnailJobClient) *HookLogService {
	hs := &HookLogService{
		Cache:     sync.Map{},
		MessageCh: make(chan *dto.JobLogTask),
		client:    client,
	}
	return hs
}

func (s *HookLogService) Init() {
	for {
		select {
		case msg := <-s.MessageCh:
			s.Wg.Add(1)
			go func() {
				defer s.Wg.Done()
				var items []*dto.JobLogTask
				items = append(items, msg)
				s.client.SendBatchLogReport(items)
			}()
		}
	}
}
