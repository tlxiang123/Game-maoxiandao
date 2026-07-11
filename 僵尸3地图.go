package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Dasongzi1366/AutoGo/motion"
)

const (
	僵尸3绳子X      = 214
	僵尸3一层左边     = 18
	僵尸3一层右边     = 231
	僵尸3二层左边     = 164
	僵尸3二层右边     = 228
	僵尸3二层可走X左   = 170
	僵尸3二层可走X右   = 221
	僵尸3去三层触发X   = 253
	僵尸3去三层准备左   = 248
	僵尸3去三层准备右   = 252
	僵尸3左侧进三层目标X = 15
	僵尸3三层下绳X    = 120
	僵尸3三层下绳左界   = 114
	僵尸3三层下绳容差   = 3
	僵尸3绳子X容差    = 3
	僵尸3移动到位容差   = 3
	僵尸3移动间隔     = 220 * time.Millisecond
	僵尸3绳子目标左    = 214
	僵尸3绳子目标右    = 215
	僵尸3爬绳最短持续   = 4000 * time.Millisecond
	僵尸3爬绳随机增加   = 500 * time.Millisecond
	僵尸3爬绳后检查时间  = 1200 * time.Millisecond
	僵尸3爬绳启动检测时间 = 1 * time.Second
	僵尸3爬绳检测间隔   = 100 * time.Millisecond
	僵尸3爬绳上升Y阈值  = 3
	僵尸3二层巡逻次数   = 3
	僵尸3下绳无换层检测  = 2 * time.Second
	僵尸3下绳重试间隔   = 500 * time.Millisecond
	僵尸3未知位置等待   = 300 * time.Millisecond
	僵尸3靠近绳子慢走距离 = 3
	僵尸3绳子对齐短按   = 35 * time.Millisecond
	僵尸3绳子对齐等待   = 80 * time.Millisecond
	僵尸3绳子校正空格次数 = 3
	僵尸3绳子校正最短   = 30 * time.Millisecond
	僵尸3绳子校正最长   = 140 * time.Millisecond
	僵尸3绳子到位稳定等待 = 500 * time.Millisecond
	僵尸3击退反向容差   = 1
	僵尸3三层换层等待   = 1500 * time.Millisecond
	僵尸3三层下绳检测   = 1 * time.Second
	僵尸3三层左下下绳超时 = 12 * time.Second
	僵尸3左侧进三层超时  = 18 * time.Second
	僵尸3卡住检测时间   = 5 * time.Second
	僵尸3卡住报警时间   = time.Minute
	僵尸3卡住坐标容差   = 1
	僵尸3一层爬绳最短间隔 = 3 * time.Minute
	僵尸3一层爬绳最长间隔 = 4 * time.Minute
)

var 僵尸3层Y配置表 = []层Y配置{
	{层: 1, Y: 171},
	{层: 2, Y: 137},
	{层: 3, Y: 114},
}

var (
	僵尸3输出锁     sync.Mutex
	僵尸3输出行     = []string{"僵尸3：等待查找黄点"}
	僵尸3输出最大行数  = 60
	僵尸3需要滚动到底  bool
	僵尸3卖杂物     = true
	僵尸3绳子学习锁   sync.Mutex
	僵尸3绳子毫秒每像素 = 35.0
	僵尸3绳子校正倍率  = 1.0
	僵尸3绳子偏好目标X = 僵尸3绳子目标右
	僵尸3爬绳结果窗口  []int
	僵尸3目标214积分 float64
	僵尸3目标215积分 float64
)

type 僵尸3卡住监控 struct {
	有上次坐标 bool
	上次X   int
	上次Y   int
	固定开始  time.Time
	已脱离   bool
	已报警   bool
	下次方向  int
}

func 读取僵尸3输出文本() string {
	僵尸3输出锁.Lock()
	defer 僵尸3输出锁.Unlock()
	return strings.Join(僵尸3输出行, "\n")
}

func 消耗僵尸3输出滚动请求() bool {
	僵尸3输出锁.Lock()
	defer 僵尸3输出锁.Unlock()
	needScroll := 僵尸3需要滚动到底
	僵尸3需要滚动到底 = false
	return needScroll
}

func 设置僵尸3层输出(format string, args ...any) {
	text := fmt.Sprintf(format, args...)
	僵尸3输出锁.Lock()
	僵尸3输出行 = append(僵尸3输出行, time.Now().Format("15:04:05  ")+text)
	if len(僵尸3输出行) > 僵尸3输出最大行数 {
		僵尸3输出行 = 僵尸3输出行[len(僵尸3输出行)-僵尸3输出最大行数:]
	}
	僵尸3需要滚动到底 = true
	僵尸3输出锁.Unlock()
	输出("僵尸3", text)
}

func 僵尸3查找黄点并更新输出() {
	ok, x, y := 查找僵尸3黄点坐标()
	if !ok {
		设置僵尸3层输出("没找到黄点，区域=(10,97,260,203)")
		return
	}
	层, matched := 识别僵尸3层数(y)
	if !matched {
		设置僵尸3层输出("找到黄点但层数未匹配：x=%d y=%d", x, y)
		return
	}
	位置 := 层位置{层: 层, X: x, Y: y}
	设置僵尸3层输出("找到黄点：层=%d x=%d y=%d", 位置.层, 位置.X, 位置.Y)
}

func 僵尸3当前层位置() (层位置, bool) {
	ok, x, y := 查找僵尸3黄点坐标()
	if !ok {
		return 层位置{}, false
	}
	层, matched := 识别僵尸3层数(y)
	if !matched {
		return 层位置{}, false
	}
	return 层位置{层: 层, X: x, Y: y}, true
}

