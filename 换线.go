package main

import (
	"sync"
	"time"

	"github.com/Dasongzi1366/AutoGo/motion"
)

var BOSS红金特征 = &FMColor{Name: "BOSS红金特征", X1: 1, Y1: 256, X2: 1269, Y2: 522, MainColor: "D3B354-202020", OffsetColor: "1,0,F5D36D-202020,11,3,DEC063-202020,-6,7,C8A639-202020,4,10,FB6D6D-202020,11,7,CDAB3D-202020,-3,21,E5C356-202020,4,14,A10000-202020,11,14,D7B546-202020", Sim: 0.90, Dir: 0}
var BOSS金红特征 = &FMColor{Name: "BOSS金红特征", X1: 1, Y1: 256, X2: 1269, Y2: 522, MainColor: "F6D470-202020", OffsetColor: "7,0,040505-202020,15,0,110D04-202020,2,12,18140D-202020,7,6,FDB7B7-202020,15,9,D9B742-202020,0,13,E3C154-202020,9,13,170000-202020,10,19,EBC959-202020", Sim: 0.90, Dir: 0}
var 僵尸3地图有人 = &FMColor{Name: "僵尸3地图有人", X1: 8, Y1: 96, X2: 262, Y2: 203, MainColor: "FA0000-202020", OffsetColor: "2,0,FA0000-202020,4,1,FF0000-202020,0,2,FA0000-202020,2,2,FA0000-202020,0,4,F40202-202020,1,4,FF0000-202020,4,4,343334-202020", Sim: 0.90, Dir: 0}
var 打开菜单 = &FMColor{Name: "打开菜单", X1: 913, Y1: 668, X2: 984, Y2: 714, MainColor: "774400-000000", OffsetColor: "7,0,FFFFFF-000000,13,-7,FFAA22-000000,14,-2,895E12-000000,15,0,774400-000000,15,3,0D4066-000000,11,6,FFFFFF-000000,4,5,004477-000000,3,2,513719-000000", Sim: 0.90, Dir: 0}
var 已经成功打开菜单 = &FMColor{Name: "已经成功打开菜单", X1: 898, Y1: 424, X2: 1009, Y2: 507, MainColor: "BFE5F7-000000", OffsetColor: "0,-1,BFDDEE-000000,-2,-6,0077BB-000000,11,-7,FFFFFF-000000,7,0,008CD5-000000,13,5,008CD0-000000,7,5,0077BB-000000,1,5,C0DEEE-000000,-5,5,FFFFFF-000000", Sim: 0.90, Dir: 0}
var 点击换线 = &FMColor{Name: "点击换线", X1: 898, Y1: 424, X2: 1009, Y2: 507, MainColor: "BFE5F7-000000", OffsetColor: "0,-1,BFDDEE-000000,-2,-6,0077BB-000000,11,-7,FFFFFF-000000,7,0,008CD5-000000,13,5,008CD0-000000,7,5,0077BB-000000,1,5,C0DEEE-000000,-5,5,FFFFFF-000000", Sim: 0.90, Dir: 0}
var 成功打开换线界面 = &FMColor{Name: "成功打开换线界面", X1: 338, Y1: 207, X2: 472, Y2: 273, MainColor: "303030-000000", OffsetColor: "-1,-2,000000-000000,11,-1,404040-000000,23,-2,000000-000000,17,0,000000-000000,25,3,000000-000000,15,3,FFFFFF-000000,7,3,FFFFFF-000000,-5,3,404040-000000", Sim: 0.90, Dir: 0}
var 当前线路 = &FMColor{Name: "当前线路", X1: 350, Y1: 308, X2: 913, Y2: 453, MainColor: "569B61-000000", OffsetColor: "3,0,5BAD6C-000000,6,0,FFFFFF-000000,1,4,5AAC6B-000000,3,4,54A665-000000,6,4,FFFFFF-000000,1,7,4C9E5D-000000,2,7,4EA05F-000000,6,6,FFFFFF-000000", Sim: 0.90, Dir: 0}
var 选中需要换线 = &FMColor{Name: "选中需要换线", X1: 350, Y1: 308, X2: 913, Y2: 453, MainColor: "26A6D9-000000", OffsetColor: "6,0,33AACC-000000,7,0,2299CC-000000,2,4,33AACC-000000,6,4,33AADD-000000,10,2,3399D0-000000,2,7,3296CF-000000,3,7,339BC9-000000,8,6,33A9DA-000000", Sim: 0.90, Dir: 0}
var 已选中需要换线 = &FMColor{Name: "已选中需要换线", X1: 350, Y1: 308, X2: 913, Y2: 453, MainColor: "26A6D9-000000", OffsetColor: "6,0,33AACC-000000,7,0,2299CC-000000,2,4,33AACC-000000,6,4,33AADD-000000,10,2,3399D0-000000,2,7,3296CF-000000,3,7,339BC9-000000,8,6,33A9DA-000000", Sim: 0.90, Dir: 0}
var 确定换线按钮 = &FMColor{Name: "确定换线按钮", X1: 800, Y1: 463, X2: 878, Y2: 495, MainColor: "FFDDCC-000000", OffsetColor: "10,0,FFDDCC-000000,11,0,FFDDCC-000000,0,1,AA6622-000000,10,4,FFDDCC-000000,11,4,FFDDCC-000000,-5,5,FFDDCC-000000,10,9,FFDDCC-000000,16,9,FFDDCC-000000", Sim: 0.90, Dir: 0}
var 再次确认按钮 = &FMColor{Name: "再次确认按钮", X1: 695, Y1: 421, X2: 749, Y2: 449, MainColor: "FBD6C1-000000", OffsetColor: "10,-2,FFDDCC-000000,16,-2,C36A22-000000,0,1,FBD6C1-000000,5,1,EEC0A2-000000,16,2,C86F26-000000,-5,6,FFEECC-000000,5,5,BF7733-000000,11,6,CC7733-000000", Sim: 0.90, Dir: 0}
var 正在换线界面 = &FMColor{Name: "正在换线界面", X1: 546, Y1: 119, X2: 748, Y2: 166, MainColor: "00FF8A-000000", OffsetColor: "5,0,00FF8A-000000,15,0,00FF8A-000000,0,7,00FF8A-000000,5,4,00FE89-000000,15,7,00FF8A-000000,0,8,00BD66-000000,5,8,004324-000000,15,8,00FF8A-000000", Sim: 0.90, Dir: 0}
var 僵尸3地图 = &FMColor{Name: "僵尸3地图", X1: 65, Y1: 34, X2: 200, Y2: 81, MainColor: "FBFBFB-000000", OffsetColor: "1,0,737779-000000,12,0,5A676E-000000,-5,9,72838B-000000,6,5,E0E9ED-000000,17,9,09090A-000000,-5,10,4F6169-000000,1,14,5C717B-000000,17,10,060607-000000", Sim: 0.90, Dir: 0}
var 僵尸4地图 = &FMColor{Name: "僵尸4地图", X1: 65, Y1: 37, X2: 200, Y2: 84, MainColor: "99BBCC-000000", OffsetColor: "1,0,99BBCC-000000,14,0,99BBCC-000000,-4,9,283135-000000,5,9,21282C-000000,14,1,FFFFFF-000000,-4,10,FEFEFE-000000,1,10,F6F6F6-000000,14,14,3C4A50-000000", Sim: 0.90, Dir: 0}
var 已重置workspace = &FMColor{Name: "系统应用", X1: 516, Y1: 531, X2: 731, Y2: 660, MainColor: "373737-202020", OffsetColor: "25,0,949494-202020,26,0,000000-202020,6,7,949494-202020,25,1,949494-202020,32,7,555555-202020,0,11,767676-202020,25,11,8D8D8D-202020,26,11,8D8D8D-202020", Sim: 0.90, Dir: 0}

