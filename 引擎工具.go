package main

import (
	"math"
	"strconv"
	"strings"

	"github.com/Dasongzi1366/AutoGo/device"
	"github.com/Dasongzi1366/AutoGo/ppocr"
)

func splitTextCandidates(text string) []string {
	parts := strings.Split(text, "|")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}

func buildDetectColorsFromCColor(feature *CColor, offsetX, offsetY, displayID int) string {
	if feature == nil {
		return ""
	}
	parts := strings.Split(feature.Color, ",")
	if len(parts) == 0 {
		return ""
	}

	currentWidth := displayWidth(displayID)
	currentHeight := displayHeight(displayID)
	x := scaleForDisplay(feature.X+offsetX, GlobalBaseWidth, currentWidth)
	y := scaleForDisplay(feature.Y+offsetY, GlobalBaseHeight, currentHeight)
	result := []string{strconv.Itoa(x), strconv.Itoa(y), strings.TrimSpace(parts[0])}

	for i := 1; i+2 < len(parts); i += 3 {
		dx, errX := strconv.Atoi(strings.TrimSpace(parts[i]))
		dy, errY := strconv.Atoi(strings.TrimSpace(parts[i+1]))
		if errX != nil || errY != nil {
			continue
		}
		result = append(result,
			strconv.Itoa(x+scaleForDisplay(dx, GlobalBaseWidth, currentWidth)),
			strconv.Itoa(y+scaleForDisplay(dy, GlobalBaseHeight, currentHeight)),
			strings.TrimSpace(parts[i+2]),
		)
	}
	return strings.Join(result, ",")
}

func containsAny(text string, candidates []string) bool {
	for _, candidate := range candidates {
		if strings.Contains(text, candidate) {
			return true
		}
	}
	return false
}

func ppocrCenter(result ppocr.Result) (int, int) {
	if result.CenterX != 0 || result.CenterY != 0 {
		return result.CenterX, result.CenterY
	}
	return result.X + result.Width/2, result.Y + result.Height/2
}

func offsetDetectColors(colors string, offsetX, offsetY, displayID int) string {
	colors = strings.TrimSpace(colors)
	if colors == "" {
		return colors
	}

	parts := strings.Split(colors, ",")
	if len(parts) < 3 {
		return colors
	}

	currentWidth := displayWidth(displayID)
	currentHeight := displayHeight(displayID)
	result := make([]string, len(parts))
	copy(result, parts)
	for i := 0; i+2 < len(result); i += 3 {
		x, errX := strconv.Atoi(strings.TrimSpace(result[i]))
		y, errY := strconv.Atoi(strings.TrimSpace(result[i+1]))
		if errX != nil || errY != nil {
			continue
		}
		result[i] = strconv.Itoa(scaleForDisplay(x+offsetX, GlobalBaseWidth, currentWidth))
		result[i+1] = strconv.Itoa(scaleForDisplay(y+offsetY, GlobalBaseHeight, currentHeight))
	}
	return strings.Join(result, ",")
}

func firstDetectColorPoint(colors string, displayID int) (int, int) {
	parts := strings.Split(colors, ",")
	if len(parts) < 2 {
		return -1, -1
	}
	x, errX := strconv.Atoi(strings.TrimSpace(parts[0]))
	y, errY := strconv.Atoi(strings.TrimSpace(parts[1]))
	if errX != nil || errY != nil {
		return -1, -1
	}
	return x, y
}

func displayWidth(displayID int) int {
	width, _, _, _ := device.GetDisplayInfo(displayID)
	if width <= 0 {
		return GlobalBaseWidth
	}
	return width
}

func displayHeight(displayID int) int {
	_, height, _, _ := device.GetDisplayInfo(displayID)
	if height <= 0 {
		return GlobalBaseHeight
	}
	return height
}

func scaleForDisplay(value, base, current int) int {
	if base <= 0 || current <= 0 {
		return value
	}
	return int(math.Round(float64(value) * float64(current) / float64(base)))
}