func 查找僵尸3黄点坐标() (bool, int, int) {
	if 引擎 == nil {
		return false, -1, -1
	}
	result := 扫描小地图黄点(僵尸3小地图黄点区域, false)
	if !result.Ok {
		return false, -1, -1
	}
	return true, result.X, result.Y
}

func 识别僵尸3层数(y int) (int, bool) {
	bestLayer := 0
	bestDiff := 层Y容差 + 1
	for _, 配置 := range 僵尸3层Y配置表 {
		diff := absInt(y - 配置.Y)
		if diff < bestDiff {
			bestLayer = 配置.层
			bestDiff = diff
		}
	}
	return bestLayer, bestDiff <= 层Y容差
}

func 僵尸3爬绳子() {
	僵尸3爬绳直到成功(func() bool {
		return !程序退出中.Load()
	})
}

func 僵尸3执行爬绳动作(起始位置 层位置, 横向键 int, shouldContinue func() bool) bool {
	displayID := 当前显示ID()
	持续 := 僵尸3爬绳最短持续 + time.Duration(移动随机.Int63n(int64(僵尸3爬绳随机增加)+1))
	设置僵尸3层输出("已到绳子下方，%s后持续按上%dms", 僵尸3爬绳动作文本(横向键), 僵尸3爬绳最短持续.Milliseconds())

	motion.KeyActionDown(motion.KEYCODE_DPAD_UP, displayID)
	if 横向键 != 0 {
		motion.KeyActionDown(横向键, displayID)
	}
	motion.KeyActionDown(motion.KEYCODE_Z, displayID)
	time.Sleep(time.Duration(方向键按下毫秒) * time.Millisecond)
	motion.KeyActionUp(motion.KEYCODE_Z, displayID)
	if 横向键 != 0 {
		motion.KeyActionUp(横向键, displayID)
	}

	上键按下 := true
	松开上键 := func() {
		if 上键按下 {
			motion.KeyActionUp(motion.KEYCODE_DPAD_UP, displayID)
			上键按下 = false
		}
	}
	defer 松开上键()

	start := time.Now()
	deadline := start.Add(持续)
	minHoldDeadline := start.Add(僵尸3爬绳最短持续)
	已到二层 := false
	for shouldContinue() && time.Now().Before(deadline) {
		now := time.Now()
		位置, ok := 僵尸3当前层位置()
		if ok {
			if 位置.层 == 2 {
				if !已到二层 {
					设置僵尸3层输出("已检测到2层，继续按上到4秒")
				}
				已到二层 = true
			}
		}
		if 已到二层 && !now.Before(minHoldDeadline) {
			设置僵尸3层输出("爬绳成功，已持续按上%dms", 僵尸3爬绳最短持续.Milliseconds())
			return true
		}
		time.Sleep(僵尸3爬绳检测间隔)
	}

	位置, ok := 僵尸3当前层位置()
	if ok && 位置.层 == 2 {
		设置僵尸3层输出("爬绳成功")
		return true
	}
	if ok {
		设置僵尸3层输出("爬绳结束未到2层：层=%d x=%d y=%d", 位置.层, 位置.X, 位置.Y)
	}
	return false
}

func 僵尸3爬绳动作文本(横向键 int) string {
	if 横向键 == 0 {
		return "上+Z"
	}
	return "上+" + 键名(横向键) + "+Z"
}

func 僵尸3爬绳直到成功(shouldContinue func() bool) bool {
	设置僵尸3层输出("开始爬绳：目标X=%d", 僵尸3绳子X)
	卡住监控 := &僵尸3卡住监控{}
	for shouldContinue() {
		卡住监控.检查("爬绳卡住检测")
		位置, ok := 僵尸3当前层位置()
		if !ok {
			设置僵尸3层输出("爬绳放弃：无法识别当前位置")
			return false
		}
		if 位置.层 == 2 {
			设置僵尸3层输出("爬绳成功")
			return true
		}
		if 位置.层 != 1 {
			time.Sleep(僵尸3未知位置等待)
			return false
		}
		位置, ok = 僵尸3移动到绳子正下方(shouldContinue)
		if !ok {
			设置僵尸3层输出("移动到绳子下方失败，放弃本次爬绳")
			return false
		}
		设置僵尸3层输出("爬绳尝试1")
		if 僵尸3执行爬绳动作(位置, 僵尸3爬绳横向键(位置), shouldContinue) {
			return true
		}
		设置僵尸3层输出("爬绳检测失败，脱离梯子去其它位置刷")
		僵尸3脱离梯子一次()
		点按空格()
		return false
	}
	return false
}

func (监控 *僵尸3卡住监控) 检查(场景 string) {
	ok, x, y := 查找僵尸3黄点坐标()
	if !ok {
		return
	}

	now := time.Now()
	if !监控.有上次坐标 || absInt(x-监控.上次X) > 僵尸3卡住坐标容差 || absInt(y-监控.上次Y) > 僵尸3卡住坐标容差 {
		监控.有上次坐标 = true
		监控.上次X = x
		监控.上次Y = y
		监控.固定开始 = now
		监控.已脱离 = false
		监控.已报警 = false
		return
	}
	if 监控.固定开始.IsZero() {
		监控.固定开始 = now
		return
	}

	固定时长 := now.Sub(监控.固定开始)
	if 固定时长 >= 僵尸3卡住检测时间 && !监控.已脱离 {
		方向键 := 监控.下一个脱离方向()
		设置僵尸3层输出("%s：坐标%d秒未动 x=%d y=%d，按%s+Z", 场景, int(固定时长.Seconds()), x, y, 键名(方向键))
		按组合键不空格(方向键, motion.KEYCODE_Z, 方向键按下毫秒)
		监控.已脱离 = true
	}
	if 固定时长 >= 僵尸3卡住报警时间 && !监控.已报警 {
		设置僵尸3层输出("%s：卡住超过%d秒，发送钉钉报警", 场景, int(僵尸3卡住报警时间.Seconds()))
		发送钉钉文本("卡住了")
		监控.已报警 = true
	}
}

