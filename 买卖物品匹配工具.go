package main

import "time"

type 买卖物品特征候选 struct {
	标签 string
	特征 interface{}
}

func 新买卖物品特征候选(label string, feature interface{}) 买卖物品特征候选 {
	return 买卖物品特征候选{标签: label, 特征: feature}
}

func 查找任一买卖物品特征(candidates ...买卖物品特征候选) (bool, int, int, 买卖物品特征候选) {
	for _, candidate := range candidates {
		if candidate.特征 == nil {
			continue
		}
		if found, x, y := 查找买卖物品特征(candidate.特征); found {
			if candidate.标签 == "" {
				candidate.标签 = 买卖物品特征名(candidate.特征)
			}
			return true, x, y, candidate
		}
	}
	return false, -1, -1, 买卖物品特征候选{}
}

func 等待任一买卖物品特征(timeout time.Duration, candidates ...买卖物品特征候选) (bool, int, int, 买卖物品特征候选) {
	deadline := time.Now().Add(timeout)
	for {
		if found, x, y, candidate := 查找任一买卖物品特征(candidates...); found {
			return true, x, y, candidate
		}
		if time.Now().After(deadline) {
			return false, -1, -1, 买卖物品特征候选{}
		}
		time.Sleep(80 * time.Millisecond)
	}
}

func 点击任一买卖物品特征(logicName string, candidates ...买卖物品特征候选) (bool, int, int, 买卖物品特征候选) {
	found, x, y, candidate := 查找任一买卖物品特征(candidates...)
	if !found {
		输出("买卖物品 找不到任一", "逻辑=", logicName, "候选=", 买卖物品候选标签文本(candidates))
		return false, -1, -1, 买卖物品特征候选{}
	}
	输出("买卖物品 点击任一", "逻辑=", logicName, "候选=", candidate.标签, "x=", x, "y=", y)
	点击买卖物品坐标(x, y, 1)
	return true, x, y, candidate
}

func 买卖物品候选标签文本(candidates []买卖物品特征候选) string {
	text := ""
	for i, candidate := range candidates {
		label := 买卖物品候选描述(candidate)
		if i > 0 {
			text += "|"
		}
		text += label
	}
	return text
}

func 买卖物品候选描述(candidate 买卖物品特征候选) string {
	label := candidate.标签
	if label == "" {
		label = 买卖物品特征名(candidate.特征)
	}
	switch f := candidate.特征.(type) {
	case *FMColor:
		if f != nil {
			return label + "(" + f.Name + " " + 格式化区域颜色(f.X1, f.Y1, f.X2, f.Y2, f.MainColor) + ")"
		}
	case *FColor:
		if f != nil {
			return label + "(" + f.Name + " " + 格式化区域颜色(f.X1, f.Y1, f.X2, f.Y2, f.Color) + ")"
		}
	case *CColor:
		if f != nil {
			return label + "(" + f.Name + " x=" + itoa(f.X) + " y=" + itoa(f.Y) + " color=" + f.Color + ")"
		}
	}
	return label
}

func 格式化区域颜色(x1, y1, x2, y2 int, color string) string {
	return "rect=" + itoa(x1) + "," + itoa(y1) + "," + itoa(x2) + "," + itoa(y2) + " color=" + color
}

func itoa(v int) string {
	if v == 0 {
		return "0"
	}
	negative := v < 0
	if negative {
		v = -v
	}
	buf := [20]byte{}
	i := len(buf)
	for v > 0 {
		i--
		buf[i] = byte('0' + v%10)
		v /= 10
	}
	if negative {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}
