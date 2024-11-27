package job

import (
	"sync"
	"time"

	"opensnail.com/snail-job/snail-job-go/dto"
)

// 线程安全缓冲区
type SafeBuffer struct {
	mu     sync.Mutex
	buffer []*dto.JobLogTask
}

// NewSafeBuffer 返回一个初始化的 SafeBuffer 实例
func NewSafeBuffer() *SafeBuffer {
	return &SafeBuffer{
		buffer: []*dto.JobLogTask{},
	}
}

// Add 向缓冲区添加一条数据
func (sb *SafeBuffer) Add(entry *dto.JobLogTask) {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	sb.buffer = append(sb.buffer, entry)
}

// GetAll 获取缓冲区中的所有数据并清空缓冲区
func (sb *SafeBuffer) GetAll() []*dto.JobLogTask {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	batch := sb.buffer
	sb.buffer = nil // 清空缓冲区
	return batch
}

// Len 返回缓冲区中数据的数量
func (sb *SafeBuffer) Len() int {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	return len(sb.buffer)
}

type HookLogService struct {
	Wg         sync.WaitGroup
	LogEntryCh chan *dto.JobLogTask
	client     SnailJobClient
	buffer     *SafeBuffer // 使用 SafeBuffer 替代原来的 buffer
}

func NewHookLogService(client SnailJobClient) *HookLogService {
	return &HookLogService{
		LogEntryCh: make(chan *dto.JobLogTask),
		client:     client,
		buffer:     NewSafeBuffer(),
	}
}

func (hls *HookLogService) Init() {
	const (
		maxBatchSize = 10               // 最大批量大小
		maxWaitTime  = 10 * time.Second // 最大等待时间
	)

	timer := time.NewTimer(maxWaitTime)
	defer timer.Stop()

	for {
		select {
		case entry, ok := <-hls.LogEntryCh:
			if !ok {
				hls.flushBuffer()
				return
			}
			hls.buffer.Add(entry)

			// 如果达到批量大小，立即发送
			if hls.buffer.Len() >= maxBatchSize {
				hls.flushBuffer()
			}
		case <-timer.C:
			hls.flushBuffer()
			timer.Reset(maxWaitTime)
		}
	}
}

// flushBuffer 清空缓冲区并异步发送日志
func (hls *HookLogService) flushBuffer() {
	if hls.buffer.Len() == 0 {
		return
	}

	// 获取并清空缓冲区
	batch := hls.buffer.GetAll()

	// 异步发送数据
	hls.Wg.Add(1)
	go func(batch []*dto.JobLogTask) {
		defer hls.Wg.Done()
		hls.client.SendBatchLogReport(batch)
	}(batch)
}
