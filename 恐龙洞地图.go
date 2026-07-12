package main

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Dasongzi1366/AutoGo/motion"
)

const (
	恐龙洞一层左边    = 23
	恐龙洞一层右边    = 148
	恐龙洞一层识别左边  = 15
	恐龙洞一层识别右边  = 165
	恐龙洞二层左边    = 74
	恐龙洞二层右边    = 132
	恐龙洞三层左边    = 82
	恐龙洞三层右边    = 136
	恐龙洞四层左边    = 87
	恐龙洞四层右边    = 123
	恐龙洞四层右侧掉落X = 123
	恐龙洞三层右侧掉落X = 137

	恐龙洞一层上梯X = 88
	恐龙洞二层上梯X = 96
	恐龙洞三层上梯X = 101
	恐龙洞四层上梯X = 113

	恐龙洞梯子左侧允许   = 0
	恐龙洞方向X移动估计  = 14
	恐龙洞下楼避梯误差   = 5
	恐龙洞层识别Y容差   = 12
	恐龙洞移动到位容差   = 3
	恐龙洞移动间隔     = 180 * time.Millisecond
	恐龙洞未知位置等待   = 200 * time.Millisecond
	恐龙洞梯子对齐短按   = 60 * time.Millisecond
	恐龙洞梯子对齐等待   = 80 * time.Millisecond
	恐龙洞梯子稳定确认   = 50 * time.Millisecond
	恐龙洞短层爬梯持续   = 2000 * time.Millisecond
	恐龙洞长层爬梯持续   = 2600 * time.Millisecond
	恐龙洞爬梯随机增加   = 400 * time.Millisecond
	恐龙洞上梯启动检查   = 500 * time.Millisecond
	恐龙洞上梯最小位移   = 2
	恐龙洞上梯最大尝试   = 2
	恐龙洞上梯后强制检查  = 3 * time.Second
	恐龙洞一层上梯后等待  = 300 * time.Millisecond
	恐龙洞爬梯稳定成功时间 = 1 * time.Second
	恐龙洞上梯检测间隔   = 50 * time.Millisecond
	恐龙洞换层确认等待   = 1400 * time.Millisecond
	恐龙洞上下楼超时    = 18 * time.Second
	恐龙洞层内移动超时   = 30 * time.Second
	恐龙洞返回一层超时   = 60 * time.Second
	恐龙洞下楼单次确认   = 2 * time.Second
	恐龙洞下梯重试等待   = 400 * time.Millisecond
	恐龙洞刷怪卡住检测时间 = 3 * time.Second
	恐龙洞卡住坐标容差   = 1
	恐龙洞四层按F持续   = 5 * time.Second
	恐龙洞四层F按下时长  = 50 * time.Millisecond
	恐龙洞四层F抬起等待  = 50 * time.Millisecond

	恐龙洞高层巡逻最短间隔 = 118 * time.Second
	恐龙洞高层巡逻最长间隔 = 233 * time.Second
)

type 恐龙洞层配置 struct {
	层     int
	左边    int
	右边    int
	上梯子左X int
	上梯子右X int
}

var 恐龙洞层配置表 = []恐龙洞层配置{
	{层: 1, 左边: 恐龙洞一层左边, 右边: 恐龙洞一层右边, 上梯子左X: 恐龙洞一层上梯X, 上梯子右X: 恐龙洞一层上梯X},
	{层: 2, 左边: 恐龙洞二层左边, 右边: 恐龙洞二层右边, 上梯子左X: 恐龙洞二层上梯X, 上梯子右X: 恐龙洞二层上梯X},
	{层: 3, 左边: 恐龙洞三层左边, 右边: 恐龙洞三层右边, 上梯子左X: 恐龙洞三层上梯X, 上梯子右X: 恐龙洞三层上梯X},
	{层: 4, 左边: 恐龙洞四层左边, 右边: 恐龙洞四层右边, 上梯子左X: 恐龙洞四层上梯X, 上梯子右X: 恐龙洞四层上梯X},
	{层: 5},
}

type 恐龙洞层识别区域 struct {
	层  int
	X1 int
	X2 int
	Y  int
}

var 恐龙洞层识别区域表 = []恐龙洞层识别区域{
	{层: 1, X1: 15, X2: 58, Y: 169},
	{层: 1, X1: 59, X2: 恐龙洞一层识别右边, Y: 185},
	{层: 2, X1: 74, X2: 132, Y: 170},
	{层: 3, X1: 82, X2: 133, Y: 149},
	{层: 4, X1: 87, X2: 121, Y: 125},
	{层: 5, X1: 10, X2: 201, Y: 112},
}

var (
	恐龙洞输出锁      sync.Mutex
	恐龙洞输出行      = []string{"恐龙洞：等待查找黄点"}
	恐龙洞输出最大行数   = 300
	恐龙洞输出显示行数   = 4
	恐龙洞需要滚动到底   bool
	恐龙洞卖物品最短分钟  int32 = 8
	恐龙洞卖物品最长分钟  int32 = 13
	恐龙洞上梯子测试执行中 atomic.Bool
	恐龙洞已记录层数    atomic.Int32
	恐龙洞到达四层次数   atomic.Int32
	恐龙洞四层F已处理次数 atomic.Int32
	恐龙洞未匹配恢复锁   sync.Mutex
	恐龙洞未匹配上次X   int
	恐龙洞未匹配上次Y   int
	恐龙洞未匹配上次方向  int
	恐龙洞未匹配有记录   bool
)

type 恐龙洞刷怪卡住监控 struct {
	有上次坐标  bool
	上次X    int
	上次Y    int
	固定开始   time.Time
	下次脱困方向 int
}

func 读取恐龙洞输出文本() string {
	恐龙洞输出锁.Lock()
	defer 恐龙洞输出锁.Unlock()
	start := len(恐龙洞输出行) - 恐龙洞输出显示行数
	if start < 0 {
		start = 0
	}
	return strings.Join(恐龙洞输出行[start:], "\n")
}

func 读取恐龙洞全部输出文本() string {
	恐龙洞输出锁.Lock()
	defer 恐龙洞输出锁.Unlock()
	return strings.Join(恐龙洞输出行, "\n")
}

