package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Dasongzi1366/AutoGo/motion"
)

type 层位置 struct {
	层 int
	X int
	Y int
}

type 层Y配置 struct {
	层 int
	Y int
}

var 层Y配置表 = []层Y配置{
	{层: 4, Y: 135},
	{层: 3, Y: 150},
	{层: 2, Y: 165},
	{层: 1, Y: 180},
}

var 移动随机 = rand.New(rand.NewSource(time.Now().UnixNano()))

const (
	跳跃移动像素       = 17
	移动到位容差       = 2
	移动校正最多次数     = 18
	移动校正间隔       = 160 * time.Millisecond
	失败等待移动间隔     = 450 * time.Millisecond
	层切换等待        = 900 * time.Millisecond
	爬绳最多持续       = 2500 * time.Millisecond
	爬绳检查间隔       = 180 * time.Millisecond
	爬绳最大尝试次数     = 2
	爬绳X偏移容差      = 5
	爬绳基准像素       = 6
	爬绳基准延迟毫秒     = 200
	爬绳人物稳定间隔     = 10 * time.Millisecond
	爬绳人物稳定容差     = 2
	爬绳人物稳定次数     = 10
	爬绳失败等待       = 3 * time.Second
	上楼两次间隔       = 250 * time.Millisecond
	上楼第二次持续      = 500 * time.Millisecond
	二层到三层上楼第二次持续 = 上楼第二次持续 * 2
)

func 当前层位置() (层位置, bool) {
	ok, x, y := 查找FEFE24坐标()
	if !ok {
		输出("未找到小地图黄点，无法判断当前层")
		return 层位置{}, false
	}
	层, ok := 识别层数(y)
	if !ok {
		输出("小地图黄点Y轴未匹配层数", "x=", x, "y=", y)
		return 层位置{}, false
	}
	return 层位置{层: 层, X: x, Y: y}, true
}

func 识别层数(y int) (int, bool) {
	bestLayer := 0
	bestDiff := 层Y容差 + 1
	for _, 配置 := range 层Y配置表 {
		diff := absInt(y - 配置.Y)
		if diff < bestDiff {
			bestLayer = 配置.层
			bestDiff = diff
		}
	}
	return bestLayer, bestDiff <= 层Y容差
}

func 当前层移动配置(层 int) (层移动配置, bool) {
	for _, 配置 := range 层移动配置表 {
		if 配置.层 == 层 {
			return 配置, true
		}
	}
	return 层移动配置{}, false
}

func 移动到当前层限制范围(层 int) bool {
	配置, ok := 当前层移动配置(层)
	if !ok {
		输出("当前层未配置移动范围", "层=", 层)
		return false
	}
	return 移动到X范围(配置.左边, 配置.右边, 配置.左边, 配置.右边)
}

func 移动到当前层限制范围换层(层 int) bool {
	配置, ok := 当前层移动配置(层)
	if !ok {
		输出("当前层未配置移动范围", "层=", 层)
		return false
	}
	return 移动到X范围换层(配置.左边, 配置.右边, 配置.左边, 配置.右边)
}

func 移动到固定X(目标X int, 层 int) bool {
	限制左, 限制右 := 0, 200
	if 配置, ok := 当前层移动配置(层); ok {
		限制左, 限制右 = 配置.左边, 配置.右边
	}
	return 移动到X范围(目标X-移动到位容差, 目标X+移动到位容差, 限制左, 限制右)
}

func 移动到固定X换层(目标X int, 层 int) bool {
	限制左, 限制右 := 0, 200
	if 配置, ok := 当前层移动配置(层); ok {
		限制左, 限制右 = 配置.左边, 配置.右边
	}
	return 移动到X范围换层(目标X-移动到位容差, 目标X+移动到位容差, 限制左, 限制右)
}

func 移动到X范围(目标左, 目标右, 限制左, 限制右 int) bool {
	for i := 0; i < 移动校正最多次数; i++ {
		位置, ok := 当前层位置()
		if !ok {
			return false
		}
		if 位置.X >= 目标左 && 位置.X <= 目标右 {
			return true
		}
		if 位置.X < 目标左 {
			随机向右移动(位置.X, 目标左, 目标右, 限制左, 限制右)
		} else {
			随机向左移动(位置.X, 目标左, 目标右, 限制左, 限制右)
		}
		time.Sleep(移动校正间隔)
	}
	位置, ok := 当前层位置()
	if ok {
		输出("移动到目标X范围失败", "层=", 位置.层, "x=", 位置.X, "目标=", 目标左, "-", 目标右)
	}
	return false
}

