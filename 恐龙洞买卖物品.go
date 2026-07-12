package main

import (
	"sync"
	"sync/atomic"
	"time"
)

var (
	恐龙洞卖物品测试锁    sync.Mutex
	恐龙洞卖物品测试已启动  bool
	恐龙洞卖物品测试当前逻辑 int
	恐龙洞卖物品下一步执行中 atomic.Bool
)

func 执行恐龙洞完整买卖物品流程(shouldContinue func() bool) bool {
	if shouldContinue == nil {
		shouldContinue = func() bool { return true }
	}
	恐龙洞买卖物品流程中.Store(true)
	defer 恐龙洞买卖物品流程中.Store(false)
	开始买卖物品流程计时()
	defer 清除买卖物品流程计时()
	for index := 0; index < len(买卖物品逻辑表); {
		if 买卖物品流程已超时() {
			执行恐龙洞买卖物品超时收尾()
			return true
		}
		if !shouldContinue() {
			设置恐龙洞输出("买卖物品流程中断")
			return false
		}
		logic := 买卖物品逻辑表[index]
		ok, next := 执行恐龙洞买卖物品逻辑并返回下一步(index, logic)
		if !ok {
			设置恐龙洞输出("买卖物品失败：%s，继续重试直到超时", logic.名称)
			time.Sleep(500 * time.Millisecond)
			continue
		}
		if next <= index {
			设置恐龙洞输出("买卖物品失败：%s 下一步异常", logic.名称)
			time.Sleep(500 * time.Millisecond)
			continue
		}
		index = next
	}
	设置恐龙洞输出("买卖物品流程结束，耗时%d秒", int(买卖物品流程已用时()/time.Second))
	return true
}

func 执行恐龙洞买卖物品超时收尾() {
	设置恐龙洞输出("卖物品超时：关闭商店并按I关闭背包")
	尝试点击买卖物品特征("恐龙洞超时关闭商店", 关闭商店)
	恐龙洞按I关闭背包并确认("卖物品超时关闭背包", nil)
}

func 执行恐龙洞卖物品测试下一步() {
	if 程序退出中.Load() {
		设置恐龙洞输出("卖物品下一步失败：程序正在退出")
		return
	}
	if 脚本运行中.Load() {
		设置恐龙洞输出("卖物品下一步失败：脚本运行中，请先点结束")
		return
	}
	if !恐龙洞卖物品下一步执行中.CompareAndSwap(false, true) {
		设置恐龙洞输出("卖物品下一步执行中，请稍等")
		return
	}
	go func() {
		defer 恐龙洞卖物品下一步执行中.Store(false)
		if !恐龙洞卖物品测试流程运行中() {
			启动恐龙洞卖物品测试流程()
		}
		if 买卖物品流程已超时() {
			执行恐龙洞卖物品测试超时收尾()
			return
		}
		index, logic, ok := 取恐龙洞卖物品测试当前逻辑()
		if !ok {
			return
		}
		设置恐龙洞输出("卖物品下一步：第%d步 %s", index+1, logic.名称)
		if success, next := 执行恐龙洞买卖物品逻辑并返回下一步(index, logic); success {
			设置恐龙洞卖物品测试下一逻辑(next)
		} else {
			设置恐龙洞输出("卖物品下一步失败：第%d步 %s", index+1, logic.名称)
		}
		if 买卖物品流程已超时() {
			执行恐龙洞卖物品测试超时收尾()
		}
	}()
}

func 执行恐龙洞卖物品测试超时收尾() {
	执行恐龙洞买卖物品超时收尾()
	恐龙洞卖物品测试锁.Lock()
	恐龙洞卖物品测试已启动 = false
	恐龙洞卖物品测试当前逻辑 = 0
	恐龙洞卖物品测试锁.Unlock()
	清除买卖物品流程计时()
	设置恐龙洞输出("卖东西单步已因超时结束")
}

func 启动恐龙洞卖物品测试流程() {
	恐龙洞卖物品测试锁.Lock()
	恐龙洞卖物品测试已启动 = true
	恐龙洞卖物品测试当前逻辑 = 0
	恐龙洞卖物品测试锁.Unlock()
	开始买卖物品流程计时()
	设置恐龙洞输出("恐龙洞卖物品测试开始")
}

func 恐龙洞卖物品测试流程运行中() bool {
	恐龙洞卖物品测试锁.Lock()
	defer 恐龙洞卖物品测试锁.Unlock()
	return 恐龙洞卖物品测试已启动
}

