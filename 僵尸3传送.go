package main

import (
	"sync"
	"sync/atomic"
	"time"
)

var 传送被传出来了 = &FMColor{Name: "被传出来了", X1: 42, Y1: 30, X2: 174, Y2: 83, MainColor: "FFFFFF-000000", OffsetColor: "5,2,4B565C-000000,16,0,FFFFFF-000000,4,7,303B40-000000,5,3,4B565C-000000,20,5,2E383D-000000,4,8,3F4D54-000000,8,8,C6C6C6-000000,20,13,3C4A50-000000", Sim: 0.90, Dir: 0}
var 传送商店按钮 = &FMColor{Name: "商店按钮", X1: 887, Y1: 635, X2: 920, Y2: 666, MainColor: "0077CC-000000", OffsetColor: "1,0,AADDEE-000000,9,0,0066CC-000000,-4,4,0077CC-000000,6,4,AADDEE-000000,12,4,AADDEE-000000,-2,11,EEAA33-000000,3,9,995500-000000,9,13,DD8833-000000", Sim: 0.90, Dir: 0}
var 传送第五页面灰色 = &FMColor{Name: "第五页面灰色", X1: 1229, Y1: 77, X2: 1276, Y2: 185, MainColor: "878787-000000", OffsetColor: "6,0,8E8E8E-000000,22,0,737373-000000,0,7,FBFBFB-000000,16,4,F7F7F7-000000,22,1,E9E9E9-000000,5,11,E9E9E9-000000,16,8,7F7F7F-000000,22,11,E9E9E9-000000", Sim: 0.90, Dir: 0}
var 传送第五页面灰色备用 = &FMColor{Name: "第五页面灰色", X1: 1002, Y1: 179, X2: 1126, Y2: 221, MainColor: "878787-000000", OffsetColor: "6,0,8E8E8E-000000,22,0,737373-000000,0,7,FBFBFB-000000,16,4,F7F7F7-000000,22,1,E9E9E9-000000,5,11,E9E9E9-000000,16,8,7F7F7F-000000,22,11,E9E9E9-000000", Sim: 0.90, Dir: 0}
var 传送第五页粉色 = &FMColor{Name: "第五页粉色", X1: 1229, Y1: 77, X2: 1276, Y2: 185, MainColor: "FDEFF2-000000", OffsetColor: "11,0,E6DDDF-000000,12,0,FAF0F3-000000,-5,3,FFF8FA-000000,1,3,FFFBFC-000000,17,6,8F3D52-000000,0,9,FFFFFF-000000,1,9,F1DCE1-000000,17,7,8F3D52-000000", Sim: 0.80, Dir: 0}
var 传送第五页粉色备用 = &FMColor{Name: "第五页粉色", X1: 1002, Y1: 179, X2: 1126, Y2: 221, MainColor: "FDEFF2-000000", OffsetColor: "11,0,E6DDDF-000000,12,0,FAF0F3-000000,-5,3,FFF8FA-000000,1,3,FFFBFC-000000,17,6,8F3D52-000000,0,9,FFFFFF-000000,1,9,F1DCE1-000000,17,7,8F3D52-000000", Sim: 0.80, Dir: 0}
var 传送已经点开传送石 = &FMColor{Name: "已经点开传送石", X1: 528, Y1: 353, X2: 650, Y2: 490, MainColor: "FF9911-000000", OffsetColor: "1,0,445577-000000,16,-3,4477BB-000000,0,1,FF9900-000000,10,4,D0DDEA-000000,11,4,FFFFFF-000000,0,14,F24800-000000,10,8,D0D9E1-000000,11,8,EEEEEE-000000", Sim: 0.90, Dir: 0}
var 传送下一页 = &FMColor{Name: "下一页", X1: 714, Y1: 357, X2: 753, Y2: 412, MainColor: "2277AA-000000", OffsetColor: "6,0,226699-000000,10,0,226699-000000,0,5,226699-000000,1,5,2A6F9D-000000,10,5,135689-000000,0,8,226699-000000,6,8,174A7D-000000,10,8,114477-000000", Sim: 0.90, Dir: 0}
var 传送需要传送的地图 = &FMColor{Name: "需要传送的地图", X1: 602, Y1: 256, X2: 666, Y2: 396, MainColor: "232623-000000", OffsetColor: "1,3,414141-000000,12,0,232323-000000,-6,7,000000-000000,4,10,000000-000000,12,7,EEEEFB-000000,0,14,000000-000000,4,14,000000-000000,9,11,414142-000000", Sim: 0.90, Dir: 0}
var 传送选中僵尸3地图 = &FMColor{Name: "选中僵尸3地图", X1: 602, Y1: 256, X2: 666, Y2: 396, MainColor: "FFFFFF-000000", OffsetColor: "13,-2,C3CDD5-000000,18,-4,4287BC-000000,4,5,FFFFFF-000000,13,3,84B1D4-000000,18,1,4287BC-000000,4,8,FFFFFF-000000,13,11,BDD5E7-000000,18,6,4287BC-000000", Sim: 0.90, Dir: 0}
var 传送点击传送 = &FMColor{Name: "点击传送", X1: 687, Y1: 517, X2: 744, Y2: 597, MainColor: "FFDDDD-000000", OffsetColor: "18,0,FFDDDD-000000,19,0,FFDDDD-000000,0,3,AA2255-000000,12,7,FFDDEE-000000,19,7,FFDDEE-000000,0,10,B25D87-000000,18,8,883355-000000,19,8,883355-000000", Sim: 0.90, Dir: 0}
var 传送点击确认 = &FMColor{Name: "点击确认", X1: 692, Y1: 403, X2: 755, Y2: 488, MainColor: "AA5511-000000", OffsetColor: "15,0,FFDDCC-000000,16,0,FFDDCC-000000,5,8,FFDDCC-000000,6,3,BB6622-000000,21,3,BB6622-000000,0,11,FFEECC-000000,6,11,FFEECC-000000,16,11,CC7733-000000", Sim: 0.90, Dir: 0}