func 移动到X范围换层(目标左, 目标右, 限制左, 限制右 int) bool {
	for i := 0; i < 移动校正最多次数; i++ {
		位置, ok := 当前层位置()
		if !ok {
			return false
		}
		if YOLO当前层附近有怪物(位置, 近怪距离阈值) {
			打怪动作(位置)
			time.Sleep(近怪攻击间隔)
			continue
		}
		if 位置.X >= 目标左 && 位置.X <= 目标右 {
			return true
		}
		if 位置.X < 目标左 {
			随机向右移动不空格(位置.X, 目标左, 目标右, 限制左, 限制右)
		} else {
			随机向左移动不空格(位置.X, 目标左, 目标右, 限制左, 限制右)
		}
		time.Sleep(移动校正间隔)
	}
	位置, ok := 当前层位置()
	if ok {
		输出("移动到目标X范围失败", "层=", 位置.层, "x=", 位置.X, "目标=", 目标左, "-", 目标右)
	}
	return false
}

func 随机向左移动(x, 目标左, 目标右, 限制左, 限制右 int) {
	随机方向移动(motion.KEYCODE_DPAD_LEFT, x, 目标左, 目标右, 限制左, 限制右)
}

func 随机向右移动(x, 目标左, 目标右, 限制左, 限制右 int) {
	随机方向移动(motion.KEYCODE_DPAD_RIGHT, x, 目标左, 目标右, 限制左, 限制右)
}

func 随机向左移动不空格(x, 目标左, 目标右, 限制左, 限制右 int) {
	随机方向移动不空格(motion.KEYCODE_DPAD_LEFT, x, 目标左, 目标右, 限制左, 限制右)
}

func 随机向右移动不空格(x, 目标左, 目标右, 限制左, 限制右 int) {
	随机方向移动不空格(motion.KEYCODE_DPAD_RIGHT, x, 目标左, 目标右, 限制左, 限制右)
}

func 随机方向移动(方向键 int, x, 目标左, 目标右, 限制左, 限制右 int) {
	用跳跃 := 移动随机.Intn(2) == 0 && 可用跳跃移动(方向键, x, 目标左, 目标右, 限制左, 限制右)
	if 用跳跃 {
		按组合键不空格(方向键, motion.KEYCODE_X, 方向键按下毫秒)
		按空格群攻()
		return
	}
	按方向键(方向键, 方向键按下毫秒)
	time.Sleep(50 * time.Millisecond)
	按空格群攻()
}

func 随机方向移动不空格(方向键 int, x, 目标左, 目标右, 限制左, 限制右 int) {
	用跳跃 := 移动随机.Intn(2) == 0 && 可用跳跃移动(方向键, x, 目标左, 目标右, 限制左, 限制右)
	if 用跳跃 {
		按组合键不空格(方向键, motion.KEYCODE_X, 方向键按下毫秒)
		return
	}
	按方向键(方向键, 方向键按下毫秒)
}

func 可用跳跃移动(方向键 int, x, 目标左, 目标右, 限制左, 限制右 int) bool {
	nextX := x
	switch 方向键 {
	case motion.KEYCODE_DPAD_LEFT:
		nextX = x - 跳跃移动像素
	case motion.KEYCODE_DPAD_RIGHT:
		nextX = x + 跳跃移动像素
	default:
		return false
	}
	if nextX < 限制左 || nextX > 限制右 {
		return false
	}
	if 方向键 == motion.KEYCODE_DPAD_LEFT {
		return nextX >= 目标左
	}
	return nextX <= 目标右
}

func 三层到四层() bool {
	return 跳跃换层(3, 4, 96, 100, motion.KEYCODE_DPAD_UP, motion.KEYCODE_X)
}

func 四层到三层() bool {
	return 跳跃换层(4, 3, 96, 100, motion.KEYCODE_DPAD_DOWN, motion.KEYCODE_X)
}

func 三层到二层() bool {
	位置, ok := 当前层位置()
	if !ok || 位置.层 != 3 {
		输出("三层到二层失败：当前不在3层", "层=", 位置.层)
		return false
	}
	if !移动到当前层限制范围换层(3) {
		return false
	}
	输出("3层到2层：按下键+X")
	按组合键不空格(motion.KEYCODE_DPAD_DOWN, motion.KEYCODE_X, 方向键按下毫秒)
	return 等待到层(2, 层切换等待)
}