func 消耗恐龙洞输出滚动请求() bool {
	恐龙洞输出锁.Lock()
	defer 恐龙洞输出锁.Unlock()
	needScroll := 恐龙洞需要滚动到底
	恐龙洞需要滚动到底 = false
	return needScroll
}

func 设置恐龙洞输出(format string, args ...any) {
	text := fmt.Sprintf(format, args...)
	恐龙洞输出锁.Lock()
	恐龙洞输出行 = append(恐龙洞输出行, time.Now().Format("15:04:05  ")+text)
	if len(恐龙洞输出行) > 恐龙洞输出最大行数 {
		恐龙洞输出行 = 恐龙洞输出行[len(恐龙洞输出行)-恐龙洞输出最大行数:]
	}
	恐龙洞需要滚动到底 = true
	恐龙洞输出锁.Unlock()
	输出("恐龙洞", text)
}

func 恐龙洞当前层位置() (层位置, bool) {
	ok, x, y := 查找恐龙洞黄点坐标()
	if !ok {
		return 层位置{}, false
	}
	if cachedLayer := int(恐龙洞已记录层数.Load()); cachedLayer >= 1 && cachedLayer <= 5 {
		if cachedLayer != 1 {
			if firstLayer, matched := 识别恐龙洞层数(x, y); matched && firstLayer == 1 {
				恐龙洞记录当前层(1)
				cachedLayer = 1
			}
		}
		return 层位置{层: cachedLayer, X: x, Y: y}, true
	}
	return 恐龙洞根据黄点识别并记录层(x, y)
}

func 恐龙洞强制识别当前层位置() (层位置, bool) {
	ok, x, y := 查找恐龙洞黄点坐标()
	if !ok {
		return 层位置{}, false
	}
	return 恐龙洞根据黄点识别并记录层(x, y)
}

func 恐龙洞根据黄点识别并记录层(x, y int) (层位置, bool) {
	layer, _, _, _, matched := 识别恐龙洞黄点所在层(x, y)
	if !matched {
		return 层位置{}, false
	}
	恐龙洞记录当前层(layer)
	return 层位置{层: layer, X: x, Y: y}, true
}

func 恐龙洞记录当前层(layer int) {
	if layer >= 1 && layer <= 5 {
		恐龙洞已记录层数.Store(int32(layer))
	}
}

func 查找恐龙洞黄点坐标() (bool, int, int) {
	if 引擎 == nil {
		return false, -1, -1
	}
	result := 扫描恐龙洞小地图黄点(恐龙洞小地图黄点区域, false)
	if !result.Ok {
		return false, -1, -1
	}
	return true, result.X, result.Y
}

func 识别恐龙洞层数(x, y int) (int, bool) {
	if x >= 恐龙洞一层识别左边 && x <= 恐龙洞一层识别右边 && y > 170 {
		return 1, true
	}
	for _, region := range 恐龙洞层识别区域表 {
		if region.层 != 1 {
			continue
		}
		if x < region.X1 || x > region.X2 {
			continue
		}
		diff := absInt(y - region.Y)
		if diff <= 恐龙洞层识别Y容差 {
			return 1, true
		}
	}
	return 0, false
}

func 取恐龙洞层配置(layer int) (恐龙洞层配置, bool) {
	for _, config := range 恐龙洞层配置表 {
		if config.层 == layer {
			return config, true
		}
	}
	return 恐龙洞层配置{}, false
}

func 运行恐龙洞循环(runID int64) {
	恐龙洞已记录层数.Store(0)
	恐龙洞到达四层次数.Store(0)
	恐龙洞四层F已处理次数.Store(0)
	最短分钟, 最长分钟 := 读取恐龙洞卖物品分钟范围()
	设置恐龙洞输出("恐龙洞开始：循环刷1-4层，4层右侧掉3层、3层右侧掉1层，%d~%d分钟到5层卖物品", 最短分钟, 最长分钟)
	启动N键守护(runID)
	卖物品截止 := 新恐龙洞卖物品截止()

	for 脚本仍应运行(runID) {
		位置, ok := 恐龙洞当前层位置()
		if !ok {
			恐龙洞恢复未匹配位置("主循环")
			continue
		}

		if 恐龙洞已到卖物品时间(&卖物品截止) {
			恐龙洞前往五层卖物品后回一层(runID, &卖物品截止)
			continue
		}

		switch 位置.层 {
		case 1, 2, 3, 4:
			设置恐龙洞输出("开始本轮1-4层循环：当前位于%d层", 位置.层)
			恐龙洞执行刷怪周期(runID)
		case 5:
			设置恐龙洞输出("5层未到卖物品时间：直接逐层回1层")
			恐龙洞逐层回一层(runID)
		default:
			time.Sleep(恐龙洞未知位置等待)
		}
	}
	设置恐龙洞输出("恐龙洞流程已停止")
}

func 恐龙洞巡逻一层一轮(runID int64) bool {
	if !恐龙洞层内走到X(runID, 1, 恐龙洞一层左边) {
		return false
	}
	return 恐龙洞层内走到X(runID, 1, 恐龙洞一层右边)
}

func 恐龙洞执行刷怪周期(runID int64) bool {
	位置, ok := 恐龙洞当前层位置()
	if !ok || 位置.层 < 1 || 位置.层 > 4 {
		return false
	}

	for layer := 位置.层; layer <= 4 && 脚本仍应运行(runID); layer++ {
		位置, ok = 恐龙洞当前层位置()
		if !ok {
			return false
		}
		if 位置.层 > layer {
			continue
		}
		if 位置.层 < layer && !恐龙洞上到下一层直到成功(runID, 位置.层) {
			return false
		}
		if layer == 4 && !恐龙洞四层左侧连续按F(runID) {
			return false
		}
		rounds := 恐龙洞当前层巡逻次数(layer)
		if !恐龙洞巡逻指定层(runID, layer, rounds) {
			if layer == 2 {
				当前位置, found := 恐龙洞强制识别当前层位置()
				if found && 当前位置.层 == 1 {
					设置恐龙洞输出("1层上2层后又掉回1层：跳过1层巡逻，立即重新爬2层")
					if !恐龙洞上到下一层直到成功(runID, 1) {
						return false
					}
					layer--
					continue
				}
			}
			return false
		}
		if layer < 4 && !恐龙洞上到下一层直到成功(runID, layer) {
			return false
		}
	}
	if !恐龙洞从右侧掉落直到层(runID, 4, 恐龙洞四层右侧掉落X, 3) {
		return false
	}
	位置, ok = 恐龙洞当前层位置()
	if ok && 位置.层 == 1 {
		设置恐龙洞输出("4层已直接掉到1层，本轮返回完成")
		return true
	}
	if !ok || 位置.层 != 3 {
		return false
	}
	return 恐龙洞从右侧掉落直到层(runID, 3, 恐龙洞三层右侧掉落X, 1)
}