const (
	僵尸3回图巡检间隔    = 10 * time.Second
	僵尸3传送等待超时    = 4 * time.Second
	僵尸3传送翻页最大次数  = 20
	僵尸3换线后巡检暂停   = 12 * time.Second
	僵尸3传送石固定双击X  = 1077
	僵尸3传送石固定双击Y  = 160
	僵尸3传送查色后最小等待 = 220 * time.Millisecond
	僵尸3传送查色后最大等待 = 520 * time.Millisecond
	僵尸3传送操作后最小等待 = 420 * time.Millisecond
	僵尸3传送操作后最大等待 = 880 * time.Millisecond
)

type 僵尸3传送特征候选 struct {
	标签 string
	特征 any
}

var 僵尸3传送第五页灰色候选 = []僵尸3传送特征候选{
	{标签: "第五页面灰色", 特征: 传送第五页面灰色},
	{标签: "第五页面灰色备用", 特征: 传送第五页面灰色备用},
}

var 僵尸3传送第五页粉色候选 = []僵尸3传送特征候选{
	{标签: "第五页粉色", 特征: 传送第五页粉色},
	{标签: "第五页粉色备用", 特征: 传送第五页粉色备用},
}

var 僵尸3传送石固定双击点列表 = []struct {
	标签 string
	X  int
	Y  int
}{
	{标签: "主点", X: 僵尸3传送石固定双击X, Y: 僵尸3传送石固定双击Y},
	{标签: "备用点", X: 881, Y: 242},
}

var (
	僵尸3传送锁      sync.Mutex
	僵尸3传送流程运行中  atomic.Bool
	僵尸3地图巡检暂停到  atomic.Int64
	僵尸3传送测试锁    sync.Mutex
	僵尸3传送测试已启动  bool
	僵尸3传送测试当前逻辑 int
	僵尸3传送下一步执行中 atomic.Bool
	僵尸3左进三层执行中  atomic.Bool
)

