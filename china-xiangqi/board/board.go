package board

import (
	"china-xiangqi/pieces"
)

// ChessBoard 表示中国象棋棋盘
type ChessBoard struct {
	Grid [9][10]*pieces.Piece // 9列10行的棋盘网格
}

// GetPiece 实现pieces.Board接口
func (b *ChessBoard) GetPiece(x, y int) *pieces.Piece {
	if x < 0 || x >= 9 || y < 0 || y >= 10 {
		return nil
	}
	return b.Grid[x][y]
}

// NewChessBoard 初始化新棋盘
func NewChessBoard() *ChessBoard {
	board := &ChessBoard{}
	// 初始化棋子位置
	board.SetupPieces()
	return board
}

// SetupPieces 设置初始棋子位置
func (b *ChessBoard) SetupPieces() {
	// 红方棋子
	b.Grid[0][0] = &pieces.Piece{Type: pieces.Rook, Color: pieces.Red}
	b.Grid[1][0] = &pieces.Piece{Type: pieces.Knight, Color: pieces.Red}
	b.Grid[2][0] = &pieces.Piece{Type: pieces.Bishop, Color: pieces.Red}
	b.Grid[3][0] = &pieces.Piece{Type: pieces.Advisor, Color: pieces.Red}
	b.Grid[4][0] = &pieces.Piece{Type: pieces.King, Color: pieces.Red}
	b.Grid[5][0] = &pieces.Piece{Type: pieces.Advisor, Color: pieces.Red}
	b.Grid[6][0] = &pieces.Piece{Type: pieces.Bishop, Color: pieces.Red}
	b.Grid[7][0] = &pieces.Piece{Type: pieces.Knight, Color: pieces.Red}
	b.Grid[8][0] = &pieces.Piece{Type: pieces.Rook, Color: pieces.Red}

	for i := 0; i < 9; i += 2 {
		b.Grid[i][2] = &pieces.Piece{Type: pieces.Cannon, Color: pieces.Red}
	}

	for i := 0; i < 9; i++ {
		b.Grid[i][3] = &pieces.Piece{Type: pieces.Soldier, Color: pieces.Red}
	}

	// 黑方棋子
	b.Grid[0][9] = &pieces.Piece{Type: pieces.Rook, Color: pieces.Black}
	b.Grid[1][9] = &pieces.Piece{Type: pieces.Knight, Color: pieces.Black}
	b.Grid[2][9] = &pieces.Piece{Type: pieces.Bishop, Color: pieces.Black}
	b.Grid[3][9] = &pieces.Piece{Type: pieces.Advisor, Color: pieces.Black}
	b.Grid[4][9] = &pieces.Piece{Type: pieces.King, Color: pieces.Black}
	b.Grid[5][9] = &pieces.Piece{Type: pieces.Advisor, Color: pieces.Black}
	b.Grid[6][9] = &pieces.Piece{Type: pieces.Bishop, Color: pieces.Black}
	b.Grid[7][9] = &pieces.Piece{Type: pieces.Knight, Color: pieces.Black}
	b.Grid[8][9] = &pieces.Piece{Type: pieces.Rook, Color: pieces.Black}

	for i := 0; i < 9; i += 2 {
		b.Grid[i][7] = &pieces.Piece{Type: pieces.Cannon, Color: pieces.Black}
	}

	for i := 0; i < 9; i++ {
		b.Grid[i][6] = &pieces.Piece{Type: pieces.Soldier, Color: pieces.Black}
	}
}

// MovePiece 移动棋子
func (b *ChessBoard) MovePiece(fromX, fromY, toX, toY int) bool {
	// 检查坐标有效性
	if !isValidPosition(fromX, fromY) || !isValidPosition(toX, toY) {
		return false
	}

	piece := b.Grid[fromX][fromY]
	if piece == nil {
		return false
	}

	// 验证移动是否合法
	if !piece.ValidMove(b, fromX, fromY, toX, toY) {
		return false
	}

	// 执行移动
	b.Grid[toX][toY] = piece
	b.Grid[fromX][fromY] = nil
	return true
}

// isValidPosition 检查坐标是否在棋盘范围内
func isValidPosition(x, y int) bool {
	return x >= 0 && x < 9 && y >= 0 && y < 10
}