func 恐龙洞四层左侧连续按F(runID int64) bool {
	到达次数 := 恐龙洞到达四层次数.Load()
	if 到达次数 <= 0 {
		设置恐龙洞输出("当前4层不是由3层爬梯到达：不计数，不执行连续按F")
		return 脚本仍应运行(runID)
	}
	if 恐龙洞四层F已处理次数.Load() == 到达次数 {
		return 脚本仍应运行(runID)
	}
	if 到达次数%2 != 0 {
		恐龙洞四层F已处理次数.Store(到达次数)
		设置恐龙洞输出("第%d次到达4层：本次跳过左侧连续按F", 到达次数)
		return 脚本仍应运行(runID)
	}
	设置恐龙洞输出("第%d次到达4层：执行左侧连续按F", 到达次数)
	设置恐龙洞输出("到达4层：先移动到最左边x=%d", 恐龙洞四层左边)
	if !恐龙洞层内走到X(runID, 4, 恐龙洞四层左边) {
		return false
	}

	displayID := 当前显示ID()
	deadline := time.Now().Add(恐龙洞四层按F持续)
	设置恐龙洞输出("4层已到最左边，连续点按F共%d秒", int(恐龙洞四层按F持续.Seconds()))
	for 脚本仍应运行(runID) && time.Now().Before(deadline) {
		motion.KeyActionDown(motion.KEYCODE_F, displayID)
		time.Sleep(恐龙洞四层F按下时长)
		motion.KeyActionUp(motion.KEYCODE_F, displayID)
		time.Sleep(恐龙洞四层F抬起等待)
	}
	motion.KeyActionUp(motion.KEYCODE_F, displayID)
	if !脚本仍应运行(runID) {
		return false
	}
	恐龙洞四层F已处理次数.Store(到达次数)
	设置恐龙洞输出("4层连续点按F完成，继续原巡逻流程")
	return true
}

func 恐龙洞巡逻指定层(runID int64, layer, rounds int) bool {
	config, ok := 取恐龙洞层配置(layer)
	if !ok || config.左边 >= config.右边 {
		return false
	}
	for round := 1; round <= rounds && 脚本仍应运行(runID); round++ {
		设置恐龙洞输出("%d层巡逻%d/%d：%d到%d", layer, round, rounds, config.左边, config.右边)
		if !恐龙洞层内走到X(runID, layer, config.左边) {
			return false
		}
		if !恐龙洞层内走到X(runID, layer, config.右边) {
			return false
		}
	}
	return 脚本仍应运行(runID)
}

func 恐龙洞当前层巡逻次数(layer int) int {
	if layer == 2 || layer == 3 {
		return 1 + 移动随机.Intn(2)
	}
	return 1
}

func 恐龙洞恢复未匹配位置(scene string) {
	ok, x, y := 查找恐龙洞黄点坐标()
	if !ok {
		time.Sleep(恐龙洞未知位置等待)
		return
	}
	direction := motion.KEYCODE_DPAD_LEFT
	if x < (恐龙洞一层识别左边+恐龙洞一层识别右边)/2 {
		direction = motion.KEYCODE_DPAD_RIGHT
	}
	恐龙洞未匹配恢复锁.Lock()
	if 恐龙洞未匹配有记录 && absInt(x-恐龙洞未匹配上次X) <= 恐龙洞卡住坐标容差 && absInt(y-恐龙洞未匹配上次Y) <= 恐龙洞卡住坐标容差 {
		if 恐龙洞未匹配上次方向 == motion.KEYCODE_DPAD_LEFT {
			direction = motion.KEYCODE_DPAD_RIGHT
		} else {
			direction = motion.KEYCODE_DPAD_LEFT
		}
	}
	恐龙洞未匹配上次X = x
	恐龙洞未匹配上次Y = y
	恐龙洞未匹配上次方向 = direction
	恐龙洞未匹配有记录 = true
	恐龙洞未匹配恢复锁.Unlock()
	设置恐龙洞输出("%s：黄点=(%d,%d)但楼层未匹配，按%s+X返回可识别区", scene, x, y, 键名(direction))
	按组合键不空格(direction, motion.KEYCODE_X, 方向键按下毫秒)
	按空格群攻()
	time.Sleep(恐龙洞移动间隔)
}

func 恐龙洞层内走到X(runID int64, layer, targetX int) bool {
	config, ok := 取恐龙洞层配置(layer)
	if !ok {
		return false
	}
	monitor := &恐龙洞刷怪卡住监控{}
	deadline := time.Now().Add(恐龙洞层内移动超时)
	for 脚本仍应运行(runID) && time.Now().Before(deadline) {
		位置, found := 恐龙洞当前层位置()
		if !found {
			恐龙洞恢复未匹配位置(fmt.Sprintf("%d层刷怪", layer))
			continue
		}
		if 位置.层 != layer {
			设置恐龙洞输出("%d层刷怪过程中已掉到%d层，按当前层继续循环", layer, 位置.层)
			return false
		}
		if monitor.检查并执行打怪(位置) {
			continue
		}
		if 恐龙洞巡逻目标已到位(位置.X, targetX, config.左边, config.右边) {
			return true
		}

		direction := motion.KEYCODE_DPAD_RIGHT
		if 位置.X > targetX {
			direction = motion.KEYCODE_DPAD_LEFT
		}
		isPatrolEdge := targetX == config.左边 || targetX == config.右边
		if isPatrolEdge || 恐龙洞可以跳跃移动(位置.X, targetX, direction, config.左边, config.右边) {
			按组合键不空格(direction, motion.KEYCODE_X, 方向键按下毫秒)
		} else {
			按方向键(direction, 方向键按下毫秒)
		}
		按空格群攻()
		time.Sleep(恐龙洞移动间隔)
	}
	设置恐龙洞输出("%d层移动到x=%d超时，交回主循环重新识别", layer, targetX)
	return false
}

