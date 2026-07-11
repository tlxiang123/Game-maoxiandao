package main

import (
	"time"

	"github.com/Dasongzi1366/AutoGo/motion"
)

const (
	驻留移动间隔   = 450 * time.Millisecond
	近怪攻击间隔   = 300 * time.Millisecond
	Del键间隔   = 3 * time.Second
	怪物识别间隔   = 1 * time.Second
	离层怪物阈值   = 5
	离层倒计时    = 8 * time.Second
	近怪距离阈值   = 200
	最大驻留时间   = 20 * time.Second
	打怪空格次数   = 3
	打怪N键最短间隔 = 30 * time.Second
	打怪N键最长间隔 = 120 * time.Second
	N键按下最短时长 = 120 * time.Millisecond
	N键按下最长时长 = 220 * time.Millisecond
	换层稳定确认时间 = 300 * time.Millisecond
)

var (
	层内移动方向 = map[int]int{}
	层访问时间  = map[int]time.Time{}
)

type 离层原因 int

const (
	离层原因无 离层原因 = iota
	离层原因空层
	离层原因低怪倒计时
	离层原因最大驻留
)

type 层稳定状态 struct {
	稳定层  int
	候选层  int
	候选开始 time.Time
}

func 运行图色循环(runID int64, 启动先四层买卖 bool) {
	启动怪物识别后台(runID)
	刷怪截止 := 新刷怪周期截止()
	本次层已到右边 := map[int]bool{}
	稳定状态 := 层稳定状态{}
	输出("开始脚本循环：1-3层刷怪，到周期后4层买卖买药")
	for 脚本仍应运行(runID) {
		位置, ok := 当前层位置()
		if !ok {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		if !确认稳定层变化(位置.层, 本次层已到右边, &稳定状态) {
			time.Sleep(200 * time.Millisecond)
			continue
		}

		if 启动先四层买卖 {
			启动先四层买卖 = false
			if 位置.层 == 4 {
				输出("买卖物品 启动时在4层，优先执行买卖买药")
				if 执行四层买卖买药后回三层(func() bool { return 脚本仍应运行(runID) }) {
					刷怪截止 = 新刷怪周期截止()
				} else {
					time.Sleep(1 * time.Second)
				}
				continue
			}
			输出("买卖物品 启动时4层优先买卖取消：当前不在4层", "层=", 位置.层)
		}

		if 位置.层 >= 1 && 位置.层 <= 3 && !time.Now().Before(刷怪截止) {
			输出("怪物 刷怪周期到，前往4层买卖买药")
			if 前往四层(func() bool { return 脚本仍应运行(runID) }) && 执行四层买卖买药后回三层(func() bool { return 脚本仍应运行(runID) }) {
				刷怪截止 = 新刷怪周期截止()
			} else {
				time.Sleep(1 * time.Second)
			}
			continue
		}

		switch 位置.层 {
		case 4:
			if 执行四层买卖买药后回三层(func() bool { return 脚本仍应运行(runID) }) {
				刷怪截止 = 新刷怪周期截止()
			} else {
				time.Sleep(1 * time.Second)
			}
		case 1:
			记录层访问(1)
			if !确保本次层已走到最右边(runID, 1, 本次层已到右边) {
				time.Sleep(500 * time.Millisecond)
				continue
			}
			原因, ok := 驻留当前层直到可离开(runID)
			if !ok {
				break
			}
			if 准备换层前确认(位置, 原因) {
				一层到二层()
			}
		case 2:
			记录层访问(2)
			if !确保本次层已走到最右边(runID, 2, 本次层已到右边) {
				time.Sleep(500 * time.Millisecond)
				continue
			}
			原因, ok := 驻留当前层直到可离开(runID)
			if !ok {
				break
			}
			if 准备换层前确认(位置, 原因) {
				if 二层下一目标层() == 3 {
					二层到三层()
				} else {
					二层到一层()
				}
			}
		case 3:
			记录层访问(3)
			原因, ok := 驻留当前层直到可离开(runID)
			if !ok {
				break
			}
			if 准备换层前确认(位置, 原因) {
				三层到二层()
			}
		default:
			输出("当前层不在1-4层巡逻范围内", "层=", 位置.层, "x=", 位置.X, "y=", 位置.Y)
			time.Sleep(1 * time.Second)
		}
	}
	if 脚本运行序号.Load() == runID {
		脚本运行中.Store(false)
	}
	释放所有按键()
	输出("脚本循环已停止")
}

func 脚本仍应运行(runID int64) bool {
	return 脚本运行中.Load() && !程序退出中.Load() && 脚本运行序号.Load() == runID
}

func 确认稳定层变化(当前层 int, 已到右边 map[int]bool, 状态 *层稳定状态) bool {
	now := time.Now()
	if 状态.稳定层 == 0 {
		状态.稳定层 = 当前层
		重置首次走右边(当前层, 已到右边, "怪物 初始层确认，重置首次走右边")
		return true
	}
	if 当前层 == 状态.稳定层 {
		if 状态.候选层 != 0 {
			输出("怪物 换层未稳定，视为上楼/换层未成功", "稳定层=", 状态.稳定层, "候选层=", 状态.候选层)
			状态.候选层 = 0
			状态.候选开始 = time.Time{}
		}
		return true
	}
	if 状态.候选层 != 当前层 {
		状态.候选层 = 当前层
		状态.候选开始 = now
		输出("怪物 检测到换层，等待稳定确认", "当前稳定层=", 状态.稳定层, "候选层=", 当前层, "毫秒=", 换层稳定确认时间.Milliseconds())
		return false
	}
	if now.Sub(状态.候选开始) < 换层稳定确认时间 {
		return false
	}

	确认稳定换层(当前层, 已到右边, 状态)
	return true
}

func 确认稳定换层(当前层 int, 已到右边 map[int]bool, 状态 *层稳定状态) {
	原层 := 状态.稳定层
	状态.稳定层 = 当前层
	状态.候选层 = 0
	状态.候选开始 = time.Time{}
	重置首次走右边(当前层, 已到右边, "怪物 换层已稳定，重置首次走右边")
	输出("怪物 稳定换层确认", "从=", 原层, "到=", 当前层)
}

func 重置首次走右边(层 int, 已到右边 map[int]bool, 文本 string) bool {
	if 层 != 1 && 层 != 2 {
		return false
	}
	已到右边[层] = false
	输出(文本, "层=", 层)
	return true
}

func 确保本次层已走到最右边(runID int64, 层 int, 已到右边 map[int]bool) bool {
	if 层 != 1 && 层 != 2 {
		return true
	}
	if 已到右边[层] {
		return true
	}
	配置, ok := 当前层移动配置(层)
	if !ok {
		输出("怪物 首次走右边失败：当前层未配置移动范围", "层=", 层)
		return false
	}

	目标左 := 配置.右边 - 移动到位容差
	if 目标左 < 配置.左边 {
		目标左 = 配置.左边
	}
	输出("怪物 当前层首次到达，先走到最右边", "层=", 层, "目标=", 目标左, "-", 配置.右边)
	for i := 0; i < 移动校正最多次数 && 脚本仍应运行(runID); i++ {
		位置, ok := 当前层位置()
		if !ok {
			time.Sleep(移动校正间隔)
			continue
		}
		if 位置.层 != 层 {
			输出("怪物 首次走右边中断：层已变化", "目标层=", 层, "当前层=", 位置.层)
			return false
		}
		if 位置.X >= 目标左 {
			已到右边[层] = true
			输出("怪物 当前层已走到最右边", "层=", 层, "x=", 位置.X)
			return true
		}
		随机向右移动(位置.X, 目标左, 配置.右边, 配置.左边, 配置.右边)
		time.Sleep(移动校正间隔)
	}
	位置, ok := 当前层位置()
	if ok {
		输出("怪物 首次走右边失败", "层=", 层, "x=", 位置.X, "目标=", 目标左, "-", 配置.右边)
	}
	return false
}

func 准备换层前确认(位置 层位置, 原因 离层原因) bool {
	最新位置, ok := 当前层位置()
	if ok {
		位置 = 最新位置
	}
	当前层怪物数, ok := 打印怪物层统计并取当前层数量(位置)
	if ok && 当前层怪物数 <= 0 {
		return true
	}
	if 原因 == 离层原因空层 {
		return false
	}
	return 原因 == 离层原因最大驻留 || 原因 == 离层原因低怪倒计时
}

func 记录层访问(层 int) {
	if 层 < 1 || 层 > 3 {
		return
	}
	层访问时间[层] = time.Now()
}

func 二层下一目标层() int {
	一层时间 := 层访问时间[1]
	三层时间 := 层访问时间[3]
	if 三层时间.IsZero() {
		return 3
	}
	if 一层时间.IsZero() {
		return 1
	}
	if 三层时间.Before(一层时间) {
		return 3
	}
	return 1
}

func 等待YOLO加载完成(runID int64) bool {
	nextLog := time.Now()
	for 脚本仍应运行(runID) {
		if YOLO已加载完成() {
			输出("怪物 YOLO已加载 4层准备下跳")
			return true
		}
		if !time.Now().Before(nextLog) {
			输出("怪物 4层等待YOLO加载")
			nextLog = time.Now().Add(怪物识别间隔)
		}
		time.Sleep(100 * time.Millisecond)
	}
	return false
}

func 驻留当前层直到可离开(runID int64) (离层原因, bool) {
	deadline := time.Now().Add(最大驻留时间)
	nextMove := time.Now()
	nextDel := time.Now()
	nextMonster := time.Now()
	var 低怪开始 time.Time
	var 空层开始 time.Time
	for 脚本仍应运行(runID) {
		now := time.Now()
		if !now.Before(deadline) {
			return 离层原因最大驻留, true
		}
		位置, ok := 当前层位置()
		if !ok {
			time.Sleep(200 * time.Millisecond)
			continue
		}

		if !now.Before(nextMove) {
			if YOLO当前层附近有怪物(位置, 近怪距离阈值) {
				打怪动作(位置)
				nextMove = now.Add(近怪攻击间隔)
			} else {
				层内追怪或单向游走(位置)
				nextMove = now.Add(驻留移动间隔)
			}
		}
		if !now.Before(nextDel) {
			按Del键()
			nextDel = now.Add(Del键间隔)
		}
		if !now.Before(nextMonster) {
			当前层怪物数, ok := 打印怪物层统计并取当前层数量(位置)
			if ok && 当前层怪物数 <= 0 {
				if 空层开始.IsZero() {
					空层开始 = now
					输出("怪物 当前层暂未识别到怪，继续搜索确认", "层=", 位置.层, "秒=", int(离层倒计时.Seconds()))
				} else if now.Sub(空层开始) >= 离层倒计时 {
					输出("怪物 当前层连续未识别到怪，准备换层", "层=", 位置.层)
					return 离层原因空层, true
				}
			} else {
				空层开始 = time.Time{}
				if !ok || YOLO当前层附近有怪物(位置, 近怪距离阈值) || 当前层怪物数 >= 离层怪物阈值 {
					低怪开始 = time.Time{}
				} else if 低怪开始.IsZero() {
					低怪开始 = now
				} else if now.Sub(低怪开始) >= 离层倒计时 {
					return 离层原因低怪倒计时, true
				}
			}
			nextMonster = now.Add(怪物识别间隔)
		}
		time.Sleep(100 * time.Millisecond)
	}
	return 离层原因无, 脚本仍应运行(runID)
}

func 层内追怪或单向游走(位置 层位置) {
	配置, ok := 当前层移动配置(位置.层)
	if !ok {
		return
	}
	if 位置.X <= 配置.左边+跳跃移动像素 {
		层内移动方向[位置.层] = 1
		随机向右移动(位置.X, 配置.左边, 配置.右边, 配置.左边, 配置.右边)
		return
	}
	if 位置.X >= 配置.右边-跳跃移动像素 {
		层内移动方向[位置.层] = -1
		随机向左移动(位置.X, 配置.左边, 配置.右边, 配置.左边, 配置.右边)
		return
	}
	方向 := YOLO当前层怪物方向(位置)
	if 方向 != 0 {
		层内移动方向[位置.层] = 方向
	} else {
		方向 = 层内移动方向[位置.层]
	}
	if 方向 == 0 {
		if 移动随机.Intn(2) == 0 {
			方向 = -1
		} else {
			方向 = 1
		}
		层内移动方向[位置.层] = 方向
	}
	if 方向 < 0 {
		随机向左移动(位置.X, 配置.左边, 配置.右边, 配置.左边, 配置.右边)
		return
	}
	随机向右移动(位置.X, 配置.左边, 配置.右边, 配置.左边, 配置.右边)
}

func 打怪动作(位置 层位置) {
	方向 := YOLO当前层怪物方向(位置)
	if 方向 == 0 {
		方向 = 层内移动方向[位置.层]
	}
	if 方向 == 0 {
		if 移动随机.Intn(2) == 0 {
			方向 = -1
		} else {
			方向 = 1
		}
	}
	层内移动方向[位置.层] = 方向

	方向键 := motion.KEYCODE_DPAD_RIGHT
	if 方向 < 0 {
		方向键 = motion.KEYCODE_DPAD_LEFT
	}
	按组合键不空格(方向键, motion.KEYCODE_X, 方向键按下毫秒)
	按空格群攻()
}

func 按空格群攻() {
	for i := 0; i < 打怪空格次数; i++ {
		点按空格()
		time.Sleep(40 * time.Millisecond)
	}
}

func 点按空格() {
	displayID := 当前显示ID()
	motion.KeyActionDown(motion.KEYCODE_SPACE, displayID)
	time.Sleep(50 * time.Millisecond)
	motion.KeyActionUp(motion.KEYCODE_SPACE, displayID)
}

func 按Del键() {
	点按键(motion.KEYCODE_DEL, 当前显示ID())
}

func 按N键() {
	displayID := 当前显示ID()
	motion.KeyActionDown(motion.KEYCODE_N, displayID)
	time.Sleep(随机N键按下时长())
	motion.KeyActionUp(motion.KEYCODE_N, displayID)
}

func 随机N键按下时长() time.Duration {
	span := int64(N键按下最长时长 - N键按下最短时长)
	if span <= 0 {
		return N键按下最短时长
	}
	return N键按下最短时长 + time.Duration(移动随机.Int63n(span+1))
}

func 启动N键守护(runID int64) {
	go N键守护循环(runID)
}

func N键守护循环(runID int64) {
	nextN := time.Now().Add(随机打怪N键间隔())
	for 脚本仍应运行(runID) {
		wait := time.Until(nextN)
		if wait > 200*time.Millisecond {
			time.Sleep(200 * time.Millisecond)
			continue
		}
		if wait > 0 {
			time.Sleep(wait)
			continue
		}
		if !脚本仍应运行(runID) {
			return
		}
		按N键()
		nextN = time.Now().Add(随机打怪N键间隔())
	}
}

func 随机打怪N键间隔() time.Duration {
	span := int64(打怪N键最长间隔 - 打怪N键最短间隔)
	if span <= 0 {
		return 打怪N键最短间隔
	}
	return 打怪N键最短间隔 + time.Duration(移动随机.Int63n(span))
}
