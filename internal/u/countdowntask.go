package u

import (
	"fmt"
	"sync"
	"time"
)

// CountdownTask 泛型倒计时任务管理器
type CountdownTask[T any] struct {
	duration     time.Duration  // 倒计时总时长
	task         func(data T)   // 要执行的任务函数（带泛型参数）
	startTime    time.Time      // 任务开始时间
	timer        *time.Timer    // 倒计时定时器
	triggered    bool           // 是否已被外部触发执行
	completed    bool           // 是否已完成
	triggerChan  chan T         // 外部触发通道（带泛型数据）
	cancelChan   chan struct{}  // 取消任务通道
	callbackChan chan struct{}  // 完成回调通道
	mutex        sync.RWMutex   // 读写锁
	wg           sync.WaitGroup // 等待组
	data         T              // 默认任务数据
	useDefault   bool           // 是否使用默认数据
}

// NewCountdownTask 创建新的泛型倒计时任务
func NewCountdownTask[T any](duration time.Duration, task func(data T)) *CountdownTask[T] {
	return &CountdownTask[T]{
		duration:     duration,
		task:         task,
		triggerChan:  make(chan T, 1), // 带缓冲的通道
		cancelChan:   make(chan struct{}, 1),
		callbackChan: make(chan struct{}, 1),
		useDefault:   false,
	}
}

// NewCountdownTaskWithData 创建带默认数据的泛型倒计时任务
func NewCountdownTaskWithData[T any](duration time.Duration, task func(data T), defaultData T) *CountdownTask[T] {
	return &CountdownTask[T]{
		duration:     duration,
		task:         task,
		triggerChan:  make(chan T, 1),
		cancelChan:   make(chan struct{}, 1),
		callbackChan: make(chan struct{}, 1),
		data:         defaultData,
		useDefault:   true,
	}
}

// SetData 设置默认任务数据（在任务开始前）
func (ct *CountdownTask[T]) SetData(data T) {
	ct.mutex.Lock()
	defer ct.mutex.Unlock()
	ct.data = data
	ct.useDefault = true
}

// Start 启动倒计时任务（非阻塞）
func (ct *CountdownTask[T]) Start() {
	ct.mutex.Lock()
	defer ct.mutex.Unlock()

	if ct.timer != nil && !ct.completed {
		return // 已经启动且未完成，避免重复启动
	}

	// 重置所有状态
	ct.startTime = time.Now()
	ct.triggered = false
	ct.completed = false
	ct.timer = time.NewTimer(ct.duration)

	// 清空通道（避免旧数据）
	select {
	case <-ct.cancelChan:
	default:
	}
	select {
	case <-ct.callbackChan:
	default:
	}
	select {
	case <-ct.triggerChan:
	default:
	}

	ct.wg.Add(1)
	go ct.run()

	fmt.Printf("[%s] 倒计时开始: %s 后执行任务\n",
		time.Now().Format("15:04:05"), ct.duration)
}

// run 执行倒计时
func (ct *CountdownTask[T]) run() {
	defer ct.wg.Done()
	defer func() {
		// 确保定时器被停止
		if !ct.timer.Stop() {
			select {
			case <-ct.timer.C:
			default:
			}
		}
	}()

	select {
	case <-ct.timer.C: // 倒计时自然结束
		ct.mutex.Lock()
		if !ct.triggered && !ct.completed {
			var data T
			if ct.useDefault {
				data = ct.data
			}
			ct.executeTask("倒计时结束", data)
		}
		ct.mutex.Unlock()

	case data := <-ct.triggerChan: // 外部触发执行（带数据）
		ct.mutex.Lock()
		if !ct.completed {
			ct.executeTask("外部触发", data)
		}
		ct.mutex.Unlock()

	case <-ct.cancelChan: // 取消任务
		ct.mutex.Lock()
		if !ct.completed {
			fmt.Println("任务已被取消")
			ct.completed = true
			ct.callbackChan <- struct{}{} // 通知任务取消
		}
		ct.mutex.Unlock()
	}
}