func 启动僵尸3地图巡检(runID int64) {
	go func() {
		僵尸3地图巡检一次(runID, "启动检查")
		ticker := time.NewTicker(僵尸3回图巡检间隔)
		defer ticker.Stop()
		for 脚本仍应运行(runID) {
			<-ticker.C
			if !脚本仍应运行(runID) {
				return
			}
			僵尸3地图巡检一次(runID, "10秒检查")
		}
	}()
}

func 僵尸3地图巡检一次(runID int64, reason string) {
	if 引擎 == nil || !脚本仍应运行(runID) {
		return
	}
	if 剩余 := 僵尸3地图巡检暂停剩余时间(); 剩余 > 0 {
		设置僵尸3层输出("%s：换线后巡检暂停中，剩余%d秒", reason, int(剩余.Seconds())+1)
		return
	}
	if 僵尸3当前在地图内() {
		if reason == "启动检查" {
			设置僵尸3层输出("%s：当前在僵尸3地图", reason)
		}
		return
	}
	设置僵尸3层输出("%s：未检测到僵尸3地图，开始传送回图", reason)
	僵尸3执行传送回图流程(func() bool { return 脚本仍应运行(runID) }, reason)
}

func 僵尸3当前在地图内() bool {
	if 引擎 == nil {
		return false
	}
	ok, _, _ := 引擎.FindFeature(僵尸3地图)
	return ok
}

func 暂停僵尸3地图巡检(duration time.Duration, reason string) {
	if duration <= 0 {
		return
	}
	deadline := time.Now().Add(duration).UnixNano()
	for {
		old := 僵尸3地图巡检暂停到.Load()
		if old >= deadline {
			return
		}
		if 僵尸3地图巡检暂停到.CompareAndSwap(old, deadline) {
			设置僵尸3层输出("%s：暂停地图巡检%d秒", reason, int(duration.Seconds()))
			return
		}
	}
}

func 僵尸3地图巡检暂停剩余时间() time.Duration {
	until := 僵尸3地图巡检暂停到.Load()
	if until <= 0 {
		return 0
	}
	剩余 := time.Until(time.Unix(0, until))
	if 剩余 < 0 {
		return 0
	}
	return 剩余
}

func 僵尸3执行传送回图流程(shouldContinue func() bool, reason string) bool {
	if shouldContinue == nil {
		shouldContinue = func() bool { return true }
	}
	if !僵尸3传送锁.TryLock() {
		设置僵尸3层输出("传送回图：已有流程运行中，跳过%s", reason)
		return false
	}
	defer 僵尸3传送锁.Unlock()
	僵尸3传送流程运行中.Store(true)
	defer 僵尸3传送流程运行中.Store(false)

	if 僵尸3当前在地图内() {
		设置僵尸3层输出("传送回图：已经在僵尸3地图")
		return true
	}
	for i := 0; i < 僵尸3传送逻辑数量(); i++ {
		if !shouldContinue() {
			设置僵尸3层输出("传送回图：流程中断")
			return false
		}
		if ok := 僵尸3执行传送逻辑(i, shouldContinue); !ok {
			设置僵尸3层输出("传送回图：第%d个逻辑失败", i+1)
			return false
		}
		if 僵尸3当前在地图内() {
			设置僵尸3层输出("传送回图：已回到僵尸3地图")
			return true
		}
	}
	if ok, _, _ := 僵尸3等待传送特征(僵尸3地图, 15*time.Second, 500*time.Millisecond, shouldContinue); ok {
		设置僵尸3层输出("传送回图：确认已回到僵尸3地图")
		return true
	}
	设置僵尸3层输出("传送回图：传送后仍未检测到僵尸3地图")
	return false
}

