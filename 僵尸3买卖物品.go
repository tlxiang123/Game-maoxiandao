package main

import (
	"sync"
	"sync/atomic"
	"time"
)

var 僵尸3点开背包按钮 = &FMColor{Name: "点开背包按钮", X1: 877, Y1: 630, X2: 931, Y2: 668, MainColor: "AADDEE-000000", OffsetColor: "10,-2,AADDEE-000000,14,2,AADDEE-000000,-2,3,AADDEE-000000,7,3,FFCC55-000000,14,3,DD9933-000000,-2,12,AADDEE-000000,4,8,EE9933-000000,11,10,EE9933-000000", Sim: 0.90, Dir: 0}
var 僵尸3第五页面灰色 = &FMColor{Name: "第五页面灰色", X1: 1069, Y1: 69, X2: 1276, Y2: 148, MainColor: "E4E4E4-000000", OffsetColor: "11,0,F0F0F0-000000,17,0,848484-000000,0,1,E4E4E4-000000,6,3,828282-000000,17,6,FCFCFC-000000,0,9,707070-000000,16,7,FAFAFA-000000,17,9,7D7D7D-000000", Sim: 0.90, Dir: 0}
var 僵尸3第五页面灰色备用 = &FMColor{Name: "第五页面灰色", X1: 1002, Y1: 179, X2: 1126, Y2: 221, MainColor: "E4E4E4-000000", OffsetColor: "11,0,F0F0F0-000000,17,0,848484-000000,0,1,E4E4E4-000000,6,3,828282-000000,17,6,FCFCFC-000000,0,9,707070-000000,16,7,FAFAFA-000000,17,9,7D7D7D-000000", Sim: 0.90, Dir: 0}
var 僵尸3第五页面粉色 = &FMColor{Name: "第五页面粉色", X1: 1069, Y1: 69, X2: 1276, Y2: 148, MainColor: "F5D2DC-202020", OffsetColor: "11,-1,FBDDE5-202020,18,0,E6DDDF-202020,5,1,CEAFB7-202020,17,4,FEF7F9-202020,18,4,FEF7F9-202020,0,8,994752-202020,17,5,FEF7F9-202020,23,6,E5BEC8-202020", Sim: 0.80, Dir: 0}
var 僵尸3第五页面粉色备用 = &FMColor{Name: "第五页面粉色", X1: 1002, Y1: 179, X2: 1126, Y2: 221, MainColor: "F5D2DC-202020", OffsetColor: "11,-1,FBDDE5-202020,18,0,E6DDDF-202020,5,1,CEAFB7-202020,17,4,FEF7F9-202020,18,4,FEF7F9-202020,0,8,994752-202020,17,5,FEF7F9-202020,23,6,E5BEC8-202020", Sim: 0.80, Dir: 0}

const (
	僵尸3猫猫固定X = 936
	僵尸3猫猫固定Y = 242
)

var 僵尸3第五页灰色候选 = []买卖物品特征候选{
	新买卖物品特征候选("僵尸3第五页面灰色", 僵尸3第五页面灰色),
	新买卖物品特征候选("僵尸3第五页面灰色备用", 僵尸3第五页面灰色备用),
}

var 僵尸3第五页粉色候选 = []买卖物品特征候选{
	新买卖物品特征候选("僵尸3第五页面粉色", 僵尸3第五页面粉色),
	新买卖物品特征候选("僵尸3第五页面粉色备用", 僵尸3第五页面粉色备用),
}

var 僵尸3点开背包候选 = []买卖物品特征候选{
	新买卖物品特征候选("僵尸3点开背包按钮", 僵尸3点开背包按钮),
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
	僵尸3买卖物品流程中.Store(true)
	defer 僵尸3买卖物品流程中.Store(false)
	开始买卖物品流程计时()
	defer 清除买卖物品流程计时()
	for index := 0; index < len(买卖物品逻辑表); {
		if 买卖物品流程已超时() {
			执行僵尸3买卖物品超时收尾()
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
		next = 调整僵尸3卖杂物下一逻辑(logic.名称, next)
		index = next
		if 买卖物品流程已超时() {
			执行僵尸3买卖物品超时收尾()
			return true
		}
	}
	输出("僵尸3买卖物品流程结束", "耗时秒=", int(买卖物品流程已用时()/time.Second))
	return true
}

func 执行僵尸3买卖物品超时收尾() {
	输出("僵尸3买卖物品流程超时，关闭商店并按I关闭背包", "耗时秒=", int(买卖物品流程已用时()/time.Second), "超时秒=", int(买卖物品流程超时时长/time.Second))
	设置僵尸3层输出("卖物品超时：关闭商店并按I关闭背包")
	尝试点击买卖物品特征("超时关闭商店", 关闭商店)
	僵尸3按I键关闭背包()
}

func 僵尸3按I键关闭背包() bool {
	return 僵尸3按I关闭背包并确认("卖物品关闭背包", nil)
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
			next = 调整僵尸3卖杂物下一逻辑(logic.名称, next)
			设置僵尸3卖物品测试下一逻辑(next)
		} else {
			设置僵尸3层输出("卖物品下一步失败：第%d步 %s", index+1, logic.名称)
		}
	}()
}

