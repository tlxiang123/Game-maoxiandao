package main

const 小地图黄点颜色 = "FFFF06-101010|FBED0A-101010|FDF301-101010|FFFF00-101010|FEF008-101010|E7E31C-101010"

var 海盗小地图黄点区域 = &FColor{Name: "海盗小地图黄点区域", X1: 26, Y1: 97, X2: 169, Y2: 204, Color: 小地图黄点颜色, Sim: 0.70, Dir: 0}
var 僵尸3小地图黄点区域 = &FColor{Name: "僵尸3小地图黄点区域", X1: 10, Y1: 97, X2: 260, Y2: 203, Color: 小地图黄点颜色, Sim: 0.70, Dir: 0}
var 下拉对话框 = &FMColor{Name: "下拉对话框", X1: 683, Y1: 631, X2: 748, Y2: 667, MainColor: "25BAFF-000000", OffsetColor: "5,0,249FDC-000000,10,0,25A6E0-000000,0,4,3FBEFF-000000,1,7,2CAFF6-000000,10,7,4EB1EC-000000,0,11,6BEAFF-000000,9,11,40ECFF-000000,10,11,4AE2FF-000000", Sim: 0.90, Dir: 0}

var 小地图黄点候选区域 = []*FColor{
	海盗小地图黄点区域,
}
