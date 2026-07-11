package main

import (
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Dasongzi1366/AutoGo/motion"
)

var 脚本运行中 atomic.Bool
var 脚本已暂停 atomic.Bool
var 程序退出中 atomic.Bool
var 脚本运行序号 atomic.Int64
var 当前动作内容 atomic.Value

var (
	UI输出锁    sync.Mutex
	UI输出行    []string
	UI输出最大行数 = 80
)

func 启动脚本() {
	启动海盗脚本()
}

func 启动地图1脚本() {
	启动海盗脚本()
}

func 启动海盗脚本() {
	if 程序退出中.Load() {
		输出("程序正在退出，忽略启动")
		return
	}
	if !脚本运行中.CompareAndSwap(false, true) {
		输出("脚本已经在运行")
		return
	}

	脚本已暂停.Store(false)
	输出("海盗 开始识别")
	runID := 脚本运行序号.Add(1)
	启动先四层买卖 := 启动时应立即执行四层买卖()
	启动N键守护(runID)
	go 运行图色循环(runID, 启动先四层买卖)
}

func 启动僵尸3脚本() {
	if 程序退出中.Load() {
		输出("程序正在退出，忽略启动")
		return
	}
	if !脚本运行中.CompareAndSwap(false, true) {
		输出("脚本已经在运行")
		return
	}

	脚本已暂停.Store(false)
	输出("僵尸3 开始识别")
	runID := 脚本运行序号.Add(1)
	启动先三层卖物品 := 启动时应立即执行僵尸3三层卖物品()
	go 运行僵尸3循环(runID, 启动先三层卖物品)
}

func 启动时应立即执行四层买卖() bool {
	位置, ok := 当前层位置()
	if !ok {
		return false
	}
	return 位置.层 == 4
}

func 启动时应立即执行僵尸3三层卖物品() bool {
	位置, ok := 僵尸3当前层位置()
	if !ok {
		return false
	}
	return 位置.层 == 3
}

func 停止脚本() {
	释放所有按键()
	if !脚本运行中.Swap(false) {
		输出("脚本未运行")
		return
	}

	脚本已暂停.Store(false)
	输出("停止脚本信号已发送")
}

func 暂停脚本() {
	释放所有按键()
	if !脚本运行中.Swap(false) {
		输出("脚本未运行")
		return
	}

	脚本已暂停.Store(true)
	输出("暂停脚本信号已发送")
}

func 设置当前动作(text string) {
	text = strings.TrimSpace(text)
	if text == "" {
		return
	}
	当前动作内容.Store(text)
}

func 当前动作文本() string {
	value := 当前动作内容.Load()
	if text, ok := value.(string); ok && text != "" {
		return text
	}
	return "等待动作"
}

func 追加UI输出(text string) {
	text = strings.TrimSpace(text)
	if text == "" {
		return
	}
	UI输出锁.Lock()
	defer UI输出锁.Unlock()

	UI输出行 = append(UI输出行, time.Now().Format("15:04:05  ")+text)
	if len(UI输出行) > UI输出最大行数 {
		UI输出行 = UI输出行[len(UI输出行)-UI输出最大行数:]
	}
}

func 读取UI输出() []string {
	UI输出锁.Lock()
	defer UI输出锁.Unlock()

	lines := make([]string, len(UI输出行))
	copy(lines, UI输出行)
	return lines
}

func 请求退出程序() {
	释放所有按键()
	firstExit := !程序退出中.Swap(true)
	if 脚本运行中.Load() {
		停止脚本()
	}
	脚本已暂停.Store(false)
	if firstExit {
		输出("退出信号已发送")
	}
}

func 释放所有按键() {
	displayID := 当前显示ID()
	keys := []int{
		motion.KEYCODE_DPAD_LEFT,
		motion.KEYCODE_DPAD_RIGHT,
		motion.KEYCODE_DPAD_UP,
		motion.KEYCODE_DPAD_DOWN,
		motion.KEYCODE_X,
		motion.KEYCODE_Z,
		motion.KEYCODE_N,
		motion.KEYCODE_SPACE,
		motion.KEYCODE_DEL,
	}
	for _, key := range keys {
		motion.KeyActionUp(key, displayID)
	}
	输出("键盘", "释放全部")
}

func 等待脚本停止() {
	for 脚本运行中.Load() {
		time.Sleep(200 * time.Millisecond)
	}
}

func 当前脚本状态文本() string {
	if 脚本运行中.Load() {
		return "运行中"
	}
	if 脚本已暂停.Load() {
		return "已暂停"
	}
	return "已停止"
}
