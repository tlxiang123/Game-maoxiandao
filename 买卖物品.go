package main

import (
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/Dasongzi1366/AutoGo/images"
	"github.com/Dasongzi1366/AutoGo/motion"
)

var 背包特征 = &FMColor{Name: "背包特征", X1: 853, Y1: 613, X2: 944, Y2: 675, MainColor: "EEAA33-000000", OffsetColor: "10,-8,BB8833-000000,11,-12,0066BB-000000,2,-1,EE9933-000000,9,-13,0066CC-000000,9,-12,AADDEE-000000,10,0,EE9933-000000,12,-9,0066BB-000000,12,0,BB7722-000000", Sim: 0.90, Dir: 0}
var 打开背包成功 = &FMColor{Name: "打开背包成功", X1: 530, Y1: 74, X2: 1276, Y2: 191, MainColor: "000000-000000", OffsetColor: "-8,1,DDAA22-000000,-10,-1,000000-000000,-3,-5,CC3322-000000,0,-4,000000-000000,-5,-5,9999AA-000000,-9,4,EE6644-000000,-8,6,FFFFFF-000000,-6,-4,000000-000000", Sim: 0.90, Dir: 0}
var 第五页灰色 = &FMColor{Name: "第五页灰色", X1: 610, Y1: 77, X2: 1274, Y2: 192, MainColor: "C7C7C7-000000", OffsetColor: "-1,-2,6C6C6C-000000,-18,-12,F8F8F8-000000,-12,2,B4B4B4-000000,-23,-9,808080-000000,-9,-4,F3F3F3-000000,-23,-2,777777-000000,-14,1,9F9F9F-000000,-3,-12,A3A3A3-000000", Sim: 0.90, Dir: 0}
var 第五页面粉色 = &FMColor{Name: "第五页面粉色", X1: 636, Y1: 77, X2: 1277, Y2: 256, MainColor: "8B3B4F-000000", OffsetColor: "-25,0,AE5269-000000,-18,7,FEF5F8-000000,-26,15,F56680-000000,-21,1,FDEFF2-000000,-13,10,C28C9A-000000,-4,4,A5475E-000000,-1,-1,E96183-000000,0,8,EBE5E6-000000", Sim: 0.90, Dir: 0}
var 找到第五页面 = &FMColor{Name: "找到第五页面", X1: 636, Y1: 77, X2: 1277, Y2: 256, MainColor: "8B3B4F-000000", OffsetColor: "-25,0,AE5269-000000,-18,7,FEF5F8-000000,-26,15,F56680-000000,-21,1,FDEFF2-000000,-13,10,C28C9A-000000,-4,4,A5475E-000000,-1,-1,E96183-000000,0,8,EBE5E6-000000", Sim: 0.90, Dir: 0}
var 双击猫猫图标 = &FMColor{Name: "双击猫猫图标", X1: 751, Y1: 74, X2: 1274, Y2: 458, MainColor: "EEEEDD-000000", OffsetColor: "11,-9,C7A562-000000,13,-7,967749-000000,22,18,000000-000000,8,-12,EFEFDC-000000,10,-13,E9E9D8-000000,-1,11,000000-000000,22,21,000000-000000,30,21,FFCC00-000000", Sim: 0.90, Dir: 0}
var 找到全卖按钮 = &FMColor{Name: "找到全卖按钮", X1: 321, Y1: 111, X2: 937, Y2: 588, MainColor: "CB6E1A-000000", OffsetColor: "-35,3,CE7022-000000,-26,13,F5D1B2-000000,-2,10,FF8C44-000000,-13,0,CB6E1A-000000,-10,12,FFEEDD-000000,-32,15,DB7E3A-000000,-9,16,FF9955-000000,-34,13,DB7E38-000000", Sim: 0.90, Dir: 0}
var 全卖页面灰色 = &FMColor{Name: "系统应用", X1: 498, Y1: 219, X2: 954, Y2: 346, MainColor: "868686-000000", OffsetColor: "7,0,868686-000000,19,0,8C8C8C-000000,0,1,D6D6D6-000000,4,7,DFDFDF-000000,15,4,F4F4F4-000000,3,14,EDEDED-000000,4,14,FFFFFF-000000,15,8,ECECEC-000000", Sim: 0.90, Dir: 0}
var 全卖页面粉色 = &FMColor{Name: "系统应用", X1: 498, Y1: 219, X2: 954, Y2: 346, MainColor: "FFFFFF-000000", OffsetColor: "14,0,FFFFFF-000000,28,0,FFFFFF-000000,13,16,D3C2C6-000000,14,12,E7D2D7-000000,28,12,F7F1F2-000000,0,25,CC0000-000000,14,25,CC0000-000000,28,25,CC0000-000000", Sim: 0.90, Dir: 0}
var 找到全卖页面 = &FMColor{Name: "找到全卖页面", X1: 321, Y1: 111, X2: 937, Y2: 588, MainColor: "EE6688-000000", OffsetColor: "7,-8,FFFFFF-000000,10,-14,BB4B64-000000,3,-4,AF4B64-000000,-1,-8,FF6688-000000,29,-4,AF4B64-000000,-2,-13,FF6688-000000,19,-16,E1667B-000000,34,-8,FF6688-000000", Sim: 0.90, Dir: 0}
var 空包袱 = &FMColor{Name: "系统应用", X1: 714, Y1: 257, X2: 944, Y2: 340, MainColor: "FFFFFF-000000", OffsetColor: "4,0,888888-000000,7,0,888888-000000,0,5,888888-000000,3,3,8E8E8E-000000,10,3,888888-000000,0,8,888888-000000,6,10,7A7A7A-000000,10,8,7A7A7A-000000", Sim: 0.90, Dir: 0}
var 全卖空包袱 = &FMColor{Name: "系统应用", X1: 714, Y1: 257, X2: 944, Y2: 340, MainColor: "FFFFFF-000000", OffsetColor: "4,0,888888-000000,7,0,888888-000000,0,5,888888-000000,3,3,8E8E8E-000000,10,3,888888-000000,0,8,888888-000000,6,10,7A7A7A-000000,10,8,7A7A7A-000000", Sim: 0.90, Dir: 0}
var 包袱不是空的 = &FMColor{Name: "包袱不是空的", X1: 321, Y1: 111, X2: 937, Y2: 588, MainColor: "2277AA-000000", OffsetColor: "7,16,83A7C4-000000,-3,-2,7D8B9C-000000,-2,12,979EAF-000000,8,11,EEEEFF-000000,-1,21,6699BB-000000,6,7,226088-000000,3,18,6688AA-000000,2,15,EEEEFF-000000", Sim: 0.90, Dir: 0}
var 全卖页面 = &FMColor{Name: "全卖页面", X1: 321, Y1: 111, X2: 937, Y2: 588, MainColor: "FF6688-000000", OffsetColor: "13,17,B87B83-000000,-1,12,EE6688-000000,5,9,AE586E-000000,3,-1,E5667B-000000,24,11,AF4B64-000000,1,4,FF6688-000000,6,5,B2475F-000000,10,2,E4C5CC-000000", Sim: 0.90, Dir: 0}
var 可以点击全卖 = &FMColor{Name: "可以点击全卖", X1: 321, Y1: 111, X2: 937, Y2: 588, MainColor: "FF6688-000000", OffsetColor: "13,17,B87B83-000000,-1,12,EE6688-000000,5,9,AE586E-000000,3,-1,E5667B-000000,24,11,AF4B64-000000,1,4,FF6688-000000,6,5,B2475F-000000,10,2,E4C5CC-000000", Sim: 0.90, Dir: 0}
var 点击全卖按钮 = &FMColor{Name: "点击全卖按钮", X1: 321, Y1: 111, X2: 937, Y2: 588, MainColor: "DF823E-000000", OffsetColor: "-9,-13,E7B688-000000,-5,-12,D68D4D-000000,-19,3,664D3B-000000,7,-10,D77022-000000,2,-15,FF9955-000000,-17,-5,DB8A4E-000000,-20,-6,D3762A-000000,-7,-7,EE8833-000000", Sim: 0.90, Dir: 0}
var 全卖确定按钮 = &FMColor{Name: "全卖确定按钮", X1: 321, Y1: 111, X2: 937, Y2: 588, MainColor: "FF9544-000000", OffsetColor: "-4,4,333333-000000,0,-4,EE882F-000000,27,-10,AA5511-000000,30,-7,AB5915-000000,9,3,D07B44-000000,19,1,D97740-000000,7,3,FF9955-000000,-2,-9,DD7711-000000", Sim: 0.90, Dir: 0}
var 点击全卖确定 = &FMColor{Name: "点击全卖确定", X1: 321, Y1: 111, X2: 937, Y2: 588, MainColor: "FF9544-000000", OffsetColor: "-4,4,333333-000000,0,-4,EE882F-000000,27,-10,AA5511-000000,30,-7,AB5915-000000,9,3,D07B44-000000,19,1,D97740-000000,7,3,FF9955-000000,-2,-9,DD7711-000000", Sim: 0.90, Dir: 0}
var 全卖确定消失 = &FMColor{Name: "全卖确定消失", X1: 321, Y1: 111, X2: 937, Y2: 588, MainColor: "FF9544-000000", OffsetColor: "-4,4,333333-000000,0,-4,EE882F-000000,27,-10,AA5511-000000,30,-7,AB5915-000000,9,3,D07B44-000000,19,1,D97740-000000,7,3,FF9955-000000,-2,-9,DD7711-000000", Sim: 0.90, Dir: 0}
var 单卖页面灰色 = &FMColor{Name: "单卖页面灰色", X1: 335, Y1: 123, X2: 943, Y2: 588, MainColor: "929292-000000", OffsetColor: "15,3,BBBBBB-000000,15,11,BBBBBB-000000,-4,14,969696-000000,7,3,808080-000000,-17,11,BBBBBB-000000,1,6,FFFFFF-000000,12,3,929292-000000,7,5,7F7F7F-000000", Sim: 0.90, Dir: 0}
var 已经切换单卖界面 = &FMColor{Name: "已经切换单卖界面", X1: 335, Y1: 123, X2: 943, Y2: 588, MainColor: "FF6688-000000", OffsetColor: "-18,7,FFFFFF-000000,-11,3,9F5364-000000,-4,9,EE6688-000000,-1,1,FF6688-000000,-12,4,A84C60-000000,-17,5,C9B0B6-000000,-8,-1,C8506B-000000,-23,3,CBAFB5-000000", Sim: 0.90, Dir: 0}
var 已切换单卖界面 = &FMColor{Name: "已切换单卖界面", X1: 335, Y1: 123, X2: 943, Y2: 588, MainColor: "FF6688-000000", OffsetColor: "-18,7,FFFFFF-000000,-11,3,9F5364-000000,-4,9,EE6688-000000,-1,1,FF6688-000000,-12,4,A84C60-000000,-17,5,C9B0B6-000000,-8,-1,C8506B-000000,-23,3,CBAFB5-000000", Sim: 0.90, Dir: 0}
var 单卖第一个格有物品 = &FMColor{Name: "单卖第一个格有物品", X1: 647, Y1: 297, X2: 914, Y2: 564, MainColor: "EE9922-000000", OffsetColor: "11,0,EE9922-000000,23,0,EE9922-000000,-5,9,EE9922-000000,6,9,EE9922-000000,18,9,EE9922-000000,0,19,EE9922-000000,11,19,EE9922-000000,23,19,EE9922-000000", Sim: 0.90, Dir: 0}
var 单卖确认卖出 = &FMColor{Name: "单卖确认卖出", X1: 1112, Y1: 611, X2: 1242, Y2: 705, MainColor: "212121-000000", OffsetColor: "27,-6,FFFFFF-000000,34,0,FBFBFB-000000,0,7,212121-000000,14,7,212121-000000,41,13,212121-000000,13,14,2E2E2E-000000,14,14,3E3E3E-000000,41,14,212121-000000", Sim: 0.90, Dir: 0}
var 单卖确认卖出备用 = &FMColor{Name: "单卖确认卖出备用", X1: 694, Y1: 405, X2: 750, Y2: 453, MainColor: "333333-000000", OffsetColor: "7,0,333333-000000,15,0,333333-000000,6,6,F3C7AD-000000,7,6,FFDDCC-000000,22,11,FFDDCC-000000,6,12,FFDDCC-000000,10,12,FFDDCC-000000,22,17,D99567-000000", Sim: 0.90, Dir: 0}
var 双击单卖之后的确定单击按钮 = &FMColor{Name: "双击单卖之后的确定单击按钮", X1: 691, Y1: 418, X2: 750, Y2: 468, MainColor: "AA5511-000000", OffsetColor: "4,0,AA5511-000000,14,0,FFDDCC-000000,-2,1,FFDDCC-000000,4,1,FFDDCC-000000,14,1,AA6611-000000,3,5,FFDDCC-000000,4,5,EDBE9F-000000,14,9,FFDDCC-000000", Sim: 0.90, Dir: 0}
var 单卖空格子颜色特征 = &FMColor{Name: "单卖空格子颜色特征", MainColor: "D9DDEA-000000", OffsetColor: "10,0,CCDDDD-000000,29,0,CCDDDD-000000,0,1,F5F7FA-000000,14,1,F2F7F7-000000,24,1,F2F7F7-000000,4,17,CFD3E8-000000,10,12,CCCCEE-000000,29,12,DDDDDD-000000", Sim: 0.90, Dir: 0}
var 关闭商店 = &FMColor{Name: "关闭商店", X1: 519, Y1: 146, X2: 631, Y2: 179, MainColor: "333333-000000", OffsetColor: "6,0,333333-000000,22,0,333333-000000,5,9,9DAE33-000000,16,4,FFFFDD-000000,17,4,889922-000000,-5,13,8A8A8A-000000,11,16,FFFFEE-000000,17,10,9DAA33-000000", Sim: 0.90, Dir: 0}
var 关闭背包 = &FMColor{Name: "关闭背包", X1: 1249, Y1: 67, X2: 1274, Y2: 119, MainColor: "DDEEFF-000000", OffsetColor: "6,0,DDEEFF-000000,10,-3,082637-000000,0,4,DDFFFF-000000,4,1,DDFFFF-000000,10,1,112F3B-000000,1,5,DDFFFF-000000,6,5,DDFFFF-000000,10,8,153340-000000", Sim: 0.90, Dir: 0}