func 取恐龙洞卖物品测试当前逻辑() (int, 买卖物品逻辑, bool) {
	恐龙洞卖物品测试锁.Lock()
	defer 恐龙洞卖物品测试锁.Unlock()
	if !恐龙洞卖物品测试已启动 {
		return 0, 买卖物品逻辑{}, false
	}
	if 恐龙洞卖物品测试当前逻辑 >= len(买卖物品逻辑表) {
		恐龙洞卖物品测试已启动 = false
		设置恐龙洞输出("恐龙洞卖物品测试结束")
		return 0, 买卖物品逻辑{}, false
	}
	return 恐龙洞卖物品测试当前逻辑, 买卖物品逻辑表[恐龙洞卖物品测试当前逻辑], true
}

func 设置恐龙洞卖物品测试下一逻辑(next int) {
	恐龙洞卖物品测试锁.Lock()
	恐龙洞卖物品测试当前逻辑 = next
	done := next >= len(买卖物品逻辑表)
	if done {
		恐龙洞卖物品测试已启动 = false
	}
	恐龙洞卖物品测试锁.Unlock()
	if done {
		清除买卖物品流程计时()
		设置恐龙洞输出("恐龙洞卖物品测试结束")
	}
}

func 执行恐龙洞买卖物品逻辑并返回下一步(index int, logic 买卖物品逻辑) (bool, int) {
	switch logic.名称 {
	case "第1个逻辑":
		return 执行恐龙洞打开背包逻辑(index)
	case "第2个逻辑":
		return 执行恐龙洞第五页切换逻辑(index)
	case "第3个逻辑":
		return 执行恐龙洞猫猫双击逻辑(index, logic)
	case "第12个逻辑":
		return 执行恐龙洞关闭背包逻辑(index)
	default:
		return 执行买卖物品逻辑并返回下一步(index, logic)
	}
}

func 执行恐龙洞打开背包逻辑(index int) (bool, int) {
	next := index + 1
	设置恐龙洞输出("卖物品第1步：查找并点击点开背包按钮")
	if found, x, y, candidate := 点击任一买卖物品特征("恐龙洞第1个逻辑", 僵尸3点开背包候选...); found {
		设置恐龙洞输出("卖物品第1步：点击%s x=%d y=%d", candidate.标签, x, y)
		return true, next
	}
	设置恐龙洞输出("卖物品第1步失败：未找到点开背包按钮")
	return false, next
}

func 执行恐龙洞第五页切换逻辑(index int) (bool, int) {
	next := index + 1
	设置恐龙洞输出("卖物品第2步：检查第五页，灰色点击、粉色通过")
	if found, x, y, candidate := 点击任一买卖物品特征("恐龙洞第2个逻辑", 僵尸3第五页灰色候选...); found {
		设置恐龙洞输出("卖物品第2步：点击灰色%s x=%d y=%d", candidate.标签, x, y)
		if ok, px, py, pink := 等待任一买卖物品特征(买卖物品猫猫检查等待, 僵尸3第五页粉色候选...); ok {
			设置恐龙洞输出("卖物品第2步：第五页已变粉色 %s x=%d y=%d", pink.标签, px, py)
			return true, next
		}
		设置恐龙洞输出("卖物品第2步失败：点击灰色后未识别粉色")
		return false, next
	}
	if found, x, y, pink := 查找任一买卖物品特征(僵尸3第五页粉色候选...); found {
		设置恐龙洞输出("卖物品第2步：第五页已经是粉色 %s x=%d y=%d", pink.标签, x, y)
		return true, next
	}
	设置恐龙洞输出("卖物品第2步失败：第五页既未识别灰色，也未识别粉色")
	return false, next
}

func 执行恐龙洞猫猫双击逻辑(index int, logic 买卖物品逻辑) (bool, int) {
	next := index + 1
	设置恐龙洞输出("卖物品第3步：固定双击猫猫 x=%d y=%d", 僵尸3猫猫固定X, 僵尸3猫猫固定Y)
	快速双击买卖物品坐标(僵尸3猫猫固定X, 僵尸3猫猫固定Y, 买卖物品猫猫双击间隔)
	if 买卖物品等待检查全部通过(logic.名称, logic.检查, 买卖物品猫猫检查等待) {
		设置恐龙洞输出("卖物品第3步：猫猫打开商店成功")
		return true, next
	}
	设置恐龙洞输出("卖物品第3步失败：双击猫猫后未出现全卖按钮")
	return false, next
}

func 执行恐龙洞关闭背包逻辑(index int) (bool, int) {
	next := index + 1
	设置恐龙洞输出("卖物品第12步：按I关闭背包")
	if !恐龙洞按I关闭背包并确认("卖物品关闭背包", nil) {
		设置恐龙洞输出("卖物品第12步失败：背包特征仍存在")
		return false, next
	}
	设置恐龙洞输出("卖物品第12步：背包特征已消失")
	return true, next
}
