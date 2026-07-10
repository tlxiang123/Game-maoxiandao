package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"time"
)

var 发现测谎 = &FMColor{Name: "发现测谎", X1: 298, Y1: 255, X2: 470, Y2: 368, MainColor: "0044BB-000000", OffsetColor: "8,0,0044BB-000000,16,0,0044BB-000000,0,3,0044BB-000000,5,5,1155BB-000000,16,5,1155BB-000000,4,8,0044BB-000000,5,6,0044BB-000000,13,6,0044BB-000000", Sim: 0.90, Dir: 0}

var 默认钉钉Webhook列表 = []string{
	"https://oapi.dingtalk.com/robot/send?access_token=6c833071fed0c97bfd0fcc1933e5239d6aa139015342e51efc71ca12445305d5",
	"https://oapi.dingtalk.com/robot/send?access_token=7c1e139f47536ae6709c68ccf680808a7712c9e97c5c5a63638a22ff614e521a",
}

var (
	测谎检测已启动 atomic.Bool
	测谎发送计数  atomic.Int64
)

func 启动测谎检测后台() {
	if !测谎检测已启动.CompareAndSwap(false, true) {
		return
	}
	go 测谎检测循环()
}

func 测谎检测循环() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	上次已发现 := false
	for !程序退出中.Load() {
		<-ticker.C
		if 引擎 == nil {
			continue
		}
		found, name, x, y := 查找测谎钉钉触发特征()
		if found {
			if !上次已发现 {
				count := 测谎发送计数.Add(1)
				输出("发现测谎触发", "count=", count, "特征=", name, "x=", x, "y=", y)
				发送测谎钉钉消息(count, name)
			}
			上次已发现 = true
			continue
		}
		上次已发现 = false
	}
}

func 查找测谎钉钉触发特征() (bool, string, int, int) {
	if found, x, y := 引擎.FindFeature(发现测谎); found {
		return true, 发现测谎.Name, x, y
	}
	if found, x, y := 引擎.FindFeature(MS系统应用); found {
		return true, MS系统应用.Name, x, y
	}
	return false, "", -1, -1
}

func 测试发送钉钉() {
	count := 测谎发送计数.Add(1)
	输出("测试钉钉发送", "count=", count)
	发送测谎钉钉消息(count, "测试")
}

func 发送测谎钉钉消息(count int64, featureName string) bool {
	return 发送钉钉文本(fmt.Sprintf("冒险岛发现测谎 %d：%s", count, featureName))
}

func 发送钉钉文本(content string) bool {
	payload := map[string]any{
		"msgtype": "text",
		"text": map[string]string{
			"content": content,
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		输出("钉钉发送失败：JSON错误", err)
		return false
	}

	client := &http.Client{Timeout: 5 * time.Second}
	allOK := true
	for index, webhook := range 默认钉钉Webhook列表 {
		resp, err := client.Post(webhook, "application/json", bytes.NewReader(body))
		if err != nil {
			输出("钉钉发送失败：请求错误", "序号=", index+1, err)
			allOK = false
			continue
		}

		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		_ = resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			输出("钉钉发送失败：HTTP", "序号=", index+1, "code=", resp.StatusCode, string(respBody))
			allOK = false
			continue
		}
		输出("钉钉发送成功", "序号=", index+1, "内容=", content)
	}
	return allOK
}