type 买卖物品操作 string

const (
	买卖物品操作查找 买卖物品操作 = "find"
	买卖物品操作点击 买卖物品操作 = "click"
	买卖物品操作双击 买卖物品操作 = "double-click"
)

type 买卖物品动作 struct {
	特征 interface{}
	操作 买卖物品操作
}

type 买卖物品检查 struct {
	特征   interface{}
	应该找到 bool
}

type 买卖物品逻辑 struct {
	名称 string
	当前 买卖物品动作
	或者 *买卖物品动作
	检查 []买卖物品检查
}

var (
	买卖物品锁      sync.Mutex
	买卖物品已启动    bool
	买卖物品当前逻辑   int
	买卖物品流程截止时间 time.Time
	卖物品测试已启动   bool
	卖物品测试当前逻辑  int
	买卖物品触摸按下时长 = 35 * time.Millisecond
	买卖物品连点间隔   = 80 * time.Millisecond
	买卖物品双击间隔   = 350 * time.Millisecond
	买卖物品猫猫双击间隔 = 120 * time.Millisecond
	买卖物品猫猫检查等待 = 1200 * time.Millisecond
	买卖物品点击后等待  = 260 * time.Millisecond
	买卖物品流程超时时长 = 2 * time.Minute
	单卖格子固定X    = 856
	单卖第一格Y     = 327
	单卖格子间隔     = 53
	单卖格子数量     = 5
	单卖空格扫描X1   = 641
	单卖空格扫描Y1   = 295
	单卖空格扫描X2   = 701
	单卖空格扫描Y2   = 574
	单卖空格检测半高   = 29
	单卖空包判断空格数  = 4
	单卖双击后确认前等待 = 200 * time.Millisecond
	单卖确认等待     = 800 * time.Millisecond
	单卖确认后等待    = 450 * time.Millisecond
	单卖循环最大次数   = 50
)

