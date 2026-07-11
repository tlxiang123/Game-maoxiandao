package main

import (
	"sync"
	"sync/atomic"
	"time"
)

var 僵尸3点开背包按钮 = &FMColor{Name: "点开背包按钮", X1: 877, Y1: 630, X2: 931, Y2: 668, MainColor: "AADDEE-000000", OffsetColor: "10,-2,AADDEE-000000,14,2,AADDEE-000000,-2,3,AADDEE-000000,7,3,FFCC55-000000,14,3,DD9933-000000,-2,12,AADDEE-000000,4,8,EE9933-000000,11,10,EE9933-000000", Sim: 0.90, Dir: 0}
var 僵尸3第五页面灰色 = &FMColor{Name: "第五页面灰色", X1: 1069, Y1: 69, X2: 1276, Y2: 148, MainColor: "E4E4E4-000000", OffsetColor: "11,0,F0F0F0-000000,17,0,848484-000000,0,1,E4E4E4-000000,6,3,828282-000000,17,6,FCFCFC-000000,0,9,707070-000000,16,7,FAFAFA-000000,17,9,7D7D7D-000000", Sim: 0.90, Dir: 0}
var 僵尸3第五页面已经点开 = &FMColor{Name: "第五页面已经点开", X1: 1069, Y1: 69, X2: 1276, Y2: 148, MainColor: "F5D2DC-000000", OffsetColor: "11,-1,FBDDE5-000000,18,0,E6DDDF-000000,5,1,CEAFB7-000000,17,4,FEF7F9-000000,18,4,FEF7F9-000000,0,8,994752-000000,17,5,FEF7F9-000000,23,6,E5BEC8-000000", Sim: 0.90, Dir: 0}
var 僵尸3第五页面粉色 = &FMColor{Name: "第五页面粉色", X1: 1069, Y1: 69, X2: 1276, Y2: 148, MainColor: "F5D2DC-000000", OffsetColor: "11,-1,FBDDE5-000000,18,0,E6DDDF-000000,5,1,CEAFB7-000000,17,4,FEF7F9-000000,18,4,FEF7F9-000000,0,8,994752-000000,17,5,FEF7F9-000000,23,6,E5BEC8-000000", Sim: 0.90, Dir: 0}
var 僵尸3双击猫猫按钮 = &FMColor{Name: "双击猫猫按钮", X1: 1055, Y1: 139, X2: 1263, Y2: 389, MainColor: "544831-000000", OffsetColor: "13,-3,000000-000000,20,3,3C2814-000000,-6,10,000000-000000,1,7,76342E-000000,20,4,5F4223-000000,0,14,020201-000000,1,17,000000-000000,14,17,87100B-000000", Sim: 0.90, Dir: 0}
var 僵尸3关闭背包备用 = &FMColor{Name: "关闭背包", X1: 1226, Y1: 65, X2: 1274, Y2: 130, MainColor: "0C70A2-000000", OffsetColor: "1,0,328CB5-000000,12,0,0C70A2-000000,0,4,2EAAEE-000000,6,4,227CBF-000000,7,4,227CBF-000000,0,10,55CCFF-000000,1,15,FFFFFF-000000,7,10,44AADD-000000", Sim: 0.90, Dir: 0}

var 僵尸3第五页灰色候选 = []买卖物品特征候选{
	新买卖物品特征候选("僵尸3第五页面灰色", 僵尸3第五页面灰色),
}

var 僵尸3第五页已打开候选 = []买卖物品特征候选{
	新买卖物品特征候选("僵尸3第五页面已经点开", 僵尸3第五页面已经点开),
	新买卖物品特征候选("僵尸3第五页面粉色", 僵尸3第五页面粉色),
}

var 僵尸3第五页粉色候选 = []买卖物品特征候选{
	新买卖物品特征候选("僵尸3第五页面粉色", 僵尸3第五页面粉色),
	新买卖物品特征候选("僵尸3第五页面已经点开", 僵尸3第五页面已经点开),
}

