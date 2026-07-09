package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"math"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Dasongzi1366/AutoGo/console"
	"github.com/Dasongzi1366/AutoGo/device"
	"github.com/Dasongzi1366/AutoGo/dotocr"
	"github.com/Dasongzi1366/AutoGo/images"
	"github.com/Dasongzi1366/AutoGo/motion"
	"github.com/Dasongzi1366/AutoGo/ppocr"
	"github.com/Dasongzi1366/AutoGo/utils"
)

var GlobalLog = false
var GlobalLogSuccessOnly = false
var GlobalDisplayId = 0
var GlobalToast = false
var GlobalConsole = false

// 0 表示不做坐标缩放，所有坐标按脚本里的绝对值执行。
var GlobalBaseWidth = 0
var GlobalBaseHeight = 0

var globalConsole *console.Console

type Zg struct {
	offsetX      int
	offsetY      int
	log          bool
	imgMap       map[string]*[]byte
	displayId    int
	hasDisplayId bool
	ppocrEngine  *ppocr.Ppocr
	ppocrVersion string
}

func New(res embed.FS) *Zg {
	return &Zg{imgMap: loadImages(res)}
}

func SetGlobalLog(enable bool) {
	GlobalLog = enable
}

func SetGlobalLogSuccessOnly(enable bool) {
	GlobalLogSuccessOnly = enable
}

func SetGlobalDisplayId(displayId int) {
	GlobalDisplayId = displayId
	log.Printf("[Display] global displayId=%d", displayId)
}

func SetGlobalToast(enable bool) {
	GlobalToast = enable
}

func SetGlobalConsole(enable bool) {
	GlobalConsole = enable
}

func SetGlobalBaseSize(width, height int) {
	GlobalBaseWidth = width
	GlobalBaseHeight = height
}

func (z *Zg) SetLog(enable bool) {
	z.log = enable
}

func (z *Zg) SetDisplayId(displayId int) {
	z.displayId = displayId
	z.hasDisplayId = true
	log.Printf("[Display] engine displayId=%d", displayId)
}

func (z *Zg) ClearDisplayId() {
	z.displayId = 0
	z.hasDisplayId = false
}

func (z *Zg) SetOffset(x, y int) {
	z.offsetX = x
	z.offsetY = y
}

func (z *Zg) markCColorFeature(feature *CColor) {
	if z == nil || feature == nil {
		return
	}

	x := feature.X + z.offsetX
	y := feature.Y + z.offsetY
	minX, maxX := x, x
	minY, maxY := y, y

	parts := strings.Split(feature.Color, ",")
	for i := 1; i+2 < len(parts); i += 3 {
		dx, errX := strconv.Atoi(strings.TrimSpace(parts[i]))
		dy, errY := strconv.Atoi(strings.TrimSpace(parts[i+1]))
		if errX != nil || errY != nil {
			continue
		}
		px := feature.X + dx + z.offsetX
		py := feature.Y + dy + z.offsetY
		if px < minX {
			minX = px
		}
		if px > maxX {
			maxX = px
		}
		if py < minY {
			minY = py
		}
		if py > maxY {
			maxY = py
		}
	}

	const padding = 12
	z.markScreenRect(
		z.scaleX(minX-padding),
		z.scaleY(minY-padding),
		z.scaleX(maxX+padding),
		z.scaleY(maxY+padding),
	)
}

func (z *Zg) CmpColor(feature *CColor) (bool, int, int) {
	if feature == nil {
		return false, -1, -1
	}
	x := feature.X + z.offsetX
	y := feature.Y + z.offsetY
	z.markCColorFeature(feature)
	ok := false
	if strings.Contains(feature.Color, ",") {
		colors := buildDetectColorsFromCColor(feature, z.offsetX, z.offsetY, z.displayID())
		ok = images.DetectsMultiColors(colors, defaultSim(feature.Sim), z.displayID())
	} else {
		ok = images.CmpColor(z.scaleX(x), z.scaleY(y), feature.Color, defaultSim(feature.Sim), z.displayID())
	}
	name := colorName(feature.Name, feature.Color)
	z.printLog("CmpColor", ok, x, y, name)
	z.toastSuccess(ok, name)
	z.consoleSuccess(ok, name)
	if ok {
		return true, x, y
	}
	return false, -1, -1
}