func 执行僵尸3传送测试下一步() {
	if 程序退出中.Load() {
		设置僵尸3层输出("传送下一步失败：程序正在退出")
		return
	}
	if !僵尸3传送下一步执行中.CompareAndSwap(false, true) {
		设置僵尸3层输出("传送下一步执行中，请稍等")
		return
	}
	go func() {
		defer 僵尸3传送下一步执行中.Store(false)
		if !僵尸3传送锁.TryLock() {
			设置僵尸3层输出("传送下一步失败：已有传送流程运行中")
			return
		}
		defer 僵尸3传送锁.Unlock()
		僵尸3传送流程运行中.Store(true)
		defer 僵尸3传送流程运行中.Store(false)

		if !僵尸3传送测试流程运行中() {
			启动僵尸3传送测试流程()
		}
		index, ok := 取僵尸3传送测试当前逻辑()
		if !ok {
			return
		}
		设置僵尸3层输出("传送下一步：第%d个逻辑", index+1)
		if 僵尸3执行传送逻辑(index, func() bool { return !程序退出中.Load() }) {
			设置僵尸3传送测试下一逻辑(index + 1)
		} else {
			设置僵尸3层输出("传送下一步失败：第%d个逻辑", index+1)
		}
	}()
}

func 启动僵尸3传送测试流程() {
	僵尸3传送测试锁.Lock()
	僵尸3传送测试已启动 = true
	僵尸3传送测试当前逻辑 = 0
	僵尸3传送测试锁.Unlock()
	设置僵尸3层输出("僵尸3传送测试开始")
}

func 僵尸3传送测试流程运行中() bool {
	僵尸3传送测试锁.Lock()
	defer 僵尸3传送测试锁.Unlock()
	return 僵尸3传送测试已启动
}

func 取僵尸3传送测试当前逻辑() (int, bool) {
	僵尸3传送测试锁.Lock()
	defer 僵尸3传送测试锁.Unlock()
	if !僵尸3传送测试已启动 {
		return 0, false
	}
	if 僵尸3传送测试当前逻辑 >= 僵尸3传送逻辑数量() {
		僵尸3传送测试已启动 = false
		设置僵尸3层输出("僵尸3传送测试结束")
		return 0, false
	}
	return 僵尸3传送测试当前逻辑, true
}

func 设置僵尸3传送测试下一逻辑(next int) {
	僵尸3传送测试锁.Lock()
	僵尸3传送测试当前逻辑 = next
	done := 僵尸3传送测试当前逻辑 >= 僵尸3传送逻辑数量()
	if done {
		僵尸3传送测试已启动 = false
	}
	僵尸3传送测试锁.Unlock()
	if done {
		设置僵尸3层输出("僵尸3传送测试结束")
	}
}

func 僵尸3传送逻辑数量() int {
	return 5
}

func 僵尸3执行传送逻辑(index int, shouldContinue func() bool) bool {
	switch index {
	case 0:
		return 僵尸3传送逻辑1(shouldContinue)
	case 1:
		return 僵尸3传送逻辑2(shouldContinue)
	case 2:
		return 僵尸3传送逻辑3(shouldContinue)
	case 3:
		return 僵尸3传送逻辑4(shouldContinue)
	case 4:
		return 僵尸3传送逻辑5(shouldContinue)
	default:
		return false
	}
}