var 僵尸3双击猫猫候选 = []买卖物品特征候选{
	新买卖物品特征候选("僵尸3双击猫猫按钮", 僵尸3双击猫猫按钮),
}

var 僵尸3点开背包候选 = []买卖物品特征候选{
	新买卖物品特征候选("僵尸3点开背包按钮", 僵尸3点开背包按钮),
}

var 僵尸3关闭背包候选 = []买卖物品特征候选{
	新买卖物品特征候选("公共关闭背包", 关闭背包),
	新买卖物品特征候选("僵尸3关闭背包备用", 僵尸3关闭背包备用),
}

var (
	僵尸3卖物品测试锁    sync.Mutex
	僵尸3卖物品测试已启动  bool
	僵尸3卖物品测试当前逻辑 int
	僵尸3卖物品下一步执行中 atomic.Bool
)

func 执行僵尸3完整买卖物品流程(shouldContinue func() bool) bool {
	if shouldContinue == nil {
		shouldContinue = func() bool { return true }
	}
	开始买卖物品流程计时()
	defer 清除买卖物品流程计时()
	for index := 0; index < len(买卖物品逻辑表); {
		if 买卖物品流程已超时() {
			执行买卖物品超时收尾()
			return true
		}
		if !shouldContinue() {
			输出("僵尸3买卖物品流程中断")
			return false
		}
		logic := 买卖物品逻辑表[index]
		ok, next := 执行僵尸3买卖物品逻辑并返回下一步(index, logic)
		if !ok {
			输出("僵尸3买卖物品流程失败", "逻辑=", logic.名称)
			设置僵尸3层输出("僵尸3买卖物品失败：%s，继续重试直到超时", logic.名称)
			time.Sleep(500 * time.Millisecond)
			continue
		}
		if next <= index {
			输出("僵尸3买卖物品流程失败", "逻辑=", logic.名称, "原因=下一步异常", "当前=", index, "下一步=", next)
			设置僵尸3层输出("僵尸3买卖物品失败：%s 下一步异常，继续重试直到超时", logic.名称)
			time.Sleep(500 * time.Millisecond)
			continue
		}
		index = next
		if 买卖物品流程已超时() {
			执行买卖物品超时收尾()
			return true
		}
	}
	输出("僵尸3买卖物品流程结束", "耗时秒=", int(买卖物品流程已用时()/time.Second))
	return true
}

func 执行僵尸3卖物品测试下一步() {
	if 程序退出中.Load() {
		设置僵尸3层输出("卖物品下一步失败：程序正在退出")
		return
	}
	if 脚本运行中.Load() {
		设置僵尸3层输出("卖物品下一步失败：脚本运行中，请先点结束")
		return
	}
	if !僵尸3卖物品下一步执行中.CompareAndSwap(false, true) {
		设置僵尸3层输出("卖物品下一步执行中，请稍等")
		return
	}
	go func() {
		defer 僵尸3卖物品下一步执行中.Store(false)
		if !僵尸3卖物品测试流程运行中() {
			启动僵尸3卖物品测试流程()
		}
		index, logic, ok := 取僵尸3卖物品测试当前逻辑()
		if !ok {
			return
		}
		设置僵尸3层输出("卖物品下一步：第%d步 %s", index+1, logic.名称)
		if ok, next := 执行僵尸3买卖物品逻辑并返回下一步(index, logic); ok {
			设置僵尸3卖物品测试下一逻辑(next)
		} else {
			设置僵尸3层输出("卖物品下一步失败：第%d步 %s", index+1, logic.名称)
		}
	}()
}

func 启动僵尸3卖物品测试流程() {
	僵尸3卖物品测试锁.Lock()
	僵尸3卖物品测试已启动 = true
	僵尸3卖物品测试当前逻辑 = 0
	僵尸3卖物品测试锁.Unlock()
	设置僵尸3层输出("僵尸3卖物品测试开始")
}