func 恐龙洞巡逻目标已到位(x, targetX, left, right int) bool {
	if targetX == left {
		return x <= left+恐龙洞移动到位容差
	}
	if targetX == right {
		return x >= right-恐龙洞移动到位容差
	}
	return absInt(x-targetX) <= 恐龙洞移动到位容差
}

func 恐龙洞从右侧掉落直到层(runID int64, startLayer, edgeX, targetLayer int) bool {
	deadline := time.Now().Add(恐龙洞上下楼超时)
	dropStarted := false
	for 脚本仍应运行(runID) && time.Now().Before(deadline) {
		位置 := 层位置{}
		found := false
		if dropStarted {
			位置, found = 恐龙洞强制识别当前层位置()
		} else {
			位置, found = 恐龙洞当前层位置()
		}
		if !found {
			dropStarted = true
			按组合键不空格(motion.KEYCODE_DPAD_RIGHT, motion.KEYCODE_X, 方向键按下毫秒)
			time.Sleep(恐龙洞未知位置等待)
			continue
		}
		if 位置.层 == targetLayer || startLayer == 4 && 位置.层 == 1 {
			恐龙洞记录当前层(位置.层)
			设置恐龙洞输出("%d层右侧掉落成功：已到%d层", startLayer, 位置.层)
			return true
		}
		if 位置.层 != startLayer {
			设置恐龙洞输出("%d层右侧掉落中：当前识别为%d层，继续按右+X防止卡点", startLayer, 位置.层)
			dropStarted = true
			按组合键不空格(motion.KEYCODE_DPAD_RIGHT, motion.KEYCODE_X, 方向键按下毫秒)
			time.Sleep(恐龙洞移动间隔)
			continue
		}
		设置恐龙洞输出("%d层右侧掉落：当前x=%d，目标边缘x=%d，直接按右+X", startLayer, 位置.X, edgeX)
		dropStarted = true
		按组合键不空格(motion.KEYCODE_DPAD_RIGHT, motion.KEYCODE_X, 方向键按下毫秒)
		time.Sleep(恐龙洞移动间隔)
	}
	设置恐龙洞输出("%d层右侧掉落到%d层超时", startLayer, targetLayer)
	return false
}

func 恐龙洞可以跳跃移动(x, targetX, direction, left, right int) bool {
	nextX := x
	if direction == motion.KEYCODE_DPAD_LEFT {
		nextX -= 跳跃移动像素
	} else {
		nextX += 跳跃移动像素
	}
	if nextX < left || nextX > right {
		return false
	}
	if direction == motion.KEYCODE_DPAD_LEFT {
		return nextX >= targetX
	}
	return nextX <= targetX
}

func (monitor *恐龙洞刷怪卡住监控) 检查并执行打怪(位置 层位置) bool {
	now := time.Now()
	if !monitor.有上次坐标 || absInt(位置.X-monitor.上次X) > 恐龙洞卡住坐标容差 || absInt(位置.Y-monitor.上次Y) > 恐龙洞卡住坐标容差 {
		monitor.有上次坐标 = true
		monitor.上次X = 位置.X
		monitor.上次Y = 位置.Y
		monitor.固定开始 = now
		return false
	}
	if monitor.固定开始.IsZero() || now.Sub(monitor.固定开始) < 恐龙洞刷怪卡住检测时间 {
		return false
	}
	direction := monitor.下一个脱困方向()
	设置恐龙洞输出("刷怪卡住检测：坐标%d秒未动，按%s+Z脱困 层=%d x=%d y=%d", int(恐龙洞刷怪卡住检测时间.Seconds()), 键名(direction), 位置.层, 位置.X, 位置.Y)
	按组合键不空格(direction, motion.KEYCODE_Z, 方向键按下毫秒)
	按空格群攻()
	time.Sleep(近怪攻击间隔)
	monitor.固定开始 = time.Now()
	return true
}

func (monitor *恐龙洞刷怪卡住监控) 下一个脱困方向() int {
	if monitor.下次脱困方向 == 0 {
		if 移动随机.Intn(2) == 0 {
			monitor.下次脱困方向 = motion.KEYCODE_DPAD_LEFT
		} else {
			monitor.下次脱困方向 = motion.KEYCODE_DPAD_RIGHT
		}
	}
	direction := monitor.下次脱困方向
	if direction == motion.KEYCODE_DPAD_LEFT {
		monitor.下次脱困方向 = motion.KEYCODE_DPAD_RIGHT
	} else {
		monitor.下次脱困方向 = motion.KEYCODE_DPAD_LEFT
	}
	return direction
}

func 恐龙洞上到下一层直到成功(runID int64, startLayer int) bool {
	ok := 恐龙洞上到下一层直到成功条件(startLayer, func() bool { return 脚本仍应运行(runID) })
	if ok || !脚本仍应运行(runID) {
		return ok
	}
	设置恐龙洞输出("%d层上%d层连续%d次未成功，返回1层", startLayer, startLayer+1, 恐龙洞上梯最大尝试)
	if 位置, found := 恐龙洞强制识别当前层位置(); found && 位置.层 == 1 {
		设置恐龙洞输出("爬梯连续失败后已经在1层")
		return false
	}
	if !恐龙洞逐层回一层(runID) {
		设置恐龙洞输出("爬梯连续失败后返回1层未完成，交回主循环继续处理")
	}
	return false
}

func 恐龙洞上到下一层直到成功条件(startLayer int, shouldContinue func() bool) bool {
	return 恐龙洞上到下一层直到成功条件从位置(startLayer, nil, shouldContinue)
}

