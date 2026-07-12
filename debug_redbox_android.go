//go:build android

package main

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/Dasongzi1366/AutoGo/device"
	"github.com/Dasongzi1366/AutoGo/imgui"
)

type debugRedBox struct {
	x1    int
	y1    int
	x2    int
	y2    int
	until time.Time
}

type debugTargetPoint struct {
	x           int
	y           int
	singlePixel bool
	until       time.Time
}

var (
	debugRedBoxMu       sync.Mutex
	debugRedBoxItems    []debugRedBox
	debugTargetPoints   []debugTargetPoint
	debugRedBoxSuppress atomic.Int32
)

const (
	debugDrawingEnabled = false
	debugRedBoxDuration = 180 * time.Millisecond
	debugRedBoxHalfSize = 28
	debugRedBoxMaxItems = 80
	debugTargetDuration = 8 * time.Second
	debugTargetMaxItems = 8
)

func addDebugRedBox(x1, y1, x2, y2 int) {
	if !debugDrawingEnabled || debugRedBoxSuppress.Load() > 0 {
		return
	}
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}
	if x2-x1 < 6 {
		x1 -= 3
		x2 += 3
	}
	if y2-y1 < 6 {
		y1 -= 3
		y2 += 3
	}

	debugRedBoxMu.Lock()
	defer debugRedBoxMu.Unlock()

	if len(debugRedBoxItems) >= debugRedBoxMaxItems {
		copy(debugRedBoxItems, debugRedBoxItems[len(debugRedBoxItems)-debugRedBoxMaxItems+1:])
		debugRedBoxItems = debugRedBoxItems[:debugRedBoxMaxItems-1]
	}
	debugRedBoxItems = append(debugRedBoxItems, debugRedBox{
		x1:    x1,
		y1:    y1,
		x2:    x2,
		y2:    y2,
		until: time.Now().Add(debugRedBoxDuration),
	})
}

func addDebugPointBox(x, y int) {
	addDebugRedBox(
		x-debugRedBoxHalfSize,
		y-debugRedBoxHalfSize,
		x+debugRedBoxHalfSize,
		y+debugRedBoxHalfSize,
	)
}

func addDebugTargetPoint(x, y int) {
	addDebugMarkerPoint(x, y, false)
}

func addDebugSinglePixel(x, y int) {
	addDebugMarkerPoint(x, y, true)
}

func addDebugMarkerPoint(x, y int, singlePixel bool) {
	if !debugDrawingEnabled {
		return
	}
	debugRedBoxMu.Lock()
	defer debugRedBoxMu.Unlock()
	if len(debugTargetPoints) >= debugTargetMaxItems {
		copy(debugTargetPoints, debugTargetPoints[len(debugTargetPoints)-debugTargetMaxItems+1:])
		debugTargetPoints = debugTargetPoints[:debugTargetMaxItems-1]
	}
	debugTargetPoints = append(debugTargetPoints, debugTargetPoint{x: x, y: y, singlePixel: singlePixel, until: time.Now().Add(debugTargetDuration)})
}

func 暂停调试红框() {
	debugRedBoxSuppress.Add(1)
}

func 恢复调试红框() {
	if debugRedBoxSuppress.Add(-1) < 0 {
		debugRedBoxSuppress.Store(0)
	}
}

func renderDebugRedBoxes() {
	if !debugDrawingEnabled {
		return
	}
	now := time.Now()

	debugRedBoxMu.Lock()
	active := debugRedBoxItems[:0]
	drawItems := make([]debugRedBox, 0, len(debugRedBoxItems))
	for _, item := range debugRedBoxItems {
		if now.Before(item.until) {
			active = append(active, item)
			drawItems = append(drawItems, item)
		}
	}
	debugRedBoxItems = active
	activeTargets := debugTargetPoints[:0]
	drawTargets := make([]debugTargetPoint, 0, len(debugTargetPoints))
	for _, item := range debugTargetPoints {
		if now.Before(item.until) {
			activeTargets = append(activeTargets, item)
			drawTargets = append(drawTargets, item)
		}
	}
	debugTargetPoints = activeTargets
	debugRedBoxMu.Unlock()

	if len(drawItems) == 0 && len(drawTargets) == 0 {
		return
	}

	red := imgui.NewColor(1, 0, 0, 1).Pack()
	fill := imgui.NewColor(1, 0, 0, 0.10).Pack()
	flags := imgui.DrawFlags(0)

	width, height, _, _ := device.GetDisplayInfo(屏幕ID)
	if width <= 0 || height <= 0 {
		width, height = debugRedBoxBounds(drawItems)
	}

	imgui.SetNextWindowPosV(imgui.Vec2{X: 0, Y: 0}, imgui.CondAlways, imgui.Vec2{X: 0, Y: 0})
	imgui.SetNextWindowSizeV(imgui.Vec2{X: float32(width), Y: float32(height)}, imgui.CondAlways)
	windowFlags := imgui.WindowFlagsNoDecoration |
		imgui.WindowFlagsNoInputs |
		imgui.WindowFlagsNoSavedSettings |
		imgui.WindowFlagsNoBackground |
		imgui.WindowFlagsNoMove |
		imgui.WindowFlagsNoBringToFrontOnFocus
	if imgui.BeginV("debug_red_box_overlay", nil, windowFlags) {
		if window := imgui.InternalCurrentWindow(); window != nil {
			drawList := window.DrawList()
			if drawList != nil {
				drawList.PushClipRectFullScreen()
				for _, item := range drawItems {
					p1 := imgui.Vec2{X: float32(item.x1), Y: float32(item.y1)}
					p2 := imgui.Vec2{X: float32(item.x2), Y: float32(item.y2)}
					drawList.AddRectFilledV(p1, p2, fill, 0, flags)
					drawList.AddRectV(p1, p2, red, 0, flags, 4)
				}
				for _, item := range drawTargets {
					center := imgui.Vec2{X: float32(item.x), Y: float32(item.y)}
					if item.singlePixel {
						drawList.AddRectFilledV(center, imgui.Vec2{X: center.X + 1, Y: center.Y + 1}, red, 0, flags)
						continue
					}
					drawList.AddCircleFilledV(center, 5, red, 16)
					drawList.AddLineV(imgui.Vec2{X: center.X - 12, Y: center.Y}, imgui.Vec2{X: center.X + 12, Y: center.Y}, red, 3)
					drawList.AddLineV(imgui.Vec2{X: center.X, Y: center.Y - 12}, imgui.Vec2{X: center.X, Y: center.Y + 12}, red, 3)
				}
				drawList.PopClipRect()
			}
		}
	}
	imgui.End()
}

func debugRedBoxBounds(items []debugRedBox) (int, int) {
	width := 1
	height := 1
	for _, item := range items {
		if item.x2 > width {
			width = item.x2
		}
		if item.y2 > height {
			height = item.y2
		}
	}
	return width, height
}

func (z *Zg) markScreenRect(x1, y1, x2, y2 int) {
	addDebugRedBox(x1, y1, x2, y2)
}

func (z *Zg) markScreenPoint(x, y int) {
	addDebugPointBox(x, y)
}

func (z *Zg) markClickPoint(x, y int) {
	if z == nil {
		return
	}
	addDebugPointBox(z.scaleX(x+z.offsetX), z.scaleY(y+z.offsetY))
}