func 二层到一层() bool {
	位置, ok := 当前层位置()
	if !ok || 位置.层 != 2 {
		输出("二层到一层失败：当前不在2层", "层=", 位置.层)
		return false
	}
	if !移动到当前层限制范围换层(2) {
		return false
	}
	输出("2层到1层：按下键+X")
	按组合键不空格(motion.KEYCODE_DPAD_DOWN, motion.KEYCODE_X, 方向键按下毫秒)
	return 等待到层(1, 层切换等待)
}

func 二层到三层() bool {
	if 随机尝试上楼方式(
		func() bool { return 上楼梯到层(2, 3, 68) },
		func() bool { return 爬绳到层(2, 3, 105) },
	) {
		return true
	}
	输出("2层到3层失败：上楼梯/爬绳均未成功")
	return false
}

func 一层到二层() bool {
	if 随机尝试上楼方式(
		func() bool { return 上楼梯到层(1, 2, 52) },
		func() bool { return 爬绳到层(1, 2, 随机选择坐标点(123, 72)) },
	) {
		return true
	}
	输出("1层到2层失败：上楼梯/爬绳均未成功，短暂移动等待后再尝试")
	失败等待并短暂移动(爬绳失败等待)
	return false
}

func 随机尝试上楼方式(方式A, 方式B func() bool) bool {
	if 移动随机.Intn(2) == 0 {
		if 方式A() {
			return true
		}
		return 方式B()
	}
	if 方式B() {
		return true
	}
	return 方式A()
}

func 上楼梯到层(起始层, 目标层, 目标X int) bool {
	位置, ok := 当前层位置()
	if !ok || 位置.层 != 起始层 {
		输出("上楼梯失败：当前层不匹配", "当前层=", 位置.层, "起始层=", 起始层, "目标层=", 目标层)
		return false
	}
	if !移动到固定X换层(目标X, 起始层) {
		输出("上楼梯移动到X失败", "当前X=", 位置.X, "目标X=", 目标X)
		return false
	}

	第二次持续 := 上楼第二次持续
	if 起始层 == 2 && 目标层 == 3 {
		第二次持续 = 二层到三层上楼第二次持续
	}
	按上楼上Z两次(第二次持续)
	return 等待到层(目标层, 层切换等待)
}

func 失败等待并短暂移动(duration time.Duration) {
	deadline := time.Now().Add(duration)
	for time.Now().Before(deadline) {
		位置, ok := 当前层位置()
		if !ok {
			time.Sleep(失败等待移动间隔)
			continue
		}
		配置, ok := 当前层移动配置(位置.层)
		if !ok {
			time.Sleep(失败等待移动间隔)
			continue
		}

		if 位置.X <= 配置.左边+跳跃移动像素 {
			随机向右移动(位置.X, 配置.左边, 配置.右边, 配置.左边, 配置.右边)
		} else if 位置.X >= 配置.右边-跳跃移动像素 {
			随机向左移动(位置.X, 配置.左边, 配置.右边, 配置.左边, 配置.右边)
		} else if 移动随机.Intn(2) == 0 {
			随机向左移动(位置.X, 配置.左边, 配置.右边, 配置.左边, 配置.右边)
		} else {
			随机向右移动(位置.X, 配置.左边, 配置.右边, 配置.左边, 配置.右边)
		}
		time.Sleep(失败等待移动间隔)
	}
}

func 跳跃换层(起始层, 目标层, 目标左, 目标右, 方向键, 动作键 int) bool {
	位置, ok := 当前层位置()
	if !ok || 位置.层 != 起始层 {
		输出("跳跃换层失败：当前层不匹配", "当前层=", 位置.层, "起始层=", 起始层, "目标层=", 目标层)
		return false
	}

	限制左, 限制右 := 0, 200
	if 配置, ok := 当前层移动配置(起始层); ok {
		限制左, 限制右 = 配置.左边, 配置.右边
	}
	if !移动到X范围换层(目标左, 目标右, 限制左, 限制右) {
		return false
	}

	输出("跳跃换层：按组合键", "起始层=", 起始层, "目标层=", 目标层)
	按组合键后处理空格(方向键, 动作键, 方向键按下毫秒, !(方向键 == motion.KEYCODE_DPAD_DOWN && 动作键 == motion.KEYCODE_X))
	return 等待到层(目标层, 层切换等待)
}