var 买卖物品逻辑表 = []买卖物品逻辑{
	{
		名称: "第1个逻辑",
		当前: 买卖物品动作{特征: 背包特征, 操作: 买卖物品操作点击},
		检查: []买卖物品检查{{特征: 打开背包成功, 应该找到: true}},
	},
	{
		名称: "第2个逻辑",
		当前: 买卖物品动作{特征: 第五页灰色, 操作: 买卖物品操作点击},
		或者: &买卖物品动作{特征: 第五页面粉色, 操作: 买卖物品操作查找},
		检查: []买卖物品检查{{特征: 找到第五页面, 应该找到: true}},
	},
	{
		名称: "第3个逻辑",
		当前: 买卖物品动作{特征: 双击猫猫图标, 操作: 买卖物品操作双击},
		检查: []买卖物品检查{{特征: 找到全卖按钮, 应该找到: true}},
	},
	{
		名称: "第4个逻辑",
		当前: 买卖物品动作{特征: 全卖页面灰色, 操作: 买卖物品操作点击},
		或者: &买卖物品动作{特征: 全卖页面粉色, 操作: 买卖物品操作查找},
		检查: []买卖物品检查{{特征: 找到全卖页面, 应该找到: true}},
	},
	{
		名称: "第5个逻辑",
		当前: 买卖物品动作{特征: 空包袱, 操作: 买卖物品操作查找},
		检查: []买卖物品检查{{特征: 全卖空包袱, 应该找到: true}},
	},
	{
		名称: "第6个逻辑",
		当前: 买卖物品动作{特征: 包袱不是空的, 操作: 买卖物品操作点击},
		或者: &买卖物品动作{特征: 全卖页面, 操作: 买卖物品操作查找},
		检查: []买卖物品检查{{特征: 可以点击全卖, 应该找到: true}},
	},
	{
		名称: "第7个逻辑",
		当前: 买卖物品动作{特征: 点击全卖按钮, 操作: 买卖物品操作点击},
		检查: []买卖物品检查{{特征: 全卖确定按钮, 应该找到: true}},
	},
	{
		名称: "第8个逻辑",
		当前: 买卖物品动作{特征: 点击全卖确定, 操作: 买卖物品操作点击},
		检查: []买卖物品检查{{特征: 全卖确定消失, 应该找到: false}},
	},
	{
		名称: "第9个逻辑",
		当前: 买卖物品动作{特征: 单卖页面灰色, 操作: 买卖物品操作点击},
		或者: &买卖物品动作{特征: 已经切换单卖界面, 操作: 买卖物品操作查找},
		检查: []买卖物品检查{{特征: 已切换单卖界面, 应该找到: true}},
	},
	{
		名称: "第10个逻辑",
		当前: 买卖物品动作{特征: 单卖第一个格有物品, 操作: 买卖物品操作查找},
	},
	{
		名称: "第11个逻辑",
		当前: 买卖物品动作{特征: 关闭商店, 操作: 买卖物品操作点击},
	},
	{
		名称: "第12个逻辑",
		当前: 买卖物品动作{特征: 关闭背包, 操作: 买卖物品操作点击},
	},
}