func 恐龙洞上到下一层直到成功条件从位置(startLayer int, initial *层位置, shouldContinue func() bool) bool {
	if startLayer < 1 || startLayer >= 5 {
		return false
	}
	config, ok := 取恐龙洞层配置(startLayer)
	if !ok {
		return false
	}
	targetLayer := startLayer + 1
	for climbAttempt := 1; climbAttempt <= 恐龙洞上梯最大尝试 && shouldContinue(); climbAttempt++ {
		位置 := 层位置{}
		found := false
		if initial != nil {
			位置 = *initial
			found = true
			initial = nil
		} else {
			位置, found = 恐龙洞强制识别当前层位置()
		}
		if !found {
			设置恐龙洞输出("%d层上%d层：第%d/%d次尝试前未识别到当前层", startLayer, targetLayer, climbAttempt, 恐龙洞上梯最大尝试)
			continue
		}
		if 位置.层 == targetLayer {
			恐龙洞完成上梯(targetLayer)
			return true
		}
		if 位置.层 != startLayer {
			设置恐龙洞输出("%d层上%d层：第%d/%d次尝试前实际位于%d层", startLayer, targetLayer, climbAttempt, 恐龙洞上梯最大尝试, 位置.层)
			continue
		}
		位置, found = 恐龙洞对齐梯子(位置, startLayer, config.上梯子左X, config.上梯子右X, shouldContinue)
		if !found {
			设置恐龙洞输出("%d层上%d层：第%d/%d次梯子对齐失败", startLayer, targetLayer, climbAttempt, 恐龙洞上梯最大尝试)
			continue
		}
		设置恐龙洞输出("%d层上%d层：第%d/%d次，已在梯子正下方x=%d，立即按上+Z", startLayer, targetLayer, climbAttempt, 恐龙洞上梯最大尝试, 位置.X)
		attemptStarted := time.Now()
		动作判断成功 := 恐龙洞执行快速上梯动作(startLayer, targetLayer, shouldContinue)
		检查时间 := attemptStarted.Add(恐龙洞上梯后强制检查)
		for shouldContinue() && time.Now().Before(检查时间) {
			time.Sleep(恐龙洞上梯检测间隔)
		}
		if !shouldContinue() {
			return false
		}
		current, currentOK := 恐龙洞强制识别当前层位置()
		if currentOK && current.层 == targetLayer {
			if startLayer == 3 && targetLayer == 4 {
				到达次数 := 恐龙洞到达四层次数.Add(1)
				设置恐龙洞输出("3层上4层确认成功：本次计为第%d次到达4层", 到达次数)
			}
			设置恐龙洞输出("%d层上%d层：约%dms强制识别成功 x=%d y=%d", startLayer, targetLayer, 恐龙洞上梯后强制检查.Milliseconds(), current.X, current.Y)
			恐龙洞完成上梯(targetLayer)
			return true
		}
		if currentOK {
			设置恐龙洞输出("%d层上%d层：第%d/%d次失败，动作判断=%t，约%dms实际识别为%d层 x=%d y=%d", startLayer, targetLayer, climbAttempt, 恐龙洞上梯最大尝试, 动作判断成功, 恐龙洞上梯后强制检查.Milliseconds(), current.层, current.X, current.Y)
		} else {
			设置恐龙洞输出("%d层上%d层：第%d/%d次失败，动作判断=%t，约%dms仍未识别到实际楼层", startLayer, targetLayer, climbAttempt, 恐龙洞上梯最大尝试, 动作判断成功, 恐龙洞上梯后强制检查.Milliseconds())
		}
		if climbAttempt < 恐龙洞上梯最大尝试 {
			设置恐龙洞输出("%d层上%d层：重新对齐，开始最后1次尝试", startLayer, targetLayer)
		}
	}
	设置恐龙洞输出("%d层上%d层：连续%d次强制识别均未成功", startLayer, targetLayer, 恐龙洞上梯最大尝试)
	return false
}

func 恐龙洞完成上梯(targetLayer int) {
	点按空格()
	恐龙洞记录当前层(targetLayer)
	if targetLayer == 2 {
		time.Sleep(恐龙洞一层上梯后等待)
		设置恐龙洞输出("上梯成功：已到2层，先按空格1次并等待%dms", 恐龙洞一层上梯后等待.Milliseconds())
		return
	}
	设置恐龙洞输出("上梯成功：已到%d层，已先按空格1次", targetLayer)
}

func 恐龙洞对齐梯子(initial 层位置, layer, ladderLeft, ladderRight int, shouldContinue func() bool) (层位置, bool) {
	monitor := &恐龙洞刷怪卡住监控{}
	位置 := initial
	useInitial := true
	acceptedLeft := ladderLeft - 恐龙洞梯子左侧允许
	acceptedRight := ladderRight
	deadline := time.Now().Add(恐龙洞上下楼超时)
	for shouldContinue() && time.Now().Before(deadline) {
		if !useInitial {
			ok, x, y := 查找恐龙洞黄点坐标()
			if !ok {
				time.Sleep(恐龙洞未知位置等待)
				continue
			}
			位置 = 层位置{层: layer, X: x, Y: y}
		}
		useInitial = false
		if monitor.检查并执行打怪(位置) {
			continue
		}
		diff := 恐龙洞到X范围距离(位置.X, acceptedLeft, acceptedRight)
		if diff == 0 {
			time.Sleep(恐龙洞梯子稳定确认)
			ok, x, y := 查找恐龙洞黄点坐标()
			if ok && x == 位置.X && y == 位置.Y {
				设置恐龙洞输出("%d层梯子已对齐：梯子x=%d-%d 允许区=%d-%d 坐标=(%d,%d) 50ms稳定", layer, ladderLeft, ladderRight, acceptedLeft, acceptedRight, 位置.X, 位置.Y)
				return 位置, true
			}
			if ok {
				位置 = 层位置{层: layer, X: x, Y: y}
				useInitial = true
			}
			continue
		}
		direction := motion.KEYCODE_DPAD_LEFT
		if diff < 0 {
			direction = motion.KEYCODE_DPAD_RIGHT
		}
		if absInt(diff) <= 3 {
			恐龙洞按移动并检测反向击退快速(位置, direction, 恐龙洞梯子对齐等待, func() {
				按方向键(direction, int(恐龙洞梯子对齐短按/time.Millisecond))
			})
			continue
		}
		恐龙洞按移动并检测反向击退快速(位置, direction, 恐龙洞移动间隔, func() {
			if absInt(diff) > 恐龙洞方向X移动估计+1 {
				按组合键不空格(direction, motion.KEYCODE_X, 方向键按下毫秒)
				return
			}
			按方向键(direction, 方向键按下毫秒)
		})
	}
	设置恐龙洞输出("%d层梯子对齐超时，交回上梯流程重试", layer)
	return 层位置{}, false
}

