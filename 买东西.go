package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/Dasongzi1366/AutoGo/ime"
)

var ms买东西灰色 = &CColor{Name: "ms买东西灰色", X: 359, Y: 282, Color: "D5D5D5-000000,7,-5,FFFFFF-000000,12,-10,FFFFFF-000000,13,-3,FAFAFA-000000,17,1,FFFFFF-000000,25,7,CC0000-000000,14,7,CC0000-000000,-2,7,CC0000-000000,-10,7,CC0000-000000", Sim: 0.90}
var nt买东西红色 = &CColor{Name: "nt买东西红色", X: 356, Y: 279, Color: "FFFFFF-000000,2,-3,D7B1BA-000000,5,-6,EAD0D6-000000,12,-6,DEBDC4-000000,23,-7,FFFFFF-000000,20,2,FDF5F7-000000,27,10,CC0000-000000,15,10,CC0000-000000,5,6,F5EBED-000000", Sim: 0.90}

var ms下滑键 = &CColor{Name: "ms下滑键", X: 616, Y: 550, Color: "92B9D9-000000,4,-1,226FAA-000000,7,-1,226699-000000,12,0,7EA3C7-000000,11,3,37699B-000000,11,5,125588-000000,7,9,114578-000000,4,9,1A5588-000000,2,9,225999-000000", Sim: 0.90}
var ck红药水 = &CColor{Name: "ck红药水", X: 351, Y: 431, Color: "000000-000000,1,-15,887744-000000,47,-18,596161-000000,59,-14,33333B-000000,72,-11,000000-000000,69,11,000000-000000,65,18,CCD0EA-000000,22,16,736B69-000000,9,14,000000-000000", Sim: 0.82}
var ck蓝盐水 = &CColor{Name: "ck蓝盐水", X: 348, Y: 378, Color: "000000-000000,-4,-19,373735-000000,22,-18,3A3736-000000,56,-10,000000-000000,82,-9,333737-000000,81,9,333738-000000,69,15,000000-000000,27,18,787975-000000,2,17,898075-000000", Sim: 0.82}

var MS系统应用 = &FMColor{Name: "系统应用", X1: 834, Y1: 180, X2: 897, Y2: 224, MainColor: "DCDCDC-000000", OffsetColor: "7,-3,DCDCDC-000000,27,-3,D4D4D4-000000,0,8,DCDCDC-000000,7,4,C9C9C9-000000,34,1,DCDCDC-000000,6,12,D9D9D9-000000,13,16,DCDCDC-000000,34,9,DCDCDC-000000", Sim: 0.90, Dir: 0}

type 买东西步骤 struct {
	名称 string
	执行 func() string
}

var (
	买东西步骤表 = []买东西步骤{
		{名称: "买东西页面", 执行: 执行买东西页面步骤},
		{名称: "买红药水", 执行: func() string { return 执行购买药水步骤("红药水", ck红药水) }},
		{名称: "买蓝盐水", 执行: func() string { return 执行购买药水步骤("蓝盐水", ck蓝盐水) }},
	}
	买东西锁       sync.Mutex
	买东西已启动     bool
	买东西当前步骤    int
	买东西页签等待    = 550 * time.Millisecond
	买东西下滑前等待   = 450 * time.Millisecond
	买东西下滑间隔    = 700 * time.Millisecond
	买东西双击后等待   = 650 * time.Millisecond
	买东西输入前等待   = 650 * time.Millisecond
	买东西输入后等待   = 650 * time.Millisecond
	买东西确认后等待   = 650 * time.Millisecond
	买东西找商品下滑次数 = 50
	买东西购买数量    = "300"
)

func 启动买东西流程() {
	if 引擎 == nil {
		输出("买东西启动失败：引擎未初始化")
		return
	}

	买东西锁.Lock()
	买东西已启动 = true
	买东西当前步骤 = 0
	买东西锁.Unlock()
}