func 调整僵尸3卖杂物下一逻辑(logicName string, next int) int {
	if 僵尸3卖杂物 {
		return next
	}

	closeIndex := 买卖物品逻辑索引("第11个逻辑", next)
	if logicName == "第5个逻辑" && next == 买卖物品逻辑索引("第9个逻辑", next) {
		输出("僵尸3买卖物品 卖杂物未勾选：空包袱后跳过单卖，进入关闭流程", "下一步=", closeIndex+1)
		设置僵尸3层输出("卖杂物未勾选：跳过单卖，下一步关闭商店")
		return closeIndex
	}
	if logicName == "第8个逻辑" {
		输出("僵尸3买卖物品 卖杂物未勾选：全卖后跳过单卖，进入关闭流程", "下一步=", closeIndex+1)
		设置僵尸3层输出("卖杂物未勾选：跳过单卖，下一步关闭商店")
		return closeIndex
	}
	return next
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
	输出("僵尸3买卖物品 执行", logic.名称, "规则=第五页灰色则点击，粉色则直接通过")
	设置僵尸3层输出("卖物品第2步：检查第五页")

	found, x, y, candidate := 点击任一买卖物品特征(logic.名称, 僵尸3第五页灰色候选...)
	if found {
		设置僵尸3层输出("卖物品第2步：点击第五页灰色 %s x=%d y=%d", candidate.标签, x, y)
		if found, vx, vy, pink := 等待任一买卖物品特征(买卖物品猫猫检查等待, 僵尸3第五页粉色候选...); found {
			输出("僵尸3买卖物品 第五页切换粉色成功", "候选=", pink.标签, "x=", vx, "y=", vy)
			设置僵尸3层输出("卖物品第2步：第五页已变粉色 %s x=%d y=%d", pink.标签, vx, vy)
			return true, next
		}
		输出("僵尸3买卖物品 第五页切换粉色失败")
		设置僵尸3层输出("卖物品第2步失败：点击灰色后未识别粉色")
		return false, next
	}

	if found, px, py, pink := 查找任一买卖物品特征(僵尸3第五页粉色候选...); found {
		输出("僵尸3买卖物品 第五页已经是粉色", "候选=", pink.标签, "x=", px, "y=", py)
		设置僵尸3层输出("卖物品第2步：第五页已经是粉色 %s x=%d y=%d", pink.标签, px, py)
		return true, next
	}

	输出("僵尸3买卖物品 第五页识别失败", "灰色候选=", 买卖物品候选标签文本(僵尸3第五页灰色候选), "粉色候选=", 买卖物品候选标签文本(僵尸3第五页粉色候选))
	设置僵尸3层输出("卖物品第2步失败：第五页既未识别到灰色，也未识别到粉色")
	return false, next
}

func 执行僵尸3猫猫双击逻辑(index int, logic 买卖物品逻辑) (bool, int) {
	next := index + 1
	输出("僵尸3买卖物品 双击猫猫", "逻辑=", logic.名称, "固定x=", 僵尸3猫猫固定X, "固定y=", 僵尸3猫猫固定Y)
	设置僵尸3层输出("卖物品第3步：固定双击猫猫 x=%d y=%d", 僵尸3猫猫固定X, 僵尸3猫猫固定Y)
	快速双击买卖物品坐标(僵尸3猫猫固定X, 僵尸3猫猫固定Y, 买卖物品猫猫双击间隔)
	if 买卖物品等待检查全部通过(logic.名称, logic.检查, 买卖物品猫猫检查等待) {
		输出("僵尸3买卖物品 逻辑成功", logic.名称)
		设置僵尸3层输出("卖物品第3步：猫猫打开商店成功")
		return true, next
	}

	输出("僵尸3买卖物品 逻辑失败", "逻辑=", logic.名称, "原因=双击猫猫后未出现全卖按钮")
	设置僵尸3层输出("卖物品第3步失败：双击猫猫后未出现全卖按钮")
	return false, next
}

func 执行僵尸3关闭背包逻辑(index int, logic 买卖物品逻辑) (bool, int) {
	next := index + 1
	输出("僵尸3买卖物品 执行", logic.名称, "规则=按I关闭背包")
	设置僵尸3层输出("卖物品第12步：按I关闭背包")
	if !僵尸3按I键关闭背包() {
		设置僵尸3层输出("卖物品第12步失败：已经打开背包特征仍存在")
		return false, next
	}
	设置僵尸3层输出("卖物品第12步：已经打开背包特征已消失")
	return true, next
}
