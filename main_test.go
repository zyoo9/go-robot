package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func Test_main(t *testing.T) {

	msg := "这是第一行\n这是第二行" // 示例字符串

	// 创建数据结构
	data := map[string]interface{}{
		"schema": "2.0",
		"elements": []map[string]interface{}{
			{
				"tag":     "markdown",
				"content": strings.ReplaceAll(msg, "\n", "\\n"),
			},
		},
	}

	// 转换为 JSON 字符串
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	// 打印转义后的 JSON 字符串
	// 使用 string(jsonData) 可以直接输出 JSON 字符串
	escJSON := strings.ReplaceAll(string(jsonData), `"`, `\"`)
	fmt.Println(escJSON)
}
