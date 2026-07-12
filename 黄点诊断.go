package main

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync/atomic"
	"time"

	"github.com/Dasongzi1366/AutoGo/images"
)

type 黄点扫描结果 struct {
	Ok         bool
	X          int
	Y          int
	Detail     string
	Screenshot string
}

type 黄点连通区域 struct {
	Count     int
	SumX      int
	SumY      int
	MinX      int
	MinY      int
	MaxX      int
	MaxY      int
	BestScore int
	BestX     int
	BestY     int
	BestColor string
}

type 黄点颜色统计项 struct {
	Color string
	Count int
}

const 小地图黄点最小连通像素 = 8

var 恐龙洞黄点测试执行中 atomic.Bool

func 执行恐龙洞黄点测试() {
	if !恐龙洞黄点测试执行中.CompareAndSwap(false, true) {
		设置恐龙洞输出("测试黄点执行中，请稍等")
		return
	}
	go func() {
		defer 恐龙洞黄点测试执行中.Store(false)
		if 引擎 == nil {
			设置恐龙洞输出("测试黄点失败：引擎未初始化")
			return
		}
		设置恐龙洞输出("测试黄点：区域=(10,96,201,231)")
		result := 扫描恐龙洞小地图黄点(恐龙洞小地图黄点区域, true)
		输出("恐龙洞黄点测试详情", result.Detail)
		if result.Screenshot != "" {
			输出("恐龙洞黄点测试截图", result.Screenshot)
		}
		if !result.Ok {
			设置恐龙洞输出("测试黄点失败：未找到，%s", result.Detail)
			return
		}
		标记恐龙洞找到的黄点(result.X, result.Y)
		if layer, baseX, baseY, diff, ok := 识别恐龙洞黄点所在层(result.X, result.Y); ok {
			设置恐龙洞输出("测试黄点成功：层=%d 黄点=(%d,%d) 基准=(%d,%d) 相对差值=%d", layer, result.X, result.Y, baseX, baseY, diff)
			return
		}
		设置恐龙洞输出("测试黄点找到坐标但楼层未匹配：x=%d y=%d", result.X, result.Y)
	}()
}

func 诊断小地图黄点() {
	if 引擎 == nil {
		输出("黄点诊断失败：引擎未初始化")
		return
	}

	matchedAny := false
	for index, 区域 := range 小地图黄点候选区域 {
		输出("黄点诊断区域", index+1, 区域.Name, "rect=", fmt.Sprintf("(%d,%d,%d,%d)", 区域.X1, 区域.Y1, 区域.X2, 区域.Y2))

		result := 扫描小地图黄点(区域, true)
		if result.Ok {
			matchedAny = true
		}
		rgbLayerText := 黄点层判断文本(result.Ok, result.Y)
		输出("黄点诊断 RGB中心", "区域=", 区域.Name, "ok=", result.Ok, "x=", result.X, "y=", result.Y, "判断=", rgbLayerText)
		输出("黄点诊断 RGB详情", result.Detail)
		if result.Screenshot != "" {
			输出("黄点诊断截图", result.Screenshot)
		}
	}

	if !matchedAny {
		输出("黄点诊断建议：两个候选框都没有 RGB 黄点。请打开诊断截图确认小地图是否在框内；如果黄点在截图里但 yellowPixels=0，就放宽黄点颜色规则或采集实际颜色。")
	}
}

func 黄点层判断文本(ok bool, y int) string {
	if !ok {
		return "未找到"
	}
	if layer, matched := 识别层数(y); matched {
		return fmt.Sprintf("%d层", layer)
	}
	return "未匹配"
}

func 扫描小地图黄点(feature *FColor, saveCrop bool) 黄点扫描结果 {
	return 扫描小地图黄点使用规则(feature, saveCrop, 是小地图黄点颜色, nil)
}

func 扫描恐龙洞小地图黄点(feature *FColor, saveCrop bool) 黄点扫描结果 {
	return 扫描小地图黄点使用规则(feature, saveCrop, 是恐龙洞玩家黄点颜色, 是恐龙洞玩家黄点区域)
}

