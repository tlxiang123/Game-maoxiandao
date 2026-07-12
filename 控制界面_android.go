//go:build android

package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Dasongzi1366/AutoGo/device"
	"github.com/Dasongzi1366/AutoGo/ime"
	"github.com/Dasongzi1366/AutoGo/imgui"
	"github.com/Dasongzi1366/AutoGo/utils"
)

var 控制窗口打开 = true

const (
	控制日志僵尸3 = iota
	控制日志海盗
	控制日志恐龙
)

var 当前控制日志 = 控制日志恐龙

func 运行控制界面() {
	if err := imgui.Init(); err != nil {
		输出("控制界面初始化失败", err)
		return
	}

	控制窗口打开 = true
	uiDone := make(chan struct{})
	var closeOnce sync.Once
	firstFrame := make(chan struct{})
	var firstFrameOnce sync.Once
	uiClosed := false

	closeUI := func() {
		closeOnce.Do(func() {
			uiClosed = true
			控制窗口打开 = false
			imgui.Close()
			close(uiDone)
		})
	}
	defer func() {
		if !uiClosed {
			imgui.Close()
		}
	}()

	imgui.Run(func() {
		firstFrameOnce.Do(func() {
			close(firstFrame)
		})

		if 程序退出中.Load() || !控制窗口打开 {
			closeUI()
			return
		}

		width, height, _, _ := device.GetDisplayInfo(屏幕ID)
		if width <= 0 {
			width = 1280
		}
		if height <= 0 {
			height = 720
		}
		scale := 控制界面缩放(width, height)
		windowWidth := 控制窗口宽度(width, height, scale)
		windowHeight := 控制窗口高度(scale)
		windowX := (float32(width) - windowWidth) / 2
		windowY := 24 * scale
		minY := 12 * scale
		if windowY < minY {
			windowY = minY
		}

		imgui.SetNextWindowSizeV(imgui.Vec2{X: windowWidth, Y: windowHeight}, imgui.CondOnce)
		imgui.SetNextWindowPosV(imgui.Vec2{X: windowX, Y: windowY}, imgui.CondOnce, imgui.Vec2{X: 0, Y: 0})

		shouldExit := false
		flags := imgui.WindowFlagsNoResize | imgui.WindowFlagsNoSavedSettings
		if imgui.BeginV("脚本控制", &控制窗口打开, flags) {
			if window := imgui.InternalCurrentWindow(); window != nil {
				window.SetFontWindowScale(scale * 0.375)
			}

			buttonWidth := windowWidth - 40*scale
			buttonGap := 6 * scale
			buttonHeight := 24 * scale
			buttonSize := imgui.Vec2{X: (buttonWidth - buttonGap*3) / 4, Y: buttonHeight}
			if imgui.ButtonV("僵尸3开始", buttonSize) {
				当前控制日志 = 控制日志僵尸3
				启动僵尸3脚本()
			}
			imgui.SameLineV(0, buttonGap)
			if imgui.ButtonV("海盗开始", buttonSize) {
				当前控制日志 = 控制日志海盗
				启动海盗脚本()
			}
			imgui.SameLineV(0, buttonGap)
			if imgui.ButtonV("恐龙开始", buttonSize) {
				当前控制日志 = 控制日志恐龙
				启动恐龙洞脚本()
			}
			imgui.SameLineV(0, buttonGap)
			if imgui.ButtonV("结束", buttonSize) {
				shouldExit = true
			}
			imgui.Spacing()
			imgui.Checkbox("卖杂物", &僵尸3卖杂物)
			walkButtonWidth := 78 * scale
			pauseButtonWidth := 64 * scale
			imgui.SameLineV(buttonWidth-walkButtonWidth-pauseButtonWidth-buttonGap, 0)
			if imgui.ButtonV("左右刷怪", imgui.Vec2{X: walkButtonWidth, Y: buttonHeight}) {
				当前控制日志 = 控制日志僵尸3
				启动僵尸3左右刷怪()
			}
			imgui.SameLineV(0, buttonGap)
			pauseLabel := "暂停"
			if 脚本已暂停.Load() {
				pauseLabel = "恢复"
			}
			if imgui.ButtonV(pauseLabel, imgui.Vec2{X: pauseButtonWidth, Y: buttonHeight}) {
				切换脚本暂停()
			}
			imgui.Text(fmt.Sprintf("卖物品时间：%d-%d 分钟", 恐龙洞卖物品最短分钟, 恐龙洞卖物品最长分钟))
			sliderWidth := (buttonWidth-buttonGap)/2 - 42*scale
			imgui.SetNextItemWidth(sliderWidth)
			if imgui.SliderInt("最短分钟", &恐龙洞卖物品最短分钟, 1, 180) {
				if 恐龙洞卖物品最短分钟 > 恐龙洞卖物品最长分钟 {
					恐龙洞卖物品最长分钟 = 恐龙洞卖物品最短分钟
				}
				设置恐龙洞卖物品分钟范围(恐龙洞卖物品最短分钟, 恐龙洞卖物品最长分钟)
			}
			imgui.SameLineV(0, buttonGap)
			imgui.SetNextItemWidth(sliderWidth)
			if imgui.SliderInt("最长分钟", &恐龙洞卖物品最长分钟, 1, 180) {
				if 恐龙洞卖物品最长分钟 < 恐龙洞卖物品最短分钟 {
					恐龙洞卖物品最短分钟 = 恐龙洞卖物品最长分钟
				}
				设置恐龙洞卖物品分钟范围(恐龙洞卖物品最短分钟, 恐龙洞卖物品最长分钟)
			}
			outputActionButtonWidth := 78 * scale
			imgui.AlignTextToFramePadding()
			imgui.Text(当前控制日志标题())
			imgui.SameLineV(buttonWidth-outputActionButtonWidth*2-buttonGap, 0)
			if imgui.ButtonV("卖东西单步", imgui.Vec2{X: outputActionButtonWidth, Y: buttonHeight}) {
				执行恐龙洞卖物品测试下一步()
			}
			imgui.SameLineV(0, buttonGap)
			if imgui.ButtonV("测试黄点", imgui.Vec2{X: outputActionButtonWidth, Y: buttonHeight}) {
				执行恐龙洞黄点测试()
			}
			imgui.Separator()
			if imgui.ButtonV("复制全部日志", imgui.Vec2{X: buttonWidth, Y: buttonHeight}) {
				if ime.SetClipText(读取当前控制全部日志()) {
					utils.Toast("全部日志已复制，可以到文本文件中粘贴", 260, 1237, 1800)
				} else {
					utils.Toast("复制失败：无法写入系统剪贴板", 260, 1237, 2200)
				}
			}
			outputFlags := imgui.WindowFlags(0)
			if imgui.BeginChildStrV("dinosaur_cave_output", imgui.Vec2{X: buttonWidth, Y: 88 * scale}, imgui.ChildFlagsBorders, outputFlags) {
				imgui.SetWindowFontScale(scale * 0.34)
				imgui.PushTextWrapPosV(0)
				imgui.TextUnformatted(读取当前控制日志文本())
				imgui.PopTextWrapPos()
				if 消耗当前控制日志滚动请求() {
					imgui.SetScrollHereYV(1)
				}
			}
			imgui.EndChild()
		}
		imgui.End()
		renderDebugRedBoxes()

		if shouldExit || !控制窗口打开 {
			请求退出程序()
			closeUI()
		}
	})

	if !uiClosed {
		select {
		case <-uiDone:
		case <-firstFrame:
			<-uiDone
		case <-time.After(1500 * time.Millisecond):
			输出("imgui控制界面未进入渲染循环，改用系统弹窗控制")
			closeUI()
			运行弹窗控制界面()
		}
	}
}