func 恐龙洞到X范围距离(x, left, right int) int {
	if x < left {
		return x - left
	}
	if x > right {
		return x - right
	}
	return 0
}

func 恐龙洞范围内最近X(x, left, right int) int {
	if x < left {
		return left
	}
	if x > right {
		return right
	}
	return x
}

func 恐龙洞执行快速上梯动作(startLayer, targetLayer int, shouldContinue func() bool) bool {
	before, beforeOK := 读取恐龙洞层数相对位置()
	if !beforeOK {
		设置恐龙洞输出("上梯启动失败：按键前未读取到相对高度，重新对齐")
		return false
	}
	displayID := 当前显示ID()
	start := time.Now()
	minimumDuration := 恐龙洞最低爬梯持续(startLayer)
	totalDuration := minimumDuration + time.Duration(移动随机.Int63n(int64(恐龙洞爬梯随机增加)+1))
	motion.KeyActionDown(motion.KEYCODE_DPAD_UP, displayID)
	motion.KeyActionDown(motion.KEYCODE_Z, displayID)
	defer func() {
		motion.KeyActionUp(motion.KEYCODE_Z, displayID)
		motion.KeyActionUp(motion.KEYCODE_DPAD_UP, displayID)
	}()

	time.Sleep(time.Duration(方向键按下毫秒) * time.Millisecond)
	motion.KeyActionUp(motion.KEYCODE_Z, displayID)
	checkAt := start.Add(恐龙洞上梯启动检查)
	for shouldContinue() && time.Now().Before(checkAt) {
		time.Sleep(恐龙洞上梯检测间隔)
	}
	if !shouldContinue() {
		return false
	}
	after, afterOK := 读取恐龙洞层数相对位置()
	delta := after.差值 - before.差值
	if !afterOK || absInt(delta) < 恐龙洞上梯最小位移 {
		设置恐龙洞输出("上梯启动失败：500ms高度未向上移动 %d->%d，立即重新爬", before.差值, after.差值)
		return false
	}
	设置恐龙洞输出("上梯启动成功：500ms高度变化 %d->%d，继续按上至少%dms", before.差值, after.差值, minimumDuration.Milliseconds())

	deadline := start.Add(totalDuration)
	successDeadline := start.Add(minimumDuration)
	lastRelativeValue := after.差值
	lastValueChangedAt := time.Now()
	hasValueChanged := true
	successDetected := false
	for shouldContinue() && time.Now().Before(deadline) {
		if successDetected {
			if !time.Now().Before(successDeadline) {
				return true
			}
			time.Sleep(恐龙洞上梯检测间隔)
			continue
		}
		relative, relativeOK := 读取恐龙洞层数相对位置()
		if relativeOK {
			if layer, matched := 识别恐龙洞相对层数(relative.差值); matched && layer == targetLayer {
				successDetected = true
			}
			now := time.Now()
			if relative.差值 != lastRelativeValue {
				设置恐龙洞输出("爬梯高度变化：%d->%d", lastRelativeValue, relative.差值)
				lastRelativeValue = relative.差值
				lastValueChangedAt = now
				hasValueChanged = true
			} else if hasValueChanged && now.Sub(lastValueChangedAt) >= 恐龙洞爬梯稳定成功时间 {
				设置恐龙洞输出("爬梯高度值%d连续%dms未变化，判定换层成功", lastRelativeValue, 恐龙洞爬梯稳定成功时间.Milliseconds())
				successDetected = true
			}
		}
		if successDetected && !time.Now().Before(successDeadline) {
			return true
		}
		time.Sleep(恐龙洞上梯检测间隔)
	}
	位置, ok := 恐龙洞当前层位置()
	return ok && 位置.层 == targetLayer
}

func 恐龙洞按移动并检测反向击退快速(original 层位置, direction int, wait time.Duration, action func()) {
	if action != nil {
		action()
	}
	if wait > 0 {
		time.Sleep(wait)
	}
	ok, x, _ := 查找恐龙洞黄点坐标()
	if !ok {
		return
	}
	deltaX := x - original.X
	reversed := direction == motion.KEYCODE_DPAD_RIGHT && deltaX < -恐龙洞卡住坐标容差
	reversed = reversed || direction == motion.KEYCODE_DPAD_LEFT && deltaX > 恐龙洞卡住坐标容差
	if !reversed {
		return
	}
	设置恐龙洞输出("梯子对齐疑似被怪反向击退：按%s x=%d->%d，攻击1次", 键名(direction), original.X, x)
	点按空格()
}

func 恐龙洞按移动并检测反向击退(original 层位置, direction int, wait time.Duration, action func()) {
	if action != nil {
		action()
	}
	if wait > 0 {
		time.Sleep(wait)
	}
	current, ok := 恐龙洞当前层位置()
	if !ok || current.层 != original.层 {
		return
	}
	deltaX := current.X - original.X
	reversed := direction == motion.KEYCODE_DPAD_RIGHT && deltaX < -恐龙洞卡住坐标容差
	reversed = reversed || direction == motion.KEYCODE_DPAD_LEFT && deltaX > 恐龙洞卡住坐标容差
	if !reversed {
		return
	}
	设置恐龙洞输出("梯子对齐疑似被怪反向击退：按%s x=%d->%d，攻击1次", 键名(direction), original.X, current.X)
	点按空格()
}

