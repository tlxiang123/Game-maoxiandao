package main

import (
	"embed"
	"log"

	"github.com/Dasongzi1366/AutoGo/device"
	"github.com/Dasongzi1366/AutoGo/utils"
)

var 屏幕ID = 0
var 应用包名 = ""

//go:embed pic/*.png best.ncnn.param best.ncnn.bin data.yaml
var res embed.FS

var 引擎 *Zg

func main() {
	初始化日志输出()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	输出("通用自动化工程启动")
	utils.Toast("脚本开始", -1, -1, 1000)

	应用UI配置(加载UI配置骨架())
	初始化运行屏幕()

	SetGlobalLog(当前UI配置.启用日志)
	SetGlobalToast(当前UI配置.启用Toast)
	SetGlobalLogSuccessOnly(当前UI配置.只记录成功)
	SetGlobalDisplayId(屏幕ID)

	引擎 = New(res)
	引擎.SetDisplayId(屏幕ID)
	启动测谎检测后台()
	输出UI配置(当前UI配置)

	if 当前UI配置.保持屏幕常亮 {
		device.KeepScreenOn()
	}
	if 当前UI配置.启动应用 && 应用包名 != "" {
		启动应用到当前屏幕()
	}

	输出("通用层已就绪，打开控制界面")
	运行控制界面()
	请求退出程序()
	等待脚本停止()
	关闭怪物识别()
	输出("程序已退出")
}