func 爬绳到层(起始层, 目标层, 目标X int) bool {
	for attempt := 1; attempt <= 爬绳最大尝试次数; attempt++ {
		位置, ok := 当前层位置()
		if !ok || 位置.层 != 起始层 {
			输出("爬绳失败：当前层不匹配", "当前层=", 位置.层, "起始层=", 起始层, "目标层=", 目标层)
			return false
		}
		if !移动到爬绳X范围(目标X, 起始层) {
			输出("爬绳移动到范围失败", "当前X=", 位置.X, "目标X=", 目标X, "容差=", 爬绳X偏移容差)
			continue
		}

		位置, ok = 当前层位置()
		if !ok {
			continue
		}
		diff := 位置.X - 目标X
		if absInt(diff) > 爬绳X偏移容差 {
			输出("爬绳位置未进入容差", "当前X=", 位置.X, "目标X=", 目标X, "容差=", 爬绳X偏移容差)
			continue
		}
		横向键 := 0
		if diff < 0 {
			横向键 = motion.KEYCODE_DPAD_RIGHT
		} else if diff > 0 {
			横向键 = motion.KEYCODE_DPAD_LEFT
		}

		偏移像素 := absInt(diff)
		方向动作时长 := 爬绳方向动作时长(偏移像素)
		动作名 := 键名(motion.KEYCODE_Z)
		if 横向键 != 0 {
			动作名 = 键名(横向键) + "+" + 键名(motion.KEYCODE_Z)
		}
		输出(fmt.Sprintf("爬梯 距离=%d 键=%s(%dms) 上(%dms)", 偏移像素, 动作名, 方向动作时长.Milliseconds(), 爬绳最多持续.Milliseconds()))
		if 按上并校正爬绳直到到层(目标层, 目标X, 横向键, 偏移像素, 爬绳最多持续) {
			return true
		}
	}
	return false
}

func 移动到爬绳X范围(目标X int, 层 int) bool {
	限制左, 限制右 := 0, 200
	if 配置, ok := 当前层移动配置(层); ok {
		限制左, 限制右 = 配置.左边, 配置.右边
	}
	return 移动到X范围换层(目标X-爬绳X偏移容差, 目标X+爬绳X偏移容差, 限制左, 限制右)
}

func 等待到层(目标层 int, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		位置, ok := 当前层位置()
		if ok && 位置.层 == 目标层 {
			输出("已到达目标层", "层=", 目标层, "x=", 位置.X, "y=", 位置.Y)
			return true
		}
		time.Sleep(爬绳检查间隔)
	}
	位置, ok := 当前层位置()
	if ok {
		输出("等待到层超时", "目标层=", 目标层, "当前层=", 位置.层, "x=", 位置.X, "y=", 位置.Y)
	}
	return false
}

func 按上并校正爬绳直到到层(目标层, 目标X, 横向键, 偏移像素 int, timeout time.Duration) bool {
	displayID := 当前显示ID()
	动作时长 := 爬绳方向动作时长(偏移像素)
	if 横向键 != 0 {
		motion.KeyActionDown(横向键, displayID)
	}
	motion.KeyActionDown(motion.KEYCODE_Z, displayID)
	time.Sleep(动作时长)
	motion.KeyActionUp(motion.KEYCODE_Z, displayID)
	if 横向键 != 0 {
		motion.KeyActionUp(横向键, displayID)
	}
	motion.KeyActionDown(motion.KEYCODE_DPAD_UP, displayID)
	上键按下 := true
	松开上键 := func() {
		if 上键按下 {
			motion.KeyActionUp(motion.KEYCODE_DPAD_UP, displayID)
			上键按下 = false
		}
	}
	defer 松开上键()

	deadline := time.Now().Add(timeout)
	已检测到目标层 := false
	稳定基准Y := 0
	稳定次数 := 0
	for time.Now().Before(deadline) {
		位置, ok := 当前层位置()
		if ok && 位置.层 == 目标层 {
			if !已检测到目标层 {
				已检测到目标层 = true
				输出("爬绳检测到目标层，继续等待人物Y稳定", "目标层=", 目标层, "x=", 位置.X, "y=", 位置.Y)
			}
			personY, ok := YOLO人物中心Y即时()
			if ok {
				if 稳定次数 == 0 || absInt(personY-稳定基准Y) > 爬绳人物稳定容差 {
					稳定基准Y = personY
					稳定次数 = 1
				} else {
					稳定次数++
				}
				if 稳定次数 >= 爬绳人物稳定次数 {
					输出("爬绳人物Y稳定，确认换层成功", "目标层=", 目标层, "人物Y=", personY, "次数=", 稳定次数)
					return true
				}
			}
		}
		if !已检测到目标层 {
			稳定次数 = 0
			稳定基准Y = 0
		}
		time.Sleep(爬绳人物稳定间隔)
	}

	位置, ok := 当前层位置()
	if ok && 位置.层 == 目标层 {
		输出("爬绳成功", "目标层=", 目标层, "x=", 位置.X, "y=", 位置.Y)
		return true
	}
	if ok {
		输出("爬绳结束未到层", "目标层=", 目标层, "当前层=", 位置.层, "x=", 位置.X, "y=", 位置.Y)
		return false
	}
	输出("爬绳结束无法判断层", "目标层=", 目标层, "目标X=", 目标X)
	return false
}