func (监控 *僵尸3卡住监控) 下一个脱离方向() int {
	if 监控.下次方向 == 0 {
		if 移动随机.Intn(2) == 0 {
			监控.下次方向 = motion.KEYCODE_DPAD_LEFT
		} else {
			监控.下次方向 = motion.KEYCODE_DPAD_RIGHT
		}
	}
	方向键 := 监控.下次方向
	if 方向键 == motion.KEYCODE_DPAD_LEFT {
		监控.下次方向 = motion.KEYCODE_DPAD_RIGHT
	} else {
		监控.下次方向 = motion.KEYCODE_DPAD_LEFT
	}
	return 方向键
}

func 僵尸3脱离梯子一次() {
	方向键 := motion.KEYCODE_DPAD_LEFT
	if 移动随机.Intn(2) == 0 {
		方向键 = motion.KEYCODE_DPAD_RIGHT
	}
	设置僵尸3层输出("尝试脱离梯子：%s+Z", 键名(方向键))
	按组合键不空格(方向键, motion.KEYCODE_Z, 方向键按下毫秒)
	time.Sleep(僵尸3移动间隔)
}

func 僵尸3等待到层(目标层 int, timeout time.Duration, shouldContinue func() bool) bool {
	deadline := time.Now().Add(timeout)
	for shouldContinue() && time.Now().Before(deadline) {
		位置, ok := 僵尸3当前层位置()
		if ok && 位置.层 == 目标层 {
			return true
		}
		time.Sleep(100 * time.Millisecond)
	}
	return false
}

func 僵尸3移动到绳子正下方(shouldContinue func() bool) (层位置, bool) {
	for shouldContinue() {
		位置, ok := 僵尸3当前层位置()
		if !ok {
			time.Sleep(僵尸3未知位置等待)
			continue
		}
		if 位置.层 != 1 {
			return 层位置{}, false
		}
		if 僵尸3绳子位置已到位(位置.X) {
			设置僵尸3层输出("检测到绳子正下方：x=%d y=%d，等待%dms确认滑行", 位置.X, 位置.Y, 僵尸3绳子到位稳定等待.Milliseconds())
			稳定位置, stable := 僵尸3等待绳子坐标稳定(位置, shouldContinue)
			if stable {
				设置僵尸3层输出("绳子坐标已稳定：x=%d y=%d，开始爬", 稳定位置.X, 稳定位置.Y)
				return 稳定位置, true
			}
			continue
		}
		目标X := 僵尸3绳子目标左
		diff := absInt(位置.X - 目标X)
		方向键 := motion.KEYCODE_DPAD_RIGHT
		if 位置.X > 目标X {
			方向键 = motion.KEYCODE_DPAD_LEFT
		}
		if diff <= 僵尸3靠近绳子慢走距离 {
			僵尸3按移动并检测反向击退(位置, 方向键, 僵尸3绳子对齐等待, func() {
				僵尸3按方向键指定时长(方向键, 僵尸3绳子对齐短按)
			})
			continue
		} else if 僵尸3应使用海盗走X(位置.X, 目标X, 方向键, 僵尸3一层左边, 僵尸3一层右边, false) {
			僵尸3按移动并检测反向击退(位置, 方向键, 僵尸3移动间隔, func() {
				按组合键不空格(方向键, motion.KEYCODE_X, 方向键按下毫秒)
			})
		} else {
			僵尸3按移动并检测反向击退(位置, 方向键, 僵尸3移动间隔, func() {
				按方向键(方向键, 方向键按下毫秒)
			})
		}
	}
	return 层位置{}, false
}

func 僵尸3按移动并检测反向击退(原位置 层位置, 方向键 int, 等待 time.Duration, 动作 func()) {
	if 动作 != nil {
		动作()
	}
	if 等待 > 0 {
		time.Sleep(等待)
	}

	新位置, ok := 僵尸3当前层位置()
	if !ok || 新位置.层 != 原位置.层 {
		return
	}
	deltaX := 新位置.X - 原位置.X
	反向 := false
	switch 方向键 {
	case motion.KEYCODE_DPAD_RIGHT:
		反向 = deltaX < -僵尸3击退反向容差
	case motion.KEYCODE_DPAD_LEFT:
		反向 = deltaX > 僵尸3击退反向容差
	}
	if !反向 {
		return
	}

	设置僵尸3层输出("爬绳对齐疑似被怪击退：按%s x=%d->%d，空格攻击2次", 键名(方向键), 原位置.X, 新位置.X)
	僵尸3空格攻击次数(2)
}

func 僵尸3空格攻击次数(times int) {
	for i := 0; i < times; i++ {
		点按空格()
		time.Sleep(40 * time.Millisecond)
	}
}

func 僵尸3等待绳子坐标稳定(原位置 层位置, shouldContinue func() bool) (层位置, bool) {
	deadline := time.Now().Add(僵尸3绳子到位稳定等待)
	for shouldContinue() && time.Now().Before(deadline) {
		time.Sleep(僵尸3绳子对齐等待)
	}
	if !shouldContinue() {
		return 层位置{}, false
	}
	新位置, ok := 僵尸3当前层位置()
	if !ok {
		设置僵尸3层输出("绳子稳定确认失败：未识别到位置")
		return 层位置{}, false
	}
	if 新位置.层 != 原位置.层 {
		设置僵尸3层输出("绳子稳定确认中断：层变化 %d->%d", 原位置.层, 新位置.层)
		return 新位置, false
	}
	if 新位置.X == 原位置.X && 新位置.Y == 原位置.Y && 僵尸3绳子位置已到位(新位置.X) {
		return 新位置, true
	}
	设置僵尸3层输出("绳子坐标滑行：%d,%d -> %d,%d，继续校正", 原位置.X, 原位置.Y, 新位置.X, 新位置.Y)
	return 新位置, false
}

