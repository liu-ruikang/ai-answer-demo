package main

import (
	"context"
	"fmt"
	"os"

	aiagent "ai-answer-demo/ai-agent/ai-agent"
)

// usage 显示用法说明
func usage() {
	fmt.Println("Usage: agent-cli [command]")
	fmt.Println("Available commands:")
	fmt.Println("  run-agent       启动完整的 Agent 示例")
	fmt.Println("  run-encourager  启动程序员鼓励师示例")
	fmt.Println("  help            显示帮助信息")
}

func main() {
	// 检查命令行参数
	if len(os.Args) < 2 {
		fmt.Println("请指定命令")
		usage()
		os.Exit(1)
	}

	ctx := context.Background()

	switch os.Args[1] {
	case "run-agent":
		aiagent.RunAgent(ctx)
	case "run-encourager":
		aiagent.RunEncourager(ctx)
	case "help":
		usage()
	default:
		fmt.Println("未知命令：", os.Args[1])
		usage()
		os.Exit(1)
	}
}
