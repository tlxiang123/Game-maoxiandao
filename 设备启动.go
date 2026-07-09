package main

import (
	"strings"

	"github.com/Dasongzi1366/AutoGo/app"
	"github.com/Dasongzi1366/AutoGo/device"
	"github.com/Dasongzi1366/AutoGo/vdisplay"
)

var 当前虚拟屏 *vdisplay.Vdisplay

func 初始化运行屏幕() {
	if !当前UI配置.使用虚拟屏 || 当前是模拟器() {
		屏幕ID = 0
		输出("使用主屏幕", "displayId=", 屏幕ID)
		return
	}

	输出("创建虚拟屏", 当前UI配置.虚拟屏宽, 当前UI配置.虚拟屏高, 当前UI配置.虚拟屏DPI)
	当前虚拟屏 = vdisplay.Create(当前UI配置.虚拟屏宽, 当前UI配置.虚拟屏高, 当前UI配置.虚拟屏DPI)
	if 当前虚拟屏 == nil {
		屏幕ID = 0
		输出("虚拟屏创建失败，回退主屏幕")
		return
	}

	屏幕ID = 当前虚拟屏.GetDisplayId()
	当前虚拟屏.SetTitle("AutoGo Project")
	输出("虚拟屏创建成功", "displayId=", 屏幕ID)
}

func 当前是模拟器() bool {
	text := strings.ToLower(strings.Join([]string{
		device.CpuAbi,
		device.Brand,
		device.Device,
		device.Model,
		device.Product,
		device.Hardware,
		device.Fingerprint,
	}, " "))

	keywords := []string{"x86", "emulator", "sdk_gphone", "google_sdk", "generic", "genymotion", "bluestacks", "nox", "ldplayer", "mumu", "vbox"}
	for _, keyword := range keywords {
		if strings.Contains(text, keyword) {
			return true
		}
	}
	return false
}

func 启动应用到当前屏幕() bool {
	if 应用包名 == "" {
		输出("应用包名为空，跳过启动")
		return false
	}
	if 当前虚拟屏 != nil {
		输出("启动应用到虚拟屏", 应用包名, "displayId=", 屏幕ID)
		return 当前虚拟屏.LaunchApp(应用包名)
	}
	输出("启动应用到主屏幕", 应用包名, "displayId=", 屏幕ID)
	return app.Launch(应用包名, 屏幕ID)
}