func 僵尸3绳子位置已到位(x int) bool {
	return x == 僵尸3绳子目标左
}

func 僵尸3绳子目标X(x int) int {
	僵尸3绳子学习锁.Lock()
	preferred := 僵尸3绳子偏好目标X
	僵尸3绳子学习锁.Unlock()
	if preferred == 僵尸3绳子目标左 || preferred == 僵尸3绳子目标右 {
		return preferred
	}
	if x <= 僵尸3绳子目标左 {
		return 僵尸3绳子目标左
	}
	return 僵尸3绳子目标右
}

func 僵尸3学习校正到绳子(位置 层位置, 目标X, 方向键 int) {
	diff := absInt(位置.X - 目标X)
	duration := 僵尸3绳子校正时长(diff)
	设置僵尸3层输出("绳子校正：x=%d 目标=%d %s %dms", 位置.X, 目标X, 键名(方向键), duration.Milliseconds())
	僵尸3按方向键指定时长(方向键, duration)
	time.Sleep(僵尸3移动间隔)

	新位置, ok := 僵尸3当前层位置()
	if !ok || 新位置.层 != 1 {
		return
	}
	移动像素 := absInt(新位置.X - 位置.X)
	if 移动像素 <= 0 {
		return
	}
	僵尸3更新绳子学习(duration, 移动像素)
	设置僵尸3层输出("绳子校正结果：%d->%d 移动=%d 学习=%.1fms/px", 位置.X, 新位置.X, 移动像素, 僵尸3当前绳子毫秒每像素())
}

func 僵尸3绳子校正时长(diff int) time.Duration {
	if diff <= 0 {
		return 僵尸3绳子校正最短
	}
	僵尸3绳子学习锁.Lock()
	msPerPixel := 僵尸3绳子毫秒每像素
	multiplier := 僵尸3绳子校正倍率
	僵尸3绳子学习锁.Unlock()
	duration := time.Duration(msPerPixel*multiplier*float64(diff)) * time.Millisecond
	if duration < 僵尸3绳子校正最短 {
		return 僵尸3绳子校正最短
	}
	if duration > 僵尸3绳子校正最长 {
		return 僵尸3绳子校正最长
	}
	return duration
}

func 僵尸3更新绳子学习(duration time.Duration, movedPixels int) {
	if movedPixels <= 0 {
		return
	}
	sample := float64(duration.Milliseconds()) / float64(movedPixels)
	if sample < float64(僵尸3绳子校正最短.Milliseconds())/3 {
		return
	}
	僵尸3绳子学习锁.Lock()
	僵尸3绳子毫秒每像素 = 僵尸3绳子毫秒每像素*0.65 + sample*0.35
	僵尸3绳子学习锁.Unlock()
}

func 僵尸3记录爬绳结果(targetX int, attempt int) {
	if attempt < 1 {
		attempt = 1
	}
	capped := attempt
	if capped > 3 {
		capped = 3
	}

	僵尸3绳子学习锁.Lock()
	僵尸3爬绳结果窗口 = append(僵尸3爬绳结果窗口, capped)
	if len(僵尸3爬绳结果窗口) > 10 {
		僵尸3爬绳结果窗口 = 僵尸3爬绳结果窗口[len(僵尸3爬绳结果窗口)-10:]
	}

	reward := float64(4 - capped)
	switch targetX {
	case 僵尸3绳子目标左:
		僵尸3目标214积分 = 僵尸3目标214积分*0.85 + reward
	case 僵尸3绳子目标右:
		僵尸3目标215积分 = 僵尸3目标215积分*0.85 + reward
	}

	one, two, three := 0, 0, 0
	for _, result := range 僵尸3爬绳结果窗口 {
		switch result {
		case 1:
			one++
		case 2:
			two++
		default:
			three++
		}
	}

	if len(僵尸3爬绳结果窗口) >= 10 {
		if one < 6 {
			僵尸3绳子校正倍率 *= 1.03
		} else if one > 7 && three == 0 {
			僵尸3绳子校正倍率 *= 0.99
		}
		if three > 1 {
			僵尸3绳子校正倍率 *= 1.04
		}
		if 僵尸3绳子校正倍率 < 0.70 {
			僵尸3绳子校正倍率 = 0.70
		}
		if 僵尸3绳子校正倍率 > 1.45 {
			僵尸3绳子校正倍率 = 1.45
		}

		if 僵尸3目标214积分 > 僵尸3目标215积分+0.8 {
			僵尸3绳子偏好目标X = 僵尸3绳子目标左
		} else if 僵尸3目标215积分 > 僵尸3目标214积分+0.8 {
			僵尸3绳子偏好目标X = 僵尸3绳子目标右
		}
	}

	multiplier := 僵尸3绳子校正倍率
	preferred := 僵尸3绳子偏好目标X
	僵尸3绳子学习锁.Unlock()

	设置僵尸3层输出("爬绳统计10次：1次=%d 2次=%d 3次=%d 目标=%d 倍率=%.2f", one, two, three, preferred, multiplier)
}

func 僵尸3当前绳子毫秒每像素() float64 {
	僵尸3绳子学习锁.Lock()
	defer 僵尸3绳子学习锁.Unlock()
	return 僵尸3绳子毫秒每像素
}

func 僵尸3按方向键指定时长(code int, duration time.Duration) {
	displayID := 当前显示ID()
	motion.KeyActionDown(code, displayID)
	time.Sleep(duration)
	motion.KeyActionUp(code, displayID)
}