func (z *Zg) FindColor(feature *FColor) (bool, int, int) {
	if feature == nil {
		return false, -1, -1
	}
	x1, y1, x2, y2 := z.scaleRect(feature.X1, feature.Y1, feature.X2, feature.Y2)
	z.markScreenRect(x1, y1, x2, y2)
	x, y := images.FindColor(x1, y1, x2, y2, feature.Color, defaultSim(feature.Sim), feature.Dir, z.displayID())
	ok := x != -1 && y != -1
	name := colorName(feature.Name, feature.Color)
	z.printLog("FindColor", ok, x, y, name)
	z.toastSuccess(ok, name)
	z.consoleSuccess(ok, name)
	if !ok {
		return false, -1, -1
	}
	z.markScreenPoint(x, y)
	return true, z.unscaleX(x), z.unscaleY(y)
}

func (z *Zg) GetColorCountInRegion(feature *CCRegion) int {
	if feature == nil {
		return 0
	}
	x1, y1, x2, y2 := z.scaleRect(feature.X1, feature.Y1, feature.X2, feature.Y2)
	z.markScreenRect(x1, y1, x2, y2)
	count := images.GetColorCountInRegion(x1, y1, x2, y2, feature.Color, defaultSim(feature.Sim), z.displayID())
	ok := count > 0
	if z.shouldLog() && (!GlobalLogSuccessOnly || ok) {
		log.Printf("[GetColorCountInRegion] ok=%v count=%d detail=%s rect=(%d,%d,%d,%d) displayId=%d", ok, count, feature.Name, x1, y1, x2, y2, z.displayID())
	}
	return count
}

func (z *Zg) DetectsMultiColors(feature *DMColor) (bool, int, int) {
	if feature == nil {
		return false, -1, -1
	}
	colors := offsetDetectColors(feature.Colors, z.offsetX, z.offsetY, z.displayID())
	firstX, firstY := firstDetectColorPoint(colors, z.displayID())
	if firstX != -1 && firstY != -1 {
		z.markScreenPoint(firstX, firstY)
	}
	ok := images.DetectsMultiColors(colors, defaultSim(feature.Sim), z.displayID())
	if !ok {
		firstX, firstY = -1, -1
	}
	name := colorName(feature.Name, colors)
	z.printLog("DetectsMultiColors", ok, firstX, firstY, name)
	z.toastSuccess(ok, name)
	z.consoleSuccess(ok, name)
	return ok, firstX, firstY
}

func (z *Zg) FindMultiColor(feature *FMColor) (bool, int, int) {
	if feature == nil {
		return false, -1, -1
	}
	x1, y1, x2, y2 := z.scaleRect(feature.X1, feature.Y1, feature.X2, feature.Y2)
	z.markScreenRect(x1, y1, x2, y2)
	colors := buildMultiColor(feature.MainColor, feature.OffsetColor)
	x, y := images.FindMultiColors(x1, y1, x2, y2, colors, defaultSim(feature.Sim), feature.Dir, z.displayID())
	ok := x != -1 && y != -1
	name := colorName(feature.Name, colors)
	z.printLog("FindMultiColor", ok, x, y, name)
	z.toastSuccess(ok, name)
	z.consoleSuccess(ok, name)
	if !ok {
		return false, -1, -1
	}
	z.markScreenPoint(x, y)
	return true, z.unscaleX(x), z.unscaleY(y)
}