func 执行恐龙洞上梯子测试() {
	if 程序退出中.Load() {
		设置恐龙洞输出("上梯子测试失败：程序正在退出")
		return
	}
	if 脚本运行中.Load() {
		设置恐龙洞输出("上梯子测试失败：脚本运行中，请先点结束")
		return
	}
	if !恐龙洞上梯子测试执行中.CompareAndSwap(false, true) {
		设置恐龙洞输出("上梯子测试执行中：上一次点击尚未结束，不重复启动")
		return
	}
	go func() {
		defer 恐龙洞上梯子测试执行中.Store(false)
		位置, ok := 恐龙洞强制识别当前层位置()
		if !ok {
			设置恐龙洞输出("上梯子测试失败：未识别到当前楼层")
			return
		}
		if 位置.层 >= 5 {
			设置恐龙洞输出("上梯子测试失败：当前已在5层")
			return
		}
		targetLayer := 位置.层 + 1
		config, found := 取恐龙洞层配置(位置.层)
		if !found {
			设置恐龙洞输出("上梯子测试失败：未找到%d层梯子配置", 位置.层)
			return
		}
		markerX := 恐龙洞范围内最近X(位置.X, config.上梯子左X, config.上梯子右X)
		标记恐龙洞梯子目标点(markerX, 位置.Y)
		设置恐龙洞输出("上梯子测试开始：当前%d层 x=%d y=%d，梯子=%d-%d，目标%d层", 位置.层, 位置.X, 位置.Y, config.上梯子左X, config.上梯子右X, targetLayer)
		shouldContinue := func() bool { return !程序退出中.Load() && !脚本运行中.Load() }
		if !恐龙洞上到下一层直到成功条件从位置(位置.层, &位置, shouldContinue) {
			设置恐龙洞输出("上梯子测试失败：未到%d层", targetLayer)
			return
		}
		newPosition, found := 恐龙洞当前层位置()
		if found && newPosition.层 == targetLayer {
			设置恐龙洞输出("上梯子测试成功：已到%d层 x=%d y=%d", newPosition.层, newPosition.X, newPosition.Y)
			return
		}
		设置恐龙洞输出("上梯子测试失败：动作结束后未确认到%d层", targetLayer)
	}()
}

func 标记恐龙洞梯子目标点(x, y int) {
	if 引擎 == nil {
		return
	}
	screenX := 引擎.scaleX(x + 引擎.offsetX)
	screenY := 引擎.scaleY(y + 引擎.offsetY)
	addDebugTargetPoint(screenX, screenY)
}

func 标记恐龙洞找到的黄点(x, y int) {
	if 引擎 == nil {
		return
	}
	screenX := 引擎.scaleX(x + 引擎.offsetX)
	screenY := 引擎.scaleY(y + 引擎.offsetY)
	addDebugSinglePixel(screenX, screenY)
}

func 恐龙洞逐层回一层(runID int64) bool {
	deadline := time.Now().Add(恐龙洞返回一层超时)
	for 脚本仍应运行(runID) && time.Now().Before(deadline) {
		位置, ok := 恐龙洞当前层位置()
		if !ok {
			恐龙洞恢复未匹配位置("返回1层")
			continue
		}
		if 位置.层 == 1 {
			设置恐龙洞输出("已逐层回到1层")
			return true
		}
		if 位置.层 < 2 || 位置.层 > 5 {
			return false
		}
		var success bool
		switch 位置.层 {
		case 4:
			success = 恐龙洞从右侧掉落直到层(runID, 4, 恐龙洞四层右侧掉落X, 3)
		case 3:
			success = 恐龙洞从右侧掉落直到层(runID, 3, 恐龙洞三层右侧掉落X, 1)
		default:
			success = 恐龙洞下到上一层直到成功(runID, 位置.层)
		}
		if !success {
			time.Sleep(恐龙洞下梯重试等待)
		}
	}
	设置恐龙洞输出("返回1层总流程超时，交回主循环重新识别")
	return false
}

func 恐龙洞下到上一层直到成功(runID int64, startLayer int) bool {
	lowerConfig, ok := 取恐龙洞层配置(startLayer - 1)
	if !ok {
		return false
	}
	targetLayer := startLayer - 1
	deadline := time.Now().Add(恐龙洞上下楼超时)
	for 脚本仍应运行(runID) && time.Now().Before(deadline) {
		位置, found := 恐龙洞当前层位置()
		if !found {
			恐龙洞恢复未匹配位置("下楼流程")
			continue
		}
		if 位置.层 == targetLayer {
			设置恐龙洞输出("下梯成功：已到%d层", targetLayer)
			return true
		}
		if 位置.层 != startLayer {
			return false
		}
		ladderDistance := 恐龙洞到X范围距离(位置.X, lowerConfig.上梯子左X, lowerConfig.上梯子右X)
		if absInt(ladderDistance) <= 恐龙洞下楼避梯误差 {
			direction := motion.KEYCODE_DPAD_RIGHT
			if config, found := 取恐龙洞层配置(startLayer); found && 位置.X > (config.左边+config.右边)/2 {
				direction = motion.KEYCODE_DPAD_LEFT
			}
			设置恐龙洞输出("%d层下楼：当前在梯子区域x=%d，先按%s+Z离开梯子", startLayer, 位置.X, 键名(direction))
			按组合键不空格(direction, motion.KEYCODE_Z, 方向键按下毫秒)
			time.Sleep(恐龙洞移动间隔)
			continue
		}
		设置恐龙洞输出("%d层下%d层：当前x=%d不在梯子x=%d-%d附近，按下+Z", startLayer, targetLayer, 位置.X, lowerConfig.上梯子左X, lowerConfig.上梯子右X)
		按组合键同时不空格(motion.KEYCODE_DPAD_DOWN, motion.KEYCODE_Z, 方向键按下毫秒)
		if 恐龙洞换层后确认(targetLayer, func() bool { return 脚本仍应运行(runID) }) {
			设置恐龙洞输出("下梯成功：已到%d层", targetLayer)
			return true
		}
		direction := motion.KEYCODE_DPAD_RIGHT
		if 移动随机.Intn(2) == 0 {
			direction = motion.KEYCODE_DPAD_LEFT
		}
		设置恐龙洞输出("%d层下楼未换层，按%s+Z脱困后重试", startLayer, 键名(direction))
		按组合键不空格(direction, motion.KEYCODE_Z, 方向键按下毫秒)
		time.Sleep(恐龙洞下梯重试等待)
	}
	设置恐龙洞输出("%d层下%d层超时，稍后重试", startLayer, targetLayer)
	return false
}