func 僵尸3爬绳横向键(位置 层位置) int {
	if 僵尸3绳子位置已到位(位置.X) {
		return 0
	}
	if 位置.X < 僵尸3绳子X {
		return motion.KEYCODE_DPAD_RIGHT
	}
	return motion.KEYCODE_DPAD_LEFT
}

func 运行僵尸3循环(runID int64, 启动先三层卖物品 bool) {
	设置僵尸3层输出("僵尸3开始：1-2-1-右边进3层，按24~43分钟卖物品")
	启动N键守护(runID)
	卖物品截止 := 新僵尸3卖物品截止()
	一层爬绳截止 := 新僵尸3一层爬绳截止()
	首次一层 := true
	if 启动先三层卖物品 {
		设置僵尸3层输出("启动时在3层：按计时判断是否卖物品")
		if 僵尸3三层按计时处理后回一层(runID, &卖物品截止) {
			首次一层 = true
			一层爬绳截止 = 新僵尸3一层爬绳截止()
		}
	}
	for 脚本仍应运行(runID) {
		if 僵尸3检查BOSS并换线(runID) {
			首次一层 = true
			一层爬绳截止 = 新僵尸3一层爬绳截止()
			continue
		}
		位置, ok := 僵尸3当前层位置()
		if !ok {
			time.Sleep(僵尸3未知位置等待)
			continue
		}
		if 位置.层 == 2 {
			if !僵尸3二层巡逻三轮(runID) {
				continue
			}
			if !僵尸3二层下到一层直到成功(runID) {
				continue
			}
			if !僵尸3一层到三层按计时后回一层(runID, &卖物品截止) {
				continue
			}
			首次一层 = true
			一层爬绳截止 = 新僵尸3一层爬绳截止()
			continue
		}
		if 位置.层 == 3 {
			if !僵尸3三层按计时处理后回一层(runID, &卖物品截止) {
				continue
			}
			首次一层 = true
			一层爬绳截止 = 新僵尸3一层爬绳截止()
			continue
		}
		if 位置.层 != 1 {
			time.Sleep(僵尸3未知位置等待)
			continue
		}

		if 首次一层 {
			僵尸3一层走到X(runID, 僵尸3一层左边)
			僵尸3一层走到X(runID, 僵尸3一层右边)
			首次一层 = false
		} else {
			僵尸3一层走到X(runID, 僵尸3一层右边)
			僵尸3一层走到X(runID, 僵尸3一层左边)
			僵尸3一层走到X(runID, 僵尸3一层右边)
		}
		if time.Now().Before(一层爬绳截止) {
			continue
		}
		设置僵尸3层输出("1层刷怪时间到，尝试爬绳1次")
		if 位置, ok := 僵尸3当前层位置(); !ok || 位置.层 != 1 {
			if ok {
				设置僵尸3层输出("准备爬绳前发现不在1层：层=%d x=%d y=%d", 位置.层, 位置.X, 位置.Y)
			}
			一层爬绳截止 = 新僵尸3一层爬绳截止()
			continue
		}
		if !僵尸3爬绳直到成功(func() bool { return 脚本仍应运行(runID) }) && 脚本仍应运行(runID) {
			设置僵尸3层输出("本次爬绳失败，回1层其它位置继续刷")
			一层爬绳截止 = 新僵尸3一层爬绳截止()
		}
	}
	设置僵尸3层输出("僵尸3流程已停止")
}

func 新僵尸3一层爬绳截止() time.Time {
	duration := 随机僵尸3一层爬绳间隔()
	设置僵尸3层输出("1层下次爬绳：%d秒后", int(duration.Seconds()))
	return time.Now().Add(duration)
}

func 随机僵尸3一层爬绳间隔() time.Duration {
	spanSeconds := int64((僵尸3一层爬绳最长间隔 - 僵尸3一层爬绳最短间隔) / time.Second)
	if spanSeconds <= 0 {
		return 僵尸3一层爬绳最短间隔
	}
	return 僵尸3一层爬绳最短间隔 + time.Duration(移动随机.Int63n(spanSeconds+1))*time.Second
}

func 新僵尸3卖物品截止() time.Time {
	duration := 随机刷怪周期时长()
	设置僵尸3层输出("下次卖物品：%d分钟后", int(duration.Minutes()))
	return time.Now().Add(duration)
}

func 僵尸3一层到三层按计时后回一层(runID int64, 卖物品截止 *time.Time) bool {
	if !僵尸3一层到三层直到成功(runID) {
		return false
	}
	return 僵尸3三层按计时处理后回一层(runID, 卖物品截止)
}

func 僵尸3三层按计时处理后回一层(runID int64, 卖物品截止 *time.Time) bool {
	if 卖物品截止 != nil && time.Now().After(*卖物品截止) {
		if !僵尸3三层卖物品后回一层(runID) {
			return false
		}
		*卖物品截止 = 新僵尸3卖物品截止()
		return true
	}
	if 卖物品截止 != nil {
		剩余 := time.Until(*卖物品截止)
		if 剩余 < 0 {
			剩余 = 0
		}
		设置僵尸3层输出("3层未到卖物品时间，剩余%d分钟，回1层打怪", int(剩余.Minutes()))
	} else {
		设置僵尸3层输出("3层未到卖物品时间，回1层打怪")
	}
	return 僵尸3三层回一层直到成功(runID)
}

func 僵尸3三层卖物品后回一层(runID int64) bool {
	设置僵尸3层输出("3层开始卖物品：使用海盗同一套买卖物品流程")
	if !执行僵尸3完整买卖物品流程(func() bool { return 脚本仍应运行(runID) }) {
		设置僵尸3层输出("3层卖物品中断，继续打怪流程")
		return false
	}
	设置僵尸3层输出("3层卖物品完成，准备回1层")
	return 僵尸3三层回一层直到成功(runID)
}

