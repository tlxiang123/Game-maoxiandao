//go:build android

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Dasongzi1366/AutoGo/device"
	"github.com/Dasongzi1366/AutoGo/yolo"
)

const (
	怪物标签      = "怪物"
	人物标签      = "人物"
	怪物模型线程数   = 6
	怪物Y样本上限   = 240
	怪物Y最小分层间隔 = 20
	怪物识别后台间隔  = 800 * time.Millisecond
	怪物结果过期时间  = 3 * time.Second
)

type 怪物层统计 struct {
	数量 int
	怪物 []yolo.Result
}

type 怪物Y层中心 struct {
	层 int
	Y int
}

type YOLO结果缓存 struct {
	Results []yolo.Result
	Time    time.Time
}

var (
	怪物识别锁     sync.Mutex
	怪物后台锁     sync.Mutex
	YOLO结果锁   sync.RWMutex
	怪物检测器     *yolo.Yolo
	怪物模型目录    string
	怪物初始化失败   bool
	怪物后台RunID int64
	怪物屏幕Y样本   []int
	最新YOLO结果  YOLO结果缓存
)

func 初始化怪物识别() bool {
	怪物识别锁.Lock()
	defer 怪物识别锁.Unlock()

	if 怪物检测器 != nil {
		return true
	}
	if 怪物初始化失败 {
		return false
	}

	paramPath, binPath, err := 准备怪物模型文件()
	if err != nil {
		怪物初始化失败 = true
		输出("怪物 初始化失败 模型文件=", err)
		return false
	}
	labels, err := 读取YOLO标签()
	if err != nil {
		怪物初始化失败 = true
		输出("怪物 初始化失败 data.yaml=", err)
		return false
	}
	if !包含标签(labels, 怪物标签) {
		怪物初始化失败 = true
		输出("怪物 初始化失败 data.yaml缺少标签=", 怪物标签)
		return false
	}

	输出("怪物 YOLO初始化 labels=", len(labels), "目标=", 怪物标签)
	detector := yolo.New("v8", 怪物模型线程数, paramPath, binPath, strings.Join(labels, ","))
	if detector == nil {
		怪物初始化失败 = true
		输出("怪物 初始化失败 YOLO=nil")
		return false
	}
	怪物检测器 = detector
	输出("怪物 YOLO加载完成")
	return true
}

func 关闭怪物识别() {
	怪物识别锁.Lock()
	defer 怪物识别锁.Unlock()

	if 怪物检测器 != nil {
		怪物检测器.Close()
		怪物检测器 = nil
	}
}

func 启动怪物识别后台(runID int64) {
	怪物后台锁.Lock()
	if 怪物后台RunID == runID {
		怪物后台锁.Unlock()
		return
	}
	怪物后台RunID = runID
	怪物后台锁.Unlock()

	go 怪物识别后台循环(runID)
}

func 怪物识别后台循环(runID int64) {
	输出("怪物 YOLO后台启动 线程=", 怪物模型线程数)
	for 脚本仍应运行(runID) {
		if 初始化怪物识别() {
			写入最新YOLO结果(执行YOLO检测())
		}
		等待怪物后台间隔(runID)
	}
}

func 等待怪物后台间隔(runID int64) {
	deadline := time.Now().Add(怪物识别后台间隔)
	for 脚本仍应运行(runID) && time.Now().Before(deadline) {
		time.Sleep(50 * time.Millisecond)
	}
}

func 执行YOLO检测() []yolo.Result {
	怪物识别锁.Lock()
	defer 怪物识别锁.Unlock()

	if 怪物检测器 == nil {
		return nil
	}
	results := 怪物检测器.Detect(0, 0, 0, 0, 当前显示ID())
	return append([]yolo.Result(nil), results...)
}

func 写入最新YOLO结果(results []yolo.Result) {
	YOLO结果锁.Lock()
	defer YOLO结果锁.Unlock()

	最新YOLO结果 = YOLO结果缓存{
		Results: append([]yolo.Result(nil), results...),
		Time:    time.Now(),
	}
}