func 僵尸3传送逻辑1(shouldContinue func() bool) bool {
	设置僵尸3层输出("传送第1个逻辑：检查被传出来了，打开第五页")
	if ok, x, y := 僵尸3查找传送特征(传送被传出来了); ok {
		设置僵尸3层输出("传送第1步：找到被传出来了 x=%d y=%d，不点击", x, y)
	} else {
		设置僵尸3层输出("传送第1步：未找到被传出来了，检查是否已在后续界面")
	}
	if ok, x, y, candidate := 僵尸3查找任一传送特征(僵尸3传送第五页粉色候选...); ok {
		设置僵尸3层输出("传送第1步：已找到%s x=%d y=%d，不点击商店按钮", candidate.标签, x, y)
		return true
	}
	if ok, x, y, candidate := 僵尸3查找任一传送特征(僵尸3传送第五页灰色候选...); ok {
		设置僵尸3层输出("传送第1步：找到%s x=%d y=%d，点击", candidate.标签, x, y)
		僵尸3传送点击(x, y)
		return true
	}
	if ok, x, y := 僵尸3等待传送特征(传送商店按钮, 僵尸3传送等待超时, 200*time.Millisecond, shouldContinue); ok {
		设置僵尸3层输出("传送第1步：找到商店按钮 x=%d y=%d，点击", x, y)
		僵尸3传送点击(x, y)
		return true
	}
	设置僵尸3层输出("传送第1步失败：未找到第五页粉色、第五页面灰色或商店按钮")
	return false
}

func 僵尸3传送逻辑2(shouldContinue func() bool) bool {
	设置僵尸3层输出("传送第2个逻辑：打开第五页并双击传送石")
	if ok, x, y, candidate := 僵尸3查找任一传送特征(僵尸3传送第五页粉色候选...); ok {
		设置僵尸3层输出("传送第2步：第五页已打开 %s x=%d y=%d", candidate.标签, x, y)
	} else if ok, x, y, candidate := 僵尸3查找任一传送特征(僵尸3传送第五页灰色候选...); ok {
		设置僵尸3层输出("传送第2步：点击%s x=%d y=%d", candidate.标签, x, y)
		僵尸3传送点击(x, y)
		if ok, x, y, opened := 僵尸3等待任一传送特征(僵尸3传送第五页粉色候选, 僵尸3传送等待超时, 200*time.Millisecond, shouldContinue); ok {
			设置僵尸3层输出("传送第2步：第五页已切到%s x=%d y=%d", opened.标签, x, y)
		} else {
			设置僵尸3层输出("传送第2步失败：点击灰色后未切到第五页粉色")
			return false
		}
	} else {
		设置僵尸3层输出("传送第2步失败：未找到第五页面灰色或第五页粉色")
		return false
	}

	for _, point := range 僵尸3传送石固定双击点列表 {
		设置僵尸3层输出("传送第2步：固定双击传送石%s x=%d y=%d", point.标签, point.X, point.Y)
		僵尸3传送双击(point.X, point.Y)
		if ok, x, y := 僵尸3等待传送特征(传送已经点开传送石, 僵尸3传送等待超时, 200*time.Millisecond, shouldContinue); ok {
			设置僵尸3层输出("传送第2步：传送石已点开 x=%d y=%d", x, y)
			return true
		}
	}
	设置僵尸3层输出("传送第2步失败：所有固定双击点都未打开传送石")
	return false
}

func 僵尸3传送逻辑3(shouldContinue func() bool) bool {
	设置僵尸3层输出("传送第3个逻辑：点击传送石界面并点下一页")
	if ok, x, y := 僵尸3等待传送特征(传送已经点开传送石, 僵尸3传送等待超时, 200*time.Millisecond, shouldContinue); ok {
		设置僵尸3层输出("传送第3步：点击已经点开传送石 x=%d y=%d", x, y)
		僵尸3传送点击(x, y)
	} else {
		设置僵尸3层输出("传送第3步失败：未找到已经点开传送石")
		return false
	}
	if ok, x, y := 僵尸3等待传送特征(传送下一页, 僵尸3传送等待超时, 200*time.Millisecond, shouldContinue); ok {
		设置僵尸3层输出("传送第3步：点击下一页 x=%d y=%d", x, y)
		僵尸3传送点击(x, y)
		return true
	}
	设置僵尸3层输出("传送第3步失败：未找到下一页")
	return false
}