func 僵尸3一层到三层直到成功(runID int64) bool {
	设置僵尸3层输出("1层准备进3层：随机右侧触发/右到左出边界")
	shouldContinue := func() bool { return 脚本仍应运行(runID) }
	for 脚本仍应运行(runID) {
		if 僵尸3检查BOSS并换线(runID) {
			return false
		}
		位置, ok := 僵尸3当前层位置()
		if !ok {
			time.Sleep(僵尸3未知位置等待)
			continue
		}
		if 位置.层 == 3 {
			设置僵尸3层输出("已到3层")
			return true
		}
		if 位置.层 != 1 {
			return false
		}
		if 僵尸3随机尝试一层到三层方式(shouldContinue) {
			return true
		}
		time.Sleep(僵尸3移动间隔)
	}
	return false
}

func 僵尸3随机尝试一层到三层方式(shouldContinue func() bool) bool {
	if 移动随机.Intn(2) == 0 {
		if 僵尸3一层到三层右侧触发(shouldContinue) {
			return true
		}
		设置僵尸3层输出("右侧触发未成功，改用右到左出边界")
		return 僵尸3一层到三层从右刷到左(shouldContinue)
	}
	if 僵尸3一层到三层从右刷到左(shouldContinue) {
		return true
	}
	设置僵尸3层输出("右到左出边界未成功，改用右侧触发")
	return 僵尸3一层到三层右侧触发(shouldContinue)
}

func 僵尸3一层到三层右侧触发(shouldContinue func() bool) bool {
	设置僵尸3层输出("1层到3层方式：到%d左边后右+X", 僵尸3去三层触发X)
	if !僵尸3移动到X范围直到(shouldContinue, 1, 僵尸3去三层准备左, 僵尸3去三层准备右, 僵尸3一层左边, 僵尸3去三层触发X) {
		return false
	}
	位置, ok := 僵尸3当前层位置()
	if ok && 位置.层 == 3 {
		设置僵尸3层输出("1层到3层成功：右侧触发")
		return true
	}
	if !ok || 位置.层 != 1 {
		return false
	}
	按组合键不空格(motion.KEYCODE_DPAD_RIGHT, motion.KEYCODE_X, 方向键按下毫秒)
	按空格群攻()
	if 僵尸3等待到层(3, 僵尸3三层换层等待, shouldContinue) {
		设置僵尸3层输出("1层到3层成功：右侧触发")
		return true
	}
	设置僵尸3层输出("右侧触发进3层未检测成功")
	return false
}

func 僵尸3一层到三层从右刷到左(shouldContinue func() bool) bool {
	设置僵尸3层输出("1层到3层方式：从右刷到左，走到x=%d出左边界", 僵尸3左侧进三层目标X)
	if !僵尸3移动到X范围直到(shouldContinue, 1, 僵尸3去三层准备左, 僵尸3去三层准备右, 僵尸3一层左边, 僵尸3去三层触发X) {
		return false
	}

	deadline := time.Now().Add(僵尸3左侧进三层超时)
	for shouldContinue() && time.Now().Before(deadline) {
		位置, ok := 僵尸3当前层位置()
		if ok {
			if 位置.层 == 3 {
				设置僵尸3层输出("1层到3层成功：右到左出边界")
				return true
			}
			if 位置.层 != 1 {
				return false
			}
			if 位置.X <= 僵尸3左侧进三层目标X+僵尸3移动到位容差 {
				按组合键不空格(motion.KEYCODE_DPAD_LEFT, motion.KEYCODE_X, 方向键按下毫秒)
			} else if 僵尸3应使用海盗走X(位置.X, 僵尸3左侧进三层目标X, motion.KEYCODE_DPAD_LEFT, 僵尸3左侧进三层目标X, 僵尸3去三层触发X, false) {
				按组合键不空格(motion.KEYCODE_DPAD_LEFT, motion.KEYCODE_X, 方向键按下毫秒)
			} else {
				按方向键(motion.KEYCODE_DPAD_LEFT, 方向键按下毫秒)
			}
		} else {
			按组合键不空格(motion.KEYCODE_DPAD_LEFT, motion.KEYCODE_X, 方向键按下毫秒)
		}
		if 僵尸3等待到层(3, 僵尸3三层换层等待, shouldContinue) {
			设置僵尸3层输出("1层到3层成功：右到左出边界")
			return true
		}
		time.Sleep(僵尸3移动间隔)
	}
	设置僵尸3层输出("右到左出边界进3层超时")
	return false
}

func 僵尸3移动到X范围直到(shouldContinue func() bool, 层, 目标左, 目标右, 限制左, 限制右 int) bool {
	deadline := time.Now().Add(僵尸3左侧进三层超时)
	for shouldContinue() && time.Now().Before(deadline) {
		位置, ok := 僵尸3当前层位置()
		if !ok {
			time.Sleep(僵尸3未知位置等待)
			continue
		}
		if 位置.层 == 3 && 层 != 3 {
			return true
		}
		if 位置.层 != 层 {
			return false
		}
		if 位置.X >= 目标左 && 位置.X <= 目标右 {
			return true
		}

		方向键 := motion.KEYCODE_DPAD_RIGHT
		目标X := 目标左
		if 位置.X > 目标右 {
			方向键 = motion.KEYCODE_DPAD_LEFT
			目标X = 目标右
		}
		if 僵尸3应使用海盗走X(位置.X, 目标X, 方向键, 限制左, 限制右, false) {
			按组合键不空格(方向键, motion.KEYCODE_X, 方向键按下毫秒)
		} else {
			按方向键(方向键, 方向键按下毫秒)
		}
		time.Sleep(僵尸3移动间隔)
	}
	设置僵尸3层输出("移动到X范围超时：层=%d 目标=%d-%d", 层, 目标左, 目标右)
	return false
}

