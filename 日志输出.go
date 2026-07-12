package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/Dasongzi1366/AutoGo/utils"
)

var (
	上次吐司内容 string
	上次吐司时间 time.Time
	调试日志锁  sync.Mutex
	调试日志文件 *os.File
	调试日志位置 string
)

func 初始化日志输出() {
	调试日志锁.Lock()
	defer 调试日志锁.Unlock()

	if 调试日志文件 != nil {
		return
	}

	logPaths := []string{}
	for _, dir := range 查找截图目录候选() {
		logPaths = append(logPaths, filepath.Join(dir, "autogo-debug.txt"))
	}
	if exePath, err := os.Executable(); err == nil && exePath != "" {
		logPaths = append(logPaths, filepath.Join(filepath.Dir(exePath), "build", "autogo-debug.log"))
	}
	logPaths = append(logPaths,
		filepath.Join("build", "autogo-debug.log"),
		filepath.Join(os.TempDir(), "autogo-debug.log"),
	)

	for _, path := range logPaths {
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			continue
		}
		file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			continue
		}
		调试日志文件 = file
		调试日志位置 = path
		log.SetOutput(io.MultiWriter(os.Stdout, file))
		_, _ = file.WriteString("\n==== start " + time.Now().Format("2006-01-02 15:04:05") + " ====\n")
		return
	}
}

func 查找截图目录候选() []string {
	candidates := []string{
		"/sdcard/Pictures/Screenshots",
		"/sdcard/DCIM/Screenshots",
		"/storage/emulated/0/Pictures/Screenshots",
		"/storage/emulated/0/DCIM/Screenshots",
	}
	for _, parent := range []string{"/sdcard/Pictures", "/sdcard/DCIM", "/storage/emulated/0/Pictures", "/storage/emulated/0/DCIM"} {
		entries, err := os.ReadDir(parent)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if entry.IsDir() && strings.Contains(strings.ToLower(entry.Name()), "screenshot") {
				candidates = append(candidates, filepath.Join(parent, entry.Name()))
			}
		}
	}

	seen := make(map[string]bool, len(candidates))
	result := make([]string, 0, len(candidates))
	for _, dir := range candidates {
		if seen[dir] {
			continue
		}
		seen[dir] = true
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			result = append(result, dir)
		}
	}
	return result
}

func 调试日志路径() string {
	调试日志锁.Lock()
	defer 调试日志锁.Unlock()
	return 调试日志位置
}

func 输出(v ...any) {
	前缀 := time.Now().Format("15:04:05")
	位置 := 调用位置(1)
	text := strings.TrimSpace(fmt.Sprintln(v...))
	if 应显示到UI日志(text) {
		设置当前动作(text)
		追加UI输出(text)
	}
	line := strings.TrimRight(fmt.Sprintln(append([]any{"[" + 前缀 + "]", 位置}, v...)...), "\r\n")
	fmt.Println(line)
	写入调试日志(line)
}

func 应显示到UI日志(text string) bool {
	return strings.HasPrefix(text, "怪物") ||
		strings.HasPrefix(text, "卖物品测试") ||
		是买卖物品日志(text)
}

func 关键输出(v ...any) {
	输出(v...)
	输出吐司(v...)
}

func 调用位置(skip int) string {
	_, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return "[unknown:0]"
	}
	return fmt.Sprintf("[%s:%d]", filepath.Base(file), line)
}

func 写入调试日志(line string) {
	if line == "" {
		return
	}
	调试日志锁.Lock()
	defer 调试日志锁.Unlock()
	if 调试日志文件 == nil {
		return
	}
	_, _ = 调试日志文件.WriteString(line + "\n")
}

func 输出吐司(v ...any) {
	text := strings.TrimSpace(fmt.Sprintln(v...))
	if text == "" {
		return
	}
	now := time.Now()
	if text == 上次吐司内容 && now.Sub(上次吐司时间) < 2*time.Second {
		return
	}
	上次吐司内容 = text
	上次吐司时间 = now
	utils.Toast(text, 260, 1237, 1800)
}
