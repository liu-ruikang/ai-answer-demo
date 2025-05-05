package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

// 常量定义
const (
	MaxRetries    = 3               // 最大重试次数
	RetryInterval = 2 * time.Second // 每次重试间隔时间
)

func main() {
	// 创建一个基于标准输入输出的 MCP 客户端，连接本地 Server
	mcpClient, err := client.NewStdioMCPClient("./weather_client", nil)
	if err != nil {
		panic(err)
	}
	defer mcpClient.Close()

	// 创建一个带有超时的上下文（10秒）
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 初始化握手请求
	fmt.Println("正在初始化客户端...")
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "WeatherBot", // 客户端名称
		Version: "1.0",        // 客户端版本
	}

	// 发送初始化请求
	if _, err := mcpClient.Initialize(ctx, initRequest); err != nil {
		panic(err)
	}

	// 执行天气查询（示例：北京市）
	if err := queryWeather(ctx, mcpClient, "北京市"); err != nil {
		fmt.Printf("天气查询失败: %v\n", err)
	}
}

// queryWeather - 调用天气查询工具
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
		"key":        os.Getenv("AMAP_API_KEY"), // 从环境变量中读取高德 API Key
		"extensions": "base",
	}

	var result *mcp.CallToolResult
	var err error

	// 带有重试机制的调用
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

// printToolResult - 打印工具调用结果
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