func (z *Zg) FindPic(feature *Pic) (bool, int, int) {
	if feature == nil || z.imgMap == nil {
		return false, -1, -1
	}
	name := picName(feature)
	rx1, ry1, rx2, ry2 := z.offsetRect(feature.X1, feature.Y1, feature.X2, feature.Y2)
	z.markScreenRect(z.scaleX(rx1), z.scaleY(ry1), z.scaleX(rx2), z.scaleY(ry2))
	template := z.imgMap[normalizePicName(feature.PicPath)]
	if template == nil {
		z.printLog("FindPic", false, -1, -1, "template not found:"+name)
		return false, -1, -1
	}
	x, y := z.findPicPure(feature, template)
	ok := x != -1 && y != -1
	z.printLog("FindPic", ok, x, y, name)
	z.toastSuccess(ok, name)
	z.consoleSuccess(ok, name)
	if !ok {
		return false, -1, -1
	}
	z.markScreenPoint(z.scaleX(x), z.scaleY(y))
	return true, z.unscaleX(x), z.unscaleY(y)
}

func (z *Zg) FindStr(feature *FStr) (bool, int, int) {
	if feature == nil {
		return false, -1, -1
	}
	x1, y1, x2, y2 := z.scaleRect(feature.X1, feature.Y1, feature.X2, feature.Y2)
	z.markScreenRect(x1, y1, x2, y2)
	x, y := dotocr.FindStr(x1, y1, x2, y2, feature.String, feature.ColorFormat, defaultSim(feature.Sim), feature.DictName, z.displayID())
	ok := x != -1 && y != -1
	name := textName(feature.Name, feature.String)
	z.printLog("FindStr", ok, x, y, name)
	z.toastSuccess(ok, name)
	z.consoleSuccess(ok, name)
	if !ok {
		return false, -1, -1
	}
	z.markScreenPoint(x, y)
	return true, z.unscaleX(x), z.unscaleY(y)
}

func (z *Zg) Ocr(feature *SOcr) string {
	if feature == nil {
		return ""
	}
	x1, y1, x2, y2 := z.scaleRect(feature.X1, feature.Y1, feature.X2, feature.Y2)
	z.markScreenRect(x1, y1, x2, y2)
	text := dotocr.Ocr(x1, y1, x2, y2, feature.ColorFormat, defaultSim(feature.Sim), false, feature.DictName, z.displayID())
	if z.shouldLog() && (!GlobalLogSuccessOnly || strings.TrimSpace(text) != "") {
		log.Printf("[Ocr] name=%s text=%q rect=(%d,%d,%d,%d) displayId=%d", textName(feature.Name, ""), text, x1, y1, x2, y2, z.displayID())
	}
	return text
}

func (z *Zg) PPOcr(feature *PPOcrRegion) []ppocr.Result {
	if feature == nil {
		return nil
	}
	engine := z.ppocr(feature.Version)
	if engine == nil {
		z.printLog("PPOcr", false, -1, -1, "engine init failed")
		return nil
	}
	x1, y1, x2, y2 := z.scaleRect(feature.X1, feature.Y1, feature.X2, feature.Y2)
	z.markScreenRect(x1, y1, x2, y2)
	results := engine.Ocr(x1, y1, x2, y2, feature.Color, z.displayID())
	ok := len(results) > 0
	if z.shouldLog() && (!GlobalLogSuccessOnly || ok) {
		log.Printf("[PPOcr] ok=%v count=%d detail=%s rect=(%d,%d,%d,%d) displayId=%d", ok, len(results), feature.Name, x1, y1, x2, y2, z.displayID())
	}
	return results
}

func (z *Zg) PPOcrText(feature *PPOcrRegion) string {
	results := z.PPOcr(feature)
	if len(results) == 0 {
		return ""
	}
	labels := make([]string, 0, len(results))
	for _, result := range results {
		label := strings.TrimSpace(result.Label)
		if label != "" {
			labels = append(labels, label)
		}
	}
	return strings.Join(labels, " ")
}

