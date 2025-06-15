package board

import (
	"china-xiangqi/pieces"
	"fmt"
)

// RenderBoard 渲染棋盘到控制台
func RenderBoard(b *ChessBoard) {
	// 清屏
	fmt.Print("\033[H\033[2J")

	// 打印列标号
	fmt.Print("   ")
	for x := 0; x < 9; x++ {
		fmt.Printf(" %c ", 'a'+x)
	}
	fmt.Println()

	// 打印上边框
	fmt.Print("  ┌─")
	for x := 1; x < 9; x++ {
		fmt.Print("┬─")
	}
	fmt.Println("┐")

	// 打印棋盘内容
	for y := 0; y < 10; y++ {
		// 打印行号
		if y == 9 {
			fmt.Printf("%d │", 0)
		} else {
			fmt.Printf("%d │", y+1)
		}

		// 打印棋子
		for x := 0; x < 9; x++ {
			piece := b.Grid[x][y]
			if piece == nil {
				fmt.Print("  │")
			} else {
				// 根据颜色设置ANSI代码
				colorCode := ""
				resetCode := "\033[0m"
				if piece.Color == pieces.Red {
					colorCode = "\033[31m" // 红色
				} else {
					colorCode = "\033[34m" // 蓝色表示黑方
				}

				// 根据棋子类型选择Unicode符号
				symbol := ""
				switch piece.Type {
				case pieces.King:
					symbol = "♚"
				case pieces.Advisor:
					symbol = "♛"
				case pieces.Bishop:
					symbol = "♝"
				case pieces.Knight:
					symbol = "♞"
				case pieces.Rook:
					symbol = "♜"
				case pieces.Cannon:
					symbol = "炮"
				case pieces.Soldier:
					symbol = "♟"
				}

				fmt.Printf("%s%s%s│", colorCode, symbol, resetCode)
			}
		}
		fmt.Println()

		// 打印分隔线
		if y < 9 {
			fmt.Print("  ├─")
			for x := 1; x < 9; x++ {
				fmt.Print("┼─")
			}
			fmt.Println("┤")
		}
	}

	// 打印下边框
	fmt.Print("  └─")
	for x := 1; x < 9; x++ {
		fmt.Print("┴─")
	}
	fmt.Println("┘")
}
