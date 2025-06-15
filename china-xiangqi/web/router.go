package main

import (
	"china-xiangqi/board"
	"china-xiangqi/game"
	"china-xiangqi/pieces"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

// GameServer 表示Web游戏服务器
type GameServer struct {
	Engine *game.GameEngine
}

// NewGameEngine 创建新的游戏引擎实例
func NewGameEngine() *game.GameEngine {
	return &game.GameEngine{
		Board:   board.NewChessBoard(),
		Turn:    pieces.Red,
		History: make([][2][2]int, 0),
	}
}

// Start 启动Web服务器
func (s *GameServer) Start() {
	http.HandleFunc("/", s.homeHandler)
	http.HandleFunc("/move", s.moveHandler)

	println("启动中国象棋Web服务器，访问 http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// homeHandler 首页处理器
func (s *GameServer) homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	template := template.Must(template.ParseFiles("china-xiangqi/web/template.html"))

	data := s.getBoardData()
	error := template.ExecuteTemplate(w, "template.html", data)
	if error != nil {
		http.Error(w, error.Error(), http.StatusInternalServerError)
	}
}

// moveHandler 移动处理器
func (s *GameServer) moveHandler(w http.ResponseWriter, r *http.Request) {
	fromStr := r.FormValue("from")
	toStr := r.FormValue("to")

	if fromStr == "" || toStr == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// 解析坐标
	from := strings.Split(fromStr, ",")
	to := strings.Split(toStr, ",")

	fromX, _ := strconv.Atoi(from[0])
	fromY, _ := strconv.Atoi(from[1])
	toX, _ := strconv.Atoi(to[0])
	toY, _ := strconv.Atoi(to[1])

	// 执行移动
	if !s.Engine.MakeMove(fromX, fromY, toX, toY) {
		// 记录错误日志或显示错误信息
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// getBoardData 获取棋盘数据
func (s *GameServer) getBoardData() *board.ChessBoard {
	data := &board.ChessBoard{}

	// 棋子名称映射
	// colorMap := map[int]string{0: "red", 1: "black"}
	// pieceNames := map[int]string{
	// 	0: "将",
	// 	1: "士",
	// 	2: "相",
	// 	3: "马",
	// 	4: "車",
	// 	5: "炮",
	// 	6: "兵",
	// }

	for x := 0; x < 9; x++ {
		for y := 0; y < 10; y++ {
			piece := s.Engine.Board.Grid[x][y]
			if piece != nil {
				data.Grid[x][y] = &pieces.Piece{
					Type:  piece.Type,
					Color: piece.Color,
				}
			}
		}
	}

	// data.Turn = int(s.Engine.Turn)
	return data
}