func 恐龙洞等待到层(targetLayer int, timeout time.Duration, shouldContinue func() bool) bool {
	deadline := time.Now().Add(timeout)
	for shouldContinue() && time.Now().Before(deadline) {
		位置, ok := 恐龙洞强制识别当前层位置()
		if ok && 位置.层 == targetLayer {
			return true
		}
		time.Sleep(100 * time.Millisecond)
	}
	return false
}

func 恐龙洞换层后确认(targetLayer int, shouldContinue func() bool) bool {
	deadline := time.Now().Add(3 * time.Second)
	virtualLogged := false
	direction := motion.KEYCODE_DPAD_LEFT
	for shouldContinue() && time.Now().Before(deadline) {
		position, ok := 恐龙洞强制识别当前层位置()
		if ok {
			if position.层 == targetLayer {
				设置恐龙洞输出("换层确认成功：识别层=%d，与虚拟目标层对应", targetLayer)
				return true
			}
			time.Sleep(100 * time.Millisecond)
			continue
		}

		yellowOK, x, _ := 查找恐龙洞黄点坐标()
		if !virtualLogged {
			设置恐龙洞输出("换层后暂未识别楼层：虚拟当前层=%d，3秒内左右移动并打怪复检", targetLayer)
			virtualLogged = true
		}
		if config, found := 取恐龙洞层配置(targetLayer); found && yellowOK {
			if x <= config.左边+恐龙洞移动到位容差 {
				direction = motion.KEYCODE_DPAD_RIGHT
			} else if x >= config.右边-恐龙洞移动到位容差 {
				direction = motion.KEYCODE_DPAD_LEFT
			}
		}
		按组合键不空格(direction, motion.KEYCODE_X, 方向键按下毫秒)
		按空格群攻()
		if direction == motion.KEYCODE_DPAD_LEFT {
			direction = motion.KEYCODE_DPAD_RIGHT
		} else {
			direction = motion.KEYCODE_DPAD_LEFT
		}
		time.Sleep(恐龙洞移动间隔)
	}
	设置恐龙洞输出("换层确认失败：3秒内未识别到目标%d层", targetLayer)
	return false
}

func 恐龙洞持续按键(code int, duration time.Duration) {
	displayID := 当前显示ID()
	motion.KeyActionDown(code, displayID)
	time.Sleep(duration)
	motion.KeyActionUp(code, displayID)
}

func 恐龙洞最低爬梯持续(startLayer int) time.Duration {
	if startLayer == 2 || startLayer == 3 {
		return 恐龙洞长层爬梯持续
	}
	return 恐龙洞短层爬梯持续
}

func 新恐龙洞高层巡逻截止() time.Time {
	duration := 随机恐龙洞高层巡逻间隔()
	设置恐龙洞输出("下次高层巡逻：%d秒后", int(duration.Seconds()))
	return time.Now().Add(duration)
}

func 随机恐龙洞高层巡逻间隔() time.Duration {
	span := int64((恐龙洞高层巡逻最长间隔 - 恐龙洞高层巡逻最短间隔) / time.Second)
	if span <= 0 {
		return 恐龙洞高层巡逻最短间隔
	}
	return 恐龙洞高层巡逻最短间隔 + time.Duration(移动随机.Int63n(span+1))*time.Second
}

func 新恐龙洞卖物品截止() time.Time {
	duration := 随机恐龙洞卖物品周期时长()
	min, max := 读取恐龙洞卖物品分钟范围()
	设置恐龙洞输出("下次卖物品：%d分钟后（范围%d~%d分钟）", int(duration.Minutes()), min, max)
	return time.Now().Add(duration)
}

func 读取恐龙洞卖物品分钟范围() (int32, int32) {
	min := 恐龙洞卖物品最短分钟
	max := 恐龙洞卖物品最长分钟
	if min < 1 {
		min = 1
	}
	if max < min {
		max = min
	}
	return min, max
}

func 设置恐龙洞卖物品分钟范围(min, max int32) {
	if min < 1 {
		min = 1
	}
	if max < 1 {
		max = 1
	}
	if min > max {
		max = min
	}
	恐龙洞卖物品最短分钟 = min
	恐龙洞卖物品最长分钟 = max
}

func 随机恐龙洞卖物品周期时长() time.Duration {
	min, max := 读取恐龙洞卖物品分钟范围()
	span := int64(max - min)
	if span <= 0 {
		return time.Duration(min) * time.Minute
	}
	return time.Duration(min)*time.Minute + time.Duration(移动随机.Int63n(span+1))*time.Minute
}

func 恐龙洞已到卖物品时间(deadline *time.Time) bool {
	return deadline != nil && !time.Now().Before(*deadline)
}

func 恐龙洞前往五层卖物品后回一层(runID int64, sellDeadline *time.Time) bool {
	设置恐龙洞输出("到卖物品时间：依次爬梯前往5层")
	if !恐龙洞确保到五层(runID) {
		设置恐龙洞输出("前往5层卖物品失败，稍后重试")
		return false
	}
	位置, ok := 恐龙洞当前层位置()
	if !ok || 位置.层 != 5 {
		设置恐龙洞输出("卖物品取消：未确认在5层")
		return false
	}
	设置恐龙洞输出("5层开始卖物品")
	if !执行恐龙洞完整买卖物品流程(func() bool { return 脚本仍应运行(runID) }) {
		设置恐龙洞输出("5层卖物品中断，稍后重试")
		return false
	}
	if sellDeadline != nil {
		*sellDeadline = 新恐龙洞卖物品截止()
	}
	设置恐龙洞输出("5层卖物品完成，逐层回1层")
	return 恐龙洞逐层回一层(runID)
}

func 恐龙洞确保到五层(runID int64) bool {
	for 脚本仍应运行(runID) {
		位置, ok := 恐龙洞当前层位置()
		if !ok {
			恐龙洞恢复未匹配位置("前往5层")
			continue
		}
		if 位置.层 == 5 {
			return true
		}
		if 位置.层 < 1 || 位置.层 > 4 {
			return false
		}
		if !恐龙洞上到下一层直到成功(runID, 位置.层) {
			return false
		}
	}
	return false
}
