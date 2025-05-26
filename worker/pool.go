package worker

import (
	"github.com/echoH00/alert-notifier/alertstore"
	"github.com/echoH00/alert-notifier/config"
	"log"
	"sync"
	"time"
	"strings"
)

var (
	jobQueue    = make(chan alertstore.AlertJob, config.JobQueueSize)
	deadLetters = make(chan alertstore.AlertJob, config.DeadLetterSize)

	workerCount  = 0
	lastActivity = make(map[int]time.Time)
	mu           sync.Mutex
)

// 初始化最小数量的 worker
func InitPool() {
	for i := 0; i < config.MinWorkers; i++ {
		if id, ok := allocateWorkerID(); ok {
			startWorker(id)
		}
	}
}

// 分配新的 worker ID（在锁内执行）
func allocateWorkerID() (int, bool) {
	mu.Lock()
	defer mu.Unlock()

	if workerCount >= config.MaxWorkers {
		return 0, false
	}
	workerCount++
	id := workerCount
	lastActivity[id] = time.Now()
	return id, true
}

// 提交一个任务到工作队列（会自动扩容）
func EnqueueJob(job alertstore.AlertJob) {
	select {
	case jobQueue <- job:
		desp := job.Alert.Annotations["description"]
		desp = strings.ReplaceAll(desp, "\n", "")
		desp = strings.ReplaceAll(desp, "\r", "")
		log.Printf("Job enqueued: %v", desp)
		scaleUpWorker()
	default:
		log.Println("Job queue full, dropping job")
	}
}

// 自动扩容 worker（不在锁内启动 goroutine，避免死锁）
func scaleUpWorker() {
	if id, ok := allocateWorkerID(); ok {
		startWorker(id)
	}
}

// 由 worker 调用，更新时间戳
func TrackActivity(id int) {
	mu.Lock()
	lastActivity[id] = time.Now()
	mu.Unlock()
}

// 判断是否可以退出 worker（超过空闲时间 & 超过最小数量）
func CanExitWorker(id int) bool {
	mu.Lock()
	defer mu.Unlock()

	if workerCount <= config.MinWorkers {
		return false
	}

	if time.Since(lastActivity[id]) > config.IdleTimeout {
		delete(lastActivity, id)
		workerCount--
		return true
	}
	return false
}

// 提供给 worker 的通道访问器
func JobQueue() <-chan alertstore.AlertJob {
	return jobQueue
}

func DeadLetterQueue() chan alertstore.AlertJob {
	return deadLetters
}