func 僵尸3卖物品测试流程运行中() bool {
	僵尸3卖物品测试锁.Lock()
	defer 僵尸3卖物品测试锁.Unlock()
	return 僵尸3卖物品测试已启动
}

func 取僵尸3卖物品测试当前逻辑() (int, 买卖物品逻辑, bool) {
	僵尸3卖物品测试锁.Lock()
	defer 僵尸3卖物品测试锁.Unlock()
	if !僵尸3卖物品测试已启动 {
		设置僵尸3层输出("僵尸3卖物品测试未开始")
		return 0, 买卖物品逻辑{}, false
	}
	if 僵尸3卖物品测试当前逻辑 >= len(买卖物品逻辑表) {
		僵尸3卖物品测试已启动 = false
		设置僵尸3层输出("僵尸3卖物品测试结束")
		return 0, 买卖物品逻辑{}, false
	}
	return 僵尸3卖物品测试当前逻辑, 买卖物品逻辑表[僵尸3卖物品测试当前逻辑], true
}

func 设置僵尸3卖物品测试下一逻辑(next int) {
	僵尸3卖物品测试锁.Lock()
	僵尸3卖物品测试当前逻辑 = next
	done := 僵尸3卖物品测试当前逻辑 >= len(买卖物品逻辑表)
	if done {
		僵尸3卖物品测试已启动 = false
	}
	僵尸3卖物品测试锁.Unlock()
	if done {
		设置僵尸3层输出("僵尸3卖物品测试结束")
	}
}

func 执行僵尸3买卖物品逻辑并返回下一步(index int, logic 买卖物品逻辑) (bool, int) {
	switch logic.名称 {
	case "第1个逻辑":
		return 执行僵尸3打开背包逻辑(index, logic)
	case "第2个逻辑":
		return 执行僵尸3第五页切换逻辑(index, logic)
	case "第3个逻辑":
		return 执行僵尸3猫猫双击逻辑(index, logic)
	case "第12个逻辑":
		return 执行僵尸3关闭背包逻辑(index, logic)
	default:
		return 执行买卖物品逻辑并返回下一步(index, logic)
	}
}

func 执行僵尸3打开背包逻辑(index int, logic 买卖物品逻辑) (bool, int) {
	next := index + 1
	输出("僵尸3买卖物品 执行", logic.名称, "规则=打开背包")
	设置僵尸3层输出("卖物品第1步：查找并点击点开背包按钮")

	if found, x, y, candidate := 点击任一买卖物品特征(logic.名称, 僵尸3点开背包候选...); found {
		设置僵尸3层输出("卖物品第1步：点击%s x=%d y=%d", candidate.标签, x, y)
		return true, next
	}

	设置僵尸3层输出("卖物品第1步失败：未找到 %s", 买卖物品候选标签文本(僵尸3点开背包候选))
	return false, next
}

func 执行僵尸3第五页切换逻辑(index int, logic 买卖物品逻辑) (bool, int) {
	next := index + 1
	输出("僵尸3买卖物品 执行", logic.名称, "规则=第五页灰色点开，粉色则通过")
	设置僵尸3层输出("卖物品第2步：检查第五页")

	if found, x, y, candidate := 查找任一买卖物品特征(僵尸3第五页已打开候选...); found {
		输出("僵尸3买卖物品 第五页已打开", "候选=", candidate.标签, "x=", x, "y=", y)
		设置僵尸3层输出("卖物品第2步：第五页已打开 %s x=%d y=%d", candidate.标签, x, y)
		return true, next
	}

	found, x, y, candidate := 点击任一买卖物品特征(logic.名称, 僵尸3第五页灰色候选...)
	if !found {
		设置僵尸3层输出("卖物品第2步失败：未找到 %s", 买卖物品候选标签文本(僵尸3第五页灰色候选))
		return false, next
	}
	设置僵尸3层输出("卖物品第2步：点击第五页灰色 %s x=%d y=%d", candidate.标签, x, y)
	if found, vx, vy, opened := 等待任一买卖物品特征(买卖物品猫猫检查等待, 僵尸3第五页已打开候选...); found {
		输出("僵尸3买卖物品 第五页打开成功", "候选=", opened.标签, "x=", vx, "y=", vy)
		设置僵尸3层输出("卖物品第2步：第五页打开成功 %s x=%d y=%d", opened.标签, vx, vy)
		return true, next
	}
	输出("僵尸3买卖物品 第五页打开失败")
	设置僵尸3层输出("卖物品第2步失败：点击后未识别 %s", 买卖物品候选标签文本(僵尸3第五页已打开候选))
	return false, next
}

