package pkg

import (
	"sync"
	"time"
)

var (
	dedupMu       sync.Mutex
	dedupLastSeen = make(map[string]time.Time)
)

// DedupDo 在 1 秒内相同 key 的事件只执行一次。
// 返回 true 表示本次执行，false 表示被去重忽略。
func DedupDo(key string, fn func()) bool {
	return DedupDoWithin(key, time.Second, fn)
}

// DedupDoWithin 在指定时间窗口内相同 key 的事件只执行一次。
func DedupDoWithin(key string, window time.Duration, fn func()) bool {
	now := time.Now()
	dedupMu.Lock()
	if last, ok := dedupLastSeen[key]; ok && now.Sub(last) < window {
		dedupMu.Unlock()
		return false
	}
	dedupLastSeen[key] = now
	// 清理过期 key，避免内存累积
	for k, t := range dedupLastSeen {
		if now.Sub(t) > window*10 {
			delete(dedupLastSeen, k)
		}
	}
	dedupMu.Unlock()

	if fn != nil {
		fn()
	}
	return true
}