func 僵尸3三层回一层直到成功(runID int64) bool {
	for 脚本仍应运行(runID) {
		位置, ok := 僵尸3当前层位置()
		if !ok {
			time.Sleep(僵尸3未知位置等待)
			continue
		}
		if 位置.层 == 1 {
			设置僵尸3层输出("3层回1层成功")
			return true
		}
		if 位置.层 != 3 {
			time.Sleep(僵尸3未知位置等待)
			continue
		}
		if 僵尸3三层执行下绳到一层(位置, func() bool { return 脚本仍应运行(runID) }) {
			设置僵尸3层输出("3层下绳回1层成功")
			return true
		}
		time.Sleep(僵尸3下绳重试间隔)
	}
	return false
}

func 僵尸3三层执行下绳到一层(起始位置 层位置, shouldContinue func() bool) bool {
	displayID := 当前显示ID()
	设置僵尸3层输出("3层下1层：按住左+下，过x=%d前检测下降", 僵尸3三层下绳左界)
	motion.KeyActionDown(motion.KEYCODE_DPAD_LEFT, displayID)
	motion.KeyActionDown(motion.KEYCODE_DPAD_DOWN, displayID)
	defer func() {
		motion.KeyActionUp(motion.KEYCODE_DPAD_DOWN, displayID)
		motion.KeyActionUp(motion.KEYCODE_DPAD_LEFT, displayID)
	}()

	deadline := time.Now().Add(僵尸3三层左下下绳超时)
	for shouldContinue() && time.Now().Before(deadline) {
		位置, ok := 僵尸3当前层位置()
		if !ok {
			time.Sleep(僵尸3爬绳检测间隔)
			continue
		}
		if 位置.层 == 1 {
			return true
		}
		已往下移动 := 位置.层 < 起始位置.层 || 位置.Y >= 起始位置.Y+僵尸3爬绳上升Y阈值
		if 已往下移动 {
			设置僵尸3层输出("3层下1层：已下降 x=%d y=%d，左+Z脱离绳子", 位置.X, 位置.Y)
			motion.KeyActionUp(motion.KEYCODE_DPAD_DOWN, displayID)
			motion.KeyActionUp(motion.KEYCODE_DPAD_LEFT, displayID)
			按组合键不空格(motion.KEYCODE_DPAD_LEFT, motion.KEYCODE_Z, 方向键按下毫秒)
			return 僵尸3等待到层(1, 僵尸3三层换层等待, shouldContinue)
		}
		time.Sleep(僵尸3爬绳检测间隔)
	}
	设置僵尸3层输出("3层下1层：左+下超时未检测到下降，重试")
	return false
}

func 僵尸3一层走到X(runID int64, 目标X int) bool {
	return 僵尸3层内走到X(runID, 1, 目标X, 僵尸3一层左边, 僵尸3一层右边, false)
}

func 僵尸3二层巡逻三轮(runID int64) bool {
	for i := 1; i <= 僵尸3二层巡逻次数 && 脚本仍应运行(runID); i++ {
		设置僵尸3层输出("2层巡逻%d/%d：左到右", i, 僵尸3二层巡逻次数)
		if !僵尸3层内走到X(runID, 2, 僵尸3二层可走X左, 僵尸3二层可走X左, 僵尸3二层可走X右, true) {
			return false
		}
		if !僵尸3层内走到X(runID, 2, 僵尸3二层可走X右, 僵尸3二层可走X左, 僵尸3二层可走X右, true) {
			return false
		}
	}
	return 脚本仍应运行(runID)
}

func 僵尸3层内走到X(runID int64, 层, 目标X, 左边, 右边 int, 二层允许X bool) bool {
	for 脚本仍应运行(runID) {
		if 僵尸3检查BOSS并换线(runID) {
			return false
		}
		位置, ok := 僵尸3当前层位置()
		if !ok {
			time.Sleep(僵尸3未知位置等待)
			continue
		}
		if 位置.层 != 层 {
			return false
		}
		diff := absInt(位置.X - 目标X)
		if diff <= 僵尸3移动到位容差 {
			return true
		}

		方向键 := motion.KEYCODE_DPAD_RIGHT
		if 位置.X > 目标X {
			方向键 = motion.KEYCODE_DPAD_LEFT
		}
		if 僵尸3应使用海盗走X(位置.X, 目标X, 方向键, 左边, 右边, 二层允许X) {
			按组合键不空格(方向键, motion.KEYCODE_X, 方向键按下毫秒)
		} else {
			按空格群攻()
			return true
		}
		按空格群攻()
		time.Sleep(僵尸3移动间隔)
	}
	return false
}

func 僵尸3移动到X范围(runID int64, 层, 目标左, 目标右, 限制左, 限制右 int) bool {
	for 脚本仍应运行(runID) {
		if 僵尸3检查BOSS并换线(runID) {
			return false
		}
		位置, ok := 僵尸3当前层位置()
		if !ok {
			time.Sleep(僵尸3未知位置等待)
			continue
		}
		if 位置.层 != 层 {
			return false
		}
		if 位置.X >= 目标左 && 位置.X <= 目标右 {
			return true
		}

		方向键 := motion.KEYCODE_DPAD_RIGHT
		目标X := 目标左
		if 位置.X > 目标右 {
			方向键 = motion.KEYCODE_DPAD_LEFT
			目标X = 目标右
		}
		if 僵尸3应使用海盗走X(位置.X, 目标X, 方向键, 限制左, 限制右, false) {
			按组合键不空格(方向键, motion.KEYCODE_X, 方向键按下毫秒)
		} else {
			按方向键(方向键, 方向键按下毫秒)
		}
		按空格群攻()
		time.Sleep(僵尸3移动间隔)
	}
	return false
}