func 扫描小地图黄点使用规则(feature *FColor, saveCrop bool, colorMatcher func(uint8, uint8, uint8) bool, componentMatcher func(黄点连通区域) bool) 黄点扫描结果 {
	if 引擎 == nil {
		return 黄点扫描结果{Ok: false, X: -1, Y: -1, Detail: "引擎未初始化"}
	}
	if feature == nil {
		return 黄点扫描结果{Ok: false, X: -1, Y: -1, Detail: "颜色区域为空"}
	}

	x1, y1, x2, y2 := 引擎.scaleRect(feature.X1, feature.Y1, feature.X2, feature.Y2)
	if x2 < x1 {
		x1, x2 = x2, x1
	}
	if y2 < y1 {
		y1, y2 = y2, y1
	}

	img := images.CaptureScreen(x1, y1, x2, y2, 引擎.displayID())
	if img == nil {
		return 黄点扫描结果{
			Ok:     false,
			X:      -1,
			Y:      -1,
			Detail: fmt.Sprintf("截图失败 rect=(%d,%d,%d,%d) displayId=%d", x1, y1, x2, y2, 引擎.displayID()),
		}
	}

	screenshot := ""
	if saveCrop {
		screenshot = 保存黄点诊断截图(img)
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	if width <= 0 || height <= 0 {
		return 黄点扫描结果{Ok: false, X: -1, Y: -1, Detail: "截图尺寸为空", Screenshot: screenshot}
	}

	mask := make([]bool, width*height)
	scores := make([]int, width*height)
	colors := make([]string, width*height)
	colorCounts := map[string]int{}
	yellowPixels := 0

	for py := bounds.Min.Y; py < bounds.Max.Y; py++ {
		for px := bounds.Min.X; px < bounds.Max.X; px++ {
			offset := img.PixOffset(px, py)
			if offset+3 >= len(img.Pix) {
				continue
			}
			r := img.Pix[offset]
			g := img.Pix[offset+1]
			b := img.Pix[offset+2]
			a := img.Pix[offset+3]
			if a < 32 || colorMatcher == nil || !colorMatcher(r, g, b) {
				continue
			}

			localX := px - bounds.Min.X
			localY := py - bounds.Min.Y
			index := localY*width + localX
			colorText := fmt.Sprintf("%02X%02X%02X", r, g, b)

			mask[index] = true
			scores[index] = 小地图黄点颜色分数(r, g, b)
			colors[index] = colorText
			colorCounts[colorText]++
			yellowPixels++
		}
	}

	best := 最大黄点连通区域使用规则(mask, scores, colors, width, height, componentMatcher)
	detailPrefix := fmt.Sprintf("rect=(%d,%d,%d,%d) capture=%dx%d yellowPixels=%d colors=%s",
		x1, y1, x2, y2, width, height, yellowPixels, 前N个黄点颜色(colorCounts, 6))
	if best.Count <= 0 {
		return 黄点扫描结果{Ok: false, X: -1, Y: -1, Detail: detailPrefix, Screenshot: screenshot}
	}

	componentDetail := fmt.Sprintf("%s component=%d bbox=(%d,%d,%d,%d) best=%s@(%d,%d)",
		detailPrefix,
		best.Count,
		引擎.unscaleX(x1+best.MinX),
		引擎.unscaleY(y1+best.MinY),
		引擎.unscaleX(x1+best.MaxX),
		引擎.unscaleY(y1+best.MaxY),
		best.BestColor,
		引擎.unscaleX(x1+best.BestX),
		引擎.unscaleY(y1+best.BestY),
	)
	if best.Count < 小地图黄点最小连通像素 {
		return 黄点扫描结果{
			Ok:         false,
			X:          -1,
			Y:          -1,
			Detail:     fmt.Sprintf("%s rejected=component<%d", componentDetail, 小地图黄点最小连通像素),
			Screenshot: screenshot,
		}
	}

	localX := int(math.Round(float64(best.SumX) / float64(best.Count)))
	localY := int(math.Round(float64(best.SumY) / float64(best.Count)))
	screenX := x1 + localX
	screenY := y1 + localY
	resultX := 引擎.unscaleX(screenX)
	resultY := 引擎.unscaleY(screenY)

	return 黄点扫描结果{Ok: true, X: resultX, Y: resultY, Detail: componentDetail, Screenshot: screenshot}
}

func 是小地图黄点颜色(r, g, b uint8) bool {
	ri := int(r)
	gi := int(g)
	bi := int(b)
	if ri < 150 || gi < 140 || bi > 110 {
		return false
	}
	if absInt(ri-gi) > 100 {
		return false
	}
	return ri > bi+70 && gi > bi+65 && ri+gi-bi >= 280
}

func 是恐龙洞玩家黄点颜色(r, g, b uint8) bool {
	ri := int(r)
	gi := int(g)
	bi := int(b)
	if ri < 220 || gi < 200 || bi > 100 {
		return false
	}
	if absInt(ri-gi) > 60 {
		return false
	}
	return ri > bi+130 && gi > bi+120
}

func 是恐龙洞玩家黄点区域(component 黄点连通区域) bool {
	if component.Count < 小地图黄点最小连通像素 || component.Count > 260 {
		return false
	}
	width := component.MaxX - component.MinX + 1
	height := component.MaxY - component.MinY + 1
	if width < 4 || height < 4 || width > 24 || height > 24 {
		return false
	}
	if width*2 < height || height*2 < width {
		return false
	}
	return component.Count*4 >= width*height
}

func 小地图黄点颜色分数(r, g, b uint8) int {
	ri := int(r)
	gi := int(g)
	bi := int(b)
	return ri + gi - 2*bi - absInt(ri-gi)
}

func 最大黄点连通区域(mask []bool, scores []int, colors []string, width, height int) 黄点连通区域 {
	return 最大黄点连通区域使用规则(mask, scores, colors, width, height, nil)
}

func 最大黄点连通区域使用规则(mask []bool, scores []int, colors []string, width, height int, componentMatcher func(黄点连通区域) bool) 黄点连通区域 {
	visited := make([]bool, len(mask))
	best := 黄点连通区域{}
	for index := range mask {
		if !mask[index] || visited[index] {
			continue
		}
		component := 扫描黄点连通区域(index, mask, visited, scores, colors, width, height)
		if componentMatcher != nil && !componentMatcher(component) {
			continue
		}
		if 黄点区域更好(component, best) {
			best = component
		}
	}
	return best
}

func 扫描黄点连通区域(start int, mask []bool, visited []bool, scores []int, colors []string, width, height int) 黄点连通区域 {
	component := 黄点连通区域{
		MinX:      width,
		MinY:      height,
		MaxX:      -1,
		MaxY:      -1,
		BestScore: math.MinInt,
	}
	queue := []int{start}
	visited[start] = true

	for len(queue) > 0 {
		index := queue[0]
		queue = queue[1:]
		x := index % width
		y := index / width

		component.Count++
		component.SumX += x
		component.SumY += y
		if x < component.MinX {
			component.MinX = x
		}
		if y < component.MinY {
			component.MinY = y
		}
		if x > component.MaxX {
			component.MaxX = x
		}
		if y > component.MaxY {
			component.MaxY = y
		}
		if scores[index] > component.BestScore {
			component.BestScore = scores[index]
			component.BestX = x
			component.BestY = y
			component.BestColor = colors[index]
		}

		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				if dx == 0 && dy == 0 {
					continue
				}
				nextX := x + dx
				nextY := y + dy
				if nextX < 0 || nextY < 0 || nextX >= width || nextY >= height {
					continue
				}
				next := nextY*width + nextX
				if !mask[next] || visited[next] {
					continue
				}
				visited[next] = true
				queue = append(queue, next)
			}
		}
	}
	return component
}

