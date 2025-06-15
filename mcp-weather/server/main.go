package main

import (
	"fmt"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// 创建一个新的 MCP 服务实例
	// 参数：工具名称、版本号
	s := server.NewMCPServer("amap_weather", "1.0",
		server.WithLogging(),  // 启用日志功能
		server.WithRecovery(), // 启用 panic 恢复机制
	)

	// 注册天气查询工具（定义在 weather.go 中）
	RegisterWeatherTool(s)

	// 启动基于标准输入输出的 MCP 服务
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