func (z *Zg) FindPPOcrText(feature *PPOcrRegion) (bool, int, int) {
	if feature == nil {
		return false, -1, -1
	}
	contains := splitTextCandidates(feature.Contains)
	results := z.PPOcr(feature)
	for _, result := range results {
		label := strings.TrimSpace(result.Label)
		if len(contains) > 0 && !containsAny(label, contains) {
			continue
		}
		if feature.MinScore > 0 && result.Score < feature.MinScore {
			continue
		}
		x, y := ppocrCenter(result)
		z.markScreenPoint(x, y)
		z.printLog("FindPPOcrText", true, x, y, textName(feature.Name, label))
		z.toastSuccess(true, textName(feature.Name, label))
		z.consoleSuccess(true, textName(feature.Name, label))
		return true, z.unscaleX(x), z.unscaleY(y)
	}
	z.printLog("FindPPOcrText", false, -1, -1, textName(feature.Name, strings.Join(contains, "|")))
	return false, -1, -1
}

func (z *Zg) ClickPPOcrText(feature *PPOcrRegion) bool {
	ok, x, y := z.FindPPOcrText(feature)
	return z.ClickResult(ok, x, y)
}

func (z *Zg) ClickPPOcrTextOffset(feature *PPOcrRegion, offsetX, offsetY int) bool {
	ok, x, y := z.FindPPOcrText(feature)
	return z.ClickResultOffset(ok, x, y, offsetX, offsetY)
}

func (z *Zg) ClosePPOcr() {
	if z.ppocrEngine == nil {
		return
	}
	z.ppocrEngine.Close()
	z.ppocrEngine = nil
	z.ppocrVersion = ""
}

func (z *Zg) FindFeature(feature any) (bool, int, int) {
	switch f := feature.(type) {
	case *CColor:
		return z.CmpColor(f)
	case CColor:
		return z.CmpColor(&f)
	case *FColor:
		return z.FindColor(f)
	case FColor:
		return z.FindColor(&f)
	case *DMColor:
		return z.DetectsMultiColors(f)
	case DMColor:
		return z.DetectsMultiColors(&f)
	case *FMColor:
		return z.FindMultiColor(f)
	case FMColor:
		return z.FindMultiColor(&f)
	case *Pic:
		return z.FindPic(f)
	case Pic:
		return z.FindPic(&f)
	case *FStr:
		return z.FindStr(f)
	case FStr:
		return z.FindStr(&f)
	case *PPOcrRegion:
		return z.FindPPOcrText(f)
	case PPOcrRegion:
		return z.FindPPOcrText(&f)
	default:
		z.printLog("FindFeature", false, -1, -1, fmt.Sprintf("unsupported feature type %T", feature))
		return false, -1, -1
	}
}

func (z *Zg) Find(feature any) *FindResult {
	ok, x, y := z.FindFeature(feature)
	return &FindResult{engine: z, ok: ok, x: x, y: y}
}

type FindResult struct {
	engine *Zg
	ok     bool
	x      int
	y      int
}

func (r *FindResult) Result() (bool, int, int) {
	if r == nil {
		return false, -1, -1
	}
	return r.ok, r.x, r.y
}

func (r *FindResult) Click() bool {
	if r == nil || r.engine == nil {
		return false
	}
	return r.engine.ClickResult(r.ok, r.x, r.y)
}

func (r *FindResult) ClickOffset(offsetX, offsetY int) bool {
	if r == nil || r.engine == nil {
		return false
	}
	return r.engine.ClickResultOffset(r.ok, r.x, r.y, offsetX, offsetY)
}

func (r *FindResult) ClickFixed(x, y int) bool {
	if r == nil || r.engine == nil || !r.ok {
		return false
	}
	r.engine.Click(x, y)
	return true
}

func (z *Zg) ClickFeature(feature any) bool {
	ok, x, y := z.FindFeature(feature)
	return z.ClickResult(ok, x, y)
}

func (z *Zg) ClickFeatureOffset(feature any, offsetX, offsetY int) bool {
	ok, x, y := z.FindFeature(feature)
	return z.ClickResultOffset(ok, x, y, offsetX, offsetY)
}

