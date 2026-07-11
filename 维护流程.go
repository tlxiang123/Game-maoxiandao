package main

import (
	"sync/atomic"
	"time"
)

const (
	刷怪周期最短 = 24 * time.Minute
	刷怪周期最长 = 43 * time.Minute
	前往四层重试 = 8
)

var (
	买卖物品下一步执行中  atomic.Bool
	卖物品测试下一步执行中 atomic.Bool
)

func 随机刷怪周期时长() time.Duration {
	spanMinutes := int64((刷怪周期最长 - 刷怪周期最短) / time.Minute)
	if spanMinutes <= 0 {
		return 刷怪周期最短
	}
	return 刷怪周期最短 + time.Duration(移动随机.Int63n(spanMinutes+1))*time.Minute
}

func 新刷怪周期截止() time.Time {
	duration := 随机刷怪周期时长()
	输出("怪物 新刷怪周期", "分钟=", int(duration.Minutes()))
	return time.Now().Add(duration)
}

func 前往四层(shouldContinue func() bool) bool {
	if shouldContinue == nil {
		shouldContinue = func() bool { return true }
	}
	for attempt := 1; attempt <= 前往四层重试 && shouldContinue(); attempt++ {
		位置, ok := 当前层位置()
		if !ok {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		switch 位置.层 {
		case 4:
			return true
		case 3:
			输出("买卖物品 前往4层：3层到4层")
			if 三层到四层() {
				return true
			}
		case 2:
			输出("买卖物品 前往4层：2层到3层")
			二层到三层()
		case 1:
			输出("买卖物品 前往4层：1层到2层")
			一层到二层()
		default:
			输出("买卖物品 前往4层失败：当前层异常", "层=", 位置.层)
			return false
		}
		time.Sleep(500 * time.Millisecond)
	}
	return false
}

func 执行四层买卖买药后回三层(shouldContinue func() bool) bool {
	if shouldContinue == nil {
		shouldContinue = func() bool { return true }
	}
	位置, ok := 当前层位置()
	if !ok || 位置.层 != 4 {
		输出("买卖物品 维护失败：当前不在4层")
		return false
	}
	输出("买卖物品 4层开始卖东西和买东西")
	if !执行完整买卖物品流程(shouldContinue) {
		return false
	}
	if !shouldContinue() {
		return false
	}
	输出("买卖物品 完成，准备4层到3层")
	if 四层到三层() {
		输出("怪物 已回到3层，继续打怪")
		return true
	}
	输出("买卖物品 4层到3层失败")
	return false
}

func 启动测试三四层上下() {
	if 脚本运行中.Load() {
		输出("买卖物品 测试上下失败：脚本运行中，请先暂停")
		return
	}
	go func() {
		位置, ok := 当前层位置()
		if !ok {
			输出("买卖物品 测试上下失败：无法判断当前层")
			return
		}
		switch 位置.层 {
		case 3:
			输出("买卖物品 测试上下：3层到4层，再4层到3层")
			if 三层到四层() {
				time.Sleep(600 * time.Millisecond)
				四层到三层()
			}
		case 4:
			输出("买卖物品 测试上下：4层到3层，再3层到4层")
			if 四层到三层() {
				time.Sleep(600 * time.Millisecond)
				三层到四层()
			}
		default:
			输出("买卖物品 测试上下失败：请先站在3层或4层", "当前层=", 位置.层)
		}
	}()
}

func 执行手动买卖物品下一步() {
	if 程序退出中.Load() {
		输出("买卖物品 下一步失败：程序正在退出")
		return
	}
	if 脚本运行中.Load() {
		输出("买卖物品 下一步失败：脚本运行中，请先暂停")
		return
	}
	if !买卖物品下一步执行中.CompareAndSwap(false, true) {
		输出("买卖物品 下一步执行中，请稍等")
		return
	}
	go func() {
		defer 买卖物品下一步执行中.Store(false)
		if !买卖物品流程运行中() {
			启动买卖物品流程()
		}
		买卖物品下一步()
	}()
}

func 执行测试卖物品下一步() {
	if 程序退出中.Load() {
		输出("卖物品测试 下一步失败：程序正在退出")
		return
	}
	if 脚本运行中.Load() {
		输出("卖物品测试 下一步失败：脚本运行中，请先暂停")
		return
	}
	if !卖物品测试下一步执行中.CompareAndSwap(false, true) {
		输出("卖物品测试 下一步执行中，请稍等")
		return
	}
	go func() {
		defer 卖物品测试下一步执行中.Store(false)
		if !卖物品测试流程运行中() {
			启动卖物品测试流程()
		}
		卖物品测试下一步()
	}()
}

func 启动手动买卖买药后打怪() {
	if 程序退出中.Load() {
		输出("买卖物品 启动失败：程序正在退出")
		return
	}
	go func() {
		脚本原本运行中 := 脚本运行中.Load()
		if 脚本原本运行中 {
			停止脚本()
			等待脚本停止()
			time.Sleep(500 * time.Millisecond)
		}
		shouldContinue := func() bool { return !程序退出中.Load() }
		输出("买卖物品 当前位置开始卖东西和买东西")
		if !执行完整买卖物品流程(shouldContinue) {
			输出("买卖物品 启动失败：买卖买药流程未完成")
			return
		}
		输出("买卖物品 当前位置买卖买药完成")
		if 脚本原本运行中 {
			启动脚本()
		}
	}()
}