func 读取最新YOLO结果() ([]yolo.Result, time.Duration, bool) {
	YOLO结果锁.RLock()
	defer YOLO结果锁.RUnlock()

	if 最新YOLO结果.Time.IsZero() {
		return nil, 0, false
	}
	return append([]yolo.Result(nil), 最新YOLO结果.Results...), time.Since(最新YOLO结果.Time), true
}

func 准备怪物模型文件() (string, string, error) {
	if 怪物模型目录 == "" {
		baseDir, err := os.UserCacheDir()
		if err != nil || baseDir == "" {
			baseDir = os.TempDir()
		}
		怪物模型目录 = filepath.Join(baseDir, "autogo-yolo")
	}
	if err := os.MkdirAll(怪物模型目录, 0755); err != nil {
		return "", "", err
	}

	paramPath := filepath.Join(怪物模型目录, "best.ncnn.param")
	binPath := filepath.Join(怪物模型目录, "best.ncnn.bin")
	if err := 写入嵌入文件("best.ncnn.param", paramPath); err != nil {
		return "", "", err
	}
	if err := 写入嵌入文件("best.ncnn.bin", binPath); err != nil {
		return "", "", err
	}
	return paramPath, binPath, nil
}

func 写入嵌入文件(name, path string) error {
	data, err := res.ReadFile(name)
	if err != nil {
		return err
	}
	if info, err := os.Stat(path); err == nil && info.Size() == int64(len(data)) {
		return nil
	}
	return os.WriteFile(path, data, 0644)
}

func 读取YOLO标签() ([]string, error) {
	data, err := res.ReadFile("data.yaml")
	if err != nil {
		return nil, err
	}
	text := string(data)
	labels := 解析DataYamlNames(text)
	if len(labels) == 0 {
		return nil, fmt.Errorf("names为空")
	}
	if nc, ok := 解析DataYamlNC(text); ok && nc != len(labels) {
		return nil, fmt.Errorf("nc=%d names=%d 不一致", nc, len(labels))
	}
	return labels, nil
}

func 解析DataYamlNC(text string) (int, bool) {
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(strings.TrimRight(line, "\r"))
		if !strings.HasPrefix(trimmed, "nc:") {
			continue
		}
		value := strings.TrimSpace(strings.TrimPrefix(trimmed, "nc:"))
		nc, err := strconv.Atoi(value)
		return nc, err == nil
	}
	return 0, false
}

func 解析DataYamlNames(text string) []string {
	lines := strings.Split(text, "\n")
	labels := []string{}
	inNames := false
	for _, line := range lines {
		raw := strings.TrimRight(line, "\r")
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		if strings.HasPrefix(trimmed, "names:") {
			inNames = true
			value := strings.TrimSpace(strings.TrimPrefix(trimmed, "names:"))
			if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
				return 解析内联标签(value)
			}
			continue
		}
		if !inNames {
			continue
		}
		if strings.HasPrefix(trimmed, "-") {
			label := 清理YAML标签(strings.TrimSpace(strings.TrimPrefix(trimmed, "-")))
			if label != "" {
				labels = append(labels, label)
			}
			continue
		}
		if strings.Contains(trimmed, ":") {
			parts := strings.SplitN(trimmed, ":", 2)
			label := 清理YAML标签(parts[1])
			if label != "" {
				labels = append(labels, label)
				continue
			}
		}
		if len(labels) > 0 {
			break
		}
	}
	return labels
}

func 解析内联标签(value string) []string {
	value = strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(value, "["), "]"))
	if value == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	labels := make([]string, 0, len(parts))
	for _, part := range parts {
		label := 清理YAML标签(part)
		if label != "" {
			labels = append(labels, label)
		}
	}
	return labels
}

func 清理YAML标签(value string) string {
	value = strings.TrimSpace(value)
	value = strings.Trim(value, `"'`)
	return strings.TrimSpace(value)
}

