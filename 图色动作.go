package main

import (
	"time"

	"github.com/Dasongzi1366/AutoGo/motion"
)

type 层移动配置 struct {
	层  int
	Y  int
	左边 int
	右边 int
}

var 层移动配置表 = []层移动配置{
	{层: 3, Y: 150, 左边: 72, 右边: 112},
	{层: 2, Y: 165, 左边: 58, 右边: 128},
	{层: 1, Y: 180, 左边: 40, 右边: 138},
}

const (
	层Y容差    = 10
	方向键按下毫秒 = 180
)

func 找到颜色就点击(name string, x, y int, color string) bool {
	输出("比色:", name, "x=", x, "y=", y, "color=", color)
	ok, rx, ry := 引擎.CmpColor(&CColor{
		Name:  name,
		X:     x,
		Y:     y,
		Color: color,
		Sim:   0.90,
	})
	if !ok {
		输出("比色失败，不点击:", name)
		return false
	}
	输出("比色成功，点击:", name, "x=", rx, "y=", ry)
	return 引擎.ClickResult(ok, rx, ry)
}

func 查找FEFE24坐标() (bool, int, int) {
	if 引擎 == nil {
		输出("查找FEFE24坐标失败：引擎未初始化")
		return false, -1, -1
	}

	for _, 区域 := range 小地图黄点候选区域 {
		result := 扫描小地图黄点(区域, false)
		if result.Ok {
			return true, result.X, result.Y
		}
	}
	return false, -1, -1
}

func 按当前层范围移动() bool {
	位置, ok := 当前层位置()
	if !ok {
		return false
	}

	配置, ok := 当前层移动配置(位置.层)
	if !ok {
		输出("当前层没有移动范围配置", "层=", 位置.层, "x=", 位置.X, "y=", 位置.Y)
		return false
	}

	if 位置.X < 配置.左边 {
		输出("当前位置偏左，向右移动", "层=", 配置.层, "x=", 位置.X, "范围=", 配置.左边, "-", 配置.右边)
		随机向右移动(位置.X, 配置.左边, 配置.右边, 配置.左边, 配置.右边)
		return true
	}
	if 位置.X > 配置.右边 {
		输出("当前位置偏右，向左移动", "层=", 配置.层, "x=", 位置.X, "范围=", 配置.左边, "-", 配置.右边)
		随机向左移动(位置.X, 配置.左边, 配置.右边, 配置.左边, 配置.右边)
		return true
	}

	输出("当前位置在移动范围内，不移动", "层=", 配置.层, "x=", 位置.X, "范围=", 配置.左边, "-", 配置.右边)
	return false
}

func 匹配层移动配置(y int) (层移动配置, bool) {
	var best 层移动配置
	bestDiff := 层Y容差 + 1
	for _, 配置 := range 层移动配置表 {
		diff := absInt(y - 配置.Y)
		if diff < bestDiff {
			best = 配置
			bestDiff = diff
		}
	}
	return best, bestDiff <= 层Y容差
}

func 按方向键左(ms int) {
	按方向键(motion.KEYCODE_DPAD_LEFT, ms)
}

func 按方向键右(ms int) {
	按方向键(motion.KEYCODE_DPAD_RIGHT, ms)
}

func 按方向键(code int, ms int) {
	displayID := 屏幕ID
	if 引擎 != nil {
		displayID = 引擎.displayID()
	}
	if ms <= 0 {
		点按键(code, displayID)
		return
	}
	motion.KeyActionDown(code, displayID)
	time.Sleep(time.Duration(ms) * time.Millisecond)
	motion.KeyActionUp(code, displayID)
}

func 点击(x, y int) {
	输出("点击坐标:", "x=", x, "y=", y)
	引擎.Click(x, y)
}

func 等待(ms int) {
	if ms <= 0 {
		return
	}
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func 等待特征(feature any, timeout time.Duration, interval time.Duration) (bool, int, int) {
	if interval <= 0 {
		interval = 300 * time.Millisecond
	}
	deadline := time.Now().Add(timeout)
	for {
		ok, x, y := 引擎.FindFeature(feature)
		if ok {
			return true, x, y
		}
		if timeout > 0 && time.Now().After(deadline) {
			return false, -1, -1
		}
		time.Sleep(interval)
	}
}

func 等待并点击(feature any, timeout time.Duration, interval time.Duration) bool {
	ok, x, y := 等待特征(feature, timeout, interval)
	return 引擎.ClickResult(ok, x, y)
}
