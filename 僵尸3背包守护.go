package main

import (
	"sync/atomic"
	"time"

	"github.com/Dasongzi1366/AutoGo/motion"
)

const (
	僵尸3背包检查间隔   = 5 * time.Second
	僵尸3背包关闭确认间隔 = 200 * time.Millisecond
	僵尸3背包关闭确认超时 = 2 * time.Second
	僵尸3背包关闭重试间隔 = 500 * time.Millisecond
)

var 已经打开背包 = &FMColor{
	Name:        "已经打开背包",
	X1:          965,
	Y1:          143,
	X2:          1103,
	Y2:          215,
	MainColor:   "EE6251-202020",
	OffsetColor: "-6,-9,CC6359-202020,-7,0,BFBFBF-202020,-5,-7,A4ABB4-202020,-3,-1,265799-202020,-9,-4,404040-202020,-4,-6,6D7074-202020,-2,-12,000000-202020,-12,-1,EE5544-202020",
	Sim:         0.90,
	Dir:         0,
}

var 僵尸3买卖物品流程中 atomic.Bool

func 启动僵尸3背包守护(runID int64) {
	go 僵尸3背包守护循环(runID)
}

func 僵尸3启动时确保背包已关闭(runID int64) bool {
	for 脚本仍应运行(runID) {
		ok, x, y := 查找僵尸3已打开背包()
		if !ok {
			设置僵尸3层输出("启动背包检查：未发现打开背包，开始判断楼层")
			return true
		}
		设置僵尸3层输出("启动背包检查：背包已打开 x=%d y=%d，先关闭再开始移动", x, y)
		if 僵尸3按I关闭背包并确认("启动背包检查", func() bool { return 脚本仍应运行(runID) }) {
			return true
		}
		time.Sleep(僵尸3背包关闭重试间隔)
	}
	return false
}

func 僵尸3背包守护循环(runID int64) {
	ticker := time.NewTicker(僵尸3背包检查间隔)
	defer ticker.Stop()

	for range ticker.C {
		if !脚本仍应运行(runID) {
			return
		}
		if 僵尸3买卖物品流程中.Load() || 僵尸3传送流程运行中.Load() || 引擎 == nil {
			continue
		}
		if ok, x, y := 查找僵尸3已打开背包(); ok {
			设置僵尸3层输出("走路/打怪检测到背包已打开 x=%d y=%d，按I关闭并确认", x, y)
			僵尸3按I关闭背包并确认("走路/打怪背包守护", func() bool { return 脚本仍应运行(runID) })
		}
	}
}

func 僵尸3按I关闭背包并确认(场景 string, shouldContinue func() bool) bool {
	if shouldContinue == nil {
		shouldContinue = func() bool { return true }
	}
	if ok, _, _ := 查找僵尸3已打开背包(); !ok {
		设置僵尸3层输出("%s：已经打开背包特征不存在，背包已关闭", 场景)
		return true
	}
	释放所有按键()
	点按键(motion.KEYCODE_I, 当前显示ID())
	deadline := time.Now().Add(僵尸3背包关闭确认超时)
	for shouldContinue() && time.Now().Before(deadline) {
		time.Sleep(僵尸3背包关闭确认间隔)
		if ok, _, _ := 查找僵尸3已打开背包(); !ok {
			设置僵尸3层输出("%s：已经打开背包特征已消失，关闭成功", 场景)
			return true
		}
	}
	if shouldContinue() {
		设置僵尸3层输出("%s：已经打开背包特征仍存在，关闭失败", 场景)
	}
	return false
}

func 查找僵尸3已打开背包() (bool, int, int) {
	暂停调试红框()
	defer 恢复调试红框()
	return 引擎.FindFeature(已经打开背包)
}