var BOSS换线锁 sync.Mutex

var 僵尸3换线触发特征列表 = []*FMColor{
	BOSS红金特征,
	BOSS金红特征,
	僵尸3地图有人,
}

const (
	换线区域左 = 350
	换线区域上 = 308
	换线区域右 = 913
	换线区域下 = 453
	换线格宽  = 89
	换线格高  = 26
	换线列数  = 6
	换线行数  = 5

	BOSS右冲X按下时间 = 60 * time.Millisecond
	BOSS右冲X间隔   = 40 * time.Millisecond
	BOSS右冲检测间隔  = 20 * time.Millisecond
)

func 僵尸3检查BOSS并换线(runID int64) bool {
	if 引擎 == nil || !脚本仍应运行(runID) {
		return false
	}
	ok, name, x, y := 查找僵尸3换线触发特征()
	if !ok {
		return false
	}
	if !BOSS换线锁.TryLock() {
		return true
	}
	defer BOSS换线锁.Unlock()

	for 脚本仍应运行(runID) {
		设置僵尸3层输出("发现换线特征：%s x=%d y=%d，固定右+X猛冲去3层换线", name, x, y)
		if !僵尸3BOSS右X冲到三层(runID) {
			设置僵尸3层输出("BOSS换线失败：未到3层")
			return true
		}
		if !执行BOSS换线流程(func() bool { return 脚本仍应运行(runID) }) {
			设置僵尸3层输出("BOSS换线流程失败")
			return true
		}
		设置僵尸3层输出("换线完成；新地图立即重新扫描BOSS或其他玩家")
		if nextOK, nextName, nextX, nextY := 查找僵尸3换线触发特征(); nextOK {
			设置僵尸3层输出("换线后再次发现：%s x=%d y=%d，继续换线", nextName, nextX, nextY)
			name, x, y = nextName, nextX, nextY
			continue
		}
		设置僵尸3层输出("换线后未发现BOSS或其他玩家，不再换线")
		break
	}

	if 位置, ok := 僵尸3当前层位置(); ok {
		if 位置.层 == 3 {
			设置僵尸3层输出("BOSS换线后在3层：回1层继续打怪")
			僵尸3三层回一层直到成功(runID)
		} else {
			设置僵尸3层输出("BOSS换线后当前位置：层=%d x=%d y=%d", 位置.层, 位置.X, 位置.Y)
		}
	} else {
		设置僵尸3层输出("BOSS换线后未识别当前位置，交给主循环")
	}
	return true
}

