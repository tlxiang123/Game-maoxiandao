package main

const 小地图黄点颜色 = "FFFF06-101010|FBED0A-101010|FDF301-101010|FFFF00-101010|FEF008-101010|E7E31C-101010"

var 海盗小地图黄点区域 = &FColor{Name: "海盗小地图黄点区域", X1: 26, Y1: 97, X2: 169, Y2: 204, Color: 小地图黄点颜色, Sim: 0.70, Dir: 0}
var 僵尸3小地图黄点区域 = &FColor{Name: "僵尸3小地图黄点区域", X1: 10, Y1: 97, X2: 260, Y2: 203, Color: 小地图黄点颜色, Sim: 0.70, Dir: 0}

var 小地图黄点候选区域 = []*FColor{
	海盗小地图黄点区域,
}