func 僵尸3应使用海盗走X(x, 目标X, 方向键, 左边, 右边 int, 二层限制 bool) bool {
	nextX := x
	switch 方向键 {
	case motion.KEYCODE_DPAD_LEFT:
		nextX = x - 跳跃移动像素
	case motion.KEYCODE_DPAD_RIGHT:
		nextX = x + 跳跃移动像素
	default:
		return false
	}
	if nextX < 左边 || nextX > 右边 {
		return false
	}
	if 二层限制 && (nextX < 僵尸3二层可走X左 || nextX > 僵尸3二层可走X右) {
		return false
	}
	if 方向键 == motion.KEYCODE_DPAD_LEFT {
		return nextX >= 目标X
	}
	return nextX <= 目标X
}

func 僵尸3二层可用X(x, 目标X, 方向键 int) bool {
	if x < 僵尸3二层可走X左 || x > 僵尸3二层可走X右 {
		return false
	}
	nextX := x
	if 方向键 == motion.KEYCODE_DPAD_LEFT {
		nextX -= 跳跃移动像素
		return nextX >= 僵尸3二层可走X左 && nextX >= 目标X
	}
	if 方向键 == motion.KEYCODE_DPAD_RIGHT {
		nextX += 跳跃移动像素
		return nextX <= 僵尸3二层可走X右 && nextX <= 目标X
	}
	return false
}

func 僵尸3二层下到一层直到成功(runID int64) bool {
	设置僵尸3层输出("2层准备下到1层")
	脱离方向 := motion.KEYCODE_DPAD_LEFT
	for 脚本仍应运行(runID) {
		位置, ok := 僵尸3当前层位置()
		if !ok {
			time.Sleep(僵尸3未知位置等待)
			continue
		}
		if 位置.层 == 1 {
			设置僵尸3层输出("已回到1层")
			return true
		}
		if 位置.层 != 2 {
			time.Sleep(僵尸3未知位置等待)
			continue
		}
		if 位置.X < 僵尸3二层可走X左 {
			按方向键(motion.KEYCODE_DPAD_RIGHT, 方向键按下毫秒)
		} else if 位置.X > 僵尸3二层可走X右 {
			按方向键(motion.KEYCODE_DPAD_LEFT, 方向键按下毫秒)
		} else {
			if 僵尸3执行下绳动作(位置.层, 脱离方向, func() bool { return 脚本仍应运行(runID) }) {
				设置僵尸3层输出("已回到1层")
				return true
			}
			if 脱离方向 == motion.KEYCODE_DPAD_LEFT {
				脱离方向 = motion.KEYCODE_DPAD_RIGHT
			} else {
				脱离方向 = motion.KEYCODE_DPAD_LEFT
			}
		}
		time.Sleep(僵尸3下绳重试间隔)
	}
	return false
}

func 僵尸3执行下绳动作(起始层 int, 脱离方向 int, shouldContinue func() bool) bool {
	displayID := 当前显示ID()
	设置僵尸3层输出("下绳：按住下+Z，2秒未换层则%s+Z", 键名(脱离方向))
	motion.KeyActionDown(motion.KEYCODE_DPAD_DOWN, displayID)
	motion.KeyActionDown(motion.KEYCODE_Z, displayID)
	defer func() {
		motion.KeyActionUp(motion.KEYCODE_Z, displayID)
		motion.KeyActionUp(motion.KEYCODE_DPAD_DOWN, displayID)
	}()

	start := time.Now()
	已脱离 := false
	for shouldContinue() {
		位置, ok := 僵尸3当前层位置()
		if ok && 位置.层 < 起始层 {
			return true
		}
		if !已脱离 && time.Since(start) >= 僵尸3下绳无换层检测 {
			motion.KeyActionUp(motion.KEYCODE_Z, displayID)
			motion.KeyActionUp(motion.KEYCODE_DPAD_DOWN, displayID)
			按组合键不空格(脱离方向, motion.KEYCODE_Z, 方向键按下毫秒)
			motion.KeyActionDown(motion.KEYCODE_DPAD_DOWN, displayID)
			motion.KeyActionDown(motion.KEYCODE_Z, displayID)
			已脱离 = true
			start = time.Now()
		}
		if 已脱离 && time.Since(start) >= 僵尸3下绳无换层检测 {
			return false
		}
		time.Sleep(僵尸3爬绳检测间隔)
	}
	return false
}

func 僵尸3按方向键冰面(code int, diff int) {
	displayID := 当前显示ID()
	duration := 僵尸3冰面按键时长(diff)
	motion.KeyActionDown(code, displayID)
	time.Sleep(duration)
	motion.KeyActionUp(code, displayID)
}

func 僵尸3绳子校正后打怪() {
	for i := 0; i < 僵尸3绳子校正空格次数; i++ {
		点按空格()
		time.Sleep(40 * time.Millisecond)
	}
}

func 僵尸3冰面按键时长(diff int) time.Duration {
	if diff <= 4 {
		return time.Duration(30+移动随机.Intn(11)) * time.Millisecond
	}
	if diff <= 8 {
		return time.Duration(30+移动随机.Intn(16)) * time.Millisecond
	}
	if diff <= 15 {
		return time.Duration(28+移动随机.Intn(18)) * time.Millisecond
	}
	if diff <= 30 {
		return time.Duration(45+移动随机.Intn(26)) * time.Millisecond
	}
	return time.Duration(70+移动随机.Intn(41)) * time.Millisecond
}
