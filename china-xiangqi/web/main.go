package main

func main() {
	server := &GameServer{
		Engine: NewGameEngine(),
	}

	// 初始化棋盘
	server.Engine.Board.SetupPieces()

	// 启动服务器
	server.Start()
}
