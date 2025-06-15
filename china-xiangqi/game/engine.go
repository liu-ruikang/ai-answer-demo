package game

import (
	"china-xiangqi/board"
	"china-xiangqi/pieces"
)

// GameEngine 表示游戏引擎
type GameEngine struct {
	Board      *board.ChessBoard
	Turn       pieces.Color // 当前回合方
	History    [][2][2]int  // 棋步历史记录
	IsGameOver bool         // 游戏是否结束
}

// NewGameEngine 创建新的游戏引擎实例
func NewGameEngine() *GameEngine {
	return &GameEngine{
		Board:   board.NewChessBoard(),
		Turn:    pieces.Red,
		History: make([][2][2]int, 0),
	}
}

// MakeMove 执行一次移动
func (g *GameEngine) MakeMove(fromX, fromY, toX, toY int) bool {
	if g.IsGameOver {
		return false
	}

	// 检查坐标有效性
	if !isValidPosition(fromX, fromY) || !isValidPosition(toX, toY) {
		return false
	}

	piece := g.Board.Grid[fromX][fromY]
	if piece == nil || piece.Color != g.Turn {
		return false
	}

	// 验证并执行移动
	if !piece.ValidMove(g.Board, fromX, fromY, toX, toY) {
		return false
	}

	// 记录棋步
	g.History = append(g.History, [2][2]int{{fromX, fromY}, {toX, toY}})

	// 检查是否将军
	if g.isCheckAfterMove(toX, toY) {
		// TODO: 实现将军提示和响应逻辑
	}

	// 切换回合
	g.Turn = 1 - g.Turn
	return true
}

// isCheckAfterMove 检查移动后是否将军
func (g *GameEngine) isCheckAfterMove(kingX, kingY int) bool {
	// TODO: 实现将军检测逻辑
	return false
}

// UndoLastMove 悔棋
func (g *GameEngine) UndoLastMove() bool {
	if len(g.History) == 0 {
		return false
	}

	// 获取最后一步
	lastMove := g.History[len(g.History)-1]
	from := lastMove[0]
	to := lastMove[1]

	// 移动回去
	piece := g.Board.Grid[to[0]][to[1]]
	g.Board.Grid[from[0]][from[1]] = piece
	g.Board.Grid[to[0]][to[1]] = nil

	// 更新历史记录
	g.History = g.History[:len(g.History)-1]

	// 切换回合
	g.Turn = 1 - g.Turn

	return true
}

// CheckGameStatus 检查游戏状态
func (g *GameEngine) CheckGameStatus() {
	// TODO: 实现将死、困毙等终局条件检测
}

// isValidPosition 检查坐标是否在棋盘范围内
func isValidPosition(x, y int) bool {
	return x >= 0 && x < 9 && y >= 0 && y < 10
}