func 包含标签(labels []string, target string) bool {
	for _, label := range labels {
		if strings.TrimSpace(label) == target {
			return true
		}
	}
	return false
}

func 打印怪物层统计(位置 层位置) {
	stats, results, ok := 读取怪物层统计(位置)
	if !ok {
		return
	}
	输出(当前层怪物统计文本(位置, stats, results))
}

func 打印怪物层统计并取当前层数量(位置 层位置) (int, bool) {
	stats, results, ok := 读取怪物层统计(位置)
	if !ok {
		return 0, false
	}
	输出(当前层怪物统计文本(位置, stats, results))
	return stats[位置.层].数量, true
}

func 读取怪物层统计(位置 层位置) (map[int]怪物层统计, []yolo.Result, bool) {
	results, age, ok := 读取最新YOLO结果()
	if !ok {
		输出("怪物 等待YOLO识别")
		return nil, nil, false
	}
	if age > 怪物结果过期时间 {
		输出("怪物 YOLO结果过期")
		return nil, nil, false
	}

	怪物结果 := 过滤怪物结果(results)
	记录怪物屏幕Y样本(怪物结果)
	return 统计怪物分层(位置.层, 怪物结果, 怪物屏幕Y聚类()), results, true
}

func YOLO已识别到怪物() bool {
	results, age, ok := 读取最新YOLO结果()
	if !ok || age > 怪物结果过期时间 {
		return false
	}
	return len(过滤怪物结果(results)) > 0
}

func YOLO已加载完成() bool {
	怪物识别锁.Lock()
	defer 怪物识别锁.Unlock()
	return 怪物检测器 != nil
}

func YOLO人物附近有怪物(maxDistance int) bool {
	results, age, ok := 读取最新YOLO结果()
	if !ok || age > 怪物结果过期时间 {
		return false
	}
	person, ok := 最高置信度结果(过滤人物结果(results))
	if !ok {
		return false
	}
	limit := maxDistance * maxDistance
	for _, monster := range 过滤怪物结果(results) {
		dx := 结果中心X(person) - 结果中心X(monster)
		dy := 怪物中心Y(person) - 怪物中心Y(monster)
		if dx*dx+dy*dy <= limit {
			return true
		}
	}
	return false
}

func YOLO当前层附近有怪物(位置 层位置, maxDistance int) bool {
	results, age, ok := 读取最新YOLO结果()
	if !ok || age > 怪物结果过期时间 {
		return false
	}
	person, ok := 最高置信度结果(过滤人物结果(results))
	if !ok {
		return false
	}
	monsters := 过滤怪物结果(results)
	if len(monsters) == 0 {
		return false
	}
	记录怪物屏幕Y样本(monsters)
	stats := 统计怪物分层(位置.层, monsters, 怪物屏幕Y聚类())
	limit := maxDistance * maxDistance
	for _, monster := range stats[位置.层].怪物 {
		dx := 结果中心X(person) - 结果中心X(monster)
		dy := 怪物中心Y(person) - 怪物中心Y(monster)
		if dx*dx+dy*dy <= limit {
			return true
		}
	}
	return false
}

func YOLO当前层怪物方向(位置 层位置) int {
	results, age, ok := 读取最新YOLO结果()
	if !ok || age > 怪物结果过期时间 {
		return 0
	}
	person, ok := 最高置信度结果(过滤人物结果(results))
	if !ok {
		return 0
	}
	monsters := 过滤怪物结果(results)
	if len(monsters) == 0 {
		return 0
	}
	记录怪物屏幕Y样本(monsters)
	stats := 统计怪物分层(位置.层, monsters, 怪物屏幕Y聚类())
	current := stats[位置.层].怪物
	if len(current) == 0 {
		return 0
	}
	target, ok := 最近YOLO结果(person, current)
	if !ok {
		return 0
	}
	if 结果中心X(target) < 结果中心X(person) {
		return -1
	}
	if 结果中心X(target) > 结果中心X(person) {
		return 1
	}
	return 0
}