func 启动买卖物品流程() {
	买卖物品锁.Lock()
	买卖物品已启动 = true
	买卖物品当前逻辑 = 0
	买卖物品锁.Unlock()
	输出("买卖物品流程开始", "逻辑数量=", len(买卖物品逻辑表))
}

func 买卖物品下一步() {
	index, logic, ok := 取买卖物品当前逻辑()
	if !ok {
		return
	}
	if ok, next := 执行买卖物品逻辑并返回下一步(index, logic); ok {
		设置买卖物品下一逻辑(next)
	}
}

func 买卖物品流程运行中() bool {
	买卖物品锁.Lock()
	defer 买卖物品锁.Unlock()
	return 买卖物品已启动
}

func 启动卖物品测试流程() {
	买卖物品锁.Lock()
	卖物品测试已启动 = true
	卖物品测试当前逻辑 = 0
	买卖物品锁.Unlock()
	输出("卖物品测试流程开始", "逻辑数量=", len(买卖物品逻辑表))
}

func 卖物品测试下一步() {
	index, logic, ok := 取卖物品测试当前逻辑()
	if !ok {
		return
	}
	if ok, next := 执行买卖物品逻辑并返回下一步(index, logic); ok {
		设置卖物品测试下一逻辑(next)
	}
}

func 卖物品测试流程运行中() bool {
	买卖物品锁.Lock()
	defer 买卖物品锁.Unlock()
	return 卖物品测试已启动
}

func 执行完整买卖物品流程(shouldContinue func() bool) bool {
	return 执行完整买卖物品流程选项(shouldContinue, true)
}

func 执行完整买卖物品流程选项(shouldContinue func() bool, 卖杂物 bool) bool {
	if shouldContinue == nil {
		shouldContinue = func() bool { return true }
	}
	开始买卖物品流程计时()
	defer 清除买卖物品流程计时()
	for index := 0; index < len(买卖物品逻辑表); {
		if 买卖物品流程已超时() {
			执行买卖物品超时收尾()
			return true
		}
		if !shouldContinue() {
			输出("买卖物品流程中断")
			return false
		}
		logic := 买卖物品逻辑表[index]
		ok, next := 执行买卖物品逻辑并返回下一步(index, logic)
		if !ok {
			输出("买卖物品流程失败", "逻辑=", logic.名称)
			return false
		}
		if next <= index {
			输出("买卖物品流程失败", "逻辑=", logic.名称, "原因=下一步异常", "当前=", index, "下一步=", next)
			return false
		}
		if !卖杂物 && logic.名称 == "第5个逻辑" && next == 买卖物品逻辑索引("第9个逻辑", next) {
			next = 买卖物品逻辑索引("第11个逻辑", next)
			输出("买卖物品 卖杂物未勾选：空包袱后跳过单卖，进入关闭流程", "下一步=", next+1)
		}
		if !卖杂物 && logic.名称 == "第8个逻辑" {
			next = 买卖物品逻辑索引("第11个逻辑", next)
			输出("买卖物品 卖杂物未勾选：全卖后跳过单卖，进入关闭流程", "下一步=", next+1)
		}
		index = next
		if 买卖物品流程已超时() {
			执行买卖物品超时收尾()
			return true
		}
	}
	输出("买卖物品流程结束", "耗时秒=", int(买卖物品流程已用时()/time.Second))
	return true
}

