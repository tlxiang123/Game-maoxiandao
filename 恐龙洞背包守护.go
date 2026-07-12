package main

import (
	"sync/atomic"
	"time"

	"github.com/Dasongzi1366/AutoGo/motion"
)

const (
	恐龙洞背包检查间隔   = 5 * time.Second
	恐龙洞背包关闭确认间隔 = 200 * time.Millisecond
	恐龙洞背包关闭确认超时 = 2 * time.Second
	恐龙洞背包关闭重试间隔 = 500 * time.Millisecond
)

var 恐龙洞买卖物品流程中 atomic.Bool

func 启动恐龙洞背包守护(runID int64) {
	go 恐龙洞背包守护循环(runID)
}

func 恐龙洞启动时确保背包已关闭(runID int64) bool {
	for 脚本仍应运行(runID) {
		ok, x, y := 查找僵尸3已打开背包()
		if !ok {
			设置恐龙洞输出("启动背包检查：背包已关闭，开始判断楼层")
			return true
		}
		设置恐龙洞输出("启动背包检查：背包已打开 x=%d y=%d，先关闭再移动", x, y)
		if 恐龙洞按I关闭背包并确认("启动背包检查", func() bool { return 脚本仍应运行(runID) }) {
			return true
		}
		time.Sleep(恐龙洞背包关闭重试间隔)
	}
	return false
}

func 恐龙洞背包守护循环(runID int64) {
	ticker := time.NewTicker(恐龙洞背包检查间隔)
	defer ticker.Stop()
	for range ticker.C {
		if !脚本仍应运行(runID) {
			return
		}
		if 恐龙洞买卖物品流程中.Load() || 引擎 == nil {
			continue
		}
		if 处理MS误触(设置恐龙洞输出, func() bool { return 脚本仍应运行(runID) }) {
			continue
		}
		if ok, x, y := 查找僵尸3已打开背包(); ok {
			设置恐龙洞输出("走路/打怪检测到背包已打开 x=%d y=%d，按I关闭并确认", x, y)
			恐龙洞按I关闭背包并确认("恐龙洞背包守护", func() bool { return 脚本仍应运行(runID) })
		}
	}
}

func 恐龙洞按I关闭背包并确认(scene string, shouldContinue func() bool) bool {
	if shouldContinue == nil {
		shouldContinue = func() bool { return true }
	}
	if ok, _, _ := 查找僵尸3已打开背包(); !ok {
		设置恐龙洞输出("%s：已经打开背包特征不存在，背包已关闭", scene)
		return true
	}
	点按键(motion.KEYCODE_I, 当前显示ID())
	deadline := time.Now().Add(恐龙洞背包关闭确认超时)
	for shouldContinue() && time.Now().Before(deadline) {
		time.Sleep(恐龙洞背包关闭确认间隔)
		if ok, _, _ := 查找僵尸3已打开背包(); !ok {
			设置恐龙洞输出("%s：背包特征已消失，关闭成功", scene)
			return true
		}
	}
	if shouldContinue() {
		设置恐龙洞输出("%s：背包特征仍存在，关闭失败", scene)
	}
	return false
}
