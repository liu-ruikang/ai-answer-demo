package main

import (
	"ai-answer-demo/china-xiangqi/board"
	"ai-answer-demo/china-xiangqi/game"
	"ai-answer-demo/china-xiangqi/pieces"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	engine := game.NewGameEngine()
	reader := bufio.NewReader(os.Stdin)

	for {
		// 渲染棋盘
		board.RenderBoard(engine.Board)

		// 显示游戏状态
		displayStatus(engine)

		// 检查游戏是否结束
		if engine.IsGameOver {
			fmt.Println("游戏结束！")
			break
		}

		// 获取用户输入
		fmt.Print("请输入移动指令 (例如 e3e5，或输入 'undo' 悔棋, 'exit' 退出): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))

		if input == "exit" {
			break
		} else if input == "undo" {
			if engine.UndoLastMove() {
				fmt.Println("已悔棋一步")
			} else {
				fmt.Println("无法悔棋")
			}
		} else if len(input) >= 4 {
			fromX := int(input[0] - 'a')
			fromY, err1 := strconv.Atoi(string(input[1]))
			toX := int(input[2] - 'a')
			toY, err2 := strconv.Atoi(string(input[3]))

			if fromY > 0 {
				fromY--
			}
			if toY > 0 {
				toY--
			}

			if err1 == nil && err2 == nil &&
				fromX >= 0 && fromX < 9 && fromY >= 0 && fromY < 10 &&
				toX >= 0 && toX < 9 && toY >= 0 && toY < 10 {

				if !engine.MakeMove(fromX, fromY, toX, toY) {
					fmt.Println("无效的移动，请重试")
				}
			} else {
				fmt.Println("输入格式错误，请使用类似 e3e5 的格式")
			}
		} else {
			fmt.Println("输入格式错误，请使用类似 e3e5 的格式")
		}

		// 检查游戏状态
		engine.CheckGameStatus()
	}
}

// displayStatus 显示游戏状态栏
func displayStatus(engine *game.GameEngine) {
	fmt.Println("----------------------------------------")
	if engine.Turn == pieces.Red {
		fmt.Println("当前回合: 红方 (Red)")
	} else {
		fmt.Println("当前回合: 黑方 (Black)")
	}
	fmt.Printf("历史棋步: %d 步\n", len(engine.History))
	fmt.Println("----------------------------------------")
}
