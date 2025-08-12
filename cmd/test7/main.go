package main

import (
	"fmt"
	"github.com/xxl6097/uclient/internal/u"
	"time"
)

func main() {
	// 1. 创建倒计时任务（10秒后执行）
	fmt.Println("=== 第一次倒计时测试 - 正常触发 ===")
	task := u.NewCountdownTask(5*time.Second, func(data string) {
		fmt.Printf("[%s] %s 任务执行: 执行预定操作\n", time.Now().Format("15:04:05.000"), data)
	})

	// 启动任务
	task.Start()

	// 监控倒计时
	go monitorCountdown(task)

	// 等待第一次任务完成
	task.WaitForCompletion()

	time.Sleep(5 * time.Second)
	fmt.Println("\n=== 第二次倒计时测试 - 外部触发 ===")
	// 重置并测试外部触发
	task.Reset(15 * time.Second)
	task.Start()

	// 模拟外部事件（3秒后手动触发）
	go func() {
		time.Sleep(3 * time.Second)
		fmt.Printf("\n[%s] 外部事件: 触发任务执行\n", time.Now().Format("15:04:05.000"))
		task.Trigger("abcvvvv")
	}()

	// 等待第二次任务完成
	task.WaitForCompletion()

	fmt.Println("\n=== 第三次倒计时测试 - 取消任务 ===")
	// 重置并测试取消
	task.Reset(8 * time.Second)
	task.Start()

	// 模拟外部取消
	go func() {
		time.Sleep(2 * time.Second)
		fmt.Printf("\n[%s] 外部事件: 取消任务\n", time.Now().Format("15:04:05.000"))
		task.Cancel()
	}()

	// 等待取消完成
	task.WaitForCompletion()

	fmt.Println("\n所有测试完成")
}

// 监控倒计时进度
func monitorCountdown(task *u.CountdownTask[string]) {
	for {
		time.Sleep(200 * time.Millisecond)

		if task.IsCompleted() {
			fmt.Printf("监控结束: 任务已完成\n")
			return
		}

		remaining := task.RemainingTime()
		if remaining > 0 {
			fmt.Printf("剩余: %.1f秒   \r", remaining.Seconds())
		} else {
			fmt.Println("任务应已执行       ")
		}
	}
}