func 买东西下一步() {
	index, step, ok := 取买东西当前步骤()
	if !ok {
		return
	}
	输出(step.执行())
	设置买东西下一步骤(index + 1)
}

func 取买东西当前步骤() (int, 买东西步骤, bool) {
	买东西锁.Lock()
	defer 买东西锁.Unlock()
	if !买东西已启动 {
		输出("买东西未开始")
		return 0, 买东西步骤{}, false
	}
	if 买东西当前步骤 >= len(买东西步骤表) {
		买东西已启动 = false
		输出("买东西流程结束")
		return 0, 买东西步骤{}, false
	}
	return 买东西当前步骤, 买东西步骤表[买东西当前步骤], true
}

func 设置买东西下一步骤(next int) {
	买东西锁.Lock()
	defer 买东西锁.Unlock()
	买东西当前步骤 = next
	if 买东西当前步骤 >= len(买东西步骤表) {
		买东西已启动 = false
	}
}

func 执行买东西页面步骤() string {
	if 匹配买卖物品特征(nt买东西红色) {
		return "买东西页面成功，已是红色"
	}
	if !点击买卖物品特征(ms买东西灰色) {
		return "买东西页面失败，未找到灰色按钮，也不是红色"
	}
	time.Sleep(买东西页签等待)
	if 匹配买卖物品特征(nt买东西红色) {
		return "买东西页面成功，灰色已点开为红色"
	}
	return "买东西页面失败，点击后未变红色"
}

func 执行购买药水步骤(名称 string, 商品 *CColor) string {
	if !确保买东西页面打开() {
		return fmt.Sprintf("买东西%s失败，买东西页面未打开", 名称)
	}

	x, y, 下滑次数, ok := 下滑查找买东西商品(商品)
	if !ok {
		return fmt.Sprintf("买东西%s失败，下滑%d次仍未找到%s", 名称, 下滑次数, 商品.Name)
	}

	双击买卖物品坐标(x, y)
	time.Sleep(买东西双击后等待)
	输入买东西数量(买东西购买数量)
	if !等待并点击买东西确定() {
		return fmt.Sprintf("买东西%s失败，已双击%s并输入%s，但未找到确定", 名称, 商品.Name, 买东西购买数量)
	}
	time.Sleep(买东西确认后等待)
	return fmt.Sprintf("买东西%s成功，下滑=%d，双击=(%d,%d)，数量=%s，确定成功", 名称, 下滑次数, x, y, 买东西购买数量)
}

func 确保买东西页面打开() bool {
	if 匹配买卖物品特征(nt买东西红色) {
		return true
	}
	if !点击买卖物品特征(ms买东西灰色) {
		return false
	}
	time.Sleep(买东西页签等待)
	return 匹配买卖物品特征(nt买东西红色)
}

func 下滑查找买东西商品(商品 *CColor) (int, int, int, bool) {
	for attempt := 0; attempt <= 买东西找商品下滑次数; attempt++ {
		if ok, x, y := 静默匹配买卖物品特征(商品); ok {
			return x, y, attempt, true
		}
		if attempt == 买东西找商品下滑次数 {
			break
		}
		time.Sleep(买东西下滑前等待)
		if !点击买卖物品特征(ms下滑键) {
			return -1, -1, attempt, false
		}
		time.Sleep(买东西下滑间隔)
		if ok, x, y := 静默匹配买卖物品特征(商品); ok {
			return x, y, attempt + 1, true
		}
	}
	return -1, -1, 买东西找商品下滑次数, false
}

func 输入买东西数量(text string) {
	time.Sleep(买东西输入前等待)
	ime.InputText(text)
	time.Sleep(买东西输入后等待)
}

func 等待并点击买东西确定() bool {
	deadline := time.Now().Add(1500 * time.Millisecond)
	for time.Now().Before(deadline) {
		if 匹配买卖物品特征(MS系统应用) {
			return 点击买卖物品特征(MS系统应用)
		}
		time.Sleep(100 * time.Millisecond)
	}
	return 点击买卖物品特征(MS系统应用)
}