func (z *Zg) ClickWhenFound(feature any, x, y int) bool {
	ok, _, _ := z.FindFeature(feature)
	if !ok {
		return false
	}
	z.Click(x, y)
	return true
}

func (z *Zg) Click(x, y int) {
	sx := z.scaleX(x + z.offsetX)
	sy := z.scaleY(y + z.offsetY)
	z.markScreenPoint(sx, sy)
	motion.Click(sx, sy, 0, z.displayID())
	z.afterActionDelay()
}

func (z *Zg) ClickResult(ok bool, x, y int) bool {
	if !ok {
		return false
	}
	sx := z.scaleX(x + z.offsetX)
	sy := z.scaleY(y + z.offsetY)
	z.markScreenPoint(sx, sy)
	motion.Click(sx, sy, 0, z.displayID())
	z.afterActionDelay()
	return true
}

func (z *Zg) ClickResultOffset(ok bool, x, y, offsetX, offsetY int) bool {
	if !ok {
		return false
	}
	sx := z.scaleX(x + offsetX + z.offsetX)
	sy := z.scaleY(y + offsetY + z.offsetY)
	z.markScreenPoint(sx, sy)
	motion.Click(sx, sy, 0, z.displayID())
	z.afterActionDelay()
	return true
}

func (z *Zg) Swipe(x1, y1, x2, y2, duration int) {
	motion.Swipe(z.scaleX(x1+z.offsetX), z.scaleY(y1+z.offsetY), z.scaleX(x2+z.offsetX), z.scaleY(y2+z.offsetY), duration, 0, z.displayID())
	z.afterActionDelay()
}

func (z *Zg) PressSwipe(x1, y1, x2, y2, holdMs, duration, steps int) {
	if holdMs < 0 {
		holdMs = 0
	}
	if duration < 0 {
		duration = 0
	}
	if steps <= 0 {
		steps = 12
	}
	startX := z.scaleX(x1 + z.offsetX)
	startY := z.scaleY(y1 + z.offsetY)
	endX := z.scaleX(x2 + z.offsetX)
	endY := z.scaleY(y2 + z.offsetY)
	displayId := z.displayID()

	motion.TouchDown(startX, startY, 0, displayId)
	if holdMs > 0 {
		time.Sleep(time.Duration(holdMs) * time.Millisecond)
	}
	stepDelay := time.Duration(0)
	if duration > 0 {
		stepDelay = time.Duration(duration/steps) * time.Millisecond
	}
	for i := 1; i <= steps; i++ {
		x := startX + (endX-startX)*i/steps
		y := startY + (endY-startY)*i/steps
		motion.TouchMove(x, y, 0, displayId)
		if stepDelay > 0 {
			time.Sleep(stepDelay)
		}
	}
	motion.TouchUp(endX, endY, 0, displayId)
	z.afterActionDelay()
}

func (z *Zg) offsetRect(x1, y1, x2, y2 int) (int, int, int, int) {
	return x1 + z.offsetX, y1 + z.offsetY, x2 + z.offsetX, y2 + z.offsetY
}

func (z *Zg) scaleX(x int) int {
	if GlobalBaseWidth <= 0 {
		return x
	}
	width, _, _, _ := device.GetDisplayInfo(z.displayID())
	if width <= 0 {
		width = GlobalBaseWidth
	}
	return scaleCoordinate(x, GlobalBaseWidth, width)
}

func (z *Zg) scaleY(y int) int {
	if GlobalBaseHeight <= 0 {
		return y
	}
	_, height, _, _ := device.GetDisplayInfo(z.displayID())
	if height <= 0 {
		height = GlobalBaseHeight
	}
	return scaleCoordinate(y, GlobalBaseHeight, height)
}

func (z *Zg) scaleRect(x1, y1, x2, y2 int) (int, int, int, int) {
	return z.scaleX(x1 + z.offsetX), z.scaleY(y1 + z.offsetY), z.scaleX(x2 + z.offsetX), z.scaleY(y2 + z.offsetY)
}