func 黄点区域更好(candidate, current 黄点连通区域) bool {
	if candidate.Count <= 0 {
		return false
	}
	if current.Count <= 0 {
		return true
	}
	if candidate.Count != current.Count {
		return candidate.Count > current.Count
	}
	return candidate.BestScore > current.BestScore
}

func 前N个黄点颜色(colorCounts map[string]int, n int) string {
	if len(colorCounts) == 0 {
		return "-"
	}
	items := make([]黄点颜色统计项, 0, len(colorCounts))
	for colorText, count := range colorCounts {
		items = append(items, 黄点颜色统计项{Color: colorText, Count: count})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Count == items[j].Count {
			return items[i].Color < items[j].Color
		}
		return items[i].Count > items[j].Count
	})

	if n > len(items) {
		n = len(items)
	}
	parts := make([]string, 0, n)
	for i := 0; i < n; i++ {
		parts = append(parts, fmt.Sprintf("%s:%d", items[i].Color, items[i].Count))
	}
	return strings.Join(parts, ",")
}

func 保存黄点诊断截图(img *image.NRGBA) string {
	dir := filepath.Dir(调试日志路径())
	if dir == "" || dir == "." {
		dir = "build"
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "保存失败:" + err.Error()
	}

	path := filepath.Join(dir, "yellow-dot-"+time.Now().Format("20060102-150405.000000000")+".png")
	file, err := os.Create(path)
	if err != nil {
		return "保存失败:" + err.Error()
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		return "保存失败:" + err.Error()
	}
	return path
}
