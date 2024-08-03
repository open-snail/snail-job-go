package job

import (
	"sync"
)

type HookLogService struct {
	Cache     sync.Map
	Wg        sync.WaitGroup //  优雅关闭
	MessageCh chan *JobLogTask
}

func NewHookLogService() *HookLogService {
	hs := &HookLogService{
		Cache:     sync.Map{},
		MessageCh: make(chan *JobLogTask),
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
				var items []*JobLogTask
				items = append(items, msg)
				SendBatchLogReport(items)
			}()
		}
	}
}
