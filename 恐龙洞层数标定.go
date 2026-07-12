package main

import (
	"sync/atomic"
)

var 层数基准 = &FMColor{
	Name:        "层数基准",
	X1:          11,
	Y1:          101,
	X2:          35,
	Y2:          207,
	MainColor:   "3A3D3D-202020",
	OffsetColor: "2,0,35CDDC-202020,6,0,299DC9-202020,0,2,3CCDFD-202020,2,4,64EFF0-202020,8,3,383224-202020,1,7,14C0FF-202020,3,7,4BF2FF-202020,5,5,5BC4E2-202020",
	Sim:         0.75,
	Dir:         0,
}

var 恐龙洞层数标定执行中 atomic.Bool

const 恐龙洞相对层数容差 = 8

var 恐龙洞相对层数基准表 = []struct {
	层  int
	差值 int
}{
	{层: 2, 差值: 46},
	{层: 3, 差值: 25},
	{层: 4, 差值: 0},
	{层: 5, 差值: -13},
}

type 恐龙洞层数相对位置 struct {
	黄点X int
	黄点Y int
	基准X int
	基准Y int
	差值  int
}

func 读取恐龙洞层数相对位置() (恐龙洞层数相对位置, bool) {
	if 引擎 == nil {
		return 恐龙洞层数相对位置{}, false
	}
	yellowOK, yellowX, yellowY := 查找恐龙洞黄点坐标()
	baseOK, baseX, baseY := 引擎.FindFeature(层数基准)
	if !yellowOK || !baseOK {
		return 恐龙洞层数相对位置{黄点X: yellowX, 黄点Y: yellowY, 基准X: baseX, 基准Y: baseY}, false
	}
	return 恐龙洞层数相对位置{
		黄点X: yellowX,
		黄点Y: yellowY,
		基准X: baseX,
		基准Y: baseY,
		差值:  yellowY - baseY,
	}, true
}

func 识别恐龙洞相对层数(diff int) (int, bool) {
	bestLayer := 0
	bestDiff := 恐龙洞相对层数容差 + 1
	for _, config := range 恐龙洞相对层数基准表 {
		distance := absInt(diff - config.差值)
		if distance < bestDiff {
			bestDiff = distance
			bestLayer = config.层
		}
	}
	return bestLayer, bestLayer != 0 && bestDiff <= 恐龙洞相对层数容差
}

func 识别恐龙洞黄点所在层(yellowX, yellowY int) (layer, baseX, baseY, diff int, ok bool) {
	if firstLayer, matched := 识别恐龙洞层数(yellowX, yellowY); matched && firstLayer == 1 {
		return 1, -1, -1, 0, true
	}
	if 引擎 == nil {
		return 0, -1, -1, 0, false
	}
	baseOK, baseX, baseY := 引擎.FindFeature(层数基准)
	if !baseOK {
		return 0, baseX, baseY, 0, false
	}
	diff = yellowY - baseY
	layer, ok = 识别恐龙洞相对层数(diff)
	return layer, baseX, baseY, diff, ok
}

func 执行恐龙洞楼层差值标定(layer int) {
	if layer < 2 || layer > 5 {
		设置恐龙洞输出("楼层差值标定失败：层数%d无效", layer)
		return
	}
	if 脚本运行中.Load() {
		设置恐龙洞输出("%d层差值标定失败：脚本运行中，请先点结束", layer)
		return
	}
	if !恐龙洞层数标定执行中.CompareAndSwap(false, true) {
		设置恐龙洞输出("楼层差值标定执行中，请稍等")
		return
	}
	go func() {
		defer 恐龙洞层数标定执行中.Store(false)
		position, ok := 读取恐龙洞层数相对位置()
		if !ok {
			设置恐龙洞输出("%d层差值标定失败：黄点=(%d,%d) 基准=(%d,%d)", layer, position.黄点X, position.黄点Y, position.基准X, position.基准Y)
			return
		}
		标记恐龙洞找到的黄点(position.黄点X, position.黄点Y)
		设置恐龙洞输出("%d层标定：黄点=(%d,%d) 基准=(%d,%d) 差值=黄点Y-基准Y=%d 绝对值=%d", layer, position.黄点X, position.黄点Y, position.基准X, position.基准Y, position.差值, absInt(position.差值))
	}()
}