func 执行僵尸3猫猫双击逻辑(index int, logic 买卖物品逻辑) (bool, int) {
	next := index + 1
	if found, x, y, candidate := 查找任一买卖物品特征(僵尸3第五页粉色候选...); found {
		设置僵尸3层输出("卖物品第3步：第五页已是粉色 %s x=%d y=%d", candidate.标签, x, y)
	} else {
		设置僵尸3层输出("卖物品第3步失败：未找到 %s", 买卖物品候选标签文本(僵尸3第五页粉色候选))
		return false, next
	}

	found, x, y, candidate := 查找任一买卖物品特征(僵尸3双击猫猫候选...)
	if !found {
		输出("僵尸3买卖物品 找不到", "逻辑=", logic.名称, "候选=", 买卖物品候选标签文本(僵尸3双击猫猫候选))
		设置僵尸3层输出("卖物品第3步失败：未找到 %s", 买卖物品候选标签文本(僵尸3双击猫猫候选))
		return false, next
	}
	设置僵尸3层输出("卖物品第3步：找到%s x=%d y=%d", candidate.标签, x, y)

	points := []struct {
		label string
		x     int
		y     int
	}{
		{label: "识别点", x: x, y: y},
		{label: "偏右下", x: x + 15, y: y + 8},
		{label: "偏下", x: x + 15, y: y + 18},
	}

	for _, point := range points {
		输出("僵尸3买卖物品 双击猫猫", "逻辑=", logic.名称, "点位=", point.label, "识别x=", x, "识别y=", y, "点击x=", point.x, "点击y=", point.y)
		设置僵尸3层输出("卖物品第3步：双击猫猫 x=%d y=%d", point.x, point.y)
		快速双击买卖物品坐标(point.x, point.y, 买卖物品猫猫双击间隔)
		if 买卖物品等待检查全部通过(logic.名称, logic.检查, 买卖物品猫猫检查等待) {
			输出("僵尸3买卖物品 逻辑成功", logic.名称)
			设置僵尸3层输出("卖物品第3步：猫猫打开商店成功")
			return true, next
		}
	}

	输出("僵尸3买卖物品 逻辑失败", "逻辑=", logic.名称, "原因=双击猫猫后未出现全卖按钮")
	设置僵尸3层输出("卖物品第3步失败：双击猫猫后未出现全卖按钮")
	return false, next
}

func 执行僵尸3关闭背包逻辑(index int, logic 买卖物品逻辑) (bool, int) {
	next := index + 1
	输出("僵尸3买卖物品 执行", logic.名称, "规则=关闭背包或备用关闭背包")
	设置僵尸3层输出("卖物品第12步：关闭背包")

	if found, x, y, candidate := 点击任一买卖物品特征(logic.名称, 僵尸3关闭背包候选...); found {
		设置僵尸3层输出("卖物品第12步：点击%s x=%d y=%d", candidate.标签, x, y)
		return true, next
	}

	设置僵尸3层输出("卖物品第12步失败：未找到 %s", 买卖物品候选标签文本(僵尸3关闭背包候选))
	return false, next
}
