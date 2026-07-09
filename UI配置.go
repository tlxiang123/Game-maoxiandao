package main

import "strings"

type UI配置 struct {
	应用包名    string
	启动应用    bool
	使用虚拟屏   bool
	虚拟屏宽    int
	虚拟屏高    int
	虚拟屏DPI  int
	保持屏幕常亮  bool
	启用日志    bool
	只记录成功   bool
	启用Toast bool
	点击后延迟毫秒 int
	OCR版本   string
}

var 当前UI配置 = 默认UI配置()

func 默认UI配置() UI配置 {
	return UI配置{
		应用包名:    "",
		启动应用:    false,
		使用虚拟屏:   false,
		虚拟屏宽:    1280,
		虚拟屏高:    720,
		虚拟屏DPI:  320,
		保持屏幕常亮:  true,
		启用日志:    true,
		只记录成功:   false,
		启用Toast: false,
		点击后延迟毫秒: 500,
		OCR版本:   "v5",
	}
}

func 加载UI配置骨架() UI配置 {
	cfg := 默认UI配置()
	应用UI配置(cfg)
	return cfg
}

func 应用UI配置(cfg UI配置) {
	cfg = 规范化UI配置(cfg)
	当前UI配置 = cfg
	应用包名 = cfg.应用包名
}

func 规范化UI配置(cfg UI配置) UI配置 {
	cfg.应用包名 = strings.TrimSpace(cfg.应用包名)
	if cfg.虚拟屏宽 <= 0 {
		cfg.虚拟屏宽 = 1280
	}
	if cfg.虚拟屏高 <= 0 {
		cfg.虚拟屏高 = 720
	}
	if cfg.虚拟屏DPI <= 0 {
		cfg.虚拟屏DPI = 320
	}
	if cfg.点击后延迟毫秒 < 0 {
		cfg.点击后延迟毫秒 = 0
	}
	if cfg.OCR版本 == "" {
		cfg.OCR版本 = "v5"
	}
	return cfg
}

func 输出UI配置(cfg UI配置) {
	输出("UI配置",
		"应用包名=", cfg.应用包名,
		"启动应用=", cfg.启动应用,
		"使用虚拟屏=", cfg.使用虚拟屏,
		"虚拟屏=", cfg.虚拟屏宽, "x", cfg.虚拟屏高,
		"日志=", cfg.启用日志,
		"Toast=", cfg.启用Toast,
		"OCR版本=", cfg.OCR版本,
	)
}