func YOLO人物中心Y即时() (int, bool) {
	if !初始化怪物识别() {
		return 0, false
	}
	results := 执行YOLO检测()
	写入最新YOLO结果(results)
	person, ok := 最高置信度结果(过滤人物结果(results))
	if !ok {
		return 0, false
	}
	return 怪物中心Y(person), true
}

func 最近YOLO结果(base yolo.Result, candidates []yolo.Result) (yolo.Result, bool) {
	if len(candidates) == 0 {
		return yolo.Result{}, false
	}
	best := candidates[0]
	bestDistance := YOLO结果距离(base, best)
	for _, candidate := range candidates[1:] {
		distance := YOLO结果距离(base, candidate)
		if distance < bestDistance {
			best = candidate
			bestDistance = distance
		}
	}
	return best, true
}

func 最高置信度结果(results []yolo.Result) (yolo.Result, bool) {
	if len(results) == 0 {
		return yolo.Result{}, false
	}
	best := results[0]
	for _, result := range results[1:] {
		if result.Score > best.Score {
			best = result
		}
	}
	return best, true
}

func 过滤怪物结果(results []yolo.Result) []yolo.Result {
	return 过滤标签结果(results, 怪物标签, "monster")
}

func 过滤人物结果(results []yolo.Result) []yolo.Result {
	return 过滤标签结果(results, 人物标签, "person")
}

func 过滤标签结果(results []yolo.Result, labels ...string) []yolo.Result {
	out := make([]yolo.Result, 0, len(results))
	for _, result := range results {
		label := strings.TrimSpace(strings.ToLower(result.Label))
		if label == "" || !匹配任一标签(label, labels...) {
			continue
		}
		out = append(out, result)
	}
	return out
}

func 匹配任一标签(label string, labels ...string) bool {
	for _, candidate := range labels {
		if label == strings.ToLower(strings.TrimSpace(candidate)) {
			return true
		}
	}
	return false
}

func 记录怪物屏幕Y样本(results []yolo.Result) {
	for _, result := range results {
		if y := 怪物中心Y(result); y > 0 {
			怪物屏幕Y样本 = append(怪物屏幕Y样本, y)
		}
	}
	if len(怪物屏幕Y样本) > 怪物Y样本上限 {
		怪物屏幕Y样本 = 怪物屏幕Y样本[len(怪物屏幕Y样本)-怪物Y样本上限:]
	}
}

func 怪物屏幕Y聚类() []怪物Y层中心 {
	if len(怪物屏幕Y样本) == 0 {
		return nil
	}
	values := append([]int(nil), 怪物屏幕Y样本...)
	sort.Ints(values)

	gaps := make([]int, 0, len(values)-1)
	for i := 1; i < len(values); i++ {
		if values[i]-values[i-1] >= 怪物Y最小分层间隔 {
			gaps = append(gaps, i)
		}
	}
	sort.Slice(gaps, func(i, j int) bool {
		return values[gaps[i]]-values[gaps[i]-1] > values[gaps[j]]-values[gaps[j]-1]
	})
	if len(gaps) > 2 {
		gaps = gaps[:2]
	}
	sort.Ints(gaps)

	start := 0
	centers := make([]怪物Y层中心, 0, len(gaps)+1)
	for _, split := range append(gaps, len(values)) {
		if split <= start {
			continue
		}
		centers = append(centers, 怪物Y层中心{Y: 平均值(values[start:split])})
		start = split
	}
	sort.Slice(centers, func(i, j int) bool {
		return centers[i].Y < centers[j].Y
	})
	for i := range centers {
		switch len(centers) {
		case 3:
			centers[i].层 = 3 - i
		case 2:
			centers[i].层 = 按屏幕Y默认判断层(centers[i].Y)
			if i == 1 && centers[0].层 == centers[1].层 {
				centers[0].层 = 3
				centers[1].层 = 1
			}
		default:
			centers[i].层 = 按屏幕Y默认判断层(centers[i].Y)
		}
	}
	return centers
}

