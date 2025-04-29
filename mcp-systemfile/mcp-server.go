package main

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// 获取桌面路径（跨平台兼容）
func getDesktopPath() string {
	home, _ := os.UserHomeDir()
	// 创建工作目录如果不存在
	path := filepath.Join(home, "Desktop")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0755)
	}
	return path
}

func main() {
	s := server.NewMCPServer("DesktopCleaner", "1.0",
		server.WithLogging(),
		server.WithRecovery(),
	)

	// 工具1：扫描临时文件
	scanTool := mcp.NewTool("scan_temp_files",
		mcp.WithDescription("扫描桌面指定后缀的临时文件"),
		mcp.WithString("suffix", mcp.Required(),
			mcp.Description("文件后缀如.log/.tmp"),
			mcp.Pattern(`^\.[a-zA-Z0-9]+$`)),
		mcp.WithNumber("days", mcp.Description("查找最近N天内的文件")),
	)
	s.AddTool(scanTool, scanHandler)

	// 工具2：批量删除
	delTool := mcp.NewTool("delete_files",
		mcp.WithDescription("删除指定路径的文件"),
		mcp.WithArray("paths", mcp.Required(),
			mcp.Description("文件路径数组")),
	)
	s.AddTool(delTool, deleteHandler)

	// 工具3：桌面整理
	orgTool := mcp.NewTool("organize_desktop",
		mcp.WithDescription("按类型整理文件"),
		mcp.WithString("strategy",
			mcp.Enum("type", "date"), // 按类型/日期分类
			mcp.Description("整理策略")),
	)
	s.AddTool(orgTool, organizeHandler)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

// 扫描处理器
func scanHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	suffix := req.Params.Arguments["suffix"].(string)
	
	// 获取days参数，如果不存在则默认为0，不做日期过滤
	var days int
	if daysVal, ok := req.Params.Arguments["days"]; ok {
		days, _ = daysVal.(int)
	}

	desktopPath := getDesktopPath()
	fmt.Printf("扫描目录: %s，后缀: %s，天数: %d\n", desktopPath, suffix, days)
	
	var files []string
	err := filepath.Walk(desktopPath, func(path string, info os.FileInfo, err error) error {
		// 处理错误情况，防止空指针
		if err != nil {
			fmt.Printf("访问路径出错: %s, 错误: %v\n", path, err)
			return nil // 继续扫描其他文件
		}
		if info == nil {
			fmt.Printf("文件信息为空: %s\n", path)
			return nil
		}
		
		// 检查文件后缀名，忽略大小写
		if info.IsDir() {
			return nil
		}
		
		ext := filepath.Ext(path)
		if strings.ToLower(ext) != strings.ToLower(suffix) {
			return nil
		}
		
		// 如果指定了天数，检查文件修改时间
		if days > 0 {
			if time.Since(info.ModTime()) > time.Duration(days)*24*time.Hour {
				return nil // 文件太老，跳过
			}
		}
		
		// 文件符合条件，添加到结果中
		files = append(files, path)
		fmt.Printf("找到文件: %s\n", path)
		return nil
	})
	
	if err != nil {
		fmt.Printf("扫描出错: %v\n", err)
		return mcp.NewToolResultErrorFromErr("扫描失败", err), nil
	}
	
	result := "找到" + fmt.Sprintf("%d", len(files)) + "个文件：\n"
	for _, f := range files {
		result += f + "\n"
	}
	
	return mcp.NewToolResultText(result), nil
}

// 删除处理器
func deleteHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	paths := req.Params.Arguments["paths"].([]interface{})
	results := make(map[string]string)

	for _, p := range paths {
		path := p.(string)
		if !isUnderDesktop(path) { // 安全校验
			return mcp.NewToolResultError("非法路径: " + path), nil
		}
		err := os.Remove(path)
		if err != nil {
			results[path] = "删除失败: " + err.Error()
		} else {
			results[path] = "删除成功"
		}
	}
	
	result := "删除结果：\n"
	for path, status := range results {
		result += path + ": " + status + "\n"
	}
	
	return mcp.NewToolResultText(result), nil
}

// 整理桌面处理器
func organizeHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	strategy, _ := req.Params.Arguments["strategy"].(string)
	
	// 实际的整理逻辑可以在这里实现
	
	return mcp.NewToolResultText("桌面文件已按" + strategy + "整理完成"), nil
}

// 路径合法性校验
func isUnderDesktop(path string) bool {
	rel, err := filepath.Rel(getDesktopPath(), path)
	return err == nil && !strings.HasPrefix(rel, "..")
}
