package main

import (
	"time"

	"github.com/Dasongzi1366/AutoGo/motion"
)

const (
	僵尸3左右刷怪左最短  = 8 * time.Second
	僵尸3左右刷怪左最长  = 30 * time.Second
	僵尸3左右刷怪右最短  = 8 * time.Second
	僵尸3左右刷怪右最长  = 31 * time.Second
	僵尸3左右刷怪长按时长 = 700 * time.Millisecond
	僵尸3左右刷怪攻击间隔 = 70 * time.Millisecond
)

func 启动僵尸3左右刷怪() {
	if 程序退出中.Load() {
		设置僵尸3层输出("程序正在退出，忽略左右刷怪")
		return
	}
	if !脚本运行中.CompareAndSwap(false, true) {
		设置僵尸3层输出("已有脚本正在运行，不能启动左右刷怪")
		return
	}

	脚本已暂停.Store(false)
	runID := 脚本运行序号.Add(1)
	设置僵尸3层输出("左右刷怪开始：左8~30秒，右8~31秒，长按方向+X后攻击2次")
	go 运行僵尸3左右刷怪(runID)
}

func 运行僵尸3左右刷怪(runID int64) {
	defer func() {
		释放所有按键()
		if 脚本运行序号.Load() == runID {
			脚本运行中.Store(false)
			脚本已暂停.Store(false)
		}
		设置僵尸3层输出("左右刷怪已停止")
	}()

	方向键 := motion.KEYCODE_DPAD_LEFT
	for 脚本仍应运行(runID) {
		持续 := 随机僵尸3左右刷怪时长(方向键)
		设置僵尸3层输出("左右刷怪：向%s%d秒", 键名(方向键), int(持续.Seconds()))
		截止 := time.Now().Add(持续)
		for 脚本仍应运行(runID) && time.Now().Before(截止) {
			if 僵尸3检查BOSS并换线(runID) {
				continue
			}
			按组合键同时不空格(方向键, motion.KEYCODE_X, int(僵尸3左右刷怪长按时长/time.Millisecond))
			if !脚本仍应运行(runID) || !僵尸3左右刷怪攻击两次(runID) {
				return
			}
		}
		if 方向键 == motion.KEYCODE_DPAD_LEFT {
			方向键 = motion.KEYCODE_DPAD_RIGHT
		} else {
			方向键 = motion.KEYCODE_DPAD_LEFT
		}
	}
}

func 随机僵尸3左右刷怪时长(方向键 int) time.Duration {
	最短, 最长 := 僵尸3左右刷怪右最短, 僵尸3左右刷怪右最长
	if 方向键 == motion.KEYCODE_DPAD_LEFT {
		最短, 最长 = 僵尸3左右刷怪左最短, 僵尸3左右刷怪左最长
	}
	范围秒 := int((最长 - 最短) / time.Second)
	return 最短 + time.Duration(移动随机.Intn(范围秒+1))*time.Second
}

func 僵尸3左右刷怪攻击两次(runID int64) bool {
	for i := 0; i < 2; i++ {
		if !脚本仍应运行(runID) {
			return false
		}
		点按空格()
		if i == 0 {
			time.Sleep(僵尸3左右刷怪攻击间隔)
		}
	}
	return true
}
