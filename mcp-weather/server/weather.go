package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// WeatherResponse - 天气数据结构体，用于封装返回结果
type WeatherResponse struct {
	City       string `json:"city"`
	Weather    string `json:"weather"`
	Temp       string `json:"temperature"`
	Humidity   string `json:"humidity"`
	ReportTime string `json:"report_time"`
}

// RegisterWeatherTool - 注册天气查询工具
// 设置工具描述和参数要求，并绑定处理函数
func RegisterWeatherTool(s *server.MCPServer) {
	tool := mcp.NewTool(
		"weather", // 工具名称
		mcp.WithDescription("获取城市实时天气数据（高德地图API）"),                             // 描述信息
		mcp.WithString("city", mcp.Required(), mcp.Description("城市名称（如：北京市）")), // 必填参数 city
	)

	// 将工具添加到服务器中，并指定对应的处理函数
	s.AddTool(tool, weatherHandler)
}

// weatherHandler - 天气查询请求的实际处理函数
func weatherHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 从请求参数中提取城市名
	city := req.Params.Arguments["city"].(string)

	// 构造高德天气 API 请求 URL
	url := fmt.Sprintf("https://restapi.amap.com/v3/weather/weatherInfo?city=%s&key=%s",
		city, os.Getenv("AMAP_API_KEY")) // 使用环境变量中的 API Key

	// 发起 HTTP GET 请求
	resp, err := http.Get(url)
	if err != nil {
		return mcp.NewToolResultError("API请求失败"), nil
	}
	defer resp.Body.Close()

	// 解析 JSON 响应
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// 数据解析部分
	// 先判断 'lives' 是否存在且是数组类型，并且长度 > 0
	if livesArray, ok := result["lives"].([]interface{}); ok && len(livesArray) > 0 {
		lives := livesArray[0].(map[string]interface{})
		weather := WeatherResponse{
			City:       lives["city"].(string),
			Weather:    lives["weather"].(string),
			Temp:       lives["temperature"].(string) + "℃",
			Humidity:   lives["humidity"].(string) + "%",
			ReportTime: time.Now().Format("2006-01-02 15:04"), // 当前时间格式化
		}

		// 构造文本响应内容
		text := fmt.Sprintf("城市: %s\n天气: %s\n温度: %s\n湿度: %s\n更新时间: %s",
			weather.City, weather.Weather, weather.Temp, weather.Humidity, weather.ReportTime)
		return mcp.NewToolResultText(text), nil
	} else {
		// 如果没有获取到有效数据，返回错误提示
		return mcp.NewToolResultError("无法获取天气数据，请稍后再试"), nil
	}
}