func 执行买卖物品超时收尾() {
	输出("买卖物品流程超时，关闭商店和背包，准备开始打怪", "耗时秒=", int(买卖物品流程已用时()/time.Second), "超时秒=", int(买卖物品流程超时时长/time.Second))
	尝试点击买卖物品特征("超时关闭商店", 关闭商店)
	尝试点击买卖物品特征("超时关闭背包", 关闭背包)
}

func 开始买卖物品流程计时() {
	买卖物品锁.Lock()
	买卖物品流程截止时间 = time.Now().Add(买卖物品流程超时时长)
	买卖物品锁.Unlock()
	输出("买卖物品流程计时开始", "超时秒=", int(买卖物品流程超时时长/time.Second))
}

func 清除买卖物品流程计时() {
	买卖物品锁.Lock()
	买卖物品流程截止时间 = time.Time{}
	买卖物品锁.Unlock()
}

func 买卖物品流程已超时() bool {
	买卖物品锁.Lock()
	deadline := 买卖物品流程截止时间
	买卖物品锁.Unlock()
	return !deadline.IsZero() && time.Now().After(deadline)
}

func 买卖物品流程已用时() time.Duration {
	买卖物品锁.Lock()
	deadline := 买卖物品流程截止时间
	买卖物品锁.Unlock()
	if deadline.IsZero() {
		return 0
	}
	elapsed := 买卖物品流程超时时长 - time.Until(deadline)
	if elapsed < 0 {
		return 0
	}
	return elapsed
}

func 取买卖物品当前逻辑() (int, 买卖物品逻辑, bool) {
	买卖物品锁.Lock()
	defer 买卖物品锁.Unlock()
	if !买卖物品已启动 {
		输出("买卖物品未开始")
		return 0, 买卖物品逻辑{}, false
	}
	if 买卖物品当前逻辑 >= len(买卖物品逻辑表) {
		买卖物品已启动 = false
		输出("买卖物品流程结束")
		return 0, 买卖物品逻辑{}, false
	}
	return 买卖物品当前逻辑, 买卖物品逻辑表[买卖物品当前逻辑], true
}

func 设置买卖物品下一逻辑(next int) {
	买卖物品锁.Lock()
	defer 买卖物品锁.Unlock()
	买卖物品当前逻辑 = next
	if 买卖物品当前逻辑 >= len(买卖物品逻辑表) {
		买卖物品已启动 = false
		输出("买卖物品流程结束")
	}
}

func 取卖物品测试当前逻辑() (int, 买卖物品逻辑, bool) {
	买卖物品锁.Lock()
	defer 买卖物品锁.Unlock()
	if !卖物品测试已启动 {
		输出("卖物品测试未开始")
		return 0, 买卖物品逻辑{}, false
	}
	if 卖物品测试当前逻辑 >= len(买卖物品逻辑表) {
		卖物品测试已启动 = false
		输出("卖物品测试流程结束")
		return 0, 买卖物品逻辑{}, false
	}
	return 卖物品测试当前逻辑, 买卖物品逻辑表[卖物品测试当前逻辑], true
}

func 设置卖物品测试下一逻辑(next int) {
	买卖物品锁.Lock()
	defer 买卖物品锁.Unlock()
	卖物品测试当前逻辑 = next
	if 卖物品测试当前逻辑 >= len(买卖物品逻辑表) {
		卖物品测试已启动 = false
		输出("卖物品测试流程结束")
	}
}

func 执行买卖物品逻辑(logic 买卖物品逻辑) bool {
	ok, _ := 执行买卖物品逻辑并返回下一步(-1, logic)
	return ok
}

func 执行买卖物品逻辑并返回下一步(index int, logic 买卖物品逻辑) (bool, int) {
	next := index + 1
	if index < 0 {
		next = 0
	}

	if logic.名称 == "第3个逻辑" {
		return 执行猫猫双击逻辑(index, logic)
	}
	if logic.名称 == "第5个逻辑" {
		return 执行空包袱分支逻辑(index, logic)
	}
	if logic.名称 == "第9个逻辑" {
		return 执行单卖页切换逻辑(index, logic)
	}
	if logic.名称 == "第10个逻辑" {
		return 执行单卖循环逻辑(index, logic)
	}

	输出("买卖物品 执行", logic.名称)
	if !执行买卖物品动作(logic.名称, logic.当前) {
		if logic.或者 == nil || !执行买卖物品动作(logic.名称+" 或者", *logic.或者) {
			输出("买卖物品 逻辑失败", "逻辑=", logic.名称, "原因=当前特征和或者特征都未成功")
			return false, next
		}
	}
	for _, check := range logic.检查 {
		if !执行买卖物品检查(logic.名称, check) {
			return false, next
		}
	}
	输出("买卖物品 逻辑成功", logic.名称)
	return true, next
}

func 执行空包袱分支逻辑(index int, logic 买卖物品逻辑) (bool, int) {
	输出("买卖物品 执行", logic.名称, "规则=空包袱则跳到第9个逻辑")
	found, x, y := 查找买卖物品特征(空包袱)
	if found {
		next := 买卖物品逻辑索引("第9个逻辑", index+1)
		输出("买卖物品 空包袱分支", "找到空包袱", "x=", x, "y=", y, "跳到=", next+1)
		return true, next
	}
	输出("买卖物品 空包袱分支", "未找到空包袱，继续第6个逻辑")
	return true, index + 1
}

