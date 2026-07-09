//go:build !android

package main

func 打印怪物层统计(位置 层位置) {
}

func 打印怪物层统计并取当前层数量(位置 层位置) (int, bool) {
	return 0, false
}

func YOLO已识别到怪物() bool {
	return false
}

func YOLO已加载完成() bool {
	return false
}

func YOLO人物附近有怪物(maxDistance int) bool {
	return false
}

func YOLO当前层附近有怪物(位置 层位置, maxDistance int) bool {
	return false
}

func YOLO当前层怪物方向(位置 层位置) int {
	return 0
}

func YOLO人物中心Y即时() (int, bool) {
	return 0, false
}

func 启动怪物识别后台(runID int64) {
}

func 关闭怪物识别() {
}