func 僵尸3传送逻辑4(shouldContinue func() bool) bool {
	设置僵尸3层输出("传送第4个逻辑：循环下一页直到僵尸3地图")
	if ok, x, y := 僵尸3查找传送特征(传送已经点开传送石); ok {
		设置僵尸3层输出("传送第4步：点击已经点开传送石 x=%d y=%d", x, y)
		僵尸3传送点击(x, y)
	}
	if ok, x, y := 僵尸3查找传送特征(传送需要传送的地图); ok {
		return 僵尸3传送点击需要传送的地图(x, y, shouldContinue)
	}
	if ok, x, y := 僵尸3查找传送特征(传送选中僵尸3地图); ok {
		设置僵尸3层输出("传送第4步：已找到选中僵尸3地图 x=%d y=%d", x, y)
		return true
	}
	for count := 1; count <= 僵尸3传送翻页最大次数 && shouldContinue(); count++ {
		ok, x, y := 僵尸3等待传送特征(传送下一页, 2*time.Second, 200*time.Millisecond, shouldContinue)
		if !ok {
			设置僵尸3层输出("传送第4步：第%d次未找到下一页按钮", count)
			continue
		}
		设置僵尸3层输出("传送第4步：第%d次点击下一页 x=%d y=%d", count, x, y)
		僵尸3传送点击(x, y)
		time.Sleep(僵尸3传送翻页等待())
		if ok, x, y := 僵尸3查找传送特征(传送需要传送的地图); ok {
			return 僵尸3传送点击需要传送的地图(x, y, shouldContinue)
		}
		if ok, x, y := 僵尸3查找传送特征(传送选中僵尸3地图); ok {
			设置僵尸3层输出("传送第4步：找到选中僵尸3地图 x=%d y=%d", x, y)
			return true
		}
	}
	设置僵尸3层输出("传送第4步失败：翻页%d次未找到需要传送的地图/选中僵尸3地图，发送钉钉", 僵尸3传送翻页最大次数)
	发送钉钉文本("僵尸3传送翻页20次未找到需要传送的地图或选中僵尸3地图")
	return false
}

func 僵尸3传送逻辑5(shouldContinue func() bool) bool {
	设置僵尸3层输出("传送第5个逻辑：选中僵尸3地图并确认传送")
	if ok, x, y := 僵尸3等待传送特征(传送选中僵尸3地图, 僵尸3传送等待超时, 200*time.Millisecond, shouldContinue); ok {
		设置僵尸3层输出("传送第5步：点击选中僵尸3地图 x=%d y=%d", x, y)
		僵尸3传送点击(x, y)
	} else {
		设置僵尸3层输出("传送第5步失败：未找到选中僵尸3地图")
		return false
	}
	if ok, x, y := 僵尸3等待传送特征(传送点击传送, 僵尸3传送等待超时, 200*time.Millisecond, shouldContinue); ok {
		设置僵尸3层输出("传送第5步：点击传送 x=%d y=%d", x, y)
		僵尸3传送点击(x, y)
	} else {
		设置僵尸3层输出("传送第5步失败：未找到点击传送")
		return false
	}
	if ok, x, y := 僵尸3等待传送特征(传送点击确认, 僵尸3传送等待超时, 200*time.Millisecond, shouldContinue); ok {
		设置僵尸3层输出("传送第5步：点击确认 x=%d y=%d", x, y)
		僵尸3传送点击(x, y)
		return true
	}
	设置僵尸3层输出("传送第5步失败：未找到点击确认")
	return false
}

func 僵尸3传送点击需要传送的地图(x, y int, shouldContinue func() bool) bool {
	设置僵尸3层输出("传送第4步：找到需要传送的地图 x=%d y=%d，点击", x, y)
	僵尸3传送点击(x, y)
	if ok, vx, vy := 僵尸3等待传送特征(传送选中僵尸3地图, 僵尸3传送等待超时, 200*time.Millisecond, shouldContinue); ok {
		设置僵尸3层输出("传送第4步：已选中僵尸3地图 x=%d y=%d", vx, vy)
		return true
	}
	设置僵尸3层输出("传送第4步失败：点击需要传送的地图后未变成选中僵尸3地图")
	return false
}