func 执行猫猫双击逻辑(index int, logic 买卖物品逻辑) (bool, int) {
	next := index + 1
	if index < 0 {
		next = 0
	}

	name := 买卖物品特征名(双击猫猫图标)
	found, x, y := 查找买卖物品特征(双击猫猫图标)
	if !found {
		输出("买卖物品 找不到", "逻辑=", logic.名称, "特征=", name)
		return false, next
	}

	points := []struct {
		label string
		x     int
		y     int
	}{
		{label: "图标中心", x: x + 15, y: y + 8},
		{label: "识别点", x: x, y: y},
		{label: "偏下中心", x: x + 15, y: y + 18},
	}

	for _, point := range points {
		输出("买卖物品 猫猫双击", "逻辑=", logic.名称, "点位=", point.label, "识别x=", x, "识别y=", y, "点击x=", point.x, "点击y=", point.y)
		快速双击买卖物品坐标(point.x, point.y, 买卖物品猫猫双击间隔)

		checkOK := true
		for _, check := range logic.检查 {
			if !等待执行买卖物品检查(logic.名称, check, 买卖物品猫猫检查等待) {
				checkOK = false
				break
			}
		}
		if checkOK {
			输出("买卖物品 逻辑成功", logic.名称)
			return true, next
		}
	}

	输出("买卖物品 逻辑失败", "逻辑=", logic.名称, "原因=猫猫双击后未出现全卖按钮")
	return false, next
}

func 执行单卖页切换逻辑(index int, logic 买卖物品逻辑) (bool, int) {
	next := index + 1
	if index < 0 {
		next = 0
	}

	if found, x, y := 查找买卖物品特征(已切换单卖界面); found {
		输出("买卖物品 单卖页切换", "已经变色", "x=", x, "y=", y)
		return true, next
	}
	if found, x, y := 查找买卖物品特征(已经切换单卖界面); found {
		输出("买卖物品 单卖页切换", "已经切换", "x=", x, "y=", y)
		return true, next
	}

	name := 买卖物品特征名(单卖页面灰色)
	found, x, y := 查找买卖物品特征(单卖页面灰色)
	if !found {
		输出("买卖物品 找不到", "逻辑=", logic.名称, "特征=", name)
		return false, next
	}

	输出("买卖物品 单卖页单击", "逻辑=", logic.名称, "特征=", name, "x=", x, "y=", y)
	快速触摸买卖物品坐标(x, y)
	time.Sleep(买卖物品点击后等待)
	if 买卖物品等待检查全部通过(logic.名称, logic.检查, 买卖物品猫猫检查等待) {
		输出("买卖物品 逻辑成功", logic.名称)
		return true, next
	}

	输出("买卖物品 单卖页双击前单击确认", "逻辑=", logic.名称, "x=", x, "y=", y)
	快速触摸买卖物品坐标(x, y)
	time.Sleep(买卖物品点击后等待)
	if !买卖物品等待检查全部通过(logic.名称, logic.检查, 买卖物品猫猫检查等待) {
		输出("买卖物品 单卖页双击取消", "逻辑=", logic.名称, "原因=单击后未变色")
		return false, next
	}

	输出("买卖物品 单卖页双击兜底", "逻辑=", logic.名称, "x=", x, "y=", y)
	快速双击买卖物品坐标(x, y, 买卖物品猫猫双击间隔)
	if 买卖物品等待检查全部通过(logic.名称, logic.检查, 买卖物品猫猫检查等待) {
		输出("买卖物品 逻辑成功", logic.名称)
		return true, next
	}

	输出("买卖物品 逻辑失败", "逻辑=", logic.名称, "原因=单卖页未变色")
	return false, next
}

func 执行单卖循环逻辑(index int, logic 买卖物品逻辑) (bool, int) {
	next := index + 1
	输出("买卖物品 执行", logic.名称, "规则=先扫空格子，空格>=4判断空包，只卖非空格子")

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	for attempt := 1; attempt <= 单卖循环最大次数; attempt++ {
		if 买卖物品流程已超时() {
			输出("买卖物品 单卖循环超时，跳过", "次数=", attempt-1)
			return true, next
		}
		emptyCount, sellSlots := 扫描单卖空格子()
		if emptyCount >= 单卖空包判断空格数 {
			输出("买卖物品 单卖循环成功", "空格子=", emptyCount, "/", 单卖格子数量, "非空格子=", sellSlots, "次数=", attempt-1)
			return true, next
		}
		if len(sellSlots) == 0 {
			输出("买卖物品 单卖循环失败", "未找到可卖格子", "空格子=", emptyCount, "/", 单卖格子数量)
			return false, next
		}

		slot := sellSlots[random.Intn(len(sellSlots))]
		clickX := 单卖格子固定X
		clickY := 单卖第一格Y + (slot-1)*单卖格子间隔
		输出("买卖物品 单卖随机双击", "次数=", attempt, "格子=", slot, "可卖格子=", sellSlots, "x=", clickX, "y=", clickY, "空格子=", emptyCount, "/", 单卖格子数量)
		快速双击买卖物品坐标不等待(clickX, clickY, 买卖物品猫猫双击间隔)
		time.Sleep(单卖双击后确认前等待)

		confirmFound, confirmX, confirmY, confirmName := 等待单卖确认卖出按钮(单卖确认等待)
		if !confirmFound {
			输出("买卖物品 单卖无确认", "格子=", slot, "下次重新扫描空格")
			continue
		}

		输出("买卖物品 单卖确认卖出", "格子=", slot, "特征=", confirmName, "x=", confirmX, "y=", confirmY)
		点击买卖物品坐标(confirmX, confirmY, 1)
		time.Sleep(单卖确认后等待)
	}

	emptyCount, sellSlots := 扫描单卖空格子()
	输出("买卖物品 单卖循环失败", "超过最大次数", "次数=", 单卖循环最大次数, "空格子=", emptyCount, "/", 单卖格子数量, "非空格子=", sellSlots)
	return false, next
}