func 爬绳方向动作时长(偏移像素 int) time.Duration {
	if 偏移像素 <= 0 {
		return 随机点按时长()
	}
	ms := (偏移像素*爬绳基准延迟毫秒 + 爬绳基准像素/2) / 爬绳基准像素
	if ms < 10 {
		ms = 10
	}
	return time.Duration(ms) * time.Millisecond
}

func 点按键(code, displayID int) {
	duration := 随机点按时长()
	motion.KeyActionDown(code, displayID)
	time.Sleep(duration)
	motion.KeyActionUp(code, displayID)
}

func 随机点按时长() time.Duration {
	return time.Duration(10+移动随机.Intn(21)) * time.Millisecond
}

func 按组合键(方向键, 动作键, ms int) {
	按组合键后处理空格(方向键, 动作键, ms, true)
}

func 按组合键不空格(方向键, 动作键, ms int) {
	按组合键后处理空格(方向键, 动作键, ms, false)
}

func 按组合键同时不空格(方向键, 动作键, ms int) {
	displayID := 当前显示ID()
	if ms <= 0 {
		ms = 方向键按下毫秒
	}
	motion.KeyActionDown(方向键, displayID)
	motion.KeyActionDown(动作键, displayID)
	time.Sleep(time.Duration(ms) * time.Millisecond)
	motion.KeyActionUp(动作键, displayID)
	motion.KeyActionUp(方向键, displayID)
}

func 按上楼上Z两次(第二次持续 time.Duration) {
	按组合键同时不空格(motion.KEYCODE_DPAD_UP, motion.KEYCODE_Z, 方向键按下毫秒)
	time.Sleep(上楼两次间隔)
	按三键同时不空格(motion.KEYCODE_DPAD_UP, motion.KEYCODE_Z, motion.KEYCODE_DPAD_RIGHT, int(第二次持续/time.Millisecond))
}

func 按三键同时不空格(按键A, 按键B, 按键C, ms int) {
	displayID := 当前显示ID()
	if ms <= 0 {
		ms = 方向键按下毫秒
	}
	motion.KeyActionDown(按键A, displayID)
	motion.KeyActionDown(按键B, displayID)
	motion.KeyActionDown(按键C, displayID)
	time.Sleep(time.Duration(ms) * time.Millisecond)
	motion.KeyActionUp(按键C, displayID)
	motion.KeyActionUp(按键B, displayID)
	motion.KeyActionUp(按键A, displayID)
}

func 按组合键后处理空格(方向键, 动作键, ms int, x后空格 bool) {
	displayID := 当前显示ID()
	if ms <= 0 {
		ms = 方向键按下毫秒
	}
	motion.KeyActionDown(方向键, displayID)
	time.Sleep(20 * time.Millisecond)
	motion.KeyActionDown(动作键, displayID)
	time.Sleep(time.Duration(ms) * time.Millisecond)
	motion.KeyActionUp(动作键, displayID)
	motion.KeyActionUp(方向键, displayID)
	if x后空格 && 动作键 == motion.KEYCODE_X {
		time.Sleep(50 * time.Millisecond)
		按空格群攻()
	}
}

func 键名(code int) string {
	switch code {
	case motion.KEYCODE_DPAD_LEFT:
		return "左"
	case motion.KEYCODE_DPAD_RIGHT:
		return "右"
	case motion.KEYCODE_DPAD_UP:
		return "上"
	case motion.KEYCODE_DPAD_DOWN:
		return "下"
	case motion.KEYCODE_X:
		return "X"
	case motion.KEYCODE_Z:
		return "Z"
	case motion.KEYCODE_SPACE:
		return "空格"
	case motion.KEYCODE_DEL:
		return "Del"
	default:
		return fmt.Sprintf("Key%d", code)
	}
}

func 键动作文本(name string, duration time.Duration) string {
	return fmt.Sprintf("%s(%dms)", name, duration.Milliseconds())
}

func 随机选择坐标点(points ...int) int {
	if len(points) == 0 {
		return 0
	}
	return points[移动随机.Intn(len(points))]
}

func 当前显示ID() int {
	if 引擎 != nil {
		return 引擎.displayID()
	}
	return 屏幕ID
}
