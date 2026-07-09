package main

import (
	"strconv"
	"strings"
)

func 清理OCR文本(text string) string {
	return strings.NewReplacer(" ", "", "\n", "", "\r", "", "\t", "").Replace(text)
}

func 仅保留数字(text string) string {
	var b strings.Builder
	for _, ch := range text {
		if ch >= '0' && ch <= '9' {
			b.WriteRune(ch)
		}
	}
	return b.String()
}

func OCR区域文本(name string, x1, y1, x2, y2 int, color string) string {
	text := 引擎.PPOcrText(&PPOcrRegion{
		Name:  name,
		X1:    x1,
		Y1:    y1,
		X2:    x2,
		Y2:    y2,
		Color: color,
	})
	return 清理OCR文本(text)
}

func OCR区域数字(name string, x1, y1, x2, y2 int, color string) (int, bool) {
	text := OCR区域文本(name, x1, y1, x2, y2, color)
	digits := 仅保留数字(text)
	if digits == "" {
		return 0, false
	}
	value, err := strconv.Atoi(digits)
	if err != nil {
		return 0, false
	}
	return value, true
}