func 查找僵尸3换线触发特征() (bool, string, int, int) {
	for _, feature := range 僵尸3换线触发特征列表 {
		if feature == nil {
			continue
		}
		if ok, x, y := 查找BOSS特征不画框(feature); ok {
			return true, feature.Name, x, y
		}
	}
	return false, "", -1, -1
}

func 查找BOSS特征不画框(feature *FMColor) (bool, int, int) {
	暂停调试红框()
	defer 恢复调试红框()
	return 引擎.FindFeature(feature)
}

func 僵尸3BOSS右X冲到三层(runID int64) bool {
	设置僵尸3层输出("BOSS右冲：固定路线，按住右并连点X直到3层")
	释放所有按键()
	displayID := 当前显示ID()
	motion.KeyActionDown(motion.KEYCODE_DPAD_RIGHT, displayID)
	defer func() {
		motion.KeyActionUp(motion.KEYCODE_X, displayID)
		motion.KeyActionUp(motion.KEYCODE_DPAD_RIGHT, displayID)
	}()

	nextLog := time.Now()
	nextX := time.Now()
	for 脚本仍应运行(runID) {
		位置, ok := 僵尸3当前层位置()
		if ok && 位置.层 == 3 {
			设置僵尸3层输出("BOSS右冲：已到3层")
			return true
		}
		if ok && !time.Now().Before(nextLog) {
			设置僵尸3层输出("BOSS右冲中：层=%d x=%d y=%d", 位置.层, 位置.X, 位置.Y)
			nextLog = time.Now().Add(2 * time.Second)
		}
		if !time.Now().Before(nextX) {
			motion.KeyActionDown(motion.KEYCODE_X, displayID)
			time.Sleep(BOSS右冲X按下时间)
			motion.KeyActionUp(motion.KEYCODE_X, displayID)
			nextX = time.Now().Add(BOSS右冲X间隔)
			continue
		}
		time.Sleep(BOSS右冲检测间隔)
	}
	return false
}

func 执行BOSS换线流程(shouldContinue func() bool) bool {
	if 引擎 == nil {
		return false
	}
	if !换线点击并验证("打开菜单", 打开菜单, 已经成功打开菜单, shouldContinue) {
		return false
	}
	if !换线点击并验证("点击换线", 点击换线, 成功打开换线界面, shouldContinue) {
		return false
	}
	if !换线选择线路(shouldContinue) {
		return false
	}
	if !换线点击并验证("确定换线", 确定换线按钮, 再次确认按钮, shouldContinue) {
		return false
	}
	if !等待并点击(再次确认按钮, 3*time.Second, 200*time.Millisecond) {
		设置僵尸3层输出("换线：再次确认按钮未点击")
		return false
	}
	if !换线点击正在换线并监控(shouldContinue) {
		return false
	}
	return true
}

func 换线点击正在换线并监控(shouldContinue func() bool) bool {
	clickDeadline := time.Now().Add(20 * time.Second)
	for shouldContinue() && time.Now().Before(clickDeadline) {
		if ok, x, y := 引擎.FindFeature(正在换线界面); ok {
			设置僵尸3层输出("换线：点击正在换线 x=%d y=%d，开始4分钟监控", x, y)
			引擎.ClickResult(true, x, y)
			time.Sleep(300 * time.Millisecond)
			return 换线等待地图成功或超时(shouldContinue)
		}
		time.Sleep(300 * time.Millisecond)
	}
	设置僵尸3层输出("换线：正在换线失败")
	return false
}