func 当前控制日志标题() string {
	switch 当前控制日志 {
	case 控制日志僵尸3:
		return "僵尸3输出"
	case 控制日志海盗:
		return "海盗输出"
	default:
		return "恐龙洞输出"
	}
}

func 读取当前控制日志文本() string {
	switch 当前控制日志 {
	case 控制日志僵尸3:
		return 读取僵尸3输出文本()
	case 控制日志海盗:
		lines := 读取UI输出()
		if len(lines) > 4 {
			lines = lines[len(lines)-4:]
		}
		return strings.Join(lines, "\n")
	default:
		return 读取恐龙洞输出文本()
	}
}

func 读取当前控制全部日志() string {
	switch 当前控制日志 {
	case 控制日志僵尸3:
		return 读取僵尸3全部输出文本()
	case 控制日志海盗:
		return strings.Join(读取UI输出(), "\n")
	default:
		return 读取恐龙洞全部输出文本()
	}
}

func 消耗当前控制日志滚动请求() bool {
	switch 当前控制日志 {
	case 控制日志僵尸3:
		return 消耗僵尸3输出滚动请求()
	case 控制日志海盗:
		return false
	default:
		return 消耗恐龙洞输出滚动请求()
	}
}

func 控制界面缩放(screenWidth, screenHeight int) float32 {
	shortSide := screenWidth
	if screenHeight < shortSide {
		shortSide = screenHeight
	}
	scale := float32(shortSide) / 750
	if scale < 1 {
		return 1
	}
	if scale > 2 {
		return 2
	}
	return scale
}

func 控制窗口高度(scale float32) float32 {
	return 300 * scale
}

func 控制窗口宽度(screenWidth, screenHeight int, scale float32) float32 {
	ratio := float32(0.52)
	if screenWidth > screenHeight {
		ratio = 0.30
	}

	width := float32(screenWidth) * ratio
	minWidth := 300 * scale
	maxWidth := 460 * scale
	screenLimit := float32(screenWidth) * 0.62

	if width < minWidth {
		width = minWidth
	}
	if width > maxWidth {
		width = maxWidth
	}
	if width > screenLimit {
		width = screenLimit
	}
	return width
}

func 运行弹窗控制界面() {
	for !程序退出中.Load() {
		content := fmt.Sprintf("状态：%s\n当前动作：%s", 当前脚本状态文本(), 当前动作文本())
		if 脚本运行中.Load() {
			primary := "暂停"
			if 脚本已暂停.Load() {
				primary = "恢复"
			}
			switch utils.Alert("脚本控制", content, primary, "结束") {
			case 0:
				切换脚本暂停()
			case 1:
				请求退出程序()
			}
			continue
		}

		switch utils.Alert("脚本控制", content, "开始", "结束") {
		case 0:
			启动脚本()
		case 1:
			请求退出程序()
		}
	}
}