func 等待买卖物品特征(feature interface{}, timeout time.Duration) (bool, int, int) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if found, x, y := 查找买卖物品特征(feature); found {
			return true, x, y
		}
		time.Sleep(80 * time.Millisecond)
	}
	return 查找买卖物品特征(feature)
}

func 等待单卖确认卖出按钮(timeout time.Duration) (bool, int, int, string) {
	deadline := time.Now().Add(timeout)
	for {
		if found, x, y := 查找买卖物品特征(单卖确认卖出); found {
			return true, x, y, 买卖物品特征名(单卖确认卖出)
		}
		if found, x, y := 查找买卖物品特征(单卖确认卖出备用); found {
			return true, x, y, 买卖物品特征名(单卖确认卖出备用)
		}
		if found, x, y := 查找买卖物品特征(双击单卖之后的确定单击按钮); found {
			return true, x, y, 买卖物品特征名(双击单卖之后的确定单击按钮)
		}
		if found, x, y := 查找买卖物品特征(MS系统应用); found {
			return true, x, y, 买卖物品特征名(MS系统应用)
		}
		if time.Now().After(deadline) {
			return false, -1, -1, ""
		}
		time.Sleep(80 * time.Millisecond)
	}
}

func 买卖物品等待检查全部通过(logicName string, checks []买卖物品检查, timeout time.Duration) bool {
	for _, check := range checks {
		if !等待执行买卖物品检查(logicName, check, timeout) {
			return false
		}
	}
	return true
}

func 等待执行买卖物品检查(logicName string, check 买卖物品检查, timeout time.Duration) bool {
	name := 买卖物品特征名(check.特征)
	if check.应该找到 {
		found, x, y := 等待买卖物品特征(check.特征, timeout)
		if found {
			输出("买卖物品 检查成功", "逻辑=", logicName, "找到=", name, "x=", x, "y=", y)
			return true
		}
		输出("买卖物品 检查失败", "逻辑=", logicName, "未找到=", name)
		return false
	}

	deadline := time.Now().Add(timeout)
	for {
		found, x, y := 查找买卖物品特征(check.特征)
		if !found {
			输出("买卖物品 检查成功", "逻辑=", logicName, "没有找到=", name)
			return true
		}
		if time.Now().After(deadline) {
			输出("买卖物品 检查失败", "逻辑=", logicName, "不该找到但找到了=", name, "x=", x, "y=", y)
			return false
		}
		time.Sleep(80 * time.Millisecond)
	}
}

func 统计单卖空格子数量() int {
	count, _ := 扫描单卖空格子()
	return count
}

func 扫描单卖空格子() (int, []int) {
	sellSlots := make([]int, 0, 单卖格子数量)
	if 引擎 == nil || 单卖空格子颜色特征 == nil {
		for slot := 1; slot <= 单卖格子数量; slot++ {
			sellSlots = append(sellSlots, slot)
		}
		return 0, sellSlots
	}

	x1, y1, x2, y2 := 引擎.scaleRect(单卖空格扫描X1, 单卖空格扫描Y1, 单卖空格扫描X2, 单卖空格扫描Y2)
	引擎.markScreenRect(x1, y1, x2, y2)
	colors := buildMultiColor(单卖空格子颜色特征.MainColor, 单卖空格子颜色特征.OffsetColor)
	points := images.FindMultiColorsAll(x1, y1, x2, y2, colors, defaultSim(单卖空格子颜色特征.Sim), 单卖空格子颜色特征.Dir, 引擎.displayID())

	emptySlots := make(map[int]bool, 单卖格子数量)
	for _, point := range points {
		slot := 单卖空格点所在格子(引擎.unscaleY(point.Y))
		if slot > 0 {
			emptySlots[slot] = true
		}
	}

	for slot := 1; slot <= 单卖格子数量; slot++ {
		if !emptySlots[slot] {
			sellSlots = append(sellSlots, slot)
		}
	}

	输出("买卖物品 单卖空格扫描", "空格子=", len(emptySlots), "/", 单卖格子数量, "非空格子=", sellSlots, "原始点=", len(points))
	return len(emptySlots), sellSlots
}

func 单卖空格点所在格子(y int) int {
	for slot := 1; slot <= 单卖格子数量; slot++ {
		centerY := 单卖第一格Y + (slot-1)*单卖格子间隔
		if y >= centerY-单卖空格检测半高 && y <= centerY+单卖空格检测半高 {
			return slot
		}
	}
	return 0
}

func 买卖物品逻辑索引(name string, fallback int) int {
	for i, logic := range 买卖物品逻辑表 {
		if logic.名称 == name {
			return i
		}
	}
	return fallback
}

func 执行买卖物品动作(logicName string, action 买卖物品动作) bool {
	name := 买卖物品特征名(action.特征)
	found, x, y := 查找买卖物品特征(action.特征)
	if !found {
		输出("买卖物品 找不到", "逻辑=", logicName, "特征=", name)
		return false
	}
	输出("买卖物品 找到", "逻辑=", logicName, "特征=", name, "x=", x, "y=", y, "操作=", string(action.操作))
	switch action.操作 {
	case 买卖物品操作点击:
		点击买卖物品坐标(x, y, 1)
	case 买卖物品操作双击:
		快速双击买卖物品坐标(x, y, 买卖物品双击间隔)
	}
	return true
}

func 执行买卖物品检查(logicName string, check 买卖物品检查) bool {
	name := 买卖物品特征名(check.特征)
	found, x, y := 查找买卖物品特征(check.特征)
	if check.应该找到 {
		if found {
			输出("买卖物品 检查成功", "逻辑=", logicName, "找到=", name, "x=", x, "y=", y)
			return true
		}
		输出("买卖物品 检查失败", "逻辑=", logicName, "未找到=", name)
		return false
	}
	if found {
		输出("买卖物品 检查失败", "逻辑=", logicName, "不该找到但找到了=", name, "x=", x, "y=", y)
		return false
	}
	输出("买卖物品 检查成功", "逻辑=", logicName, "没有找到=", name)
	return true
}

