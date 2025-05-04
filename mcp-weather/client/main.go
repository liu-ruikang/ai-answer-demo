package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

const (
	MaxRetries    = 3
	RetryInterval = 2 * time.Second
)

func main() {
	// 连接本地Server
	mcpClient, err := client.NewStdioMCPClient("./weather_client", nil)
	if err != nil {
		panic(err)
	}
	defer mcpClient.Close()

	// 创建超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 初始化握手
	fmt.Println("正在初始化客户端...")
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "WeatherBot",
		Version: "1.0",
	}
	if _, err := mcpClient.Initialize(ctx, initRequest); err != nil {
		panic(err)
	}

	// 执行天气查询
	if err := queryWeather(ctx, mcpClient, "北京市"); err != nil {
		fmt.Printf("天气查询失败: %v\n", err)
	}
}

func queryWeather(ctx context.Context, client *client.Client, city string) error {
	// 构造天气查询请求
	weatherReq := mcp.CallToolRequest{
		Request: mcp.Request{
			Method: "tools/call",
		},
	}
	weatherReq.Params.Name = "amap_weather"
	weatherReq.Params.Arguments = map[string]interface{}{
		"city":       city,
		"key":        "your_amap_api_key",
		"extensions": "base",
	}

	// 调用接口（带重试机制）
	var result *mcp.CallToolResult
	var err error

	for i := 0; i < MaxRetries; i++ {
		result, err = client.CallTool(ctx, weatherReq)
		if err == nil {
			break
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(RetryInterval):
			continue
		}
	}

	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}

	// 打印结果
	printToolResult(result)
	return nil
}

// 打印工具调用结果
func printToolResult(result *mcp.CallToolResult) {
	if result == nil {
		fmt.Println("结果为空")
		return
	}

	if result.IsError {
		if result.Content != nil {
			fmt.Printf("错误信息：%s\n", result.Content)
		} else {
			fmt.Println("未知错误")
		}
		return
	}

	for _, content := range result.Content {
		if textContent, ok := content.(mcp.TextContent); ok {
			fmt.Println(textContent.Text)
		} else {
			jsonBytes, _ := json.MarshalIndent(content, "", "  ")
			fmt.Println(string(jsonBytes))
		}
	}

}