func 换线等待地图成功或超时(shouldContinue func() bool) bool {
	deadline := time.Now().Add(4 * time.Minute)
	for shouldContinue() && time.Now().Before(deadline) {
		if ok, _, _ := 引擎.FindFeature(僵尸3地图); ok {
			设置僵尸3层输出("换线：已到僵尸3地图，继续流程")
			暂停僵尸3地图巡检(僵尸3换线后巡检暂停, "换线成功")
			return true
		}
		if ok, x, y := 引擎.FindFeature(已重置workspace); ok {
			设置僵尸3层输出("换线：检测到系统应用 x=%d y=%d，判断换线完成", x, y)
			暂停僵尸3地图巡检(僵尸3换线后巡检暂停, "换线成功")
			return true
		}
		if ok, _, _ := 引擎.FindFeature(僵尸4地图); ok {
			设置僵尸3层输出("换线：检测到僵尸4地图，继续等待僵尸3地图")
		}
		time.Sleep(500 * time.Millisecond)
	}
	设置僵尸3层输出("换线：4分钟未到僵尸3地图，发送钉钉")
	发送钉钉文本("切换僵尸3地图失败。")
	return false
}

func 换线点击并验证(name string, clickFeature any, verifyFeature any, shouldContinue func() bool) bool {
	return 换线点击并验证任一(name, clickFeature, []any{verifyFeature}, 5*time.Second, shouldContinue)
}

func 换线点击并验证任一(name string, clickFeature any, verifyFeatures []any, timeout time.Duration, shouldContinue func() bool) bool {
	deadline := time.Now().Add(timeout)
	clicked := false
	for shouldContinue() && time.Now().Before(deadline) {
		if !clicked {
			if ok, x, y := 引擎.FindFeature(clickFeature); ok {
				设置僵尸3层输出("换线：点击%s x=%d y=%d", name, x, y)
				引擎.ClickResult(true, x, y)
				clicked = true
				time.Sleep(300 * time.Millisecond)
			}
		}
		for _, feature := range verifyFeatures {
			if ok, x, y := 引擎.FindFeature(feature); ok {
				设置僵尸3层输出("换线：%s验证成功 x=%d y=%d", name, x, y)
				return true
			}
		}
		time.Sleep(200 * time.Millisecond)
	}
	设置僵尸3层输出("换线：%s失败", name)
	return false
}

func 换线选择线路(shouldContinue func() bool) bool {
	currentIndex := 换线当前线路格子()
	order := 移动随机.Perm(换线列数 * 换线行数)
	for _, i := range order {
		if !shouldContinue() {
			return false
		}
		if i == currentIndex {
			continue
		}
		x, y := 换线格子中心(i)
		设置僵尸3层输出("换线：选择线路格子%d x=%d y=%d", i+1, x, y)
		引擎.Click(x, y)
		time.Sleep(350 * time.Millisecond)
		if ok, vx, vy := 引擎.FindFeature(已选中需要换线); ok {
			设置僵尸3层输出("换线：线路已选中 x=%d y=%d", vx, vy)
			return true
		}
	}
	return false
}

func 换线当前线路格子() int {
	ok, x, y := 引擎.FindFeature(当前线路)
	if !ok {
		return -1
	}
	for i := 0; i < 换线列数*换线行数; i++ {
		left, top, right, bottom := 换线格子矩形(i)
		if x >= left && x <= right && y >= top && y <= bottom {
			设置僵尸3层输出("换线：当前线路格子%d", i+1)
			return i
		}
	}
	return -1
}

func 换线格子中心(index int) (int, int) {
	left, top, right, bottom := 换线格子矩形(index)
	return (left + right) / 2, (top + bottom) / 2
}

func 换线格子矩形(index int) (int, int, int, int) {
	row := index / 换线列数
	col := index % 换线列数
	left := 换线区域左
	top := 换线区域上
	if 换线列数 > 1 {
		left += (col*(换线区域右-换线区域左-换线格宽) + (换线列数-1)/2) / (换线列数 - 1)
	}
	if 换线行数 > 1 {
		top += (row*(换线区域下-换线区域上-换线格高) + (换线行数-1)/2) / (换线行数 - 1)
	}
	return left, top, left + 换线格宽, top + 换线格高
}
