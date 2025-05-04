package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer("amap_weather", "1.0",
		server.WithLogging(),
		server.WithRecovery(),
	)

	RegisterWeatherTool(s)
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

type WeatherResponse struct {
	City       string `json:"city"`
	Weather    string `json:"weather"`
	Temp       string `json:"temperature"`
	Humidity   string `json:"humidity"`
	ReportTime string `json:"report_time"`
}

func RegisterWeatherTool(s *server.MCPServer) {
	tool := mcp.NewTool(
		"weather",
		mcp.WithDescription("获取城市实时天气数据（高德地图API）"),
		mcp.WithString("city", mcp.Required(), mcp.Description("城市名称（如：北京市）")),
	)

	s.AddTool(tool, weatherHandler)
}

func weatherHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	city := req.Params.Arguments["city"].(string)

	// 调用高德API
	url := fmt.Sprintf("https://restapi.amap.com/v3/weather/weatherInfo?city=%s&key=%s",
		city, "your_api_key_here")

	resp, err := http.Get(url)
	if err != nil {
		return mcp.NewToolResultError("API请求失败"), nil
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// 数据解析
	if livesArray, ok := result["lives"].([]interface{}); ok && len(livesArray) > 0 {
		lives := livesArray[0].(map[string]interface{})
		weather := WeatherResponse{
			City:       lives["city"].(string),
			Weather:    lives["weather"].(string),
			Temp:       lives["temperature"].(string) + "℃",
			Humidity:   lives["humidity"].(string) + "%",
			ReportTime: time.Now().Format("2006-01-02 15:04"),
		}

		// 把weather转换为*mcp.CallToolResult结构
		text := fmt.Sprintf("城市: %s\n天气: %s\n温度: %s\n湿度: %s\n更新时间: %s",
			weather.City, weather.Weather, weather.Temp, weather.Humidity, weather.ReportTime)
		return mcp.NewToolResultText(text), nil
	} else {
		return mcp.NewToolResultError("无法获取天气数据，请稍后再试"), nil
	}
}