func 僵尸3查找传送特征(feature any) (bool, int, int) {
	if 引擎 == nil {
		return false, -1, -1
	}
	ok, x, y := 引擎.FindFeature(feature)
	僵尸3传送查色后等待()
	return ok, x, y
}

func 僵尸3查找任一传送特征(candidates ...僵尸3传送特征候选) (bool, int, int, 僵尸3传送特征候选) {
	for _, candidate := range candidates {
		if ok, x, y := 僵尸3查找传送特征(candidate.特征); ok {
			return true, x, y, candidate
		}
	}
	return false, -1, -1, 僵尸3传送特征候选{}
}

func 僵尸3等待传送特征(feature any, timeout, interval time.Duration, shouldContinue func() bool) (bool, int, int) {
	if interval <= 0 {
		interval = 200 * time.Millisecond
	}
	deadline := time.Now().Add(timeout)
	for shouldContinue == nil || shouldContinue() {
		if ok, x, y := 僵尸3查找传送特征(feature); ok {
			return true, x, y
		}
		if timeout > 0 && time.Now().After(deadline) {
			return false, -1, -1
		}
		time.Sleep(interval)
	}
	return false, -1, -1
}

func 僵尸3等待任一传送特征(candidates []僵尸3传送特征候选, timeout, interval time.Duration, shouldContinue func() bool) (bool, int, int, 僵尸3传送特征候选) {
	if interval <= 0 {
		interval = 200 * time.Millisecond
	}
	deadline := time.Now().Add(timeout)
	for shouldContinue == nil || shouldContinue() {
		if ok, x, y, candidate := 僵尸3查找任一传送特征(candidates...); ok {
			return true, x, y, candidate
		}
		if timeout > 0 && time.Now().After(deadline) {
			return false, -1, -1, 僵尸3传送特征候选{}
		}
		time.Sleep(interval)
	}
	return false, -1, -1, 僵尸3传送特征候选{}
}

func 僵尸3传送翻页等待() time.Duration {
	return time.Duration(700+time.Now().Nanosecond()%501) * time.Millisecond
}

func 僵尸3传送查色后等待() {
	僵尸3传送随机等待(僵尸3传送查色后最小等待, 僵尸3传送查色后最大等待)
}

func 僵尸3传送操作后等待() {
	僵尸3传送随机等待(僵尸3传送操作后最小等待, 僵尸3传送操作后最大等待)
}

func 僵尸3传送点击(x, y int) {
	引擎.ClickResult(true, x, y)
	僵尸3传送操作后等待()
}

func 僵尸3传送双击(x, y int) {
	快速双击买卖物品坐标(x, y, 买卖物品猫猫双击间隔)
	僵尸3传送操作后等待()
}

func 僵尸3传送随机等待(min, max time.Duration) {
	if max <= min {
		time.Sleep(min)
		return
	}
	span := int64(max - min)
	time.Sleep(min + time.Duration(time.Now().UnixNano()%span))
}

func 执行僵尸3左进三层测试() {
	if 程序退出中.Load() {
		设置僵尸3层输出("左进3层测试失败：程序正在退出")
		return
	}
	if 脚本运行中.Load() {
		设置僵尸3层输出("左进3层测试失败：脚本运行中，请先点结束")
		return
	}
	if !僵尸3左进三层执行中.CompareAndSwap(false, true) {
		设置僵尸3层输出("左进3层测试执行中，请稍等")
		return
	}
	go func() {
		defer 僵尸3左进三层执行中.Store(false)
		设置僵尸3层输出("左进3层测试开始：从1层左侧出边界进3层")
		if 僵尸3一层到三层左侧出边界(func() bool { return !程序退出中.Load() }) {
			设置僵尸3层输出("左进3层测试成功")
		} else {
			设置僵尸3层输出("左进3层测试失败")
		}
	}()
}
