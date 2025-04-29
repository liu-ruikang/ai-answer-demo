package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"strings"
	"time"
)

func main() {
	// 连接本地Server
	mcpClient, err := client.NewStdioMCPClient("./desktop_cleaner", nil)
	if err != nil {
		panic(err)
	}
	defer mcpClient.Close()

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 初始化握手
	fmt.Println("Initializing client...")
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "CleanBot",
		Version: "1.0",
	}
	if _, err := mcpClient.Initialize(ctx, initRequest); err != nil {
		panic(err)
	}

	// 执行工具链
	if err := cleanTempFiles(ctx, mcpClient); err != nil {
		fmt.Printf("清理失败: %v\n", err)
	}
}

func cleanTempFiles(ctx context.Context, client *client.Client) error {
	// 步骤1：扫描.log文件
	scanReq := mcp.CallToolRequest{
		Request: mcp.Request{
			Method: "tools/call",
		},
	}
	scanReq.Params.Name = "scan_temp_files"
	scanReq.Params.Arguments = map[string]interface{}{
		"suffix": ".log",
		"days":   7,
	}

	scanRes, err := client.CallTool(ctx, scanReq)
	if err != nil {
		return err
	}

	printToolResult(scanRes)
	fmt.Println()

	// 解析扫描结果
	files := parseFiles(scanRes)

	// 步骤2：删除文件
	delReq := mcp.CallToolRequest{
		Request: mcp.Request{
			Method: "tools/call",
		},
	}
	delReq.Params.Name = "delete_files"
	delReq.Params.Arguments = map[string]interface{}{
		"paths": files,
	}
	if _, err := client.CallTool(ctx, delReq); err != nil {
		return err
	}

	fmt.Printf("成功清理%d个文件\n", len(files))
	return nil
}

// Helper function to print tool results
func printToolResult(result *mcp.CallToolResult) {
	for _, content := range result.Content {
		if textContent, ok := content.(mcp.TextContent); ok {
			fmt.Println(textContent.Text)
		} else {
			jsonBytes, _ := json.MarshalIndent(content, "", "  ")
			fmt.Println(string(jsonBytes))
		}
	}
}

// 解析calltoolresult
func parseFiles(result *mcp.CallToolResult) (files []string) {
	for _, content := range result.Content {
		if textContent, ok := content.(mcp.TextContent); ok {
			for _, line := range strings.Split(textContent.Text, "\n") {
				if strings.HasPrefix(line, "/") {
					files = append(files, line)
				}
			}
		} else {
			jsonBytes, _ := json.MarshalIndent(content, "", "  ")
			// 解析成[]string
			var result []string
			if err := json.Unmarshal(jsonBytes, &result); err != nil {
				fmt.Println("解析失败:", err)
			}
			files = append(files, result...)
		}
	}
	return
}