func (z *Zg) unscaleX(x int) int {
	if GlobalBaseWidth <= 0 {
		return x
	}
	width, _, _, _ := device.GetDisplayInfo(z.displayID())
	if width <= 0 {
		width = GlobalBaseWidth
	}
	return scaleCoordinate(x, width, GlobalBaseWidth)
}

func (z *Zg) unscaleY(y int) int {
	if GlobalBaseHeight <= 0 {
		return y
	}
	_, height, _, _ := device.GetDisplayInfo(z.displayID())
	if height <= 0 {
		height = GlobalBaseHeight
	}
	return scaleCoordinate(y, height, GlobalBaseHeight)
}

func scaleCoordinate(value, base, current int) int {
	if base <= 0 || current <= 0 {
		return value
	}
	return int(math.Round(float64(value) * float64(current) / float64(base)))
}

func (z *Zg) shouldLog() bool {
	return z.log || GlobalLog
}

func (z *Zg) displayID() int {
	if z.hasDisplayId {
		return z.displayId
	}
	return GlobalDisplayId
}

func (z *Zg) afterActionDelay() {
	if 当前UI配置.点击后延迟毫秒 <= 0 {
		return
	}
	time.Sleep(time.Duration(当前UI配置.点击后延迟毫秒) * time.Millisecond)
}

func (z *Zg) toastSuccess(ok bool, name string) {
	if ok && GlobalToast {
		utils.Toast(name, -1, -1, 800)
	}
}

func (z *Zg) consoleSuccess(ok bool, name string) {
	if !ok || !GlobalConsole {
		return
	}
	if globalConsole == nil {
		globalConsole = console.New()
	}
	globalConsole.Println(name)
}

func (z *Zg) printLog(action string, ok bool, x, y int, detail string) {
	if !z.shouldLog() || (GlobalLogSuccessOnly && !ok) {
		return
	}
	log.Printf("[%s] ok=%v x=%d y=%d detail=%s offset=(%d,%d) displayId=%d", action, ok, x, y, detail, z.offsetX, z.offsetY, z.displayID())
}

func (z *Zg) ppocr(version string) *ppocr.Ppocr {
	if version == "" {
		version = 当前UI配置.OCR版本
	}
	if version == "" {
		version = "v5"
	}
	if z.ppocrEngine != nil && z.ppocrVersion == version {
		return z.ppocrEngine
	}
	z.ClosePPOcr()
	z.ppocrEngine = ppocr.New(version)
	z.ppocrVersion = version
	return z.ppocrEngine
}

func loadImages(res embed.FS) map[string]*[]byte {
	result := map[string]*[]byte{}
	_ = fs.WalkDir(res, "pic", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d == nil || d.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
			return nil
		}
		data, err := res.ReadFile(path)
		if err != nil {
			return nil
		}
		bytes := make([]byte, len(data))
		copy(bytes, data)
		key := normalizePicName(path)
		result[key] = &bytes
		result[normalizePicName(filepath.Base(path))] = &bytes
		return nil
	})
	return result
}

func normalizePicName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ReplaceAll(name, "\\", "/")
	name = strings.TrimPrefix(name, "pic/")
	return strings.ToLower(name)
}

func picName(feature *Pic) string {
	if feature == nil {
		return ""
	}
	if feature.Name != "" {
		return feature.Name
	}
	return feature.PicPath
}

func colorName(name, color string) string {
	if name != "" {
		return name
	}
	return color
}

func textName(name, text string) string {
	if name != "" {
		return name
	}
	return text
}

func defaultSim(sim float32) float32 {
	if sim <= 0 {
		return 0.90
	}
	return sim
}

func buildMultiColor(mainColor, offsetColor string) string {
	mainColor = strings.TrimSpace(mainColor)
	offsetColor = strings.TrimSpace(offsetColor)
	if offsetColor == "" {
		return mainColor
	}
	if mainColor == "" {
		return offsetColor
	}
	return mainColor + "," + offsetColor
}
