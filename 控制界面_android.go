//go:build android

package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/Dasongzi1366/AutoGo/device"
	"github.com/Dasongzi1366/AutoGo/imgui"
	"github.com/Dasongzi1366/AutoGo/utils"
)

var 控制窗口打开 = true

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
				window.SetFontWindowScale(scale * 0.75)
			}

			buttonWidth := windowWidth - 40*scale
			buttonGap := 6 * scale
			buttonHeight := 34 * scale
			buttonSize := imgui.Vec2{X: (buttonWidth - buttonGap) / 2, Y: buttonHeight}
			if imgui.ButtonV("地图1开始", buttonSize) {
				启动地图1脚本()
			}
			imgui.SameLineV(0, buttonGap)
			if imgui.ButtonV("结束", buttonSize) {
				shouldExit = true
			}
			imgui.Spacing()
			if imgui.ButtonV("下一步", imgui.Vec2{X: buttonWidth, Y: buttonHeight}) {
				执行测试卖物品下一步()
			}
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
	return 112 * scale
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
			switch utils.Alert("脚本控制", content, "暂停", "结束") {
			case 0:
				暂停脚本()
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
