package main

import (
	"testing"
	"time"
)

func Test是恐龙洞玩家黄点颜色(t *testing.T) {
	if !是恐龙洞玩家黄点颜色(0xFF, 0xF7, 0x29) {
		t.Fatal("玩家亮黄色核心应该匹配")
	}
	if 是恐龙洞玩家黄点颜色(0xC5, 0xB4, 0x19) {
		t.Fatal("较暗的洞穴金色地形不应该匹配")
	}
}

func Test是恐龙洞玩家黄点区域(t *testing.T) {
	if !是恐龙洞玩家黄点区域(黄点连通区域{Count: 48, MinX: 0, MinY: 0, MaxX: 8, MaxY: 8}) {
		t.Fatal("接近菱形的玩家黄点区域应该匹配")
	}
	if 是恐龙洞玩家黄点区域(黄点连通区域{Count: 60, MinX: 0, MinY: 0, MaxX: 29, MaxY: 3}) {
		t.Fatal("细长平台区域不应该匹配")
	}
}

func Test识别恐龙洞层数(t *testing.T) {
	tests := []struct {
		name    string
		x       int
		y       int
		layer   int
		matched bool
	}{
		{name: "一层左侧高地", x: 30, y: 169, layer: 1, matched: true},
		{name: "一层在二层以下", x: 57, y: 180, layer: 1, matched: true},
		{name: "一层右侧低地", x: 120, y: 185, layer: 1, matched: true},
		{name: "一层越过右侧刷怪边界", x: 155, y: 181, layer: 1, matched: true},
		{name: "上层不再使用绝对Y", x: 100, y: 149, layer: 0, matched: false},
		{name: "范围外不匹配", x: 190, y: 185, layer: 0, matched: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			layer, matched := 识别恐龙洞层数(test.x, test.y)
			if layer != test.layer || matched != test.matched {
				t.Fatalf("识别恐龙洞层数(%d, %d) = (%d, %v), want (%d, %v)", test.x, test.y, layer, matched, test.layer, test.matched)
			}
		})
	}
}

func Test识别恐龙洞相对层数(t *testing.T) {
	tests := []struct {
		diff  int
		layer int
	}{
		{diff: 46, layer: 2},
		{diff: 45, layer: 2},
		{diff: 25, layer: 3},
		{diff: 0, layer: 4},
		{diff: -1, layer: 4},
		{diff: -13, layer: 5},
	}
	for _, test := range tests {
		layer, ok := 识别恐龙洞相对层数(test.diff)
		if !ok || layer != test.layer {
			t.Fatalf("识别恐龙洞相对层数(%d) = (%d, %v), want (%d, true)", test.diff, layer, ok, test.layer)
		}
	}
}

func Test恐龙洞巡逻边界越界后直接到位(t *testing.T) {
	if !恐龙洞巡逻目标已到位(20, 23, 23, 148) {
		t.Fatal("越过左边界后应该直接到位，不应反向小步修正")
	}
	if !恐龙洞巡逻目标已到位(151, 148, 23, 148) {
		t.Fatal("越过右边界后应该直接到位，不应反向小步修正")
	}
	if 恐龙洞巡逻目标已到位(100, 148, 23, 148) {
		t.Fatal("距离右边界较远时不应该判定到位")
	}
}

func Test恐龙洞各层最低爬梯时间(t *testing.T) {
	tests := []struct {
		startLayer int
		want       time.Duration
	}{
		{startLayer: 1, want: 2000 * time.Millisecond},
		{startLayer: 2, want: 2600 * time.Millisecond},
		{startLayer: 3, want: 2600 * time.Millisecond},
		{startLayer: 4, want: 2000 * time.Millisecond},
	}
	for _, test := range tests {
		if got := 恐龙洞最低爬梯持续(test.startLayer); got != test.want {
			t.Fatalf("%d层开始爬梯最低时间=%s, want %s", test.startLayer, got, test.want)
		}
	}
}

func Test恐龙洞所有梯子使用精确X(t *testing.T) {
	wantX := map[int]int{1: 88, 2: 96, 3: 101, 4: 113}
	for layer, x := range wantX {
		config, ok := 取恐龙洞层配置(layer)
		if !ok {
			t.Fatalf("缺少%d层梯子配置", layer)
		}
		if config.上梯子左X != x || config.上梯子右X != x {
			t.Fatalf("%d层梯子范围=%d-%d, want %d-%d", layer, config.上梯子左X, config.上梯子右X, x, x)
		}
	}
}