func 平均值(values []int) int {
	if len(values) == 0 {
		return 0
	}
	total := 0
	for _, value := range values {
		total += value
	}
	return total / len(values)
}

func 统计怪物分层(当前层 int, results []yolo.Result, centers []怪物Y层中心) map[int]怪物层统计 {
	stats := map[int]怪物层统计{
		1: {},
		2: {},
		3: {},
	}
	for _, result := range results {
		layer := 判断怪物层(当前层, result, centers)
		if layer < 1 || layer > 3 {
			continue
		}
		item := stats[layer]
		item.数量++
		item.怪物 = append(item.怪物, result)
		stats[layer] = item
	}
	return stats
}

func 判断怪物层(当前层 int, result yolo.Result, centers []怪物Y层中心) int {
	y := 怪物中心Y(result)
	if len(centers) == 1 && 当前层 >= 1 && 当前层 <= 3 {
		return 当前层
	}

	best := 0
	bestDiff := 0
	for _, center := range centers {
		diff := absInt(y - center.Y)
		if best == 0 || diff < bestDiff {
			best = center.层
			bestDiff = diff
		}
	}
	if best != 0 {
		return best
	}
	return 按屏幕Y默认判断层(y)
}

func 怪物中心Y(result yolo.Result) int {
	if result.CenterY > 0 {
		return result.CenterY
	}
	return result.Y + result.Height/2
}

func 结果中心X(result yolo.Result) int {
	if result.CenterX > 0 {
		return result.CenterX
	}
	return result.X + result.Width/2
}

func 按屏幕Y默认判断层(y int) int {
	_, height, _, _ := device.GetDisplayInfo(当前显示ID())
	if height <= 0 {
		height = 720
	}
	if y <= height/3 {
		return 3
	}
	if y <= height*2/3 {
		return 2
	}
	return 1
}

func 怪物统计文本(stats map[int]怪物层统计) string {
	parts := make([]string, 0, 3)
	for layer := 1; layer <= 3; layer++ {
		stat := stats[layer]
		if stat.数量 == 0 {
			parts = append(parts, fmt.Sprintf("%d层=0只", layer))
			continue
		}
		parts = append(parts, fmt.Sprintf("%d层=%d只", layer, stat.数量))
	}
	return "怪物 " + strings.Join(parts, " | ")
}

func 当前层怪物统计文本(位置 层位置, stats map[int]怪物层统计, results []yolo.Result) string {
	stat := stats[位置.层]
	return fmt.Sprintf("怪物 当前%d层=%d只 距离=%s", 位置.层, stat.数量, 怪物距离文本(stat.怪物, results))
}

func 怪物距离文本(monsters []yolo.Result, results []yolo.Result) string {
	if len(monsters) == 0 {
		return "-"
	}
	person, ok := 最高置信度结果(过滤人物结果(results))
	if !ok {
		return "无人物"
	}
	distances := make([]int, 0, len(monsters))
	for _, monster := range monsters {
		distances = append(distances, YOLO结果距离(person, monster))
	}
	sort.Ints(distances)
	parts := make([]string, 0, len(distances))
	for _, distance := range distances {
		parts = append(parts, fmt.Sprintf("%d", distance))
	}
	return strings.Join(parts, "/")
}

func YOLO结果距离(a, b yolo.Result) int {
	dx := 结果中心X(a) - 结果中心X(b)
	dy := 怪物中心Y(a) - 怪物中心Y(b)
	return intSqrt(dx*dx + dy*dy)
}

func intSqrt(v int) int {
	if v <= 0 {
		return 0
	}
	x := v
	y := (x + 1) / 2
	for y < x {
		x = y
		y = (x + v/x) / 2
	}
	return x
}

func 置信度文本(scores []float64) string {
	if len(scores) == 0 {
		return "-"
	}
	sort.Slice(scores, func(i, j int) bool {
		return scores[i] > scores[j]
	})
	parts := make([]string, 0, len(scores))
	for _, score := range scores {
		parts = append(parts, fmt.Sprintf("%.2f", score))
	}
	return strings.Join(parts, "/")
}