// executeTask 安全执行任务
func (ct *CountdownTask[T]) executeTask(reason string, data T) {
	// 记录开始执行时间
	start := time.Now()

	elapsed := time.Since(ct.startTime)
	fmt.Printf("%s: 提前 %.2f秒执行任务 (数据: %v)\n",
		reason, float64(ct.duration-elapsed)/float64(time.Second), data)

	// 执行实际任务
	ct.task(data)

	// 标记任务完成
	ct.completed = true
	ct.callbackChan <- struct{}{} // 通知任务完成

	fmt.Printf("任务执行完成, 耗时: %v\n", time.Since(start))
}

// Trigger 外部触发立即执行任务（带泛型数据）
func (ct *CountdownTask[T]) Trigger(data T) {
	ct.mutex.RLock()
	defer ct.mutex.RUnlock()

	if ct.completed || ct.triggered {
		return // 任务已完成或已触发
	}

	// 如果倒计时还没结束，发送触发信号
	ct.triggered = true
	select {
	case ct.triggerChan <- data:
		fmt.Println("已发送触发信号和数据")
	default:
		fmt.Println("触发信号通道已满，忽略触发请求")
	}
}

func (ct *CountdownTask[T]) TriggerSign(fn func(data T)) {
	ct.mutex.RLock()
	defer ct.mutex.RUnlock()

	if ct.completed || ct.triggered {
		return // 任务已完成或已触发
	}
	if fn != nil {
		fn(ct.data)
	}

	ct.triggered = true
	select {
	case ct.triggerChan <- ct.data:
		fmt.Println("已发送触发信号和默认数据")
	default:
		fmt.Println("触发信号通道已满，忽略触发请求")
	}
}

// TriggerDefault 使用默认数据触发任务
func (ct *CountdownTask[T]) TriggerDefault() {
	ct.mutex.RLock()
	defer ct.mutex.RUnlock()

	if ct.completed || ct.triggered {
		return // 任务已完成或已触发
	}

	ct.triggered = true
	select {
	case ct.triggerChan <- ct.data:
		fmt.Println("已发送触发信号和默认数据")
	default:
		fmt.Println("触发信号通道已满，忽略触发请求")
	}
}

// Cancel 取消倒计时任务
func (ct *CountdownTask[T]) Cancel() {
	ct.mutex.RLock()
	defer ct.mutex.RUnlock()

	if ct.completed {
		return // 任务已完成
	}

	select {
	case ct.cancelChan <- struct{}{}:
		fmt.Println("已发送取消信号")
	default:
		fmt.Println("取消信号通道已满，忽略取消请求")
	}
}

// WaitForCompletion 等待任务完成
func (ct *CountdownTask[T]) WaitForCompletion() {
	<-ct.callbackChan
}

// ElapsedTime 获取自任务开始后经过的时间
func (ct *CountdownTask[T]) ElapsedTime() time.Duration {
	ct.mutex.RLock()
	defer ct.mutex.RUnlock()

	if ct.startTime.IsZero() {
		return 0
	}
	return time.Since(ct.startTime)
}

// RemainingTime 获取剩余倒计时时间
func (ct *CountdownTask[T]) RemainingTime() time.Duration {
	ct.mutex.RLock()
	defer ct.mutex.RUnlock()

	if ct.completed {
		return 0
	}

	elapsed := time.Since(ct.startTime)
	if elapsed >= ct.duration {
		return 0
	}
	return ct.duration - elapsed
}

// IsCompleted 检查任务是否已完成
func (ct *CountdownTask[T]) IsCompleted() bool {
	ct.mutex.RLock()
	defer ct.mutex.RUnlock()
	return ct.completed
}

// Reset 重置倒计时
func (ct *CountdownTask[T]) Reset(duration time.Duration) {
	ct.Cancel()
	ct.duration = duration
}
