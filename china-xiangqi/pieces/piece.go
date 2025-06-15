package pieces

import "fmt"

// Board 接口定义了棋盘的基本功能
type Board interface {
	GetPiece(x, y int) *Piece
}

// Color 表示棋子颜色（阵营）
type Color int

const (
	Red Color = iota
	Black
)

// PieceType 表示棋子类型
type PieceType int

const (
	King    PieceType = iota // 将/帅
	Advisor                  // 士
	Bishop                   // 象
	Knight                   // 马
	Rook                     // 车
	Cannon                   // 炮
	Soldier                  // 兵/卒
)

// Piece 表示棋子
type Piece struct {
	Type  PieceType
	Color Color
}

// ValidMove 验证棋子移动是否合法
func (p *Piece) ValidMove(b Board, fromX, fromY, toX, toY int) bool {
	switch p.Type {
	case King:
		return validKingMove(b, fromX, fromY, toX, toY)
	case Advisor:
		return validAdvisorMove(b, fromX, fromY, toX, toY)
	case Bishop:
		return validBishopMove(b, fromX, fromY, toX, toY)
	case Knight:
		return validKnightMove(b, fromX, fromY, toX, toY)
	case Rook:
		return validRookMove(b, fromX, fromY, toX, toY)
	case Cannon:
		return validCannonMove(b, fromX, fromY, toX, toY)
	case Soldier:
		return validSoldierMove(b, fromX, fromY, toX, toY)
	}
	fmt.Println("未知的棋子类型")
	return false
}

// validKingMove 验证将/帅移动
func validKingMove(b Board, fromX, fromY, toX, toY int) bool {
	// 实现在九宫格内的移动
	color := b.GetPiece(fromX, fromY).Color
	if color == Red && (toX < 3 || toX > 5 || toY < 0 || toY > 2) {
		return false
	}
	if color == Black && (toX < 3 || toX > 5 || toY < 7 || toY > 9) {
		return false
	}

	// 只能移动一步
	if abs(fromX-toX)+abs(fromY-toY) > 1 {
		return false
	}

	// 检查是否同线直面对
	if fromX == toX {
		for y := min(fromY, toY) + 1; y < max(fromY, toY); y++ {
			if b.GetPiece(toX, y) != nil {
				return false
			}
		}
	}
	return true
}

// validAdvisorMove 验证士的移动
func validAdvisorMove(b Board, fromX, fromY, toX, toY int) bool {
	// 实现在九宫格内的斜线移动
	color := b.GetPiece(fromX, fromY).Color
	if color == Red && (toX < 3 || toX > 5 || toY < 0 || toY > 2) {
		return false
	}
	if color == Black && (toX < 3 || toX > 5 || toY < 7 || toY > 9) {
		return false
	}

	// 必须斜线移动一步
	return abs(fromX-toX) == 1 && abs(fromY-toY) == 1
}

// validBishopMove 验证象的移动
func validBishopMove(b Board, fromX, fromY, toX, toY int) bool {
	// 飞田字且不能过河
	color := b.GetPiece(fromX, fromY).Color
	if color == Red && toY > 4 {
		return false // 红方不能过河
	}
	if color == Black && toY < 5 {
		return false // 黑方不能过河
	}

	// 检查是否田字形移动
	if abs(fromX-toX) != 2 || abs(fromY-toY) != 2 {
		return false
	}

	// 检查田字中心是否有阻挡
	midX := (fromX + toX) / 2
	midY := (fromY + toY) / 2
	return b.GetPiece(midX, midY) == nil
}

// validKnightMove 验证马的移动
func validKnightMove(b Board, fromX, fromY, toX, toY int) bool {
	dx := abs(fromX - toX)
	dy := abs(fromY - toY)

	// 检查是否日字形移动
	if !(dx == 1 && dy == 2 || dx == 2 && dy == 1) {
		return false
	}

	// 检查蹩脚
	if dx == 1 {
		// 横向移动，检查纵向是否有阻挡
		midY := (fromY + toY) / 2
		return b.GetPiece(fromX, midY) == nil
	} else {
		// 纵向移动，检查横向是否有阻挡
		midX := (fromX + toX) / 2
		return b.GetPiece(midX, fromY) == nil
	}
}

// validRookMove 验证车的移动
func validRookMove(b Board, fromX, fromY, toX, toY int) bool {
	// 必须直线移动
	if fromX != toX && fromY != toY {
		return false
	}

	// 检查路径上是否有阻挡
	if fromX == toX {
		// 垂直移动
		for y := min(fromY, toY) + 1; y < max(fromY, toY); y++ {
			if b.GetPiece(toX, y) != nil {
				return false
			}
		}
	} else {
		// 水平移动
		for x := min(fromX, toX) + 1; x < max(fromX, toX); x++ {
			if b.GetPiece(x, toY) != nil {
				return false
			}
		}
	}
	return true
}

// validCannonMove 验证炮的移动
func validCannonMove(b Board, fromX, fromY, toX, toY int) bool {
	target := b.GetPiece(toX, toY)

	// 移动规则与车相同，但吃子时需要隔山
	if fromX != toX && fromY != toY {
		return false
	}

	// 计算路径上的棋子数量
	count := 0
	if fromX == toX {
		// 垂直移动
		for y := min(fromY, toY) + 1; y < max(fromY, toY); y++ {
			if b.GetPiece(toX, y) != nil {
				count++
			}
		}
	} else {
		// 水平移动
		for x := min(fromX, toX) + 1; x < max(fromX, toX); x++ {
			if b.GetPiece(x, toY) != nil {
				count++
			}
		}
	}

	// 没有吃子时必须路径畅通
	if target == nil {
		return count == 0
	}
	// 吃子时必须有一个阻挡物
	return count == 1
}

// validSoldierMove 验证兵/卒的移动
func validSoldierMove(b Board, fromX, fromY, toX, toY int) bool {
	piece := b.GetPiece(fromX, fromY)

	// 检查移动方向
	if piece.Color == Red {
		// 红方只能向前或左右移动
		if toY-fromY < 0 || abs(toX-fromX) > 1 || toY-fromY > 1 {
			return false
		}

		// 过河前只能直行
		if fromY < 5 && toX != fromX {
			return false
		}
	} else {
		// 黑方只能向后或左右移动
		if fromY-toY < 0 || abs(toX-fromX) > 1 || fromY-toY > 1 {
			return false
		}

		// 过河前只能直行
		if fromY > 4 && toX != fromX {
			return false
		}
	}

	// 目标位置不能有同色棋子
	targetPiece := b.GetPiece(toX, toY)
	if targetPiece != nil && targetPiece.Color == piece.Color {
		return false
	}

	return true
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