func 查找买卖物品特征(feature interface{}) (bool, int, int) {
	if 引擎 == nil || feature == nil {
		return false, -1, -1
	}
	return 引擎.FindFeature(feature)
}

func 买卖物品特征名(feature interface{}) string {
	switch f := feature.(type) {
	case *FMColor:
		if f != nil {
			return f.Name
		}
	case *CColor:
		if f != nil {
			return f.Name
		}
	}
	return "未知特征"
}

func 点击买卖物品坐标(x, y, count int) {
	点击买卖物品坐标带间隔(x, y, count, 买卖物品连点间隔)
}

func 点击买卖物品坐标带间隔(x, y, count int, interval time.Duration) {
	if count < 1 {
		count = 1
	}
	for i := 0; i < count; i++ {
		if 引擎 != nil {
			引擎.ClickResult(true, x, y)
		} else {
			手动点击买卖物品坐标(x, y)
		}
		if i+1 < count {
			time.Sleep(interval)
		}
	}
	time.Sleep(买卖物品点击后等待)
}

func 匹配买卖物品特征(feature interface{}) bool {
	if c, ok := feature.(*CColor); ok {
		ok, _, _ := 静默匹配买卖物品特征(c)
		return ok
	}
	ok, _, _ := 查找买卖物品特征(feature)
	return ok
}

func 点击买卖物品特征(feature interface{}) bool {
	if 引擎 == nil || feature == nil {
		return false
	}
	if c, ok := feature.(*CColor); ok {
		found, x, y := 静默匹配买卖物品特征(c)
		if !found {
			return false
		}
		手动点击买卖物品坐标(x, y)
		return true
	}
	ok, x, y := 查找买卖物品特征(feature)
	if !ok {
		return false
	}
	手动点击买卖物品坐标(x, y)
	return true
}

func 尝试点击买卖物品特征(label string, feature interface{}) bool {
	name := 买卖物品特征名(feature)
	found, x, y := 查找买卖物品特征(feature)
	if !found {
		输出("买卖物品", label, "未找到，跳过", "特征=", name)
		return false
	}
	输出("买卖物品", label, "点击", "特征=", name, "x=", x, "y=", y)
	点击买卖物品坐标(x, y, 1)
	return true
}

func 双击买卖物品坐标(x, y int) {
	快速双击买卖物品坐标(x, y, 买卖物品双击间隔)
}

func 快速双击买卖物品坐标(x, y int, interval time.Duration) {
	快速双击买卖物品坐标不等待(x, y, interval)
	time.Sleep(买卖物品点击后等待)
}

func 快速双击买卖物品坐标不等待(x, y int, interval time.Duration) {
	快速触摸买卖物品坐标(x, y)
	time.Sleep(interval)
	快速触摸买卖物品坐标(x, y)
}

func 快速触摸买卖物品坐标(x, y int) {
	if 引擎 == nil {
		return
	}
	displayID := 引擎.displayID()
	sx := 引擎.scaleX(x + 引擎.offsetX)
	sy := 引擎.scaleY(y + 引擎.offsetY)
	引擎.markScreenPoint(sx, sy)
	motion.TouchDown(sx, sy, 0, displayID)
	time.Sleep(买卖物品触摸按下时长)
	motion.TouchUp(sx, sy, 0, displayID)
}

func 手动点击买卖物品坐标(x, y int) {
	if 引擎 == nil {
		return
	}
	displayID := 引擎.displayID()
	sx := 引擎.scaleX(x + 引擎.offsetX)
	sy := 引擎.scaleY(y + 引擎.offsetY)
	引擎.markScreenPoint(sx, sy)
	motion.TouchDown(sx, sy, 0, displayID)
	time.Sleep(买卖物品触摸按下时长)
	motion.TouchUp(sx, sy, 0, displayID)
	time.Sleep(买卖物品连点间隔)
}

func 手动点击买卖物品绝对坐标(x, y int) {
	if 引擎 == nil {
		return
	}
	displayID := 引擎.displayID()
	sx := 引擎.scaleX(x)
	sy := 引擎.scaleY(y)
	引擎.markScreenPoint(sx, sy)
	motion.TouchDown(sx, sy, 0, displayID)
	time.Sleep(买卖物品触摸按下时长)
	motion.TouchUp(sx, sy, 0, displayID)
	time.Sleep(买卖物品连点间隔)
}

func 静默匹配买卖物品特征(feature *CColor) (bool, int, int) {
	if 引擎 == nil || feature == nil {
		return false, -1, -1
	}
	x := feature.X + 引擎.offsetX
	y := feature.Y + 引擎.offsetY
	if strings.Contains(feature.Color, ",") {
		colors := buildDetectColorsFromCColor(feature, 引擎.offsetX, 引擎.offsetY, 引擎.displayID())
		if images.DetectsMultiColors(colors, defaultSim(feature.Sim), 引擎.displayID()) {
			return true, x, y
		}
		return false, -1, -1
	}
	if images.CmpColor(引擎.scaleX(x), 引擎.scaleY(y), feature.Color, defaultSim(feature.Sim), 引擎.displayID()) {
		return true, x, y
	}
	return false, -1, -1
}

func 是买卖物品日志(text string) bool {
	return strings.HasPrefix(text, "ms") ||
		strings.HasPrefix(text, "买卖物品") ||
		strings.HasPrefix(text, "买东西") ||
		strings.Contains(text, "第")
}
